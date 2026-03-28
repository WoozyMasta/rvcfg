package rvcfg

import (
	"context"
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

// lintRulesTestRegistrar captures registered runners for tests.
type lintRulesTestRegistrar struct {
	// runners stores registered runner instances.
	runners []lint.RuleRunner
}

// Register appends all provided runners into local test storage.
func (registrar *lintRulesTestRegistrar) Register(
	runners ...lint.RuleRunner,
) error {
	registrar.runners = append(registrar.runners, runners...)
	return nil
}

func TestRegisterLintRules(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	if err := RegisterLintRules(&registrar); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	catalog := DiagnosticCatalog()
	if len(registrar.runners) != len(catalog) {
		t.Fatalf(
			"registered runners=%d, want %d",
			len(registrar.runners),
			len(catalog),
		)
	}
}

func TestLintRulesProviderRegisterRules(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	provider := LintRulesProvider{}
	if err := provider.RegisterRules(&registrar); err != nil {
		t.Fatalf("LintRulesProvider.RegisterRules() error: %v", err)
	}

	catalog := DiagnosticCatalog()
	if len(registrar.runners) != len(catalog) {
		t.Fatalf(
			"registered runners=%d, want %d",
			len(registrar.runners),
			len(catalog),
		)
	}
}

func TestRegisterLintRulesNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRules(nil); err != ErrNilLintRuleRegistrar {
		t.Fatalf(
			"RegisterLintRules(nil) error=%v, want %v",
			err,
			ErrNilLintRuleRegistrar,
		)
	}
}

func TestRegisterLintRulesByStage(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	if err := RegisterLintRulesByStage(&registrar, StageParse); err != nil {
		t.Fatalf("RegisterLintRulesByStage() error: %v", err)
	}

	if len(registrar.runners) == 0 {
		t.Fatal("registered runners=0, want parse runners")
	}

	for index := range registrar.runners {
		if registrar.runners[index].RuleSpec().Scope != string(StageParse) {
			t.Fatalf(
				"runner[%d].scope=%q, want %q",
				index,
				registrar.runners[index].RuleSpec().Scope,
				StageParse,
			)
		}
	}

	if _, ok := findRunnerByRuleID(
		registrar.runners,
		LintRuleID(CodePPIncludeNotFound),
	); ok {
		t.Fatal("preprocess rule registered in parse-only stage filter")
	}
}

func TestRegisterLintRulesByStageNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRulesByStage(nil, StageParse); err != ErrNilLintRuleRegistrar {
		t.Fatalf(
			"RegisterLintRulesByStage(nil) error=%v, want %v",
			err,
			ErrNilLintRuleRegistrar,
		)
	}
}

func TestLintRulesProviderSupportsStageRegistration(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	provider := LintRulesProvider{}
	err := lint.RegisterRuleProvidersByStage(
		&registrar,
		[]lint.Stage{lint.Stage(StageParse)},
		provider,
	)
	if err != nil {
		t.Fatalf("RegisterRuleProvidersByStage(provider) error: %v", err)
	}

	if len(registrar.runners) == 0 {
		t.Fatal("registered runners=0, want parse runners")
	}

	for index := range registrar.runners {
		if registrar.runners[index].RuleSpec().Scope != string(StageParse) {
			t.Fatalf(
				"runner[%d].scope=%q, want %q",
				index,
				registrar.runners[index].RuleSpec().Scope,
				StageParse,
			)
		}
	}

	if _, ok := findRunnerByRuleID(
		registrar.runners,
		LintRuleID(CodePPIncludeNotFound),
	); ok {
		t.Fatal("preprocess rule registered in parse-only stage filter")
	}
}

func TestLintRuleRunnerCheck(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	if err := RegisterLintRules(&registrar); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	runner, ok := findRunnerByRuleID(
		registrar.runners,
		LintRuleID(CodeParExpectedValue),
	)
	if !ok {
		t.Fatalf(
			"runner for %q not found",
			LintRuleID(CodeParExpectedValue),
		)
	}

	runContext := lint.RunContext{
		TargetPath: "config.cpp",
	}
	AttachLintDiagnostics(&runContext, []Diagnostic{
		{
			Code:     CodeParExpectedValue,
			Message:  "expected value",
			Severity: SeverityError,
			Start: Position{
				File:   "config.cpp",
				Line:   12,
				Column: 5,
			},
		},
		{
			Code:     CodeParExpectedValue,
			Message:  "expected value again",
			Severity: SeverityError,
			Start: Position{
				File:   "config.cpp",
				Line:   18,
				Column: 2,
			},
		},
		{
			Code:     CodePPMacroRedefined,
			Message:  "macro redefined",
			Severity: SeverityWarning,
			Start: Position{
				File:   "config.cpp",
				Line:   40,
				Column: 1,
			},
		},
	})

	diagnostics := make([]lint.Diagnostic, 0, 2)
	err := runner.Check(
		context.Background(),
		&runContext,
		func(diagnostic lint.Diagnostic) {
			diagnostics = append(diagnostics, diagnostic)
		},
	)
	if err != nil {
		t.Fatalf("runner.Check() error: %v", err)
	}

	if len(diagnostics) != 2 {
		t.Fatalf("len(Diagnostics)=%d, want 2", len(diagnostics))
	}

	for index := range diagnostics {
		if diagnostics[index].RuleID !=
			LintRuleID(CodeParExpectedValue) {
			t.Fatalf("Diagnostics[%d].RuleID=%q", index, diagnostics[index].RuleID)
		}
	}
}

// findRunnerByRuleID returns runner by stable rule id.
func findRunnerByRuleID(
	runners []lint.RuleRunner,
	ruleID string,
) (lint.RuleRunner, bool) {
	for index := range runners {
		if runners[index].RuleSpec().ID == ruleID {
			return runners[index], true
		}
	}

	return nil, false
}
