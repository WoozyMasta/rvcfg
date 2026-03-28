// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"

	"github.com/woozymasta/lintkit/lint"
)

// Diagnostic describes parser/lexer/preprocessor issue.
type Diagnostic struct {
	// Code is stable machine-readable identifier.
	Code lint.Code `json:"code,omitempty" yaml:"code,omitempty"`

	// Message is technical description of issue.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Severity defines impact of diagnostic.
	Severity lint.Severity `json:"severity,omitzero" yaml:"severity,omitempty"`

	// Start is start location.
	Start lint.Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is end location.
	End lint.Position `json:"end,omitzero" yaml:"end,omitempty"`
}

// Error renders diagnostic in compact form.
func (d Diagnostic) Error() string {
	if d.Start.File == "" {
		return fmt.Sprintf("%d: %s", d.Code, d.Message)
	}

	return fmt.Sprintf(
		"%s:%d:%d: %s: %s",
		d.Start.File,
		d.Start.Line,
		d.Start.Column,
		lint.FormatCode(d.Code),
		d.Message,
	)
}

// LintDiagnostic converts one rvcfg diagnostic into lintkitmodel.
func (d Diagnostic) LintDiagnostic() lint.Diagnostic {
	start := lint.Position{
		File:   d.Start.File,
		Line:   d.Start.Line,
		Column: d.Start.Column,
		Offset: d.Start.Offset,
	}
	end := lint.Position{
		File:   d.End.File,
		Line:   d.End.Line,
		Column: d.End.Column,
		Offset: d.End.Offset,
	}

	if end.File == "" && end.Line == 0 && end.Column == 0 && end.Offset == 0 {
		end = start
	}

	path := start.File
	if path == "" {
		path = end.File
	}

	return lint.Diagnostic{
		RuleID:   LintRuleID(d.Code),
		Severity: d.Severity,
		Message:  d.Message,
		Path:     path,
		Start:    start,
		End:      end,
	}
}
