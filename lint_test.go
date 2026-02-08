package rvcfg

import (
	"testing"
)

// lintTestRule emits one issue per file for API contract verification.
type lintTestRule struct{}

// RuleID returns stable rule-local code.
func (r *lintTestRule) RuleID() string {
	return "PATH001"
}

// CheckFile emits one diagnostic for provided file.
func (r *lintTestRule) CheckFile(ctx *LintContext, file File, emit LintEmit) {
	emit(LintIssue{
		Code:     "FILE001",
		Message:  "missing file",
		Severity: SeverityError,
		Start: Position{
			Line:   7,
			Column: 3,
		},
	})
}

// lintFallbackRule emits issue without local code to test RuleID fallback.
type lintFallbackRule struct{}

// RuleID returns fallback code suffix.
func (r *lintFallbackRule) RuleID() string {
	return "SCHEMA001"
}

// CheckFile emits warning issue without explicit code.
func (r *lintFallbackRule) CheckFile(ctx *LintContext, file File, emit LintEmit) {
	emit(LintIssue{
		Message: "schema warning",
	})
}

// lintFinalizeRule emits one aggregate issue from Finalize hook.
type lintFinalizeRule struct {
	seen int
}

// RuleID returns rule code suffix.
func (r *lintFinalizeRule) RuleID() string {
	return "AGG001"
}

// CheckFile increments file counter in rule state.
func (r *lintFinalizeRule) CheckFile(ctx *LintContext, file File, emit LintEmit) {
	r.seen++
}

// Finalize emits aggregate warning with number of processed files.
func (r *lintFinalizeRule) Finalize(ctx *LintContext, emit LintEmit) {
	emit(LintIssue{
		Message: "aggregate issue",
	})
}

func TestLintEngineRegisterPrefixValidation(t *testing.T) {
	t.Parallel()

	engine := NewLintEngine(LintOptions{})

	err := engine.Register("lint-path", &lintTestRule{})
	if err == nil {
		t.Fatal("expected invalid prefix registration error")
	}
}

func TestLintEngineRegisterRuleIDValidation(t *testing.T) {
	t.Parallel()

	engine := NewLintEngine(LintOptions{})

	err := engine.Register("LINT", lintRuleInvalidID{})
	if err == nil {
		t.Fatal("expected invalid rule id registration error")
	}
}

func TestLintEngineRunFilesBuildsPrefixedDiagnostics(t *testing.T) {
	t.Parallel()

	engine := NewLintEngine(LintOptions{})
	if err := engine.Register("LINT_PATH", &lintTestRule{}); err != nil {
		t.Fatalf("Register error: %v", err)
	}

	file := File{
		Source: "vehicle.cpp",
	}

	result := engine.RunFiles(file)
	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
	}

	diagnostic := result.Diagnostics[0]
	if diagnostic.Code != DiagnosticCode("LINT_PATH_FILE001") {
		t.Fatalf("unexpected diagnostic code: %s", diagnostic.Code)
	}

	if diagnostic.Severity != SeverityError {
		t.Fatalf("unexpected diagnostic severity: %s", diagnostic.Severity)
	}

	if diagnostic.Start.File != "vehicle.cpp" {
		t.Fatalf("expected source fallback file vehicle.cpp, got %q", diagnostic.Start.File)
	}
}

func TestLintEngineRunFilesUsesRuleIDFallbackAndDefaultSeverity(t *testing.T) {
	t.Parallel()

	engine := NewLintEngine(LintOptions{})
	if err := engine.Register("LINT_CFG", &lintFallbackRule{}); err != nil {
		t.Fatalf("Register error: %v", err)
	}

	result := engine.RunFiles(File{
		Source: "config.cpp",
	})

	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(result.Diagnostics))
	}

	diagnostic := result.Diagnostics[0]
	if diagnostic.Code != DiagnosticCode("LINT_CFG_SCHEMA001") {
		t.Fatalf("unexpected fallback diagnostic code: %s", diagnostic.Code)
	}

	if diagnostic.Severity != SeverityWarning {
		t.Fatalf("expected warning default severity, got %s", diagnostic.Severity)
	}
}

func TestLintEngineRunFilesCallsFinalizeHook(t *testing.T) {
	t.Parallel()

	finalizeRule := &lintFinalizeRule{}
	engine := NewLintEngine(LintOptions{})
	if err := engine.Register("LINT_AGG", finalizeRule); err != nil {
		t.Fatalf("Register error: %v", err)
	}

	result := engine.RunFiles(
		File{Source: "a.cpp"},
		File{Source: "b.cpp"},
	)

	if finalizeRule.seen != 2 {
		t.Fatalf("expected finalize rule to process 2 files, got %d", finalizeRule.seen)
	}

	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected one finalize diagnostic, got %d", len(result.Diagnostics))
	}

	if result.Diagnostics[0].Code != DiagnosticCode("LINT_AGG_AGG001") {
		t.Fatalf("unexpected finalize diagnostic code: %s", result.Diagnostics[0].Code)
	}
}

// lintRuleInvalidID is intentionally malformed rule id for validation tests.
type lintRuleInvalidID struct{}

// RuleID returns invalid token to test registration checks.
func (r lintRuleInvalidID) RuleID() string {
	return "bad-id"
}

// CheckFile does nothing for invalid rule id case.
func (r lintRuleInvalidID) CheckFile(ctx *LintContext, file File, emit LintEmit) {}
