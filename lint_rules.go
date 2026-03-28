// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"github.com/woozymasta/lintkit/lint"
)

const (
	// LintModule is stable lint module namespace for rvcfg rules.
	LintModule = "rvcfg"
)

const (
	// lintRunValueByCodeKey stores grouped diagnostics map in run values.
	lintRunValueByCodeKey = "lint.by_code"
)

// LintRulesProvider registers rvcfg diagnostic rules into any RuleRegistrar.
type LintRulesProvider struct{}

// RegisterRules adds provider-owned rules to target registrar.
func (provider LintRulesProvider) RegisterRules(
	registrar lint.RuleRegistrar,
) error {
	return RegisterLintRules(registrar)
}

// RegisterLintRules registers stable rvcfg diagnostic rules into registrar.
func RegisterLintRules(registrar lint.RuleRegistrar) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return err
	}

	provider, err := lint.NewCodeCatalogProvider(
		lintRunValueByCodeKey,
		catalog,
		Diagnostic.LintDiagnostic,
	)
	if err != nil {
		return err
	}

	return provider.RegisterRules(registrar)
}

// LintRuleSpecs returns deterministic lint rule specs from diagnostics catalog.
func LintRuleSpecs() []lint.RuleSpec {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return nil
	}

	return catalog.RuleSpecs()
}

// AttachLintDiagnostics stores diagnostics in run context values.
func AttachLintDiagnostics(run *lint.RunContext, diagnostics []Diagnostic) {
	lint.AttachCatalogDiagnostics(
		run,
		lintRunValueByCodeKey,
		diagnostics,
		func(item Diagnostic) lint.Code {
			return item.Code
		},
	)
}
