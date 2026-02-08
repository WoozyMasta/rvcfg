package rvcfg

import (
	"errors"
	"testing"
)

func TestParseFileVehicleConfig(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "vehicle", "config.cpp")
	result, err := ParseFile(path, ParseOptions{})
	if err != nil {
		t.Fatalf("ParseFile(%s) error: %v", path, err)
	}

	if len(result.Diagnostics) > 0 {
		for _, d := range result.Diagnostics {
			if d.Severity == SeverityError {
				t.Fatalf("unexpected parse diagnostic: %s", d.Error())
			}
		}
	}

	classCount := countClasses(result.File.Statements)
	if classCount < 50 {
		t.Fatalf("expected many class declarations, got %d", classCount)
	}
}

func TestParseFileBackpacksConfig(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "backpacks", "config.cpp")
	result, err := ParseFile(path, ParseOptions{})
	if err != nil {
		t.Fatalf("ParseFile(%s) error: %v", path, err)
	}

	appendCount := countArrayAppends(result.File.Statements)
	if appendCount == 0 {
		t.Fatal("expected array append assignments in mod config sample")
	}
}

func TestParseFilePantsConfig(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "pants", "config.cpp")
	result, err := ParseFile(path, ParseOptions{})
	if err != nil {
		t.Fatalf("ParseFile(%s) error: %v", path, err)
	}

	classCount := countClasses(result.File.Statements)
	if classCount == 0 {
		t.Fatal("expected class declarations in pants sample")
	}
}

func TestParseFileModelCfg(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "parse", "model", "model.cfg")
	result, err := ParseFile(path, ParseOptions{})
	if err != nil {
		t.Fatalf("ParseFile(%s) error: %v", path, err)
	}

	classCount := countClasses(result.File.Statements)
	if classCount < 2 {
		t.Fatalf("expected model.cfg class declarations, got %d", classCount)
	}
}

func TestParseBytesSupportsTrailingCommaInNestedArray(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	class Nested
	{
		weights[] = {{1, 2,}, {3, 4,},};
	};
};
`

	result, err := ParseBytes("inline.cpp", []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes inline source error: %v", err)
	}

	arrayAssign, ok := findFirstArrayAssign(result.File.Statements)
	if !ok {
		t.Fatal("expected one array assignment in parsed AST")
	}

	if arrayAssign.Value.Kind != ValueArray || len(arrayAssign.Value.Elements) != 2 {
		t.Fatalf("unexpected array value shape: kind=%s elements=%d", arrayAssign.Value.Kind, len(arrayAssign.Value.Elements))
	}

	for idx, item := range arrayAssign.Value.Elements {
		if item.Kind != ValueArray {
			t.Fatalf("expected nested array at index %d, got %s", idx, item.Kind)
		}
	}
}

func TestParseBytesMissingSemicolon(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = 1
};
`

	result, err := ParseBytes("broken.cpp", []byte(source), ParseOptions{})
	if !errors.Is(err, ErrParse) {
		t.Fatalf("expected ErrParse, got %v", err)
	}

	foundMissing := false
	for _, d := range result.Diagnostics {
		if d.Code == CodeParMissingAssignSemicolon {
			foundMissing = true

			break
		}
	}

	if !foundMissing {
		t.Fatalf("expected PAR014 diagnostic, got %v", result.Diagnostics)
	}
}

func TestParseBytesMissingClassSemicolonAutoFixMode(t *testing.T) {
	t.Parallel()

	source := `
class Root
{
	class Child
	{
		value = 1;
	}
};
`

	_, err := ParseBytes("missing-class-semicolon-default.cpp", []byte(source), ParseOptions{})
	if !errors.Is(err, ErrParse) {
		t.Fatalf("expected ErrParse in default mode, got %v", err)
	}

	result, err := ParseBytes("missing-class-semicolon-autofix.cpp", []byte(source), ParseOptions{
		AutoFixMissingClassSemicolon: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes autofix mode error: %v", err)
	}

	if len(result.File.Statements) != 1 {
		t.Fatalf("expected one top-level class, got %d", len(result.File.Statements))
	}

	foundAutoFixWarning := false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code != CodeParAutofixClassSemicolon || diagnostic.Severity != SeverityWarning {
			continue
		}

		if diagnostic.Message == "autofix: inserted missing ';' after class declaration" {
			foundAutoFixWarning = true

			break
		}
	}

	if !foundAutoFixWarning {
		t.Fatalf("expected PAR020 autofix warning, got %v", result.Diagnostics)
	}
}

