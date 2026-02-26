// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"os"
	"strings"
)

// emitError appends error diagnostic with source location.
func (p *preprocessor) emitError(code DiagnosticCode, file string, line int, msg string) {
	p.emit(code, SeverityError, file, line, msg)
}

// emitWarning appends warning diagnostic with source location.
func (p *preprocessor) emitWarning(code DiagnosticCode, file string, line int, msg string) {
	p.emit(code, SeverityWarning, file, line, msg)
}

// emit appends diagnostic with explicit severity and source location.
func (p *preprocessor) emit(code DiagnosticCode, severity Severity, file string, line int, msg string) {
	start := Position{
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

	parts := strings.Fields(trimmed)
	if len(parts) == 1 {
		return parts[0], ""
	}

	name := parts[0]
	arg := strings.TrimSpace(trimmed[len(name):])

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

// parseQuotedInclude extracts "file" from include argument.
func parseQuotedInclude(arg string) (string, error) {
	arg = strings.TrimSpace(arg)
	if len(arg) < 2 || arg[0] != '"' || arg[len(arg)-1] != '"' {
		return "", fmt.Errorf("invalid #include path %q", arg)
	}

	return arg[1 : len(arg)-1], nil
}

// fileExists checks whether path exists and is file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
