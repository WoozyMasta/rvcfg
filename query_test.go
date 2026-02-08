package rvcfg

import (
	"reflect"
	"testing"
)

func TestFileFindClass(t *testing.T) {
	t.Parallel()

	file := parseInlineFile(t, `
class CfgVehicles
{
	class Hatchback_02
	{
		displayName = "Ada";
		hiddenSelections[] = {"camo"};
	};
};
`)

	classDecl, ok := file.FindClass("CfgVehicles", "Hatchback_02")
	if !ok {
		t.Fatal("expected nested class CfgVehicles/Hatchback_02")
	}

	property, ok := classDecl.FindProperty("displayName")
	if !ok {
		t.Fatal("expected displayName property")
	}

	if property.Value.Raw != "\"Ada\"" {
		t.Fatalf("unexpected displayName raw value: %q", property.Value.Raw)
	}

	arrayAssign, ok := classDecl.FindArrayAssign("hiddenSelections")
	if !ok {
		t.Fatal("expected hiddenSelections array assignment")
	}

	if arrayAssign.Value.Kind != ValueArray || len(arrayAssign.Value.Elements) != 1 {
		t.Fatalf("unexpected array assignment shape: kind=%s elements=%d", arrayAssign.Value.Kind, len(arrayAssign.Value.Elements))
	}
}

func TestFileWalkClasses(t *testing.T) {
	t.Parallel()

	file := parseInlineFile(t, `
class A
{
	class B {};
	class C
	{
		class D {};
	};
};
`)

	var paths []string
	file.WalkClasses(func(path []string, _ *ClassDecl) bool {
		paths = append(paths, joinPath(path))

		return true
	})

	want := []string{"A", "A/B", "A/C", "A/C/D"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("unexpected class walk order\nwant=%v\ngot=%v", want, paths)
	}
}

func TestFileWalkClassesCanStop(t *testing.T) {
	t.Parallel()

	file := parseInlineFile(t, `
class A
{
	class B {};
	class C {};
};
`)

	var paths []string
	file.WalkClasses(func(path []string, _ *ClassDecl) bool {
		paths = append(paths, joinPath(path))

		return len(paths) < 2
	})

	want := []string{"A", "A/B"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("unexpected class walk stop behavior\nwant=%v\ngot=%v", want, paths)
	}
}

func TestFileWalkStatementsWithPositions(t *testing.T) {
	t.Parallel()

	file := parseInlineFile(t, `
class A
{
	value = 1;
	class B
	{
		items[] = {1, 2};
	};
};
`)

	var refs []StatementRef
	file.WalkStatements(func(ref StatementRef) bool {
		refs = append(refs, ref)

		return true
	})

	if len(refs) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(refs))
	}

	if refs[0].Statement.Kind != NodeClass || refs[0].PathString() != "" {
		t.Fatalf("unexpected first statement ref: kind=%s path=%q", refs[0].Statement.Kind, refs[0].PathString())
	}

	if refs[1].Statement.Kind != NodeProperty || refs[1].PathString() != "A" {
		t.Fatalf("unexpected second statement ref: kind=%s path=%q", refs[1].Statement.Kind, refs[1].PathString())
	}

	if refs[2].Statement.Kind != NodeClass || refs[2].PathString() != "A" {
		t.Fatalf("unexpected third statement ref: kind=%s path=%q", refs[2].Statement.Kind, refs[2].PathString())
	}

	if refs[3].Statement.Kind != NodeArrayAssign || refs[3].PathString() != "A/B" {
		t.Fatalf("unexpected fourth statement ref: kind=%s path=%q", refs[3].Statement.Kind, refs[3].PathString())
	}

	if refs[3].Start.Line <= 0 || refs[3].Start.Column <= 0 {
		t.Fatalf("invalid start position: %+v", refs[3].Start)
	}

	if refs[3].End.Line < refs[3].Start.Line {
		t.Fatalf("invalid end position: start=%+v end=%+v", refs[3].Start, refs[3].End)
	}
}

func TestStatementRefPathString(t *testing.T) {
	t.Parallel()

	ref := StatementRef{
		ClassPath: []string{"CfgVehicles", "Car"},
	}
	if ref.PathString() != "CfgVehicles/Car" {
		t.Fatalf("unexpected path string: %q", ref.PathString())
	}
}

func TestClassFindHelpersNilReceiver(t *testing.T) {
	t.Parallel()

	var classDecl *ClassDecl
	if _, ok := classDecl.FindClass("Any"); ok {
		t.Fatal("expected FindClass on nil receiver to fail")
	}

	if _, ok := classDecl.FindProperty("Any"); ok {
		t.Fatal("expected FindProperty on nil receiver to fail")
	}

	if _, ok := classDecl.FindArrayAssign("Any"); ok {
		t.Fatal("expected FindArrayAssign on nil receiver to fail")
	}
}

func parseInlineFile(t *testing.T, source string) File {
	t.Helper()

	result, err := ParseBytes("query-test.cpp", []byte(source), ParseOptions{
		CaptureScalarRaw: true,
	})
	if err != nil {
		t.Fatalf("ParseBytes query test source error: %v", err)
	}

	return result.File
}

func joinPath(items []string) string {
	if len(items) == 0 {
		return ""
	}

	joined := items[0]
	for idx := 1; idx < len(items); idx++ {
		joined += "/" + items[idx]
	}

	return joined
}
