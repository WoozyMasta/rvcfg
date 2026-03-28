package rvcfg

import (
	"context"
	"testing"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/lintkit/linttest"
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

func TestLintRuleSpecsMatchCatalog(t *testing.T) {
	t.Parallel()

	linttest.AssertCatalogContract(
		t,
		LintModule,
		DiagnosticCatalog(),
		LintRuleSpecs(),
		LintRuleID,
	)
}

func TestRegisterLintRules(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	if err := RegisterLintRules(&registrar); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	catalog := DiagnosticCatalog()
	if len(registrar.runners) != len(catalog) {
		t.Fatalf("registered runners=%d, want %d", len(registrar.runners), len(catalog))
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
		t.Fatalf("registered runners=%d, want %d", len(registrar.runners), len(catalog))
	}
}

func TestRegisterLintRulesNilRegistrar(t *testing.T) {
	t.Parallel()

	if err := RegisterLintRules(nil); err != ErrNilLintRuleRegistrar {
		t.Fatalf("RegisterLintRules(nil) error=%v, want %v", err, ErrNilLintRuleRegistrar)
	}
}

func TestLintRuleRunnerCheck(t *testing.T) {
	t.Parallel()

	var registrar lintRulesTestRegistrar
	if err := RegisterLintRules(&registrar); err != nil {
		t.Fatalf("RegisterLintRules() error: %v", err)
	}

	runner, ok := findRunnerByRuleID(registrar.runners, LintRuleID(CodeParExpectedValue))
	if !ok {
		t.Fatalf("runner for %q not found", LintRuleID(CodeParExpectedValue))
	}

	runContext := lint.RunContext{
		TargetPath: "config.cpp",
	}
	AttachLintDiagnostics(&runContext, []Diagnostic{
		{
			Code:     CodeParExpectedValue,
			Message:  "expected value",
			Severity: lint.SeverityError,
			Start: lint.Position{
				File:   "config.cpp",
				Line:   12,
				Column: 5,
			},
		},
		{
			Code:     CodeParExpectedValue,
			Message:  "expected value again",
			Severity: lint.SeverityError,
			Start: lint.Position{
				File:   "config.cpp",
				Line:   18,
				Column: 2,
			},
		},
		{
			Code:     CodePPMacroRedefined,
			Message:  "macro redefined",
			Severity: lint.SeverityWarning,
			Start: lint.Position{
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
		if diagnostics[index].RuleID != LintRuleID(CodeParExpectedValue) {
			t.Fatalf("Diagnostics[%d].RuleID=%q", index, diagnostics[index].RuleID)
		}
	}
}

func TestDiagnosticLintDiagnostic(t *testing.T) {
	t.Parallel()

	converted := (Diagnostic{
		Code:     999999,
		Message:  "message",
		Severity: lint.SeverityInfo,
		Start: lint.Position{
			File:   "source.cpp",
			Line:   1,
			Column: 2,
		},
	}).LintDiagnostic()

	if converted.RuleID != "rvcfg.unknown" {
		t.Fatalf("RuleID=%q, want rvcfg.unknown", converted.RuleID)
	}

	if converted.Severity != lint.SeverityInfo {
		t.Fatalf("lint.Severity=%q, want %q", converted.Severity, lint.SeverityInfo)
	}

	if converted.Path != "source.cpp" {
		t.Fatalf("Path=%q, want source.cpp", converted.Path)
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
