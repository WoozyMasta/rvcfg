package rvcfg

import "testing"

func TestLexFilePlainRVMAT(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "lex", "plain.rvmat")
	tokens, diagnostics, err := LexFile(path)
	if err != nil {
		t.Fatalf("LexFile(%s) error: %v", path, err)
	}

	assertNoLexErrorDiagnostics(t, diagnostics)

	if len(tokens) == 0 {
		t.Fatal("expected non-empty token stream")
	}
}

func TestLexFileLargeConfig(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "vehicle", "config.cpp")
	tokens, diagnostics, err := LexFile(path)
	if err != nil {
		t.Fatalf("LexFile(%s) error: %v", path, err)
	}

	assertNoLexErrorDiagnostics(t, diagnostics)

	if len(tokens) < 100 {
		t.Fatalf("expected large token stream, got %d tokens", len(tokens))
	}
}

func TestLexFileModelCfg(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "model", "model.cfg")
	tokens, diagnostics, err := LexFile(path)
	if err != nil {
		t.Fatalf("LexFile(%s) error: %v", path, err)
	}

	assertNoLexErrorDiagnostics(t, diagnostics)

	if len(tokens) == 0 {
		t.Fatal("expected non-empty token stream")
	}
}

func TestLexFileLightingUnderground(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "macros", "lighting_underground.txt")
	tokens, diagnostics, err := LexFile(path)
	if err != nil {
		t.Fatalf("LexFile(%s) error: %v", path, err)
	}

	assertNoLexErrorDiagnostics(t, diagnostics)

	if len(tokens) == 0 {
		t.Fatal("expected non-empty token stream")
	}
}

// assertNoLexErrorDiagnostics fails when lexer produced at least one error-level diagnostic.
func assertNoLexErrorDiagnostics(t *testing.T, diagnostics []Diagnostic) {
	t.Helper()

	for _, diagnostic := range diagnostics {
		if diagnostic.Severity != SeverityError {
			continue
		}

		t.Fatalf("unexpected error diagnostic: %s", diagnostic.Error())
	}
}
