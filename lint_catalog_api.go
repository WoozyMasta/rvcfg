// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"github.com/woozymasta/lintkit/lint"
)

const (
	// Module is stable lint module namespace for rvcfg rules.
	Module = "rvcfg"
)

var diagnosticCodeCatalogHandle = lint.NewCodeCatalogHandle(
	lint.CodeCatalogConfig{
		Module:            Module,
		CodePrefix:        "CFG",
		ModuleName:        "Real Virtuality Configs",
		ModuleDescription: "Lint rules for Real Virtuality config lexer, parser and preprocessor flows.",
		ScopeDescriptions: map[lint.Stage]string{
			"lex":        "Lexer diagnostics.",
			"parse":      "Parser diagnostics.",
			"preprocess": "Preprocessor diagnostics.",
		},
	},
	diagnosticCatalog,
)

// getDiagnosticCodeCatalog returns lazy-initialized diagnostics catalog.
func getDiagnosticCodeCatalog() (lint.CodeCatalog, error) {
	return diagnosticCodeCatalogHandle.Catalog()
}

// DiagnosticRuleSpec converts one diagnostic spec into lint rule metadata.
func DiagnosticRuleSpec(spec lint.CodeSpec) (lint.RuleSpec, error) {
	return diagnosticCodeCatalogHandle.RuleSpec(spec)
}

// LintRuleID returns lint rule ID mapped from stable rvcfg diagnostic code.
func LintRuleID(code Code) string {
	return diagnosticCodeCatalogHandle.RuleIDOrUnknown(code)
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []lint.CodeSpec {
	return diagnosticCodeCatalogHandle.CodeSpecs()
}

// DiagnosticByCode returns diagnostic metadata for code.
func DiagnosticByCode(code Code) (lint.CodeSpec, bool) {
	return diagnosticCodeCatalogHandle.ByCode(code)
}

// LintRuleSpecs returns deterministic lint rule specs from diagnostics catalog.
func LintRuleSpecs() []lint.RuleSpec {
	return diagnosticCodeCatalogHandle.RuleSpecs()
}
