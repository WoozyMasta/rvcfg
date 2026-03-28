package rvcfg

import (
	"testing"
)

func TestAnalyzePAR026Emitted(t *testing.T) {
	t.Parallel()

	source := `
class ParentBase
{
	class DamageSystem
	{
		class DamageZones
		{
			class Torso {};
			class Head {};
		};
	};
};

class Parent: ParentBase
{
	class DamageSystem
	{
		class DamageZones: DamageZones
		{
			class Torso
			{
				class Health
				{
					hitpoints = 750000;
				};
			};
		};
	};
};
`

	parseResult := parseForAnalyze(t, "derived-nested-no-base.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	found := false
	for _, diagnostic := range diagnostics {
		if diagnostic.Code != CodeParDerivedNestedClassWithoutBase {
			continue
		}

		if diagnostic.Severity != SeverityInfo {
			t.Fatalf(
				"diagnostic severity=%q, want %q",
				diagnostic.Severity,
				SeverityInfo,
			)
		}

		found = true

		break
	}

	if !found {
		t.Fatalf("expected PAR026 info diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR026NoEmitExplicitInheritance(t *testing.T) {
	t.Parallel()

	source := `
class ParentBase
{
	class DamageSystem
	{
		class DamageZones
		{
			class Torso {};
		};
	};
};

class Parent: ParentBase
{
	class DamageSystem: DamageSystem
	{
		class DamageZones: DamageZones
		{
			class Torso: Torso
			{
				class Health: Health
				{
					hitpoints = 750000;
				};
			};
		};
	};
};
`

	parseResult := parseForAnalyze(t, "derived-nested-explicit-base.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	for _, diagnostic := range diagnostics {
		if diagnostic.Code == CodeParDerivedNestedClassWithoutBase {
			t.Fatalf("unexpected PAR026 diagnostic: %v", diagnostic)
		}
	}
}

func TestAnalyzePAR026NoEmitNonDerivedParent(t *testing.T) {
	t.Parallel()

	source := `
class Parent
{
	class DamageSystem
	{
		class DamageZones
		{
			class Torso
			{
				hitpoints = 750000;
			};
		};
	};
};
`

	parseResult := parseForAnalyze(t, "non-derived-parent.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	for _, diagnostic := range diagnostics {
		if diagnostic.Code == CodeParDerivedNestedClassWithoutBase {
			t.Fatalf("unexpected PAR026 diagnostic: %v", diagnostic)
		}
	}
}

func TestAnalyzePAR027UnsupportedScalar(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1 + 2;
};
`

	parseResult := parseForAnalyze(t, "rap-unsupported-scalar.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPUnsupportedScalar) {
		t.Fatalf("expected PAR027 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR028PrecisionLoss(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1234.56789;
};
`

	parseResult := parseForAnalyze(t, "rap-float-precision.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPFloatPrecisionLoss) {
		t.Fatalf("expected PAR028 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR029UnderflowToZero(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1e-60;
};
`

	parseResult := parseForAnalyze(t, "rap-float-underflow.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPFloatUnderflowToZero) {
		t.Fatalf("expected PAR029 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR030UnsafeEscape(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = "DZ\gear\consumables\data\\"\".rvmat";
};
`

	parseResult := parseForAnalyze(t, "rap-unsafe-string-escape.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPUnsafeStringEscape) {
		t.Fatalf("expected PAR030 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR030NoWarnBICompatibleEscape(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = "DZ\gear\consumables\data\"""".rvmat";
};
`

	parseResult := parseForAnalyze(t, "rap-bi-compatible-escape.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if hasDiagnosticCode(diagnostics, CodeParRAPUnsafeStringEscape) {
		t.Fatalf("unexpected PAR030 diagnostic for BI-compatible escape form: %v", diagnostics)
	}
}

func TestAnalyzePAR031ExtremeMagnitude(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1e30;
};
`

	parseResult := parseForAnalyze(t, "rap-extreme-float.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPExtremeFloatMagnitude) {
		t.Fatalf("expected PAR031 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR032OverflowToInf(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1e39;
};
`

	parseResult := parseForAnalyze(t, "rap-float-overflow.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPFloatOverflowToInf) {
		t.Fatalf("expected PAR032 diagnostic, got %v", diagnostics)
	}

	if hasDiagnosticCode(diagnostics, CodeParRAPUnsupportedScalar) {
		t.Fatalf("unexpected PAR027 diagnostic for float32 overflow: %v", diagnostics)
	}
}

func TestAnalyzePAR033NonFiniteFloat(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = inf;
};
`

	parseResult := parseForAnalyze(t, "rap-non-finite-float.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if !hasDiagnosticCode(diagnostics, CodeParRAPNonFiniteFloat) {
		t.Fatalf("expected PAR033 diagnostic, got %v", diagnostics)
	}
}

func TestAnalyzePAR028NoWarnHumanFloat(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 0.3;
};
`

	parseResult := parseForAnalyze(t, "rap-human-float.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if hasDiagnosticCode(diagnostics, CodeParRAPFloatPrecisionLoss) {
		t.Fatalf("unexpected PAR028 diagnostic for 0.3 literal: %v", diagnostics)
	}
}

func TestAnalyzePAR029NoWarnMinSubnormal(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1e-45;
};
`

	parseResult := parseForAnalyze(t, "rap-float-min-subnormal.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{})

	if hasDiagnosticCode(diagnostics, CodeParRAPFloatUnderflowToZero) {
		t.Fatalf("unexpected PAR029 diagnostic for float32 min subnormal: %v", diagnostics)
	}
}

func TestAnalyzeDisableOptionalPasses(t *testing.T) {
	t.Parallel()

	source := `
class ParentBase
{
	class DamageSystem
	{
		class DamageZones
		{
			class Torso {};
		};
	};
};

class Parent: ParentBase
{
	class DamageSystem
	{
		class DamageZones: DamageZones
		{
			class Torso
			{
				class Health
				{
					hitpoints = 750000;
				};
			};
		};
	};
};
`

	parseResult := parseForAnalyze(t, "analyze-disable-flags.cpp", source)
	diagnostics := AnalyzeFile(parseResult.File, []byte(source), AnalyzeOptions{
		DisableInheritanceHints: true,
		DisableRAPScalarHints:   true,
	})
	if len(diagnostics) != 0 {
		t.Fatalf("expected no diagnostics when all optional passes disabled, got %v", diagnostics)
	}
}

// parseForAnalyze parses source and fails test on parse errors.
func parseForAnalyze(t *testing.T, filename string, source string) ParseResult {
	t.Helper()

	parseResult, err := ParseBytes(filename, []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes(%s) error: %v", filename, err)
	}

	return parseResult
}

// hasDiagnosticCode reports whether diagnostics list contains diagnostic code.
func hasDiagnosticCode(diagnostics []Diagnostic, code Code) bool {
	for index := range diagnostics {
		if diagnostics[index].Code == code {
			return true
		}
	}

	return false
}
