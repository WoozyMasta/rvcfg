// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import "fmt"

// Severity is diagnostic severity level.
type Severity string

const (
	// SeverityError marks diagnostics that should fail current operation.
	SeverityError Severity = "error"

	// SeverityWarning marks diagnostics that should be reported but are non-fatal.
	SeverityWarning Severity = "warning"

	// SeverityInfo marks informational diagnostics.
	SeverityInfo Severity = "info"

	// SeverityNotice marks low-priority notices.
	SeverityNotice Severity = "notice"
)

// DiagnosticCode is stable machine-readable diagnostic identifier.
type DiagnosticCode string

// Position is a source location.
type Position struct {
	// File is source file path or logical source name.
	File string `json:"file,omitempty" yaml:"file,omitempty"`

	// Line is 1-based line index.
	Line int `json:"line,omitempty" yaml:"line,omitempty"`

	// Column is 1-based column index.
	Column int `json:"column,omitempty" yaml:"column,omitempty"`

	// Offset is 0-based absolute byte offset in source.
	Offset int `json:"offset,omitempty" yaml:"offset,omitempty"`
}

// Diagnostic describes parser/lexer/preprocessor issue.
type Diagnostic struct {
	// Code is stable machine-readable identifier.
	Code DiagnosticCode `json:"code,omitempty" yaml:"code,omitempty"`

	// Message is technical description of issue.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Severity defines impact of diagnostic.
	Severity Severity `json:"severity,omitzero" yaml:"severity,omitempty"`

	// Start is start location.
	Start Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is end location.
	End Position `json:"end,omitzero" yaml:"end,omitempty"`
}

// Error renders diagnostic in compact form.
func (d Diagnostic) Error() string {
	if d.Start.File == "" {
		return fmt.Sprintf("%s: %s", d.Code, d.Message)
	}

	return fmt.Sprintf(
		"%s:%d:%d: %s: %s",
		d.Start.File,
		d.Start.Line,
		d.Start.Column,
		d.Code,
		d.Message,
	)
}
