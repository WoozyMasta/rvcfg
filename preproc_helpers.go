// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/woozymasta/lintkit/lint"
)

// emitError appends error diagnostic with source location.
func (p *preprocessor) emitError(code lint.Code, file string, line int, msg string) {
	p.emit(code, lint.SeverityError, file, line, msg)
}

// emitWarning appends warning diagnostic with source location.
func (p *preprocessor) emitWarning(code lint.Code, file string, line int, msg string) {
	p.emit(code, lint.SeverityWarning, file, line, msg)
}

// emit appends diagnostic with explicit severity and source location.
func (p *preprocessor) emit(code lint.Code, severity lint.Severity, file string, line int, msg string) {
	start := lint.Position{
		File:   file,
		Line:   line,
		Column: 1,
	}

	p.diagnostics = append(p.diagnostics, Diagnostic{
		Code:     code,
		Message:  msg,
		Severity: severity,
		Start:    start,
		End:      start,
	})
}

// splitDirective extracts directive name and raw argument text.
func splitDirective(line string) (string, string) {
	trimmed := strings.TrimSpace(line)
	trimmed = strings.TrimPrefix(trimmed, "#")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return "", ""
	}

	nameEnd := 0
	for nameEnd < len(trimmed) && isDirectiveNameChar(trimmed[nameEnd]) {
		nameEnd++
	}

	if nameEnd == 0 {
		return "", ""
	}

	name := trimmed[:nameEnd]
	arg := strings.TrimLeft(trimmed[nameEnd:], " \t")

	return name, arg
}

// isConditionalDirective checks whether directive is condition control.
func isConditionalDirective(name string) bool {
	switch name {
	case "if", "ifdef", "ifndef", "elif", "else", "endif":
		return true
	default:
		return false
	}
}

// parseIncludePathWithTail extracts include path and raw tail after include path.
func parseIncludePathWithTail(arg string) (string, string, error) {
	arg = strings.TrimLeft(arg, " \t")
	if len(arg) < 2 {
		return "", "", fmt.Errorf("invalid #include path %q", arg)
	}

	var (
		delimStart byte
		delimEnd   byte
	)

	switch arg[0] {
	case '"':
		delimStart = '"'
		delimEnd = '"'
	case '<':
		delimStart = '<'
		delimEnd = '>'
	default:
		return "", "", fmt.Errorf("invalid #include path %q", arg)
	}

	endIdx := -1
	for idx := 1; idx < len(arg); idx++ {
		if arg[idx] == delimEnd {
			endIdx = idx
			break
		}
	}

	if endIdx < 0 || arg[0] != delimStart {
		return "", "", fmt.Errorf("invalid #include path %q", arg)
	}

	path := arg[1:endIdx]
	tail := arg[endIdx+1:]
	tail, _ = stripComments(tail, false)

	return path, tail, nil
}

// splitDirectiveHeadTail splits first argument token and raw trailing suffix.
func splitDirectiveHeadTail(arg string) (string, string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return "", ""
	}

	idx := 0
	for idx < len(arg) && !isWhitespace(arg[idx]) {
		idx++
	}

	return arg[:idx], arg[idx:]
}

// isWhitespace reports whether byte is ASCII whitespace.
func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

// isDirectiveNameChar reports whether byte is valid directive name character.
func isDirectiveNameChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch == '_')
}

// fileExists checks whether path exists and is file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// stripComments removes // and /* */ comments while preserving string literals.
// It keeps block-comment state to support comments spanning multiple lines.
func stripComments(line string, inBlockComment bool) (string, bool) {
	if line == "" {
		return line, inBlockComment
	}

	if inBlockComment && !strings.Contains(line, "*/") {
		return "", true
	}

	if !inBlockComment {
		if !strings.Contains(line, "/") {
			return line, false
		}

		if !strings.Contains(line, "//") && !strings.Contains(line, "/*") {
			return line, false
		}
	}

	var out strings.Builder
	out.Grow(len(line))
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
			out.WriteByte(line[i])
			if line[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if line[i] == '"' {
			inString = true
			out.WriteByte(line[i])
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

		out.WriteByte(line[i])
		i++
	}

	return out.String(), inBlockComment
}