func TestParseBytesExternDeleteAndForwardClass(t *testing.T) {
	t.Parallel()

	source := `
extern class CfgPatches;
delete LegacyClass;
class ForwardOnly;
`

	result, err := ParseBytes("extern.cpp", []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes extern source error: %v", err)
	}

	if len(result.File.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(result.File.Statements))
	}

	if result.File.Statements[0].Kind != NodeExtern || result.File.Statements[0].Extern == nil {
		t.Fatalf("expected first statement extern, got %+v", result.File.Statements[0])
	}

	if result.File.Statements[1].Kind != NodeDelete || result.File.Statements[1].Delete == nil {
		t.Fatalf("expected second statement delete, got %+v", result.File.Statements[1])
	}

	if result.File.Statements[2].Kind != NodeClass || result.File.Statements[2].Class == nil {
		t.Fatalf("expected third statement class, got %+v", result.File.Statements[2])
	}

	if !result.File.Statements[2].Class.Forward {
		t.Fatal("expected forward class declaration")
	}
}

func TestParseBytesCaptureScalarRawOption(t *testing.T) {
	t.Parallel()

	source := `
class CfgTest
{
	value = -1;
};
`

	defaultResult, err := ParseBytes("raw-default.cpp", []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes default options error: %v", err)
	}

	withRawResult, err := ParseBytes("raw-enabled.cpp", []byte(source), ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes with CaptureScalarRaw error: %v", err)
	}

	defaultProp, ok := findFirstProperty(defaultResult.File.Statements)
	if !ok {
		t.Fatal("expected property assignment in default parse result")
	}

	withRawProp, ok := findFirstProperty(withRawResult.File.Statements)
	if !ok {
		t.Fatal("expected property assignment in raw-enabled parse result")
	}

	if defaultProp.Value.Raw != "" {
		t.Fatalf("expected empty Value.Raw by default, got %q", defaultProp.Value.Raw)
	}

	if withRawProp.Value.Raw != "-1" {
		t.Fatalf("expected Value.Raw to be -1, got %q", withRawProp.Value.Raw)
	}
}

func TestParseBytesDigitLeadingClassName(t *testing.T) {
	t.Parallel()

	source := `
class CfgSoundSets
{
	class 1kHz_mono_1s_SoundSet
	{
		loop = 1;
	};
};
`

	result, err := ParseBytes("digit-class.cpp", []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes digit-leading class name error: %v", err)
	}

	if len(result.File.Statements) != 1 {
		t.Fatalf("expected one top-level statement, got %d", len(result.File.Statements))
	}

	root := result.File.Statements[0]
	if root.Kind != NodeClass || root.Class == nil {
		t.Fatalf("expected top-level class, got %+v", root)
	}

	if len(root.Class.Body) != 1 {
		t.Fatalf("expected one nested class, got %d", len(root.Class.Body))
	}

	nested := root.Class.Body[0]
	if nested.Kind != NodeClass || nested.Class == nil {
		t.Fatalf("expected nested class, got %+v", nested)
	}

	if nested.Class.Name != "1kHz_mono_1s_SoundSet" {
		t.Fatalf("unexpected nested class name: %q", nested.Class.Name)
	}
}

func TestParseBytesDigitLeadingClassNameStrictMode(t *testing.T) {
	t.Parallel()

	source := `
class CfgSoundSets
{
	class 1kHz_mono_1s_SoundSet
	{
		loop = 1;
	};
};
`

	result, err := ParseBytes("digit-class-strict.cpp", []byte(source), ParseOptions{
		Strict: true,
	})
	if !errors.Is(err, ErrParse) {
		t.Fatalf("expected ErrParse in strict mode, got %v", err)
	}

	found := false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code == CodeParStrictDigitLeadingClassName {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected PAR025 in strict mode diagnostics, got %v", result.Diagnostics)
	}
}

