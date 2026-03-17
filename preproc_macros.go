// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"fmt"
	"maps"
	"sort"
	"strings"
)

// defineMacro parses and stores object-like or function-like macro.
func (p *preprocessor) defineMacro(raw string, filename string, lineNo int) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		p.emitError(CodePPMissingMacroName, filename, lineNo, "missing macro name in #define")

		return ErrInvalidDirective
	}

	nameEnd := 0
	for nameEnd < len(raw) && isIdentifierPart(raw[nameEnd]) {
		nameEnd++
	}

	if nameEnd == 0 {
		p.emitError(CodePPInvalidMacroName, filename, lineNo, "invalid macro name in #define")

		return ErrInvalidDirective
	}

	name := raw[:nameEnd]
	rest := raw[nameEnd:]

	existing, exists := p.macros[name]
	if p.enableMacroRedefWarn && exists && existing.FunctionLike {
		p.emitMacroRedefinedWarning(filename, lineNo, "function-like", name)
	}

	if p.enableMacroRedefWarn && exists && !existing.FunctionLike {
		p.emitMacroRedefinedWarning(filename, lineNo, "object-like", name)
	}

	if strings.HasPrefix(rest, "(") {
		closeIdx := strings.Index(rest, ")")
		if closeIdx < 0 {
			p.emitError(CodePPUnterminatedMacroParams, filename, lineNo, "unterminated parameter list in #define")

			return ErrInvalidDirective
		}

		paramText := strings.TrimSpace(rest[1:closeIdx])
		// Keep continuation indentation and trailing separator spaces for
		// closer CfgConvert -pcpp parity.
		// Drop only single delimiter whitespace after ")" in one-line macros.
		body := rest[closeIdx+1:]
		body = trimSingleMacroBodyDelimiter(body)
		params := parseParams(paramText)

		p.macros[name] = macroDefinition{
			Name:         name,
			Params:       params,
			Body:         body,
			FunctionLike: true,
		}
		p.macroNamesDirty = true

		return nil
	}

	body := strings.TrimSpace(rest)
	p.macros[name] = macroDefinition{
		Name: name,
		Body: body,
	}
	p.macroNamesDirty = true

	return nil
}

// trimSingleMacroBodyDelimiter removes exactly one whitespace delimiter between
// function-like macro parameter list and body token.
func trimSingleMacroBodyDelimiter(body string) string {
	if len(body) < 2 {
		return body
	}

	first := body[0]
	second := body[1]
	if (first == ' ' || first == '\t') && second != ' ' && second != '\t' {
		return body[1:]
	}

	return body
}

// emitMacroRedefinedWarning emits one warning per file+kind+name redefinition key.
func (p *preprocessor) emitMacroRedefinedWarning(filename string, lineNo int, kind string, name string) {
	key := filename + "\x00" + kind + "\x00" + name
	if _, exists := p.macroRedefWarnedV0[key]; exists {
		return
	}

	p.macroRedefWarnedV0[key] = struct{}{}
	p.emitWarning(CodePPMacroRedefined, filename, lineNo, fmt.Sprintf("redefining %s macro %s", kind, name))
}

// expandLine expands function-like and object-like macros.
func (p *preprocessor) expandLine(line string) (string, error) {
	hasFunctionMacros := len(p.functionMacroNames()) > 0
	if hasFunctionMacros {
		savedMacros := cloneMacroDefinitions(p.macros)
		savedDirty := p.macroNamesDirty
		savedObjectNames := append([]string(nil), p.objectMacroNamesV0...)
		savedFunctionNames := append([]string(nil), p.functionMacroNamesV0...)
		defer func() {
			p.macros = savedMacros
			p.macroNamesDirty = savedDirty
			p.objectMacroNamesV0 = savedObjectNames
			p.functionMacroNamesV0 = savedFunctionNames
		}()
	}

	result := line

	for pass := 0; pass < p.maxExpandDepth; pass++ {
		changed := false

		next := result
		changedFunc := false
		if hasFunctionMacros {
			next, changedFunc = p.expandFunctionMacros(next)
		}

		next, changedObj := p.expandObjectMacros(next)
		next, changedStringify := collapseStringifyTokens(next)
		changed = changedFunc || changedObj || changedStringify
		result = next

		if !changed {
			result = collapseTokenPaste(result)

			return result, nil
		}
	}

	return result, fmt.Errorf("%w: expansion depth overflow", ErrMacroExpand)
}

