// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"fmt"
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
	if exists && existing.FunctionLike {
		p.emitMacroRedefinedWarning(filename, lineNo, "function-like", name)
	}

	if exists && !existing.FunctionLike {
		p.emitMacroRedefinedWarning(filename, lineNo, "object-like", name)
	}

	if strings.HasPrefix(rest, "(") {
		closeIdx := strings.Index(rest, ")")
		if closeIdx < 0 {
			p.emitError(CodePPUnterminatedMacroParams, filename, lineNo, "unterminated parameter list in #define")

			return ErrInvalidDirective
		}

		paramText := strings.TrimSpace(rest[1:closeIdx])
		body := strings.TrimSpace(rest[closeIdx+1:])
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
	result := line

	for pass := 0; pass < p.maxExpandDepth; pass++ {
		changed := false

		next, changedFunc, err := p.expandFunctionMacros(result)
		if err != nil {
			return result, fmt.Errorf("%w: %v", ErrMacroExpand, err)
		}

		next, changedObj := p.expandObjectMacros(next)
		changed = changedFunc || changedObj
		result = next

		if !changed {
			return result, nil
		}
	}

	return result, fmt.Errorf("%w: expansion depth overflow", ErrMacroExpand)
}

// expandFunctionMacros expands function-like invocations.
func (p *preprocessor) expandFunctionMacros(input string) (string, bool, error) {
	names := p.functionMacroNames()
	if len(names) == 0 {
		return input, false, nil
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
			start, args, end, ok, err := findMacroCall(output, name, searchFrom)
			if err != nil {
				return output, changedAny, err
			}

			if !ok {
				break
			}

			if len(args) != len(def.Params) {
				return output, changedAny, fmt.Errorf(
					"macro %s expects %d args, got %d",
					name,
					len(def.Params),
					len(args),
				)
			}

			replacement := def.Body
			replacement = stringifyMacroParams(replacement, def.Params, args)

			for idx := range def.Params {
				replacement = replaceIdentifierTokens(replacement, def.Params[idx], strings.TrimSpace(args[idx]))
			}

			replacement = collapseTokenPaste(replacement)
			output = output[:start] + replacement + output[end:]
			searchFrom = start + len(replacement)
			changedAny = true
		}
	}

	return output, changedAny, nil
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
	value := strings.TrimSpace(arg)
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
func findMacroCall(input string, name string, from int) (int, []string, int, bool, error) {
	if from < 0 {
		from = 0
	}

	for i := from; i < len(input); i++ {
		if !hasIdentifierAt(input, i, name) {
			continue
		}

		open := i + len(name)
		if open >= len(input) || input[open] != '(' {
			continue
		}

		args, end, err := parseMacroArgs(input, open)
		if err != nil {
			return 0, nil, 0, false, err
		}

		return i, args, end, true, nil
	}

	return 0, nil, 0, false, nil
}

// parseMacroArgs parses (...) argument list from opening parenthesis position.
func parseMacroArgs(input string, open int) ([]string, int, error) {
	if open >= len(input) || input[open] != '(' {
		return nil, 0, errors.New("macro call parse without opening parenthesis")
	}

	args := make([]string, 0, 4)
	start := open + 1
	depth := 1
	inString := false

	for i := open + 1; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '"' {
				inString = false
			}

			continue
		}

		if ch == '"' {
			inString = true
			continue
		}

		if ch == '(' {
			depth++
			continue
		}

		if ch == ')' {
			depth--
			if depth == 0 {
				arg := strings.TrimSpace(input[start:i])
				if arg != "" || len(args) > 0 {
					args = append(args, arg)
				}

				return args, i + 1, nil
			}

			continue
		}

		if ch == ',' && depth == 1 {
			arg := strings.TrimSpace(input[start:i])
			args = append(args, arg)
			start = i + 1
		}
	}

	return nil, 0, errors.New("unterminated macro argument list")
}

// collapseTokenPaste removes ## operator and adjacent whitespace.
func collapseTokenPaste(input string) string {
	for {
		idx := strings.Index(input, "##")
		if idx < 0 {
			return input
		}

		left := idx
		for left > 0 && isTokenPasteSpace(input[left-1]) {
			left--
		}

		right := idx + 2
		for right < len(input) && isTokenPasteSpace(input[right]) {
			right++
		}

		// Keep a separator for keyword->identifier edges (for example "class ##Name"),
		// because corpus macros rely on this compatibility behavior.
		if shouldKeepTokenPasteSeparator(input, left, right) {
			input = input[:left] + " " + input[right:]
			continue
		}

		input = input[:left] + input[right:]
	}
}

// isTokenPasteSpace checks collapsible whitespace around ## operator.
func isTokenPasteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

// shouldKeepTokenPasteSeparator decides whether ## collapse must keep one space.
func shouldKeepTokenPasteSeparator(input string, left int, right int) bool {
	if left <= 0 || right >= len(input) {
		return false
	}

	leftToken := tokenLeftOf(input, left)
	if leftToken == "" {
		return false
	}

	rightToken := tokenRightOf(input, right)
	if rightToken == "" {
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

// findUnresolvedMacroCalls finds unresolved macro-like NAME(...) invocations in a source line.
func (p *preprocessor) findUnresolvedMacroCalls(line string, inBlockComment bool) ([]string, bool) {
	found := make(map[string]struct{})
	inString := false

	for i := 0; i < len(line); {
		if inBlockComment {
			if i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			if line[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if line[i] == '"' {
			inString = true
			i++

			continue
		}

		if i+1 < len(line) && line[i] == '/' && line[i+1] == '/' {
			break
		}

		if i+1 < len(line) && line[i] == '/' && line[i+1] == '*' {
			inBlockComment = true
			i += 2

			continue
		}

		if !isIdentifierStart(line[i]) {
			i++

			continue
		}

		start := i
		i++
		for i < len(line) && isIdentifierPart(line[i]) {
			i++
		}

		name := line[start:i]
		if !isLikelyMacroName(name) {
			continue
		}

		if i >= len(line) || line[i] != '(' {
			continue
		}

		if p.macroExists(name) {
			continue
		}

		found[name] = struct{}{}
	}

	if len(found) == 0 {
		return nil, inBlockComment
	}

	out := make([]string, 0, len(found))
	for name := range found {
		out = append(out, name)
	}

	sort.Strings(out)

	return out, inBlockComment
}

// isLikelyMacroName checks macro-like naming style to reduce false positives on regular calls.
func isLikelyMacroName(name string) bool {
	if name == "" {
		return false
	}

	if !isIdentifierStart(name[0]) {
		return false
	}

	hasLetter := false

	for i := 0; i < len(name); i++ {
		ch := name[i]

		if ch >= 'a' && ch <= 'z' {
			return false
		}

		if ch >= 'A' && ch <= 'Z' {
			hasLetter = true

			continue
		}

		if ch == '_' || (ch >= '0' && ch <= '9') {
			continue
		}

		return false
	}

	return hasLetter
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
