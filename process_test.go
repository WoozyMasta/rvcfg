package rvcfg

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessAndParseFile_ConfigSample(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "macros", "config.cpp")
	result, err := ProcessAndParseFile(path, PreprocessOptions{}, ParseOptions{
		CaptureScalarRaw:             true,
		AutoFixMissingClassSemicolon: true,
	})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", path, err)
	}

	if len(result.Preprocess.Includes) != 3 {
		t.Fatalf("expected 3 resolved includes, got %d", len(result.Preprocess.Includes))
	}

	if len(result.Parse.File.Statements) == 0 {
		t.Fatal("expected non-empty parsed statements")
	}
}

func TestProcessAndParseFile_MacroDefinedInRoot(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
#define SET_SCOPE(V) scope = V;
class CfgVehicles
{
	class CarScript;
	class TestCar: CarScript
	{
		SET_SCOPE(2)
	};
};
`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", root, err)
	}

	assertNoProcessErrorDiagnostics(t, result.Parse.Diagnostics)

	if strings.Contains(result.Preprocess.Text, "SET_SCOPE(") {
		t.Fatalf("expected macro call expansion in processed text:\n%s", result.Preprocess.Text)
	}

	cfgVehicles, ok := result.Parse.File.FindClass("CfgVehicles")
	if !ok {
		t.Fatal("expected class CfgVehicles")
	}

	testCar, ok := cfgVehicles.FindClass("TestCar")
	if !ok {
		t.Fatal("expected class CfgVehicles/TestCar")
	}

	scope, ok := testCar.FindProperty("scope")
	if !ok {
		t.Fatal("expected scope property produced by root macro")
	}

	if scope.Value.Raw != "2" {
		t.Fatalf("expected scope raw value 2, got %q", scope.Value.Raw)
	}
}

func TestProcessAndParseFile_MacroDefinedInIncludeUsedInRoot(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
#include "defs.hpp"
class CfgVehicles
{
	class CarScript;
	class TestCar: CarScript
	{
		SET_HITPOINTS(250)
	};
};
`,
		"defs.hpp": `#define SET_HITPOINTS(V) hitpoints = V;`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", root, err)
	}

	assertNoProcessErrorDiagnostics(t, result.Parse.Diagnostics)

	if len(result.Preprocess.Includes) != 1 {
		t.Fatalf("expected 1 resolved include, got %d", len(result.Preprocess.Includes))
	}

	if strings.Contains(result.Preprocess.Text, "SET_HITPOINTS(") {
		t.Fatalf("expected include macro call expansion in processed text:\n%s", result.Preprocess.Text)
	}

	cfgVehicles, ok := result.Parse.File.FindClass("CfgVehicles")
	if !ok {
		t.Fatal("expected class CfgVehicles")
	}

	testCar, ok := cfgVehicles.FindClass("TestCar")
	if !ok {
		t.Fatal("expected class CfgVehicles/TestCar")
	}

	hitpoints, ok := testCar.FindProperty("hitpoints")
	if !ok {
		t.Fatal("expected hitpoints property produced by include macro")
	}

	if hitpoints.Value.Raw != "250" {
		t.Fatalf("expected hitpoints raw value 250, got %q", hitpoints.Value.Raw)
	}
}

func TestProcessAndParseFile_IncludePlainConfigChunk(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
class CfgVehicles
{
	class CarScript;
#include "chunk.hpp"
};
`,
		"chunk.hpp": `
class TestChunkCar: CarScript
{
	scope = 2;
};
`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", root, err)
	}

	assertNoProcessErrorDiagnostics(t, result.Parse.Diagnostics)

	if len(result.Preprocess.Includes) != 1 {
		t.Fatalf("expected 1 resolved include, got %d", len(result.Preprocess.Includes))
	}

	cfgVehicles, ok := result.Parse.File.FindClass("CfgVehicles")
	if !ok {
		t.Fatal("expected class CfgVehicles")
	}

	if _, ok := cfgVehicles.FindClass("TestChunkCar"); !ok {
		t.Fatal("expected class from plain include chunk")
	}
}

func TestProcessAndParseFile_IncludeMacroTemplatesAndInvocations(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
class cfgVehicles
{
	class Clothing;
#include "templates.hpp"
#include "items.hpp"
};
`,
		"templates.hpp": `#define DECL_ITEM(NAME, SX, SY) class NAME: Clothing { itemSize[] = {SX, SY}; };`,
		"items.hpp": `
DECL_ITEM(ItemA, 2, 3)
DECL_ITEM(ItemB, 4, 5)
`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", root, err)
	}

	assertNoProcessErrorDiagnostics(t, result.Parse.Diagnostics)

	if len(result.Preprocess.Includes) != 2 {
		t.Fatalf("expected 2 resolved includes, got %d", len(result.Preprocess.Includes))
	}

	if strings.Contains(result.Preprocess.Text, "DECL_ITEM(") {
		t.Fatalf("expected template macro expansion in processed text:\n%s", result.Preprocess.Text)
	}

	cfgVehicles, ok := result.Parse.File.FindClass("cfgVehicles")
	if !ok {
		t.Fatal("expected class cfgVehicles")
	}

	itemA, ok := cfgVehicles.FindClass("ItemA")
	if !ok {
		t.Fatal("expected class cfgVehicles/ItemA")
	}

	itemSize, ok := itemA.FindArrayAssign("itemSize")
	if !ok {
		t.Fatal("expected itemSize array assignment from template macro")
	}

	if itemSize.Value.Kind != ValueArray || len(itemSize.Value.Elements) != 2 {
		t.Fatalf("expected itemSize with 2 array elements, got kind=%s len=%d", itemSize.Value.Kind, len(itemSize.Value.Elements))
	}
}