// collapseStringifyTokens converts standalone #IDENT tokens into quoted strings.
// This keeps DayZ-compatible behavior for object-like macros such as:
//
//	#define MACRO#test
//	MACRO -> "test"
func collapseStringifyTokens(input string) (string, bool) {
	if input == "" || !strings.Contains(input, "#") {
		return input, false
	}

	var out strings.Builder
	lastWrite := 0
	replaced := false

	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); {
		if inLineComment {
			if input[i] == '\n' {
				inLineComment = false
			}

			i++

			continue
		}

		if inBlockComment {
			if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			if input[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if input[i] == '"' {
			inString = true
			i++

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			inLineComment = true
			i += 2

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '*' {
			inBlockComment = true
			i += 2

			continue
		}

		if input[i] != '#' {
			i++

			continue
		}

		if i+1 < len(input) && input[i+1] == '#' {
			i += 2

			continue
		}

		nameStart := i + 1
		for nameStart < len(input) && isTokenPasteSpace(input[nameStart]) {
			nameStart++
		}

		if nameStart >= len(input) || !isIdentifierStart(input[nameStart]) {
			i++

			continue
		}

		nameEnd := nameStart + 1
		for nameEnd < len(input) && isIdentifierPart(input[nameEnd]) {
			nameEnd++
		}

		name := input[nameStart:nameEnd]
		if !replaced {
			out.Grow(len(input))
			replaced = true
		}

		out.WriteString(input[lastWrite:i])
		out.WriteString(quoteIntrinsicString(name))
		lastWrite = nameEnd
		i = nameEnd
	}

	if !replaced {
		return input, false
	}

	out.WriteString(input[lastWrite:])

	return out.String(), true
}

// expandFunctionMacros expands function-like invocations.
func (p *preprocessor) expandFunctionMacros(input string) (string, bool) {
	names := p.functionMacroNames()
	if len(names) == 0 {
		return input, false
	}

	changedAny := false
	output := input

	for _, name := range names {
		if !strings.Contains(output, name+"(") {
			continue
		}

		def := p.macros[name]
		searchFrom := 0

		for {
			start, args, end, ok := findMacroCall(output, name, searchFrom)

			if !ok {
				break
			}

			if args == nil {
				if fallbackArgs, ok := malformedTwoArgFallback(def, output[start+len(name)+1:end]); ok {
					if end < len(output) && output[end] == ';' {
						end++
					}

					args = fallbackArgs
				}
			}

			if args == nil {
				if shouldConsumeMalformedCallSemicolon(output, start, name, end) {
					end++
				}

				// Malformed call syntax is tolerated by DayZ CfgConvert in strict
				// mode. Drop the invocation span and keep preprocessing.
				output = output[:start] + output[end:]
				searchFrom = start
				changedAny = true

				continue
			}

			if len(args) != len(def.Params) {
				// DayZ CfgConvert does not fail on arg-count mismatch for
				// function-like macro calls. It drops the whole invocation.
				output = output[:start] + output[end:]
				searchFrom = start
				changedAny = true

				continue
			}

			p.bindMacroParams(def.Params, args)

			replacement := def.Body
			replacement = stringifyMacroParams(replacement, def.Params, args)

			for idx := range def.Params {
				replacement = replaceIdentifierTokens(replacement, def.Params[idx], args[idx])
			}

			replacement, _ = p.expandObjectMacros(replacement)
			output = output[:start] + replacement + output[end:]
			searchFrom = start + len(replacement)
			changedAny = true
		}
	}

	return output, changedAny
}

// stringifyMacroParams applies #param stringification for function-like macros.
// It skips strings/comments and ignores token-paste operator (##).
func stringifyMacroParams(body string, params []string, args []string) string {
	if body == "" || len(params) == 0 || len(args) == 0 || !strings.Contains(body, "#") {
		return body
	}

	paramIndex := make(map[string]int, len(params))
	for idx := range params {
		paramIndex[params[idx]] = idx
	}

	var out strings.Builder
	lastWrite := 0
	replaced := false

	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(body); {
		if inLineComment {
			if body[i] == '\n' {
				inLineComment = false
			}

			i++

			continue
		}

		if inBlockComment {
			if body[i] == '*' && i+1 < len(body) && body[i+1] == '/' {
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			if body[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if body[i] == '"' {
			inString = true
			i++

			continue
		}

		if body[i] == '/' && i+1 < len(body) && body[i+1] == '/' {
			inLineComment = true
			i += 2

			continue
		}

		if body[i] == '/' && i+1 < len(body) && body[i+1] == '*' {
			inBlockComment = true
			i += 2

			continue
		}

		if body[i] != '#' {
			i++

			continue
		}

		// Leave token paste untouched.
		if i+1 < len(body) && body[i+1] == '#' {
			i += 2

			continue
		}

		nameStart := i + 1
		for nameStart < len(body) && isTokenPasteSpace(body[nameStart]) {
			nameStart++
		}

		if nameStart >= len(body) || !isIdentifierStart(body[nameStart]) {
			i++

			continue
		}

		nameEnd := nameStart + 1
		for nameEnd < len(body) && isIdentifierPart(body[nameEnd]) {
			nameEnd++
		}

		name := body[nameStart:nameEnd]
		argPos, ok := paramIndex[name]
		if !ok || argPos >= len(args) {
			i++

			continue
		}

		if !replaced {
			out.Grow(len(body))
			replaced = true
		}

		out.WriteString(body[lastWrite:i])
		out.WriteString(stringifyMacroArg(args[argPos]))
		lastWrite = nameEnd
		i = nameEnd
	}

	if !replaced {
		return body
	}

	out.WriteString(body[lastWrite:])

	return out.String()
}

// stringifyMacroArg converts raw macro argument text to quoted string literal.
func stringifyMacroArg(arg string) string {
	value := arg
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)

	return `"` + value + `"`
}

// expandObjectMacros expands object-like macro tokens.
func (p *preprocessor) expandObjectMacros(input string) (string, bool) {
	names := p.objectMacroNames()
	if len(names) == 0 {
		return input, false
	}

	result := input
	changedAny := false

	for _, name := range names {
		if !strings.Contains(result, name) {
			continue
		}

		def := p.macros[name]
		next := replaceIdentifierTokens(result, name, def.Body)
		if next != result {
			changedAny = true
			result = next
		}
	}

	return result, changedAny
}

// objectMacroNames returns object-like macro names sorted by length descending.
func (p *preprocessor) objectMacroNames() []string {
	p.refreshMacroNameCache()

	return p.objectMacroNamesV0
}

// functionMacroNames returns function-like macro names sorted by length descending.
func (p *preprocessor) functionMacroNames() []string {
	p.refreshMacroNameCache()

	return p.functionMacroNamesV0
}

// refreshMacroNameCache rebuilds sorted macro name lists when table changed.
func (p *preprocessor) refreshMacroNameCache() {
	if !p.macroNamesDirty {
		return
	}

	objectNames := make([]string, 0, len(p.macros))
	functionNames := make([]string, 0, len(p.macros))

	for name, def := range p.macros {
		if def.FunctionLike {
			functionNames = append(functionNames, name)
		} else {
			objectNames = append(objectNames, name)
		}
	}

	sort.Slice(objectNames, func(i int, j int) bool {
		if len(objectNames[i]) == len(objectNames[j]) {
			return objectNames[i] < objectNames[j]
		}

		return len(objectNames[i]) > len(objectNames[j])
	})

	sort.Slice(functionNames, func(i int, j int) bool {
		if len(functionNames[i]) == len(functionNames[j]) {
			return functionNames[i] < functionNames[j]
		}

		return len(functionNames[i]) > len(functionNames[j])
	})

	p.objectMacroNamesV0 = objectNames
	p.functionMacroNamesV0 = functionNames
	p.macroNamesDirty = false
}

// parseParams splits function-like macro parameters.
func parseParams(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	params := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			params = append(params, trimmed)
		}
	}

	return params
}

// bindMacroParams binds function-like arguments as temporary object-like macros.
// This mirrors DayZ CfgConvert behavior for nested expansions relying on outer args.
func (p *preprocessor) bindMacroParams(params []string, args []string) {
	if len(params) == 0 || len(args) == 0 {
		return
	}

	count := min(len(params), len(args))

	for i := range count {
		name := strings.TrimSpace(params[i])
		if name == "" {
			continue
		}

		p.macros[name] = macroDefinition{
			Name: name,
			Body: args[i],
		}
	}

	p.macroNamesDirty = true
}

// cloneMacroDefinitions copies macro table for temporary line-local mutation.
func cloneMacroDefinitions(in map[string]macroDefinition) map[string]macroDefinition {
	if len(in) == 0 {
		return map[string]macroDefinition{}
	}

	out := make(map[string]macroDefinition, len(in))
	maps.Copy(out, in)

	return out
}

// replaceIdentifierTokens replaces identifier token outside strings/comments.
func replaceIdentifierTokens(input string, name string, replacement string) string {
	if input == "" || name == "" {
		return input
	}

	var out strings.Builder
	lastWrite := 0
	replaced := false

	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); {
		if inLineComment {
			if input[i] == '\n' {
				inLineComment = false
			}

			i++
			continue
		}

		if inBlockComment {
			if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i += 2
				continue
			}

			i++
			continue
		}

		if inString {
			if input[i] == '"' {
				inString = false
			}

			i++
			continue
		}

		if input[i] == '"' {
			inString = true
			i++
			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			inLineComment = true
			i += 2
			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '*' {
			inBlockComment = true
			i += 2
			continue
		}

		if hasIdentifierAt(input, i, name) {
			if !replaced {
				out.Grow(len(input))
				replaced = true
			}

			out.WriteString(input[lastWrite:i])
			out.WriteString(replacement)
			i += len(name)
			lastWrite = i

			continue
		}

		i++
	}

	if !replaced {
		return input
	}

	out.WriteString(input[lastWrite:])

	return out.String()
}

// hasIdentifierAt checks identifier with token boundaries.
func hasIdentifierAt(input string, at int, name string) bool {
	if at < 0 || at+len(name) > len(input) {
		return false
	}

	if input[at:at+len(name)] != name {
		return false
	}

	if at > 0 && isIdentifierPart(input[at-1]) {
		return false
	}

	if at+len(name) < len(input) && isIdentifierPart(input[at+len(name)]) {
		return false
	}

	return true
}

// findMacroCall finds NAME(args...) invocation and returns parsed args and range.
func findMacroCall(input string, name string, from int) (int, []string, int, bool) {
	if from < 0 {
		from = 0
	}

	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); {
		if inLineComment {
			if input[i] == '\n' {
				inLineComment = false
			}

			i++

			continue
		}

		if inBlockComment {
			if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			if input[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if input[i] == '"' {
			inString = true
			i++

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			inLineComment = true
			i += 2

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '*' {
			inBlockComment = true
			i += 2

			continue
		}

		if i < from {
			i++

			continue
		}

		if !hasIdentifierAt(input, i, name) {
			i++

			continue
		}

		open := i + len(name)
		if open >= len(input) || input[open] != '(' {
			i++

			continue
		}

		args, end, err := parseMacroArgs(input, open)
		if err != nil {
			return i, nil, findMalformedMacroCallEnd(input, open), true
		}

		return i, args, end, true
	}

	return 0, nil, 0, false
}

// findMalformedMacroCallEnd finds safe truncate point for malformed macro call.
func findMalformedMacroCallEnd(input string, open int) int {
	if open < 0 || open >= len(input) {
		return len(input)
	}

	for i := open + 1; i < len(input); i++ {
		switch input[i] {
		case ';', '\n', '\r', '{', '}':
			return i
		}
	}

	return len(input)
}

// malformedTwoArgFallback emulates DayZ malformed two-arg behavior where
// `NAME(a,b` may still resolve as `NAME(a,a)` in strict mode.
func malformedTwoArgFallback(def macroDefinition, callBody string) ([]string, bool) {
	if len(def.Params) != 2 {
		return nil, false
	}

	if strings.Contains(callBody, `"`) || strings.Contains(callBody, "//") || strings.Contains(callBody, "/*") {
		return nil, false
	}

	depth := 0
	comma := -1
	for i := 0; i < len(callBody); i++ {
		switch callBody[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				if comma >= 0 {
					return nil, false
				}

				comma = i
			}
		}
	}

	if comma < 0 {
		return nil, false
	}

	first := strings.TrimSpace(callBody[:comma])
	if first == "" {
		return nil, false
	}

	return []string{first, first}, true
}

// shouldConsumeMalformedCallSemicolon decides whether trailing ';' should be
// consumed together with malformed invocation span.
func shouldConsumeMalformedCallSemicolon(input string, start int, name string, end int) bool {
	if end < 0 || end >= len(input) || input[end] != ';' {
		return false
	}

	if end+1 < len(input) && (input[end+1] == '\n' || input[end+1] == '\r') {
		return true
	}

	bodyStart := start + len(name) + 1
	if bodyStart < 0 || bodyStart > end {
		return false
	}

	return strings.Contains(input[bodyStart:end], `"`)
}

// parseMacroArgs parses (...) argument list from opening parenthesis position.
func parseMacroArgs(input string, open int) ([]string, int, error) {
	if open >= len(input) || input[open] != '(' {
		return nil, 0, errors.New("macro call parse without opening parenthesis")
	}

	args := make([]string, 0, 4)
	var current strings.Builder
	current.Grow(32)
	depth := 1
	inString := false

	for i := open + 1; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '"' {
				inString = false
			}

			// DayZ CfgConvert quirk: commas inside macro argument
			// double-quoted strings are removed.
			if ch != ',' {
				current.WriteByte(ch)
			}

			continue
		}

		if ch == '"' {
			inString = true
			current.WriteByte(ch)

			continue
		}

		if ch == '(' {
			depth++
			current.WriteByte(ch)

			continue
		}

		if ch == ')' {
			depth--
			if depth == 0 {
				arg := current.String()
				if strings.TrimSpace(arg) != "" || len(args) > 0 {
					args = append(args, arg)
				}

				return args, i + 1, nil
			}

			current.WriteByte(ch)

			continue
		}

		if ch == ',' && depth == 1 {
			arg := current.String()
			args = append(args, arg)
			current.Reset()

			continue
		}

		current.WriteByte(ch)
	}

	return nil, 0, errors.New("unterminated macro argument list")
}

