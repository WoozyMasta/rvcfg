// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"strings"
)

// intrinsicValue stores evaluated intrinsic expression value.
type intrinsicValue struct {
	text     string
	number   float64
	isString bool
}

// expandExecEvalIntrinsics runs compat and extended intrinsic passes.
func (p *preprocessor) expandExecEvalIntrinsics(line string, sourceFile string) (string, error) {
	out := line

	for pass := 0; pass < p.maxExpandDepth; pass++ {
		changed := false

		next, changedExec, err := p.expandExecIntrinsics(out)
		if err != nil {
			return out, err
		}

		next, changedEval, err := p.expandEvalIntrinsics(next)
		if err != nil {
			return out, err
		}

		changedExtended := false
		if p.enableExtended {
			next, changedExtended, err = p.expandExtendedIntrinsics(next, sourceFile)
			if err != nil {
				return out, err
			}
		}

		out = next
		changed = changedExec || changedEval || changedExtended
		if !changed {
			return strings.TrimSpace(out), nil
		}
	}

	return out, errors.New("intrinsic expansion depth overflow")
}

// findIntrinsicCall finds NAME(...) call and returns call body and bounds.
func findIntrinsicCall(input string, name string, from int) (int, string, int, bool, error) {
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

		body, end, err := parseIntrinsicCallBody(input, open)
		if err != nil {
			return 0, "", 0, false, err
		}

		return i, body, end, true, nil
	}

	return 0, "", 0, false, nil
}

// parseIntrinsicCallBody parses balanced (...) body from input and open index.
func parseIntrinsicCallBody(input string, open int) (string, int, error) {
	if open >= len(input) || input[open] != '(' {
		return "", 0, errors.New("intrinsic call parse without opening parenthesis")
	}

	start := open + 1
	depth := 1
	inString := false

	for i := start; i < len(input); i++ {
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
				return strings.TrimSpace(input[start:i]), i + 1, nil
			}
		}
	}

	return "", 0, errors.New("unterminated intrinsic call")
}

// parseIntrinsicArgs splits intrinsic body arguments by top-level commas.
func parseIntrinsicArgs(body string) ([]string, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, nil
	}

	args := make([]string, 0, 4)
	start := 0
	depth := 0
	inString := false
	escaped := false

	for idx := 0; idx < len(body); idx++ {
		ch := body[idx]
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

		switch ch {
		case '"':
			inString = true
		case '(':
			depth++
		case ')':
			if depth == 0 {
				return nil, errors.New("unexpected ')' in argument list")
			}

			depth--
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(body[start:idx]))
				start = idx + 1
			}
		}
	}

	if inString {
		return nil, errors.New("unterminated string in argument list")
	}

	if depth != 0 {
		return nil, errors.New("unterminated nested expression in argument list")
	}

	args = append(args, strings.TrimSpace(body[start:]))

	return args, nil
}

// quoteIntrinsicString wraps text in quotes and escapes embedded quotes.
func quoteIntrinsicString(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}

// isIntegralFloat checks whether float value can be represented as integer.
func isIntegralFloat(v float64) bool {
	const eps = 1e-9
	iv := float64(int64(v))
	if v >= 0 {
		return v-iv < eps
	}

	return iv-v < eps
}
