package rvcfg

import "testing"

func TestDiagnosticCatalogIntegrity(t *testing.T) {
	t.Parallel()

	catalog := DiagnosticCatalog()
	if len(catalog) == 0 {
		t.Fatal("expected non-empty diagnostics catalog")
	}

	seen := make(map[DiagnosticCode]struct{}, len(catalog))

	for _, spec := range catalog {
		if spec.Code == "" {
			t.Fatal("diagnostic catalog contains empty code")
		}

		if spec.Stage == "" {
			t.Fatalf("diagnostic %s has empty stage", spec.Code)
		}

		if spec.Severity == "" {
			t.Fatalf("diagnostic %s has empty severity", spec.Code)
		}

		if spec.Summary == "" {
			t.Fatalf("diagnostic %s has empty summary", spec.Code)
		}

		if _, ok := seen[spec.Code]; ok {
			t.Fatalf("diagnostic catalog has duplicate code %s", spec.Code)
		}

		seen[spec.Code] = struct{}{}

		lookup, ok := DiagnosticByCode(spec.Code)
		if !ok {
			t.Fatalf("DiagnosticByCode(%s) returned not found", spec.Code)
		}

		if lookup != spec {
			t.Fatalf("DiagnosticByCode(%s) returned different spec: %+v != %+v", spec.Code, lookup, spec)
		}
	}
}

func TestDiagnosticByCodeUnknown(t *testing.T) {
	t.Parallel()

	if _, ok := DiagnosticByCode(DiagnosticCode("UNKNOWN")); ok {
		t.Fatal("expected unknown diagnostic code lookup to fail")
	}
}