// collapseTokenPaste removes ## operator and adjacent whitespace outside
// string literals and comments.
func collapseTokenPaste(input string) string {
	if input == "" || !strings.Contains(input, "##") {
		return input
	}

	out := make([]byte, 0, len(input))
	changed := false
	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); {
		if inLineComment {
			out = append(out, input[i])
			if input[i] == '\n' {
				inLineComment = false
			}

			i++

			continue
		}

		if inBlockComment {
			out = append(out, input[i])
			if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
				out = append(out, input[i+1])
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			out = append(out, input[i])
			if input[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if input[i] == '"' {
			inString = true
			out = append(out, input[i])
			i++

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			inLineComment = true
			out = append(out, input[i], input[i+1])
			i += 2

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '*' {
			inBlockComment = true
			out = append(out, input[i], input[i+1])
			i += 2

			continue
		}

		if i+1 < len(input) && input[i] == '#' && input[i+1] == '#' {
			changed = true

			leftBoundary := len(out)
			for leftBoundary > 0 && isTokenPasteSpace(out[leftBoundary-1]) {
				leftBoundary--
			}

			leftToken := tokenLeftOf(string(out[:leftBoundary]), leftBoundary)
			leftHasToken := leftToken != ""
			if leftHasToken {
				out = out[:leftBoundary]
			}

			i += 2
			rightBoundary := i
			for rightBoundary < len(input) && isTokenPasteSpace(input[rightBoundary]) {
				rightBoundary++
			}

			rightToken := tokenRightOf(input, rightBoundary)
			if leftHasToken && rightToken != "" {
				i = rightBoundary

				if shouldKeepTokenPasteSeparatorByToken(leftToken, rightToken) {
					out = append(out, ' ')
				}

				continue
			}

			out = append(out, input[i:rightBoundary]...)
			i = rightBoundary

			continue
		}

		out = append(out, input[i])
		i++
	}

	if !changed {
		return input
	}

	return string(out)
}

// isTokenPasteSpace checks collapsible whitespace around ## operator.
func isTokenPasteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

// shouldKeepTokenPasteSeparatorByToken applies keyword boundary separator rule.
func shouldKeepTokenPasteSeparatorByToken(leftToken string, rightToken string) bool {
	if leftToken == "" || rightToken == "" {
		return false
	}

	switch leftToken {
	case "class", "delete", "extern":
		return true
	default:
		return false
	}
}

// tokenLeftOf returns identifier token ending at boundary index.
func tokenLeftOf(input string, boundary int) string {
	start := boundary
	for start > 0 && isIdentifierPart(input[start-1]) {
		start--
	}

	if start == boundary {
		return ""
	}

	return input[start:boundary]
}

// tokenRightOf returns identifier token starting at boundary index.
func tokenRightOf(input string, boundary int) string {
	end := boundary
	for end < len(input) && isIdentifierPart(input[end]) {
		end++
	}

	if end == boundary {
		return ""
	}

	return input[boundary:end]
}

// mergeLineContinuationsWithSourceLines joins lines ending with backslash.
func mergeLineContinuationsWithSourceLines(input string) []logicalLine {
	lines := strings.Split(input, "\n")
	if len(lines) == 0 {
		return nil
	}

	out := make([]logicalLine, 0, len(lines))
	var carry string
	carryLine := 0

	for idx, raw := range lines {
		lineNo := idx + 1
		line := strings.TrimSuffix(raw, "\r")
		if strings.HasSuffix(line, "\\") {
			if carryLine == 0 {
				carryLine = lineNo
			}

			carry += strings.TrimSuffix(line, "\\")

			continue
		}

		if carry != "" {
			out = append(out, logicalLine{
				Text:       carry + line,
				SourceLine: carryLine,
			})
			carry = ""
			carryLine = 0

			continue
		}

		out = append(out, logicalLine{
			Text:       line,
			SourceLine: lineNo,
		})
	}

	if carry != "" {
		out = append(out, logicalLine{
			Text:       carry,
			SourceLine: carryLine,
		})
	}

	return out
}

// normalizeLineEndings converts CRLF/CR into LF.
func normalizeLineEndings(input string) string {
	out := strings.ReplaceAll(input, "\r\n", "\n")
	out = strings.ReplaceAll(out, "\r", "\n")

	return out
}
