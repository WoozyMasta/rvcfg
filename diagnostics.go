// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
)

// Diagnostic describes parser/lexer/preprocessor issue.
type Diagnostic struct {

	// Message is technical description of issue.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Severity defines impact of diagnostic.
	Severity Severity `json:"severity,omitzero" yaml:"severity,omitempty"`

	// Start is start location.
	Start Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is end location.
	End Position `json:"end,omitzero" yaml:"end,omitempty"`
	// Code is stable machine-readable identifier.
	Code Code `json:"code,omitempty" yaml:"code,omitempty"`
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
		FormatCode(d.Code),
		d.Message,
	)
}
