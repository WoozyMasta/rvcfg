// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import lint "github.com/woozymasta/lintkit/lint"

const (
	// SeverityError marks hard diagnostics.
	SeverityError Severity = "error"

	// SeverityWarning marks non-fatal diagnostics.
	SeverityWarning Severity = "warning"

	// SeverityInfo marks informational diagnostics.
	SeverityInfo Severity = "info"

	// SeverityNotice marks low-priority diagnostics.
	SeverityNotice Severity = "notice"
)

// Severity defines normalized diagnostic severity.
type Severity = lint.Severity

// Code defines stable numeric diagnostic code.
type Code = lint.Code

// Stage defines diagnostic pipeline stage token.
type Stage = lint.Stage

// Position stores one source position in file-oriented diagnostics.
type Position = lint.Position

// FormatCode formats numeric diagnostic code token as base-10 string.
func FormatCode(code Code) string {
	return lint.FormatCode(code)
}