func TestProcessAndParseFile_RemapParseDiagnosticToIncludeFile(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
class CfgVehicles
{
#include "inc.hpp"
};
`,
		"inc.hpp": `
class Broken
{
	value = 1
};
`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{})
	if err == nil {
		t.Fatal("expected ProcessAndParseFile to fail on missing semicolon in include file")
	}

	found := false
	for _, diagnostic := range result.Parse.Diagnostics {
		if diagnostic.Code != CodeParMissingAssignSemicolon {
			continue
		}

		found = true
		if diagnostic.Start.File != "inc.hpp" {
			t.Fatalf("expected remapped diagnostic file inc.hpp, got %q", diagnostic.Start.File)
		}

		if diagnostic.Start.Line != 5 {
			t.Fatalf("expected remapped diagnostic line 5 in include file, got %d", diagnostic.Start.Line)
		}
	}

	if !found {
		t.Fatalf("expected %s diagnostic in parse diagnostics", FormatCode(CodeParMissingAssignSemicolon))
	}
}

func TestProcessAndParseFile_RemapASTPositionsToIncludeFile(t *testing.T) {
	t.Parallel()

	root := writeProcessFixture(t, map[string]string{
		"root.cpp": `
class CfgVehicles
{
#include "inc.hpp"
};
`,
		"inc.hpp": `
class TestChunkCar
{
	scope = 2;
};
`,
	})

	result, err := ProcessAndParseFile(root, PreprocessOptions{}, ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ProcessAndParseFile(%s) error: %v", root, err)
	}

	cfgVehicles, ok := result.Parse.File.FindClass("CfgVehicles")
	if !ok {
		t.Fatal("expected class CfgVehicles")
	}

	var includeClassStatement *Statement
	for index := range cfgVehicles.Body {
		statement := &cfgVehicles.Body[index]
		if statement.Kind != NodeClass || statement.Class == nil {
			continue
		}

		if statement.Class.Name != "TestChunkCar" {
			continue
		}

		includeClassStatement = statement

		break
	}

	if includeClassStatement == nil {
		t.Fatal("expected class CfgVehicles/TestChunkCar")
	}

	if includeClassStatement.Start.File != "inc.hpp" {
		t.Fatalf("expected include class start file inc.hpp, got %q", includeClassStatement.Start.File)
	}

	if includeClassStatement.Start.Line != 2 {
		t.Fatalf("expected include class start line 2, got %d", includeClassStatement.Start.Line)
	}

	scope, ok := includeClassStatement.Class.FindProperty("scope")
	if !ok {
		t.Fatal("expected scope property in include class")
	}

	if scope.Value.Start.File != "inc.hpp" {
		t.Fatalf("expected include property value file inc.hpp, got %q", scope.Value.Start.File)
	}

	if scope.Value.Start.Line != 4 {
		t.Fatalf("expected include property value line 4, got %d", scope.Value.Start.Line)
	}
}

// writeProcessFixture writes temporary fixture files and returns root.cpp path.
func writeProcessFixture(t *testing.T, files map[string]string) string {
	t.Helper()

	dir := t.TempDir()
	for relPath, content := range files {
		path := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatalf("create fixture dir for %s: %v", relPath, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write fixture file %s: %v", relPath, err)
		}
	}

	return filepath.Join(dir, "root.cpp")
}

// assertNoProcessErrorDiagnostics fails when parse diagnostics contain error severity.
func assertNoProcessErrorDiagnostics(t *testing.T, diagnostics []Diagnostic) {
	t.Helper()

	for _, diagnostic := range diagnostics {
		if diagnostic.Severity != SeverityError {
			continue
		}

		t.Fatalf("unexpected parse diagnostic: %s", diagnostic.Error())
	}
}
