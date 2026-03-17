// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// expandExtendedIntrinsics expands all extended-only intrinsics in one pass.
func (p *preprocessor) expandExtendedIntrinsics(input string, sourceFile string) (string, bool, error) {
	out := input
	changedAny := false

	out, changed, err := p.expandPathNormIntrinsics(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_TRIM", strings.TrimSpace)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_LOWER", strings.ToLower)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_UPPER", strings.ToUpper)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_PASCAL", toPascalCase)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_CAMEL", toCamelCase)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_SNAKE", toSnakeCase)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringUnaryIntrinsic(out, "__STR_CONST", toConstCase)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringQuoteIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringReplaceIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringJoinIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandStringSplitIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandFilesJoinIntrinsic(out, sourceFile)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandFilesCountIntrinsic(out, sourceFile)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandFilesGetIntrinsic(out, sourceFile)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandFilesRenderIntrinsic(out, sourceFile)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandRangeRenderIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	out, changed, err = p.expandEachRenderIntrinsic(out)
	if err != nil {
		return input, false, err
	}

	changedAny = changedAny || changed

	return out, changedAny, nil
}

// expandPathNormIntrinsics expands extended __PATH_NORM(expr) helper.
func (p *preprocessor) expandPathNormIntrinsics(input string) (string, bool, error) {
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, expr, end, ok, err := findIntrinsicCall(out, "__PATH_NORM", searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		expr = strings.TrimSpace(expr)
		if expr == "" {
			return out, changedAny, errors.New("__PATH_NORM requires one expression argument")
		}

		value, evalErr := p.evalIntrinsicExpr(expr)
		if evalErr != nil {
			return out, changedAny, fmt.Errorf("__PATH_NORM evaluate argument: %w", evalErr)
		}

		pathText := intrinsicValueToString(value)
		normalized := normalizePathSlashes(pathText)
		replacement := quoteIntrinsicString(normalized)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// normalizePathSlashes converts mixed slashes to single backslash separators.
func normalizePathSlashes(input string) string {
	if input == "" {
		return input
	}

	normalized := strings.ReplaceAll(input, "/", `\`)
	if !strings.Contains(normalized, `\\`) {
		return normalized
	}

	out := make([]byte, 0, len(normalized))
	prevSlash := false

	for idx := 0; idx < len(normalized); idx++ {
		ch := normalized[idx]
		if ch != '\\' {
			out = append(out, ch)
			prevSlash = false

			continue
		}

		if prevSlash {
			continue
		}

		out = append(out, ch)
		prevSlash = true
	}

	return string(out)
}

// expandStringUnaryIntrinsic expands one-argument string intrinsic.
func (p *preprocessor) expandStringUnaryIntrinsic(
	input string,
	name string,
	transform func(string) string,
) (string, bool, error) {
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 1 {
			return out, changedAny, fmt.Errorf("%s requires exactly one argument", name)
		}

		value, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument: %w", name, err)
		}

		replacement := quoteIntrinsicString(transform(value))
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandStringQuoteIntrinsic expands __STR_QUOTE(expr).
func (p *preprocessor) expandStringQuoteIntrinsic(input string) (string, bool, error) {
	const name = "__STR_QUOTE"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 1 {
			return out, changedAny, fmt.Errorf("%s requires exactly one argument", name)
		}

		value, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument: %w", name, err)
		}

		replacement := quoteIntrinsicString(value)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandStringReplaceIntrinsic expands __STR_REPLACE(text, old, new).
func (p *preprocessor) expandStringReplaceIntrinsic(input string) (string, bool, error) {
	const name = "__STR_REPLACE"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 3 {
			return out, changedAny, fmt.Errorf("%s requires exactly three arguments", name)
		}

		textValue, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		oldValue, err := p.evalIntrinsicArgToString(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		newValue, err := p.evalIntrinsicArgToString(args[2])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #3: %w", name, err)
		}

		replacement := quoteIntrinsicString(strings.ReplaceAll(textValue, oldValue, newValue))
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandStringJoinIntrinsic expands __STR_JOIN(delimiter, value1, ...).
func (p *preprocessor) expandStringJoinIntrinsic(input string) (string, bool, error) {
	const name = "__STR_JOIN"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) < 2 {
			return out, changedAny, fmt.Errorf("%s requires delimiter and at least one value", name)
		}

		delimiter, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate delimiter: %w", name, err)
		}

		parts := make([]string, 0, len(args)-1)
		for idx := 1; idx < len(args); idx++ {
			part, partErr := p.evalIntrinsicArgToString(args[idx])
			if partErr != nil {
				return out, changedAny, fmt.Errorf("%s evaluate argument #%d: %w", name, idx+1, partErr)
			}

			parts = append(parts, part)
		}

		replacement := quoteIntrinsicString(strings.Join(parts, delimiter))
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandStringSplitIntrinsic expands __STR_SPLIT(text, delimiter, index).
func (p *preprocessor) expandStringSplitIntrinsic(input string) (string, bool, error) {
	const name = "__STR_SPLIT"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 3 {
			return out, changedAny, fmt.Errorf("%s requires exactly three arguments", name)
		}

		textValue, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		delimiter, err := p.evalIntrinsicArgToString(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		if delimiter == "" {
			return out, changedAny, fmt.Errorf("%s delimiter must not be empty", name)
		}

		index, err := p.evalIntrinsicArgToIndex(args[2])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #3: %w", name, err)
		}

		parts := strings.Split(textValue, delimiter)
		part := ""
		if index >= 0 && index < len(parts) {
			part = parts[index]
		}

		replacement := quoteIntrinsicString(part)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandFilesJoinIntrinsic expands __FILES_JOIN(pattern, delimiter) and
// __FILES_JOIN(pattern) with default ';' delimiter.
func (p *preprocessor) expandFilesJoinIntrinsic(input string, sourceFile string) (string, bool, error) {
	const name = "__FILES_JOIN"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) < 1 || len(args) > 2 {
			return out, changedAny, fmt.Errorf("%s requires pattern and optional delimiter", name)
		}

		pattern, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		delimiter := ";"
		if len(args) == 2 {
			delimiter, err = p.evalIntrinsicArgToString(args[1])
			if err != nil {
				return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
			}
		}

		files, err := p.globFiles(pattern, sourceFile)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s: %w", name, err)
		}

		replacement := quoteIntrinsicString(strings.Join(files, delimiter))
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandFilesCountIntrinsic expands __FILES_COUNT(pattern).
func (p *preprocessor) expandFilesCountIntrinsic(input string, sourceFile string) (string, bool, error) {
	const name = "__FILES_COUNT"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 1 {
			return out, changedAny, fmt.Errorf("%s requires exactly one argument", name)
		}

		pattern, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		files, err := p.globFiles(pattern, sourceFile)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s: %w", name, err)
		}

		replacement := strconv.Itoa(len(files))
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandFilesGetIntrinsic expands __FILES_GET(pattern, index).
func (p *preprocessor) expandFilesGetIntrinsic(input string, sourceFile string) (string, bool, error) {
	const name = "__FILES_GET"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) != 2 {
			return out, changedAny, fmt.Errorf("%s requires exactly two arguments", name)
		}

		pattern, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		index, err := p.evalIntrinsicArgToIndex(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		files, err := p.globFiles(pattern, sourceFile)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s: %w", name, err)
		}

		value := ""
		if index < len(files) {
			value = files[index]
		}

		replacement := quoteIntrinsicString(value)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandFilesRenderIntrinsic expands __FILES_RENDER(pattern, template[, delimiter]).
// It emits raw joined text without wrapping resulting output in quotes.
func (p *preprocessor) expandFilesRenderIntrinsic(input string, sourceFile string) (string, bool, error) {
	const name = "__FILES_RENDER"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) < 2 || len(args) > 3 {
			return out, changedAny, fmt.Errorf("%s requires pattern, template, and optional delimiter", name)
		}

		pattern, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		templateText, err := p.evalIntrinsicArgToString(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		delimiter := "\n"
		if len(args) == 3 {
			delimiter, err = p.evalIntrinsicArgToString(args[2])
			if err != nil {
				return out, changedAny, fmt.Errorf("%s evaluate argument #3: %w", name, err)
			}
		}

		files, err := p.globFiles(pattern, sourceFile)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s: %w", name, err)
		}

		rendered := make([]string, 0, len(files))
		for idx, filePath := range files {
			value, renderErr := renderFileTemplate(templateText, filePath, idx)
			if renderErr != nil {
				return out, changedAny, fmt.Errorf("%s render item #%d: %w", name, idx+1, renderErr)
			}

			rendered = append(rendered, value)
		}

		replacement := strings.Join(rendered, delimiter)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandRangeRenderIntrinsic expands __FOR_RANGE_RENDER(start, end, template[, delimiter]).
// It emits raw joined text without wrapping resulting output in quotes.
func (p *preprocessor) expandRangeRenderIntrinsic(input string) (string, bool, error) {
	const name = "__FOR_RANGE_RENDER"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) < 3 || len(args) > 4 {
			return out, changedAny, fmt.Errorf("%s requires start, end, template, and optional delimiter", name)
		}

		startValue, err := p.evalIntrinsicArgToInt(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		endValue, err := p.evalIntrinsicArgToInt(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		templateText, err := p.evalIntrinsicArgToString(args[2])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #3: %w", name, err)
		}

		delimiter := "\n"
		if len(args) == 4 {
			delimiter, err = p.evalIntrinsicArgToString(args[3])
			if err != nil {
				return out, changedAny, fmt.Errorf("%s evaluate argument #4: %w", name, err)
			}
		}

		steps := countRangeSteps(startValue, endValue)
		if steps > p.extendedLoopMaxItems {
			return out, changedAny, fmt.Errorf("%s iteration limit exceeded (%d)", name, p.extendedLoopMaxItems)
		}

		rendered := make([]string, 0, steps)
		step := 1
		if startValue > endValue {
			step = -1
		}

		for idx, value := 0, startValue; ; idx, value = idx+1, value+step {
			item, renderErr := renderRangeTemplate(templateText, idx, value)
			if renderErr != nil {
				return out, changedAny, fmt.Errorf("%s render item #%d: %w", name, idx+1, renderErr)
			}

			rendered = append(rendered, item)
			if value == endValue {
				break
			}
		}

		replacement := strings.Join(rendered, delimiter)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// expandEachRenderIntrinsic expands
// __FOR_EACH_RENDER(template, delimiter, value1, ...).
// It emits raw joined text without wrapping resulting output in quotes.
func (p *preprocessor) expandEachRenderIntrinsic(input string) (string, bool, error) {
	const name = "__FOR_EACH_RENDER"
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, name, searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		args, err := parseIntrinsicArgs(body)
		if err != nil {
			return out, changedAny, fmt.Errorf("%s parse arguments: %w", name, err)
		}

		if len(args) < 3 {
			return out, changedAny, fmt.Errorf(
				"%s requires template, delimiter, and at least one value",
				name,
			)
		}

		templateText, err := p.evalIntrinsicArgToString(args[0])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #1: %w", name, err)
		}

		delimiter, err := p.evalIntrinsicArgToString(args[1])
		if err != nil {
			return out, changedAny, fmt.Errorf("%s evaluate argument #2: %w", name, err)
		}

		values := make([]string, 0, len(args)-2)
		for idx := 2; idx < len(args); idx++ {
			value, valueErr := p.evalIntrinsicArgToString(args[idx])
			if valueErr != nil {
				return out, changedAny, fmt.Errorf(
					"%s evaluate argument #%d: %w",
					name,
					idx+1,
					valueErr,
				)
			}

			values = append(values, value)
		}

		if len(values) > p.extendedLoopMaxItems {
			return out, changedAny, fmt.Errorf("%s iteration limit exceeded (%d)", name, p.extendedLoopMaxItems)
		}

		rendered := make([]string, 0, len(values))
		for idx, value := range values {
			item, renderErr := renderEachTemplate(templateText, idx, value)
			if renderErr != nil {
				return out, changedAny, fmt.Errorf("%s render item #%d: %w", name, idx+1, renderErr)
			}

			rendered = append(rendered, item)
		}

		replacement := strings.Join(rendered, delimiter)
		out = out[:start] + replacement + out[end:]
		searchFrom = start + len(replacement)
		changedAny = true
	}
}

// evalIntrinsicArgToString evaluates intrinsic argument expression and converts it to string.
func (p *preprocessor) evalIntrinsicArgToString(expr string) (string, error) {
	value, err := p.evalIntrinsicExpr(strings.TrimSpace(expr))
	if err != nil {
		return "", err
	}

	return intrinsicValueToString(value), nil
}

// evalIntrinsicArgToIndex evaluates intrinsic argument expression and converts it to non-negative index.
func (p *preprocessor) evalIntrinsicArgToIndex(expr string) (int, error) {
	index, err := p.evalIntrinsicArgToInt(expr)
	if err != nil {
		return 0, err
	}

	if index < 0 {
		return 0, errors.New("index must be non-negative")
	}

	return index, nil
}

// evalIntrinsicArgToInt evaluates intrinsic argument expression and converts it to integer.
func (p *preprocessor) evalIntrinsicArgToInt(expr string) (int, error) {
	value, err := p.evalIntrinsicExpr(strings.TrimSpace(expr))
	if err != nil {
		return 0, err
	}

	if value.isString {
		return 0, errors.New("value must be numeric")
	}

	if !isIntegralFloat(value.number) {
		return 0, errors.New("value must be integer")
	}

	return int(value.number), nil
}

// globFiles resolves glob pattern from source file context and returns normalized relative paths.
func (p *preprocessor) globFiles(pattern string, sourceFile string) ([]string, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return nil, errors.New("glob pattern must not be empty")
	}

	baseDir := filepath.Dir(sourceFile)
	rootDir := baseDir
	if strings.TrimSpace(p.extendedFSRoot) != "" {
		rootDir = p.extendedFSRoot
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("resolve fs root: %w", err)
	}

	pattern = strings.ReplaceAll(pattern, "/", string(os.PathSeparator))
	pattern = strings.ReplaceAll(pattern, `\`, string(os.PathSeparator))

	globPattern := pattern
	if !filepath.IsAbs(globPattern) {
		globPattern = filepath.Join(baseDir, globPattern)
	}

	absGlobPattern, err := filepath.Abs(globPattern)
	if err != nil {
		return nil, fmt.Errorf("resolve glob pattern: %w", err)
	}

	if !isPathWithinRoot(absGlobPattern, absRoot) {
		return nil, errors.New("glob pattern is outside extended fs root")
	}

	matches, err := filepath.Glob(absGlobPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}

	sort.Strings(matches)

	if len(matches) > p.extendedFSMaxItems {
		return nil, fmt.Errorf("matched items exceed limit (%d)", p.extendedFSMaxItems)
	}

	out := make([]string, 0, len(matches))
	for _, match := range matches {
		info, statErr := os.Stat(match)
		if statErr != nil || info.IsDir() {
			continue
		}

		absMatch, absErr := filepath.Abs(match)
		if absErr != nil {
			return nil, fmt.Errorf("resolve matched path: %w", absErr)
		}

		if !isPathWithinRoot(absMatch, absRoot) {
			return nil, errors.New("matched file is outside extended fs root")
		}

		rel, relErr := filepath.Rel(baseDir, absMatch)
		if relErr != nil {
			return nil, fmt.Errorf("relativize matched path: %w", relErr)
		}

		out = append(out, normalizePathSlashes(rel))
	}

	return out, nil
}

// isPathWithinRoot reports whether path is equal to root or nested under it.
func isPathWithinRoot(path string, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)

	if strings.EqualFold(path, root) {
		return true
	}

	sep := string(os.PathSeparator)
	rootPrefix := root
	if !strings.HasSuffix(rootPrefix, sep) {
		rootPrefix += sep
	}

	return strings.HasPrefix(strings.ToLower(path), strings.ToLower(rootPrefix))
}

// renderFileTemplate replaces file placeholders in template text.
func renderFileTemplate(templateText string, filePath string, index int) (string, error) {
	name := filepathBase(filePath)
	stem := trimExt(name)
	ext := filepath.Ext(name)
	values := map[string]string{
		"path":  filePath,
		"name":  name,
		"stem":  stem,
		"ext":   ext,
		"index": strconv.Itoa(index),
	}

	return renderTemplateWithFilters(templateText, values)
}

// renderRangeTemplate replaces range placeholders in template text.
func renderRangeTemplate(templateText string, index int, value int) (string, error) {
	values := map[string]string{
		"index": strconv.Itoa(index),
		"value": strconv.Itoa(value),
	}

	return renderTemplateWithFilters(templateText, values)
}

// renderEachTemplate replaces foreach placeholders in template text.
func renderEachTemplate(templateText string, index int, value string) (string, error) {
	values := map[string]string{
		"index": strconv.Itoa(index),
		"value": value,
	}

	return renderTemplateWithFilters(templateText, values)
}

// countRangeSteps returns inclusive number of values in [start..end] or [start..end] descending.
func countRangeSteps(start int, end int) int {
	if start <= end {
		return end - start + 1
	}

	return start - end + 1
}

// renderTemplateWithFilters renders placeholders like
// {stem|pascal} or {path|lower|replace(abc, xyz)}.
func renderTemplateWithFilters(templateText string, values map[string]string) (string, error) {
	if templateText == "" {
		return "", nil
	}

	var out strings.Builder
	out.Grow(len(templateText))

	for idx := 0; idx < len(templateText); {
		if templateText[idx] != '{' {
			out.WriteByte(templateText[idx])
			idx++

			continue
		}

		placeholder, ok, err := parseTemplatePlaceholder(templateText, idx)
		if err != nil {
			return "", err
		}

		if !ok {
			out.WriteByte(templateText[idx])
			idx++

			continue
		}

		value, exists := values[placeholder.name]
		if !exists {
			return "", fmt.Errorf("unknown template placeholder %q", placeholder.name)
		}

		value, err = applyTemplateFilters(value, placeholder.filters)
		if err != nil {
			return "", fmt.Errorf("template placeholder %q: %w", placeholder.name, err)
		}

		out.WriteString(value)
		idx = placeholder.end
	}

	return out.String(), nil
}

type templateFilter struct {
	name string
	args []string
}

type templatePlaceholder struct {
	name    string
	filters []templateFilter
	end     int
}

func parseTemplatePlaceholder(input string, start int) (templatePlaceholder, bool, error) {
	if start < 0 || start >= len(input) || input[start] != '{' {
		return templatePlaceholder{}, false, nil
	}

	cursor := start + 1
	cursor = skipTemplateSpaces(input, cursor)

	if cursor >= len(input) || !isTemplateIdentStart(input[cursor]) {
		return templatePlaceholder{}, false, nil
	}

	nameStart := cursor
	cursor++
	for cursor < len(input) && isTemplateIdentPart(input[cursor]) {
		cursor++
	}

	placeholder := templatePlaceholder{
		name: strings.ToLower(strings.TrimSpace(input[nameStart:cursor])),
	}

	for {
		cursor = skipTemplateSpaces(input, cursor)
		if cursor >= len(input) {
			return templatePlaceholder{}, false, errors.New("unterminated template placeholder")
		}

		if input[cursor] == '}' {
			placeholder.end = cursor + 1

			return placeholder, true, nil
		}

		if input[cursor] != '|' {
			return templatePlaceholder{}, false, nil
		}

		cursor++
		cursor = skipTemplateSpaces(input, cursor)
		if cursor >= len(input) || !isTemplateIdentStart(input[cursor]) {
			return templatePlaceholder{}, false, errors.New("template filter name expected")
		}

		filterStart := cursor
		cursor++
		for cursor < len(input) && isTemplateIdentPart(input[cursor]) {
			cursor++
		}

		filter := templateFilter{
			name: strings.ToLower(strings.TrimSpace(input[filterStart:cursor])),
		}

		cursor = skipTemplateSpaces(input, cursor)
		if cursor < len(input) && input[cursor] == '(' {
			argsBody, end, err := parseTemplateFilterArgsBody(input, cursor)
			if err != nil {
				return templatePlaceholder{}, false, err
			}

			args, err := parseIntrinsicArgs(argsBody)
			if err != nil {
				return templatePlaceholder{}, false, fmt.Errorf(
					"parse filter %q arguments: %w",
					filter.name,
					err,
				)
			}

			filter.args = make([]string, 0, len(args))
			for _, arg := range args {
				filter.args = append(filter.args, unquoteTemplateArg(arg))
			}

			cursor = end
		}

		placeholder.filters = append(placeholder.filters, filter)
	}
}

func parseTemplateFilterArgsBody(input string, open int) (string, int, error) {
	if open >= len(input) || input[open] != '(' {
		return "", 0, errors.New("template filter arguments must start with '('")
	}

	start := open + 1
	depth := 1
	inString := false
	escaped := false

	for idx := start; idx < len(input); idx++ {
		ch := input[idx]
		if inString {
			if escaped {
				escaped = false

				continue
			}

			if ch == '\\' {
				escaped = true

				continue
			}

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
				return strings.TrimSpace(input[start:idx]), idx + 1, nil
			}
		}
	}

	return "", 0, errors.New("unterminated template filter arguments")
}

func applyTemplateFilters(value string, filters []templateFilter) (string, error) {
	out := value

	for _, filter := range filters {
		switch filter.name {
		case "trim":
			if len(filter.args) != 0 {
				return "", errors.New("trim filter does not accept arguments")
			}

			out = strings.TrimSpace(out)

		case "lower":
			if len(filter.args) != 0 {
				return "", errors.New("lower filter does not accept arguments")
			}

			out = strings.ToLower(out)

		case "upper":
			if len(filter.args) != 0 {
				return "", errors.New("upper filter does not accept arguments")
			}

			out = strings.ToUpper(out)

		case "replace":
			if len(filter.args) != 2 {
				return "", errors.New("replace filter requires exactly two arguments")
			}

			out = strings.ReplaceAll(out, filter.args[0], filter.args[1])

		case "split":
			if len(filter.args) != 2 {
				return "", errors.New("split filter requires delimiter and index")
			}

			index, err := strconv.Atoi(filter.args[1])
			if err != nil {
				return "", errors.New("split filter index must be integer")
			}

			if index < 0 {
				return "", errors.New("split filter index must be >= 0")
			}

			parts := strings.Split(out, filter.args[0])
			if index >= len(parts) {
				out = ""

				continue
			}

			out = parts[index]

		case "quote":
			if len(filter.args) != 0 {
				return "", errors.New("quote filter does not accept arguments")
			}

			out = quoteIntrinsicString(out)

		case "pascal":
			if len(filter.args) != 0 {
				return "", errors.New("pascal filter does not accept arguments")
			}

			out = toPascalCase(out)

		case "camel":
			if len(filter.args) != 0 {
				return "", errors.New("camel filter does not accept arguments")
			}

			out = toCamelCase(out)

		case "snake":
			if len(filter.args) != 0 {
				return "", errors.New("snake filter does not accept arguments")
			}

			out = toSnakeCase(out)

		case "const":
			if len(filter.args) != 0 {
				return "", errors.New("const filter does not accept arguments")
			}

			out = toConstCase(out)

		case "path_norm", "slash_norm":
			if len(filter.args) != 0 {
				return "", errors.New("path_norm filter does not accept arguments")
			}

			out = normalizePathSlashes(out)

		default:
			return "", fmt.Errorf("unknown template filter %q", filter.name)
		}
	}

	return out, nil
}

func unquoteTemplateArg(value string) string {
	out := strings.TrimSpace(value)
	if len(out) < 2 {
		return out
	}

	if out[0] == '"' && out[len(out)-1] == '"' {
		out = out[1 : len(out)-1]
		out = strings.ReplaceAll(out, `\"`, `"`)
		out = strings.ReplaceAll(out, `\\`, `\`)

		return out
	}

	if out[0] == '\'' && out[len(out)-1] == '\'' {
		return out[1 : len(out)-1]
	}

	return out
}

func skipTemplateSpaces(input string, index int) int {
	for index < len(input) {
		switch input[index] {
		case ' ', '\t', '\r', '\n':
			index++
		default:
			return index
		}
	}

	return index
}

func isTemplateIdentStart(ch byte) bool {
	return ch == '_' || unicode.IsLetter(rune(ch))
}

func isTemplateIdentPart(ch byte) bool {
	return isTemplateIdentStart(ch) || unicode.IsDigit(rune(ch))
}

func toPascalCase(input string) string {
	parts := splitCaseWords(input)
	if len(parts) == 0 {
		return ""
	}

	var out strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}

		lower := strings.ToLower(part)
		runes := []rune(lower)
		runes[0] = unicode.ToUpper(runes[0])
		out.WriteString(string(runes))
	}

	return out.String()
}

func toCamelCase(input string) string {
	pascal := toPascalCase(input)
	if pascal == "" {
		return ""
	}

	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func toSnakeCase(input string) string {
	parts := splitCaseWords(input)
	if len(parts) == 0 {
		return ""
	}

	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}

		out = append(out, strings.ToLower(part))
	}

	return strings.Join(out, "_")
}

func toConstCase(input string) string {
	return strings.ToUpper(toSnakeCase(input))
}

func splitCaseWords(input string) []string {
	runes := []rune(input)
	words := make([]string, 0, 8)
	buf := make([]rune, 0, len(runes))

	flush := func() {
		if len(buf) == 0 {
			return
		}

		words = append(words, string(buf))
		buf = buf[:0]
	}

	for idx := range len(runes) {
		curr := runes[idx]
		if !unicode.IsLetter(curr) && !unicode.IsDigit(curr) {
			flush()

			continue
		}

		if len(buf) > 0 {
			prev := buf[len(buf)-1]
			next := rune(0)
			if idx+1 < len(runes) {
				next = runes[idx+1]
			}

			if isCaseBoundary(prev, curr, next) {
				flush()
			}
		}

		buf = append(buf, curr)
	}

	flush()

	return words
}

func isCaseBoundary(prev rune, curr rune, next rune) bool {
	if unicode.IsDigit(prev) != unicode.IsDigit(curr) {
		return true
	}

	if unicode.IsLower(prev) && unicode.IsUpper(curr) {
		return true
	}

	if unicode.IsUpper(prev) && unicode.IsUpper(curr) && unicode.IsLower(next) {
		return true
	}

	return false
}
