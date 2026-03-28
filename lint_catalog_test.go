package rvcfg

import (
	"testing"

	"github.com/woozymasta/lintkit/lint"
	"github.com/woozymasta/lintkit/linttest"
)

func TestLintRuleSpecsMatchCatalog(t *testing.T) {
	t.Parallel()

	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		t.Fatalf("getDiagnosticCodeCatalog() error: %v", err)
	}

	linttest.AssertCodeCatalogContract(t, catalog)
}

func TestLintRuleIDMatchesCodeCatalog(t *testing.T) {
	t.Parallel()

	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		t.Fatalf("getDiagnosticCodeCatalog() error: %v", err)
	}

	for _, spec := range DiagnosticCatalog() {
		expected, err := catalog.RuleID(spec.Code)
		if err != nil {
			t.Fatalf("catalog.RuleID(%d) error: %v", spec.Code, err)
		}

		got := LintRuleID(Code(spec.Code))
		if got != expected {
			t.Fatalf(
				"LintRuleID(%d)=%q, want %q",
				spec.Code,
				got,
				expected,
			)
		}
	}
}

func TestDiagnosticCatalogIntegrity(t *testing.T) {
	t.Parallel()

	catalog := DiagnosticCatalog()
	if len(catalog) == 0 {
		t.Fatal("expected non-empty diagnostics catalog")
	}

	seen := make(map[lint.Code]struct{}, len(catalog))

	for _, spec := range catalog {
		if spec.Code == 0 {
			t.Fatal("diagnostic catalog contains empty code")
		}

		if spec.Stage == "" {
			t.Fatalf("diagnostic %d has empty stage", spec.Code)
		}

		if spec.Severity == "" {
			t.Fatalf("diagnostic %d has empty severity", spec.Code)
		}

		if spec.Message == "" {
			t.Fatalf("diagnostic %d has empty message", spec.Code)
		}

		if _, ok := seen[spec.Code]; ok {
			t.Fatalf("diagnostic catalog has duplicate code %d", spec.Code)
		}

		seen[spec.Code] = struct{}{}

		lookup, ok := DiagnosticByCode(Code(spec.Code))
		if !ok {
			t.Fatalf("DiagnosticByCode(%d) returned not found", spec.Code)
		}

		if lookup != spec {
			t.Fatalf(
				"DiagnosticByCode(%d) returned different spec: %+v != %+v",
				spec.Code,
				lookup,
				spec,
			)
		}
	}
}

func TestDiagnosticByCodeUnknown(t *testing.T) {
	t.Parallel()

	if _, ok := DiagnosticByCode(0); ok {
		t.Fatal("expected unknown diagnostic code lookup to fail")
	}
}
