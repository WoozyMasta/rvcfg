package rvcfg

import (
	"reflect"
	"testing"
)

func TestPreprocessConfigMacroParityWithExplainedSample(t *testing.T) {
	t.Parallel()

	sourcePath := testDataPath("cases", "macros", "config.cpp")
	explainedPath := testDataPath("cases", "macros", "config-explained.cpp")

	preprocessed, err := PreprocessFile(sourcePath, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", sourcePath, err)
	}

	preParsed, err := ParseBytes(sourcePath, []byte(preprocessed.Text), ParseOptions{
		CaptureScalarRaw:             true,
		AutoFixMissingClassSemicolon: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes(preprocessed) error: %v", err)
	}

	assertNoErrorDiagnostics(t, preParsed.Diagnostics, "preprocessed parse")

	explainedData := readTestFile(t, explainedPath)
	expParsed, err := ParseBytes(explainedPath, explainedData, ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes(%s) error: %v", explainedPath, err)
	}

	assertNoErrorDiagnostics(t, expParsed.Diagnostics, "explained parse")

	gotMetrics := collectASTMetrics(preParsed.File.Statements)
	wantMetrics := collectASTMetrics(expParsed.File.Statements)
	if !reflect.DeepEqual(gotMetrics, wantMetrics) {
		t.Fatalf("macro parity mismatch: metrics differ\nwant=%+v\ngot=%+v", wantMetrics, gotMetrics)
	}
}

// assertNoErrorDiagnostics fails test when at least one error diagnostic is present.
func assertNoErrorDiagnostics(t *testing.T, diagnostics []Diagnostic, stage string) {
	t.Helper()

	for _, diagnostic := range diagnostics {
		if diagnostic.Severity != SeverityError {
			continue
		}

		t.Fatalf("%s has error diagnostic: %s", stage, diagnostic.Error())
	}
}

// astMetrics stores coarse semantic counters for macro parity comparison.
type astMetrics struct {
	Classes       int
	Forward       int
	Delete        int
	Extern        int
	Properties    int
	ArraySet      int
	ArrayAppend   int
	ArrayLiterals int
}

// collectASTMetrics gathers recursive parser counters for parity baseline.
func collectASTMetrics(statements []Statement) astMetrics {
	metrics := astMetrics{}
	collectMetricsInto(statements, &metrics)

	return metrics
}

// collectMetricsInto updates counters recursively.
func collectMetricsInto(statements []Statement, metrics *astMetrics) {
	for _, statement := range statements {
		switch statement.Kind {
		case NodeClass:
			if statement.Class == nil {
				continue
			}

			metrics.Classes++
			if statement.Class.Forward {
				metrics.Forward++
			}

			collectMetricsInto(statement.Class.Body, metrics)
		case NodeDelete:
			metrics.Delete++
		case NodeExtern:
			metrics.Extern++
		case NodeProperty:
			metrics.Properties++
			if statement.Property != nil {
				countArrayLiterals(statement.Property.Value, metrics)
			}
		case NodeArrayAssign:
			if statement.ArrayAssign == nil {
				continue
			}

			if statement.ArrayAssign.Append {
				metrics.ArrayAppend++
			} else {
				metrics.ArraySet++
			}

			countArrayLiterals(statement.ArrayAssign.Value, metrics)
		}
	}
}

// countArrayLiterals counts nested ValueArray literals recursively.
func countArrayLiterals(value Value, metrics *astMetrics) {
	if value.Kind != ValueArray {
		return
	}

	metrics.ArrayLiterals++
	for _, element := range value.Elements {
		countArrayLiterals(element, metrics)
	}
}
