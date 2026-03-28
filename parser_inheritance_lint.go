// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"

	"github.com/woozymasta/lintkit/lint"
)

// collectInheritanceHints emits info diagnostics for potentially risky inheritance overrides.
func collectInheritanceHints(file File) []Diagnostic {
	diagnostics := make([]Diagnostic, 0)
	walkDerivedNestedClassWithoutBase(file.Statements, nil, &diagnostics)

	return diagnostics
}

// walkDerivedNestedClassWithoutBase recursively walks class statements.
func walkDerivedNestedClassWithoutBase(statements []Statement, parentClass *ClassDecl, diagnostics *[]Diagnostic) {
	for index := range statements {
		stmt := &statements[index]
		if stmt.Kind != NodeClass || stmt.Class == nil || stmt.Class.Forward {
			continue
		}

		if parentClass != nil && parentClass.Base != "" && stmt.Class.Base == "" {
			*diagnostics = append(*diagnostics, Diagnostic{
				Code: CodeParDerivedNestedClassWithoutBase,
				Message: fmt.Sprintf(
					"nested class %q in derived class %q has no explicit inheritance and may replace parent subtree",
					stmt.Class.Name,
					parentClass.Name,
				),
				Severity: lint.SeverityInfo,
				Start:    stmt.Start,
				End:      stmt.End,
			})
		}

		walkDerivedNestedClassWithoutBase(
			stmt.Class.Body,
			stmt.Class,
			diagnostics,
		)
	}
}