func TestParseBytesEnumDeclaration(t *testing.T) {
	t.Parallel()

	source := `
enum ESurfaceType
{
	SURFACE_DEFAULT,
	SURFACE_DIRT = 10,
	SURFACE_WET = SURFACE_DIRT + 2,
};
`

	result, err := ParseBytes("enum.cpp", []byte(source), ParseOptions{})
	if err != nil {
		t.Fatalf("ParseBytes enum source error: %v", err)
	}

	if len(result.File.Statements) != 1 {
		t.Fatalf("expected one top-level statement, got %d", len(result.File.Statements))
	}

	stmt := result.File.Statements[0]
	if stmt.Kind != NodeEnum || stmt.Enum == nil {
		t.Fatalf("expected enum statement, got %+v", stmt)
	}

	if stmt.Enum.Name != "ESurfaceType" {
		t.Fatalf("unexpected enum name: %q", stmt.Enum.Name)
	}

	if len(stmt.Enum.Items) != 3 {
		t.Fatalf("expected 3 enum items, got %d", len(stmt.Enum.Items))
	}

	if stmt.Enum.Items[0].Name != "SURFACE_DEFAULT" || stmt.Enum.Items[0].ValueRaw != "" {
		t.Fatalf("unexpected enum item[0]: %+v", stmt.Enum.Items[0])
	}

	if stmt.Enum.Items[1].Name != "SURFACE_DIRT" || stmt.Enum.Items[1].ValueRaw != "10" {
		t.Fatalf("unexpected enum item[1]: %+v", stmt.Enum.Items[1])
	}

	if stmt.Enum.Items[2].Name != "SURFACE_WET" || stmt.Enum.Items[2].ValueRaw != "SURFACE_DIRT + 2" {
		t.Fatalf("unexpected enum item[2]: %+v", stmt.Enum.Items[2])
	}
}

func TestParseBytesEnumMissingSemicolon(t *testing.T) {
	t.Parallel()

	source := `
enum EBroken
{
	V0,
}
`

	result, err := ParseBytes("enum-missing-semicolon.cpp", []byte(source), ParseOptions{})
	if !errors.Is(err, ErrParse) {
		t.Fatalf("expected ErrParse, got %v", err)
	}

	found := false
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code == CodeParMissingEnumSemicolon {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected PAR024 diagnostic, got %v", result.Diagnostics)
	}
}

func TestParseBytesPreserveComments(t *testing.T) {
	t.Parallel()

	source := `
// top comment
class Cfg
{
	value = 1; // inline value
};
`

	result, err := ParseBytes("comments.cpp", []byte(source), ParseOptions{
		PreserveComments: true,
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes preserve comments error: %v", err)
	}

	if len(result.File.Statements) != 1 {
		t.Fatalf("expected one top-level statement, got %d", len(result.File.Statements))
	}

	root := result.File.Statements[0]
	if len(root.LeadingComments) != 1 {
		t.Fatalf("expected one leading comment on class, got %d", len(root.LeadingComments))
	}

	if root.LeadingComments[0].Text != "// top comment" {
		t.Fatalf("unexpected class leading comment text: %q", root.LeadingComments[0].Text)
	}

	if root.Kind != NodeClass || root.Class == nil || len(root.Class.Body) != 1 {
		t.Fatalf("expected class with one child statement, got %+v", root)
	}

	child := root.Class.Body[0]
	if child.TrailingComment == nil {
		t.Fatal("expected trailing comment on child statement")
	}

	if child.TrailingComment.Text != "// inline value" {
		t.Fatalf("unexpected child trailing comment text: %q", child.TrailingComment.Text)
	}
}

// countClasses counts class declarations recursively.
func countClasses(statements []Statement) int {
	count := 0
	for _, stmt := range statements {
		if stmt.Kind != NodeClass || stmt.Class == nil {
			continue
		}

		count++
		count += countClasses(stmt.Class.Body)
	}

	return count
}

// countArrayAppends counts array append assignments recursively.
func countArrayAppends(statements []Statement) int {
	count := 0
	for _, stmt := range statements {
		if stmt.Kind == NodeArrayAssign && stmt.ArrayAssign != nil && stmt.ArrayAssign.Append {
			count++
		}

		if stmt.Kind == NodeClass && stmt.Class != nil {
			count += countArrayAppends(stmt.Class.Body)
		}
	}

	return count
}

// findFirstArrayAssign returns first recursive array assignment.
func findFirstArrayAssign(statements []Statement) (ArrayAssign, bool) {
	for _, stmt := range statements {
		if stmt.Kind == NodeArrayAssign && stmt.ArrayAssign != nil {
			return *stmt.ArrayAssign, true
		}

		if stmt.Kind != NodeClass || stmt.Class == nil {
			continue
		}

		found, ok := findFirstArrayAssign(stmt.Class.Body)
		if ok {
			return found, true
		}
	}

	return ArrayAssign{}, false
}

// findFirstProperty returns first recursive property assignment.
func findFirstProperty(statements []Statement) (PropertyAssign, bool) {
	for _, stmt := range statements {
		if stmt.Kind == NodeProperty && stmt.Property != nil {
			return *stmt.Property, true
		}

		if stmt.Kind != NodeClass || stmt.Class == nil {
			continue
		}

		found, ok := findFirstProperty(stmt.Class.Body)
		if ok {
			return found, true
		}
	}

	return PropertyAssign{}, false
}
