// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"sync"

	"github.com/woozymasta/lintkit/lint"
)

const (
	// lintRunValueByCodeKey stores grouped diagnostics map in run values.
	lintRunValueByCodeKey = "lint.by_code"
)

var (
	// lintBindingState stores lazy-initialized code-catalog binding state.
	lintBindingState struct {
		// once guards one-time binding construction.
		once sync.Once

		// binding stores reusable register+attach helper.
		binding lint.CodeCatalogBinding[Diagnostic]

		// err stores binding construction error.
		err error
	}
)

// LintRulesProvider registers rvcfg diagnostic rules into any RuleRegistrar.
type LintRulesProvider struct{}

// RegisterRules adds provider-owned rules to target registrar.
func (provider LintRulesProvider) RegisterRules(
	registrar lint.RuleRegistrar,
) error {
	return RegisterLintRules(registrar)
}

// RegisterRulesByScope adds provider-owned rules filtered by scope tokens.
func (provider LintRulesProvider) RegisterRulesByScope(
	registrar lint.RuleRegistrar,
	scopes ...string,
) error {
	return RegisterLintRulesByScope(registrar, scopes...)
}

// RegisterRulesByStage adds provider-owned rules filtered by stage tokens.
func (provider LintRulesProvider) RegisterRulesByStage(
	registrar lint.RuleRegistrar,
	stages ...Stage,
) error {
	return RegisterLintRulesByStage(registrar, stages...)
}

// RegisterLintRules registers stable rvcfg diagnostic rules into registrar.
func RegisterLintRules(registrar lint.RuleRegistrar) error {
	return registerLintRulesWithBinding(registrar, func(
		binding lint.CodeCatalogBinding[Diagnostic],
	) error {
		return binding.RegisterRules(registrar)
	})
}

// RegisterLintRulesByScope registers rvcfg rules filtered by scope tokens.
func RegisterLintRulesByScope(
	registrar lint.RuleRegistrar,
	scopes ...string,
) error {
	return registerLintRulesWithBinding(registrar, func(
		binding lint.CodeCatalogBinding[Diagnostic],
	) error {
		return binding.RegisterRulesByScope(registrar, scopes...)
	})
}

// RegisterLintRulesByStage registers rvcfg rules filtered by stage tokens.
func RegisterLintRulesByStage(
	registrar lint.RuleRegistrar,
	stages ...Stage,
) error {
	return registerLintRulesWithBinding(registrar, func(
		binding lint.CodeCatalogBinding[Diagnostic],
	) error {
		return binding.RegisterRulesByStage(registrar, stages...)
	})
}

// registerLintRulesWithBinding validates registrar and executes binding callback.
func registerLintRulesWithBinding(
	registrar lint.RuleRegistrar,
	register func(binding lint.CodeCatalogBinding[Diagnostic]) error,
) error {
	if registrar == nil {
		return ErrNilLintRuleRegistrar
	}

	binding, err := getLintBinding()
	if err != nil {
		return err
	}

	return register(binding)
}

// AttachLintDiagnostics stores diagnostics in run context values.
func AttachLintDiagnostics(
	run *lint.RunContext,
	diagnostics []Diagnostic,
) {
	binding, err := getLintBinding()
	if err != nil {
		return
	}

	_ = binding.Attach(run, diagnostics)
}

// getLintBinding returns lazy-initialized code-catalog binding helper.
func getLintBinding() (lint.CodeCatalogBinding[Diagnostic], error) {
	lintBindingState.once.Do(func() {
		catalog, err := getDiagnosticCodeCatalog()
		if err != nil {
			lintBindingState.err = err
			return
		}

		lintBindingState.binding, lintBindingState.err = lint.NewCodeCatalogBinding(
			lint.CodeCatalogBindingConfig[Diagnostic]{
				RunValueKey:        lintRunValueByCodeKey,
				Catalog:            catalog,
				CodeFromDiagnostic: lintDiagnosticCode,
				DiagnosticToLint:   lintDiagnostic,
				UnknownCodePolicy:  lint.UnknownCodeDrop,
			},
		)
	})

	if lintBindingState.err != nil {
		return lint.CodeCatalogBinding[Diagnostic]{}, lintBindingState.err
	}

	return lintBindingState.binding, nil
}

// lintDiagnosticCode extracts numeric code from internal diagnostic item.
func lintDiagnosticCode(item Diagnostic) lint.Code {
	return item.Code
}

// lintDiagnostic converts one rvcfg diagnostic into lintkit diagnostic.
func lintDiagnostic(diagnostic Diagnostic) lint.Diagnostic {
	start := diagnostic.Start
	end := diagnostic.End

	if end.File == "" && end.Line == 0 && end.Column == 0 && end.Offset == 0 {
		end = start
	}

	path := start.File
	if path == "" {
		path = end.File
	}

	return lint.Diagnostic{
		Severity: diagnostic.Severity,
		Message:  diagnostic.Message,
		Path:     path,
		Start:    start,
		End:      end,
	}
}
