package rvcfg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/woozymasta/lintkit/lint"
)

// testIncludeResolver maps include names to concrete file paths for resolver tests.
type testIncludeResolver struct {
	resolved map[string]string
	calls    int
}

// Resolve returns mapped include path or ErrIncludeNotFound.
func (r *testIncludeResolver) Resolve(currentFile string, includePath string, includeDirs []string) (string, error) {
	r.calls++

	path, ok := r.resolved[includePath]
	if !ok {
		return "", ErrIncludeNotFound
	}

	return path, nil
}

func TestPreprocessFile_ConfigIncludesAndMacros(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "macros", "config.cpp")
	got, err := PreprocessFile(path, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", path, err)
	}

	if strings.Contains(got.Text, "#include") {
		t.Fatalf("expected include directives to be resolved, got text with #include")
	}

	needles := []string{
		"class BirdWalk_LookupTable",
		"class walkErc_Char_LookupTable",
		"class bodyfall_Zmb_LookupTable",
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected preprocessed output to contain %q", needle)
		}
	}

	if len(got.Diagnostics) == 0 {
		return
	}

	for _, d := range got.Diagnostics {
		if d.Severity == lint.SeverityError {
			t.Fatalf("unexpected error diagnostic: %s", d.Error())
		}
	}
}

func TestPreprocessFile_MissingInclude(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")

	source := `#include "not_found.hpp"
class Test {};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrIncludeNotFound) {
		t.Fatalf("expected ErrIncludeNotFound, got %v", err)
	}
}

func TestPreprocessFile_CustomIncludeResolver(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	external := filepath.Join(dir, "external.hpp")

	rootSource := `#include "virtual.hpp"
class Root {};
`

	includeSource := `class External {};
`

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(external, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	resolver := &testIncludeResolver{
		resolved: map[string]string{
			"virtual.hpp": external,
		},
	}

	got, err := PreprocessFile(root, PreprocessOptions{
		IncludeResolver: resolver,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	if resolver.calls == 0 {
		t.Fatal("expected custom include resolver to be called")
	}

	if strings.Contains(got.Text, "#include") {
		t.Fatalf("expected include directive to be resolved, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class External {};") {
		t.Fatalf("expected include content from custom resolver, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_IncludeWithAngleBrackets(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	include := filepath.Join(dir, "inc.hpp")

	rootSource := `#include <inc.hpp>
class Root {};
`

	includeSource := `#define FROM_INCLUDE 11
value = FROM_INCLUDE;
`

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(include, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	got, err := PreprocessFile(root, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	if !strings.Contains(got.Text, "value = 11;") {
		t.Fatalf("expected include with angle brackets to resolve, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_IncludeWithoutSpaceAfterDirective(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	include := filepath.Join(dir, "inc.hpp")

	rootSource := `#include"inc.hpp"
#include<inc.hpp>
class Root {};
`

	includeSource := `class Inc {};
`

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(include, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	got, err := PreprocessFile(root, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	if strings.Count(got.Text, "class Inc {};") != 2 {
		t.Fatalf("expected include without space to resolve for quote and angle forms, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_IncludeDirectiveTailEmittedAsSeparateLine(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	include := filepath.Join(dir, "inc.hpp")

	rootSource := `#include "inc.hpp" extra_tail
class Root {};
`

	includeSource := `class Inc {};
`

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(include, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	got, err := PreprocessFile(root, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	if !strings.Contains(got.Text, "class Inc {};\n extra_tail\nclass Root {};") {
		t.Fatalf("expected include tail as separate source line after include content, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_UnsupportedDirective(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#unknown_directive\nclass Test {};\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnsupportedDirective) {
		t.Fatalf("expected ErrUnsupportedDirective, got %v", err)
	}
}

func TestPreprocessFile_IfDirectiveUnsupportedInStrict(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#if 1\nclass A {};\n#endif\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnsupportedDirective) {
		t.Fatalf("expected ErrUnsupportedDirective for #if in strict mode, got %v", err)
	}

	found := false
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code == CodePPUnsupportedDirective {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected %d diagnostic, got: %+v", CodePPUnsupportedDirective, got.Diagnostics)
	}
}

func TestPreprocessFile_IfDirectiveEnabledInCompatMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#if 1\nclass A {};\n#endif\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeCompat,
	})
	if err != nil {
		t.Fatalf("expected #if to be enabled in compat mode, got %v", err)
	}

	if !strings.Contains(got.Text, "class A {};") {
		t.Fatalf("expected #if branch output in compat mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ElifDirectiveUnsupportedInStrict(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#ifdef A\n#elif 1\nclass B {};\n#endif\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnsupportedDirective) {
		t.Fatalf("expected ErrUnsupportedDirective for #elif in strict mode, got %v", err)
	}

	found := false
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code == CodePPUnsupportedDirective {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected %d diagnostic, got: %+v", CodePPUnsupportedDirective, got.Diagnostics)
	}
}

func TestPreprocessFile_IfDirectiveUnsupportedInsideInactiveIfdef(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#ifdef MISSING\n#if 1\nclass A {};\n#endif\n#endif\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnsupportedDirective) {
		t.Fatalf("expected ErrUnsupportedDirective for nested #if in strict mode, got %v", err)
	}
}

func TestPreprocessFile_MissingDirectiveMacroName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		source string
	}{
		{name: "undef", source: "#undef\nclass A {};\n"},
		{name: "ifdef", source: "#ifdef\nclass A {};\n#endif\n"},
		{name: "ifndef", source: "#ifndef\nclass A {};\n#endif\n"},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			file := filepath.Join(dir, "root.cpp")
			if err := os.WriteFile(file, []byte(tc.source), 0o600); err != nil {
				t.Fatalf("write temp file: %v", err)
			}

			got, err := PreprocessFile(file, PreprocessOptions{})
			if !errors.Is(err, ErrInvalidDirective) {
				t.Fatalf("expected ErrInvalidDirective, got %v", err)
			}

			found := false
			for _, diagnostic := range got.Diagnostics {
				if diagnostic.Code == CodePPMissingMacroName {
					found = true

					break
				}
			}

			if !found {
				t.Fatalf("expected %d diagnostic, got: %+v", CodePPMissingMacroName, got.Diagnostics)
			}
		})
	}
}

func TestPreprocessFile_UnexpectedConditionalDirectives(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		source       string
		expectedCode lint.Code
	}{
		{name: "else", source: "#else\nclass A {};\n", expectedCode: CodePPUnexpectedElse},
		{name: "endif", source: "#endif\nclass A {};\n", expectedCode: CodePPUnexpectedEndif},
		{name: "elif", source: "#elif 1\nclass A {};\n", expectedCode: CodePPUnsupportedDirective},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			file := filepath.Join(dir, "root.cpp")
			if err := os.WriteFile(file, []byte(tc.source), 0o600); err != nil {
				t.Fatalf("write temp file: %v", err)
			}

			got, err := PreprocessFile(file, PreprocessOptions{})
			if err == nil {
				t.Fatalf("expected preprocess error for case %s", tc.name)
			}

			found := false
			for _, diagnostic := range got.Diagnostics {
				if diagnostic.Code == tc.expectedCode {
					found = true

					break
				}
			}

			if !found {
				t.Fatalf("expected %d diagnostic, got: %+v", tc.expectedCode, got.Diagnostics)
			}
		})
	}
}

func TestPreprocessFile_IncludeCycle(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	a := filepath.Join(dir, "a.hpp")
	b := filepath.Join(dir, "b.hpp")

	if err := os.WriteFile(root, []byte(`#include "a.hpp"`+"\n"), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(a, []byte(`#include "b.hpp"`+"\n"), 0o600); err != nil {
		t.Fatalf("write include a: %v", err)
	}

	if err := os.WriteFile(b, []byte(`#include "a.hpp"`+"\n"), 0o600); err != nil {
		t.Fatalf("write include b: %v", err)
	}

	_, err := PreprocessFile(root, PreprocessOptions{})
	if !errors.Is(err, ErrIncludeNotFound) {
		t.Fatalf("expected ErrIncludeNotFound for include cycle, got %v", err)
	}
}

func TestPreprocessFile_MaxIncludeDepthExceeded(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	one := filepath.Join(dir, "one.hpp")
	two := filepath.Join(dir, "two.hpp")

	if err := os.WriteFile(root, []byte(`#include "one.hpp"`+"\n"), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(one, []byte(`#include "two.hpp"`+"\n"), 0o600); err != nil {
		t.Fatalf("write include one: %v", err)
	}

	if err := os.WriteFile(two, []byte("class TooDeep {};\n"), 0o600); err != nil {
		t.Fatalf("write include two: %v", err)
	}

	_, err := PreprocessFile(root, PreprocessOptions{
		MaxIncludeDepth: 1,
	})
	if !errors.Is(err, ErrIncludeNotFound) {
		t.Fatalf("expected ErrIncludeNotFound for include depth overflow, got %v", err)
	}
}

func TestPreprocessFile_DirectiveTailEmission(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#define A 1
#ifdef A A
class C {};
#endif 777
#undef A A
#ifdef MISSING
#undef X 999
#endif
#ifndef A A
class End {};
#else 111
class Else {};
#endif
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "1\nclass C {};") {
		t.Fatalf("expected #ifdef tail to be emitted and macro-expanded, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "\n777\n") {
		t.Fatalf("expected #endif tail to be emitted after popping frame, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "\n A\n") && !strings.HasPrefix(got.Text, " A\n") {
		t.Fatalf("expected #undef tail to be emitted after undef, got:\n%s", got.Text)
	}

	if strings.Contains(got.Text, "999") {
		t.Fatalf("did not expect tail from inactive #ifdef branch, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class End {};") {
		t.Fatalf("expected active #ifndef branch body, got:\n%s", got.Text)
	}

	if strings.Contains(got.Text, "class Else {};") || strings.Contains(got.Text, "111") {
		t.Fatalf("did not expect inactive #else tail/body, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_DirectiveTailElseActive(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#ifdef MISSING
class A {};
#else 111
class B {};
#endif
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "111\nclass B {};") {
		t.Fatalf("expected #else tail in active else branch, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class B {};") {
		t.Fatalf("expected else branch body, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_SourceMapAndIncludeMarkers(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	include := filepath.Join(dir, "inc.hpp")

	rootSource := `class RootStart {};
#include "inc.hpp"
class RootEnd {};
`

	includeSource := `class IncA {};
class IncB {};
`

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(include, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	got, err := PreprocessFile(root, PreprocessOptions{
		EmitIncludeMarkers: true,
		TrackSourceMap:     true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	absInclude, err := filepath.Abs(include)
	if err != nil {
		t.Fatalf("resolve include absolute path: %v", err)
	}

	if !strings.Contains(got.Text, `<include-start "`) || !strings.Contains(got.Text, `<include-end "`) {
		t.Fatalf("expected include markers in output, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class IncA {};") || !strings.Contains(got.Text, "class IncB {};") {
		t.Fatalf("expected include content in output, got:\n%s", got.Text)
	}

	if len(got.SourceMap) == 0 {
		t.Fatal("expected non-empty source map")
	}

	hasIncludeStart := false
	hasIncludeEnd := false
	hasIncludeSource := false

	for _, entry := range got.SourceMap {
		if entry.Kind == "include-start" && entry.IncludeFile == absInclude {
			hasIncludeStart = true
		}

		if entry.Kind == "include-end" && entry.IncludeFile == absInclude {
			hasIncludeEnd = true
		}

		if entry.Kind == "source" && (entry.SourceFile == absInclude || entry.SourceFile == "inc.hpp") {
			hasIncludeSource = true
		}
	}

	if !hasIncludeStart || !hasIncludeEnd || !hasIncludeSource {
		t.Fatalf(
			"expected include-start/include-end/source map entries for include file, got map: %+v",
			got.SourceMap,
		)
	}
}

func TestPreprocessFile_LightingUndergroundMacroExpansion(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "macros", "lighting_underground.txt")
	got, err := PreprocessFile(path, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", path, err)
	}

	if strings.Contains(got.Text, "UNUSED_PARAMS") {
		t.Fatalf("expected object-like macro expansion, got unresolved token in output:\n%s", got.Text)
	}

	needles := []string{
		"luminanceRectCoef = 0.0;",
		"swBrightness = 0;",
		"class Underground",
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected preprocessed output to contain %q", needle)
		}
	}
}

func TestPreprocessFile_ConfigDeterministicGoldenHash(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "macros", "config.cpp")
	got, err := PreprocessFile(path, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", path, err)
	}

	normalized := strings.ReplaceAll(got.Text, "\r\n", "\n")
	sum := sha256.Sum256([]byte(normalized))
	hash := hex.EncodeToString(sum[:])

	const (
		expectedHash  = "93baf439e23fef12e4ec00ffa2eb4a6f52b885df39c2159e221fc902afceb98b"
		expectedLines = 164
	)

	if hash != expectedHash {
		t.Fatalf("unexpected preprocess hash: got=%s want=%s", hash, expectedHash)
	}

	lines := strings.Count(normalized, "\n") + 1
	if lines != expectedLines {
		t.Fatalf("unexpected preprocess line count: got=%d want=%d", lines, expectedLines)
	}
}

func TestPreprocessFile_FunctionLikeMacroArgumentEdgeCases(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define MAKE_ARRAY(NAME, A, B) NAME[] = {A, B};
class CfgTest
{
	MAKE_ARRAY(nested, func(1, (2 + 3)), finalValue)
	MAKE_ARRAY(strPairs, "left,right", "x,y")
	MAKE_ARRAY(emptyLead, , 2)
	MAKE_ARRAY(emptyTail, 1, )
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		"nested[] = {func(1, (2 + 3)), finalValue};",
		`strPairs[] = {"leftright", "xy"};`,
		"emptyLead[] = {, 2};",
		"emptyTail[] = {1, };",
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)

	for _, needle := range needles {
		if !strings.Contains(normalized, stripHorizontalWhitespaceTest(needle)) {
			t.Fatalf("expected expanded output to contain %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_FunctionLikeMacroArgCountMismatch_DropsInvocation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define MAKE(A, B) A + B
class CfgTest
{
	value1 = MAKE(1);
	value2 = MAKE(1, 2, 3);
	value3 = MAKE(1,2);
	value4 = MAKE(1/*x,y*/,2);
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)
	needles := []string{
		"value1=;",
		"value2=;",
		"value3=1+2;",
		"value4=;",
	}

	for _, needle := range needles {
		if !strings.Contains(normalized, needle) {
			t.Fatalf("expected output to contain %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_FunctionLikeMacroMalformedCall_DropsInvocation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define M(A,B) A + B
class CfgTest
{
	value1 = M(1;
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("expected malformed macro calls to be tolerated, got %v", err)
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)
	if !strings.Contains(normalized, "value1=;") {
		t.Fatalf("expected malformed invocation to be dropped, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_FunctionLikeMacroMalformedCall_QuotedArg_DropsSemicolon(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := "#define M(ARG) ARG\nvalue1 = M(\"A\\\"B\");"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("expected malformed macro calls to be tolerated, got %v", err)
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)
	if !strings.Contains(normalized, "value1=") || strings.Contains(normalized, "value1=;") {
		t.Fatalf("expected malformed quoted invocation to consume ';', got:\n%s", got.Text)
	}
}

func TestPreprocessFile_FunctionLikeMacroMalformedTwoArg_DuplicatesFirstArg(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define M(A,B) A + B
class CfgTest
{
	value1 = M(1,2;
	value2 = M(3,4
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("expected malformed two-arg call to be tolerated, got %v", err)
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)
	needles := []string{
		"value1=1+1",
		"value2=3+3",
	}

	for _, needle := range needles {
		if !strings.Contains(normalized, needle) {
			t.Fatalf("expected output to contain %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_FunctionLikeMacroStringArg_CommaRemoved(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define M(ARG) ARG
value1 = M("Some, content");
value2 = M("A,B,C");
value3 = "keep,comma";
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	normalized := stripHorizontalWhitespaceTest(got.Text)
	needles := []string{
		`value1="Somecontent";`,
		`value2="ABC";`,
		`value3="keep,comma";`,
	}

	for _, needle := range needles {
		if !strings.Contains(normalized, stripHorizontalWhitespaceTest(needle)) {
			t.Fatalf("expected output to contain %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_TokenPasteWhitespaceVariants(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define JOIN_TIGHT(A, B) A##B
#define JOIN_LEFT(A, B) A ##B
#define JOIN_RIGHT(A, B) A## B
#define JOIN_WIDE(A, B) A ## B
class CfgTest
{
	value0 = JOIN_TIGHT(foo, bar);
	value1 = JOIN_LEFT(alpha, beta);
	value2 = JOIN_RIGHT(left, right);
	value3 = JOIN_WIDE(one, two);
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		"value0 = foobar;",
		"value1 = alphabeta;",
		"value2 = leftright;",
		"value3 = onetwo;",
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected token-paste expansion %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestCollapseTokenPaste_OneSidedKeepsIndentation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "leading_token_paste_keeps_prefix_tabs",
			in:   "{\t\t##ValueToken",
			want: "{\t\tValueToken",
		},
		{
			name: "trailing_token_paste_keeps_suffix_tabs",
			in:   "class Name##\t\t\t{",
			want: "class Name\t\t\t{",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := collapseTokenPaste(tc.in)
			if got != tc.want {
				t.Fatalf("collapseTokenPaste mismatch\ngot:  %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func TestTrimSingleMacroBodyDelimiter(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "single space", in: " A##B", want: "A##B"},
		{name: "single tab", in: "\tA##B", want: "A##B"},
		{name: "multiple tabs", in: "\t\t\tclass A{}", want: "\t\t\tclass A{}"},
		{name: "multiple spaces", in: "   class A{}", want: "   class A{}"},
		{name: "no leading whitespace", in: "A##B", want: "A##B"},
		{name: "empty", in: "", want: ""},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := trimSingleMacroBodyDelimiter(tc.in)
			if got != tc.want {
				t.Fatalf("trimSingleMacroBodyDelimiter(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestShouldConsumeMalformedCallSemicolon(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		in    string
		want  bool
		start int
		end   int
	}{
		{
			name:  "quoted_body_consumes",
			in:    `x = M("A\"B");`,
			want:  true,
			start: strings.Index(`x = M("A\"B");`, "M"),
			end:   strings.Index(`x = M("A\"B");`, ";"),
		},
		{
			name:  "newline_after_semicolon_consumes",
			in:    "x = M(1;\n",
			want:  true,
			start: strings.Index("x = M(1;\n", "M"),
			end:   strings.Index("x = M(1;\n", ";"),
		},
		{
			name:  "plain_malformed_keeps_semicolon",
			in:    "x = M(1;",
			want:  false,
			start: strings.Index("x = M(1;", "M"),
			end:   strings.Index("x = M(1;", ";"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := shouldConsumeMalformedCallSemicolon(tc.in, tc.start, "M", tc.end)
			if got != tc.want {
				t.Fatalf("want consume=%v, got %v for %q", tc.want, got, tc.in)
			}
		})
	}
}

func TestPreprocessFile_StringifyOperator(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define STR(X) #X
class CfgTest
{
	text0 = STR(hello_world);
	text1 = STR(1 + 2);
	text2 = STR(path\to\file);
	text3 = STR("quoted");
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		`text0 = "hello_world";`,
		`text1 = "1 + 2";`,
		`text2 = "path\\to\\file";`,
		`text3 = "\"quoted\"";`,
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected stringify expansion %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_StringifyWithTokenPaste(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define DECL(NAME) class NAME##_Node { value = #NAME; };
DECL(Foo)
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		`class Foo_Node { value = "Foo"; };`,
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected stringify+paste expansion %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_StringifyTokenPasteCfgConvertStyle(t *testing.T) {
	t.Parallel()

	path := testDataPath("preproc", "stringify_tokenpaste", "input.cpp")
	got, err := PreprocessFile(path, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", path, err)
	}

	if !strings.Contains(got.Text, `displayName = "$##Warpbox##_##  CardboardBox##_name";`) {
		t.Fatalf("expected CfgConvert-style stringify+paste result for displayName, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `damageMat = "Warpbox/assets/data/testbox##_damage.rvmat";`) {
		t.Fatalf("expected CfgConvert-style BASE##_damage preservation in string, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_TokenPasteOutputStaysParseable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#define MOD_PREFIX Warpbox
#define CFG_CAT_INNER(A, B) A##B
#define CFG_CAT(A, B) CFG_CAT_INNER(A, B)
class CFG_CAT(MOD_PREFIX, _Data) {};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile default mode error: %v", err)
	}

	if !strings.Contains(got.Text, `class Warpbox_Data {};`) {
		t.Fatalf("expected parseable default token-paste output, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_UnresolvedMacroLikeInvocation_PreservedInStrictMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgTest
{
	healthLevels[] = HL_FROM_BASE("bagpack");
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, `HL_FROM_BASE("bagpack")`) {
		t.Fatalf("expected unresolved macro call to remain untouched in strict mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_UnresolvedMacroLikeInvocation_IgnoresStringsAndComments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgTest
{
	value = func(1, 2);
	text = "HL_FROM_BASE(bagpack)";
	// HL_FROM_BASE("bagpack")
	/* HL_FROM_BASE("bagpack")
	   KEEP_SCOPE(2) */
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("expected no unresolved macro error for strings/comments/lowercase calls, got %v", err)
	}
}

func TestPreprocessFile_FunctionLikeMacro_DoesNotExpandInStringsOrComments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define F(X) X
class CfgTest
{
	value = F(7);
	text = "F(123)";
	// F(456)
	/* F(789) */
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "value = 7;") {
		t.Fatalf("expected function-like macro expansion in code, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `text = "F(123)";`) {
		t.Fatalf("expected string literal macro-like text untouched, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_MacroRedefinedWarningDedup(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define DUP 1
#define DUP 2
#define DUP 3
class CfgTest {};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	count := 0
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code != CodePPMacroRedefined {
			continue
		}

		count++
	}

	if count != 1 {
		t.Fatalf("expected one %d warning after dedup, got=%d diagnostics=%+v", CodePPMacroRedefined, count, got.Diagnostics)
	}
}

func TestPreprocessFile_MacroRedefinedWarningLineWithContinuation(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#define DUP(A, B) A \
/* keep line */ \
B
#define DUP(A, B) A
class CfgTest {};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableMacroRedefinitionWarnings: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	diagLine := 0
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code != CodePPMacroRedefined {
			continue
		}

		diagLine = diagnostic.Start.Line

		break
	}

	if diagLine == 0 {
		t.Fatalf("expected %d warning, diagnostics=%+v", CodePPMacroRedefined, got.Diagnostics)
	}

	if diagLine != 4 {
		t.Fatalf("expected %d warning on source line 4, got=%d diagnostics=%+v", CodePPMacroRedefined, diagLine, got.Diagnostics)
	}
}

func TestPreprocessFile_DocsDoubleNestedHealthLevelsMacro(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
#define CAMPING_DATA_PATH "DZ\gear\camping\data\"
#define MAT_OK(BASE) CAMPING_DATA_PATH BASE ".rvmat"
#define MAT_DMG(BASE) CAMPING_DATA_PATH BASE "_damage.rvmat"
#define MAT_DEST(BASE) CAMPING_DATA_PATH BASE "_destruct.rvmat"
#define HL_ROW(PCT, MAT) {PCT, {MAT}}
#define HL_STD(OK, DMG, DEST) \
{ \
	HL_ROW(1.0, OK), \
	HL_ROW(0.5, DMG), \
	HL_ROW(0.0, DEST) \
}
#define HL_FROM_BASE(BASE) HL_STD(MAT_OK(BASE), MAT_DMG(BASE), MAT_DEST(BASE))

healthLevels[] = HL_FROM_BASE("bagpack");
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		"healthLevels[] = {",
		`"bagpack" ".rvmat"`,
		`"bagpack" "_damage.rvmat"`,
		`"bagpack" "_destruct.rvmat"`,
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected expanded output to contain %q, got:\n%s", needle, got.Text)
		}
	}

	blockedTokens := []string{
		"HL_FROM_BASE(",
		"HL_STD(",
		"HL_ROW(",
		"MAT_OK(",
		"MAT_DMG(",
		"MAT_DEST(",
	}

	for _, token := range blockedTokens {
		if strings.Contains(got.Text, token) {
			t.Fatalf("expected docs macro token %q to be fully expanded, got:\n%s", token, got.Text)
		}
	}
}

func TestPreprocessFile_ExecEval_DefaultPreserved(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgTest
{
	class Root
	{
		value = __EVAL(1+2);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "value = __EVAL(1+2);") {
		t.Fatalf("expected __EVAL to stay untouched in strict mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ExecEval_EnabledEvalOnly(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgExecEval
{
	class Root
	{
		eval_value = __EVAL(1+2);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableExecEvalIntrinsics: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "eval_value = 3;") {
		t.Fatalf("expected __EVAL result to be numeric 3, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ExecEval_EnabledExecEval(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgExecEval
{
	class Root
	{
		__EXEC(testVar = 7)
		value = __EVAL(testVar + 5);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableExecEvalIntrinsics: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if strings.Contains(got.Text, "__EXEC(") {
		t.Fatalf("expected __EXEC to be consumed in compatibility mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "value = 12;") {
		t.Fatalf("expected __EVAL(testVar + 5) to be 12, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ExecEval_EnabledByCompatMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgExecEval
{
	class Root
	{
		__EXEC(testVar = 7)
		value = __EVAL(testVar + 5);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeCompat,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "value = 12;") {
		t.Fatalf("expected __EVAL(testVar + 5) to be 12 in compat mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_PathNorm_DefaultPreserved(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgPathNorm
{
	class Root
	{
		path = __PATH_NORM("a/b\\c//d");
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, `path = __PATH_NORM("a/b\\c//d");`) {
		t.Fatalf("expected __PATH_NORM to stay untouched in strict mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_PathNorm_EnabledByExtendedMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgPathNorm
{
	class Root
	{
		pathA = __PATH_NORM("a/b\\c//d");
		__EXEC(base = "mods/demo")
		pathB = __PATH_NORM(base + "/assets\\data//x.paa");
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if strings.Contains(got.Text, "__PATH_NORM(") {
		t.Fatalf("expected __PATH_NORM to be expanded in extended mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `pathA = "a\b\c\d";`) {
		t.Fatalf("expected pathA normalization, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `pathB = "mods\demo\assets\data\x.paa";`) {
		t.Fatalf("expected pathB normalization with __EXEC variable, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_StringIntrinsics_DefaultPreserved(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgStr
{
	class Root
	{
		v1 = __STR_TRIM("  x  ");
		v2 = __STR_LOWER("AB");
		v3 = __STR_UPPER("ab");
		v4 = __STR_REPLACE("a/b", "/", "\\");
		v5 = __STR_JOIN("_", "a", "b");
		v6 = __STR_QUOTE("demo");
		v7 = __STR_SPLIT("a_b_c", "_", 1);
		v8 = __STR_PASCAL("hello-world_demo x");
		v9 = __STR_CAMEL("hello-world_demo x");
		v10 = __STR_SNAKE("HelloWorld Demo");
		v11 = __STR_CONST("HelloWorld Demo");
		v12 = __FILES_JOIN("sounds/*.ogg", "|");
		v13 = __FILES_COUNT("sounds/*.ogg");
		v14 = __FILES_GET("sounds/*.ogg", 0);
		v15 = __FILES_RENDER("sounds/*.ogg", "X:{stem}");
		v16 = __FOR_RANGE_RENDER(1, 3, "{index}:{value}", "|");
		v17 = __FOR_EACH_RENDER("{index}:{value}", "|", "a", "b");
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		`v1 = __STR_TRIM("  x  ");`,
		`v2 = __STR_LOWER("AB");`,
		`v3 = __STR_UPPER("ab");`,
		`v4 = __STR_REPLACE("a/b", "/", "\\");`,
		`v5 = __STR_JOIN("_", "a", "b");`,
		`v6 = __STR_QUOTE("demo");`,
		`v7 = __STR_SPLIT("a_b_c", "_", 1);`,
		`v8 = __STR_PASCAL("hello-world_demo x");`,
		`v9 = __STR_CAMEL("hello-world_demo x");`,
		`v10 = __STR_SNAKE("HelloWorld Demo");`,
		`v11 = __STR_CONST("HelloWorld Demo");`,
		`v12 = __FILES_JOIN("sounds/*.ogg", "|");`,
		`v13 = __FILES_COUNT("sounds/*.ogg");`,
		`v14 = __FILES_GET("sounds/*.ogg", 0);`,
		`v15 = __FILES_RENDER("sounds/*.ogg", "X:{stem}");`,
		`v16 = __FOR_RANGE_RENDER(1, 3, "{index}:{value}", "|");`,
		`v17 = __FOR_EACH_RENDER("{index}:{value}", "|", "a", "b");`,
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected strict mode to keep %q unchanged, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_StringIntrinsics_EnabledByExtendedMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgStr
{
	class Root
	{
		__EXEC(base = "Mods/DEMO"; suffix = "x"; item1 = "mods"; item2 = "demo")
		trimmed = __STR_TRIM("  Mods/DEMO  ");
		lower = __STR_LOWER(base);
		upper = __STR_UPPER("abC");
		replaced = __STR_REPLACE("a/b/c", "/", "\\");
		joined = __STR_JOIN("\\", item1, item2, "assets", suffix);
		quoted = __STR_QUOTE(base + "/" + suffix);
		split1 = __STR_SPLIT("a_b_c", "_", 1);
		split9 = __STR_SPLIT("a_b_c", "_", 9);
		pascal = __STR_PASCAL("hello-world_demo x");
		camel = __STR_CAMEL("hello-world_demo x");
		snake = __STR_SNAKE("HelloWorld Demo");
		constv = __STR_CONST("HelloWorld Demo");
		count = __FILES_COUNT("sounds/*.ogg");
		first = __FILES_GET("sounds/*.ogg", 0);
		missing = __FILES_GET("sounds/*.ogg", 9);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		`trimmed = "Mods/DEMO";`,
		`lower = "mods/demo";`,
		`upper = "ABC";`,
		`replaced = "a\\b\\c";`,
		`joined = "mods\\demo\\assets\\x";`,
		`quoted = "Mods/DEMO/x";`,
		`split1 = "b";`,
		`split9 = "";`,
		`pascal = "HelloWorldDemoX";`,
		`camel = "helloWorldDemoX";`,
		`snake = "hello_world_demo";`,
		`constv = "HELLO_WORLD_DEMO";`,
		`count = 0;`,
		`first = "";`,
		`missing = "";`,
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected extended mode output to contain %q, got:\n%s", needle, got.Text)
		}
	}
}

func TestPreprocessFile_StringIntrinsics_InvalidArguments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgStr
{
	class Root
	{
		v = __STR_REPLACE("a", "b");
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __STR_REPLACE call, got %v", err)
	}
}

func TestPreprocessFile_StringIntrinsics_SplitInvalidIndex(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgStr
{
	class Root
	{
		v = __STR_SPLIT("a_b_c", "_", -1);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __STR_SPLIT index, got %v", err)
	}
}

func TestPreprocessFile_StringIntrinsics_InvalidTemplateFilter(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sounds"), 0o700); err != nil {
		t.Fatalf("mkdir sounds: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "sounds", "a.ogg"), []byte("x"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgStr
{
	class Root
	{
		v = __FILES_RENDER("sounds/*.ogg", "{stem|nope}", "|");
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := PreprocessFile(file, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid template filter, got %v", err)
	}
}

func TestPreprocessFile_FilesJoin_EnabledByExtendedMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sounds"), 0o700); err != nil {
		t.Fatalf("mkdir sounds: %v", err)
	}

	files := []string{
		filepath.Join(dir, "sounds", "b.ogg"),
		filepath.Join(dir, "sounds", "a.ogg"),
		filepath.Join(dir, "sounds", "ignore.txt"),
	}

	for _, filePath := range files {
		if err := os.WriteFile(filePath, []byte("x"), 0o600); err != nil {
			t.Fatalf("write fixture file %s: %v", filePath, err)
		}
	}

	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgFiles
{
	class Root
	{
		filesPipe = __FILES_JOIN("sounds/*.ogg", "|");
		filesDefault = __FILES_JOIN("sounds/*.ogg");
		filesCount = __FILES_COUNT("sounds/*.ogg");
		firstFile = __FILES_GET("sounds/*.ogg", 0);
		missingFile = __FILES_GET("sounds/*.ogg", 99);
		filesRender = __FILES_RENDER("sounds/*.ogg", "{index}:{stem}={path}", "|");
		filesRenderFiltered = __FILES_RENDER("sounds/*.ogg", "{index}:{stem|pascal}:{stem|snake}:{stem|const}={path|lower|replace(sounds, sfx)}", "|");
		filesRenderClass = __FILES_RENDER("sounds/*.ogg", "class Snd_{index} { path = {path|lower}; };", "\n");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	got, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", sourcePath, err)
	}

	if !strings.Contains(got.Text, `filesPipe = "sounds\a.ogg|sounds\b.ogg";`) {
		t.Fatalf("expected sorted and joined file list for filesPipe, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `filesDefault = "sounds\a.ogg;sounds\b.ogg";`) {
		t.Fatalf("expected default ';' delimiter for filesDefault, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `filesCount = 2;`) {
		t.Fatalf("expected filesCount=2, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `firstFile = "sounds\a.ogg";`) {
		t.Fatalf("expected firstFile with sorted index 0, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `missingFile = "";`) {
		t.Fatalf("expected missingFile to be empty string, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `filesRender = 0:a=sounds\a.ogg|1:b=sounds\b.ogg;`) {
		t.Fatalf("expected filesRender template output, got:\n%s", got.Text)
	}

	if !strings.Contains(
		got.Text,
		`filesRenderFiltered = 0:A:a:A=sfx\a.ogg|1:B:b:B=sfx\b.ogg;`,
	) {
		t.Fatalf("expected filesRenderFiltered template output, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `class Snd_0 { path = sounds\a.ogg; };`) ||
		!strings.Contains(got.Text, `class Snd_1 { path = sounds\b.ogg; };`) {
		t.Fatalf("expected class-style template rendering with literal braces, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_FilesJoin_MaxItemsExceeded(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sounds"), 0o700); err != nil {
		t.Fatalf("mkdir sounds: %v", err)
	}

	for _, name := range []string{"a.ogg", "b.ogg"} {
		if err := os.WriteFile(filepath.Join(dir, "sounds", name), []byte("x"), 0o600); err != nil {
			t.Fatalf("write fixture file %s: %v", name, err)
		}
	}

	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgFiles
{
	class Root
	{
		filesPipe = __FILES_JOIN("sounds/*.ogg", "|");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode:               PreprocessModeExtended,
		ExtendedFSMaxItems: 1,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic when files exceed limit, got %v", err)
	}
}

func TestPreprocessFile_FilesJoin_RootRestriction(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	allowedRoot := filepath.Join(dir, "allowed")
	outsideDir := filepath.Join(dir, "outside")
	if err := os.MkdirAll(filepath.Join(allowedRoot, "sounds"), 0o700); err != nil {
		t.Fatalf("mkdir allowed sounds: %v", err)
	}

	if err := os.MkdirAll(outsideDir, 0o700); err != nil {
		t.Fatalf("mkdir outside: %v", err)
	}

	outsideFile := filepath.Join(outsideDir, "x.ogg")
	if err := os.WriteFile(outsideFile, []byte("x"), 0o600); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}

	sourcePath := filepath.Join(allowedRoot, "root.cpp")
	source := `
class CfgFiles
{
	class Root
	{
		files = __FILES_JOIN("` + normalizePathSlashes(outsideFile) + `", "|");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode:           PreprocessModeExtended,
		ExtendedFSRoot: allowedRoot,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for outside root file glob, got %v", err)
	}
}

func TestPreprocessFile_FilesGet_InvalidIndex(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "sounds"), 0o700); err != nil {
		t.Fatalf("mkdir sounds: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "sounds", "a.ogg"), []byte("x"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgFiles
{
	class Root
	{
		file = __FILES_GET("sounds/*.ogg", -1);
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __FILES_GET index, got %v", err)
	}
}

func TestPreprocessFile_FilesRender_InvalidArguments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgFiles
{
	class Root
	{
		file = __FILES_RENDER("sounds/*.ogg");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __FILES_RENDER args, got %v", err)
	}
}

func TestPreprocessFile_ForRangeRender_EnabledByExtendedMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForRange
{
	class Root
	{
		asc = __FOR_RANGE_RENDER(1, 3, "A{index}:{value}", "|");
		desc = __FOR_RANGE_RENDER(3, 1, "{value}", ",");
		pascal = __FOR_RANGE_RENDER(1, 3, "{value|pascal}", "|");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	got, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", sourcePath, err)
	}

	if !strings.Contains(got.Text, `asc = A0:1|A1:2|A2:3;`) {
		t.Fatalf("expected ascending __FOR_RANGE_RENDER output, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `desc = 3,2,1;`) {
		t.Fatalf("expected descending __FOR_RANGE_RENDER output, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `pascal = 1|2|3;`) {
		t.Fatalf("expected filtered __FOR_RANGE_RENDER output, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ForRangeRender_InvalidArguments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForRange
{
	class Root
	{
		asc = __FOR_RANGE_RENDER(1, 3);
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __FOR_RANGE_RENDER args, got %v", err)
	}
}

func TestPreprocessFile_ForRangeRender_MaxItemsExceeded(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForRange
{
	class Root
	{
		asc = __FOR_RANGE_RENDER(1, 5, "{value}", "|");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode:                 PreprocessModeExtended,
		ExtendedLoopMaxItems: 3,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for __FOR_RANGE_RENDER loop limit, got %v", err)
	}
}

func TestPreprocessFile_ForEachRender_EnabledByExtendedMode(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForEach
{
	class Root
	{
		__EXEC(a = "alpha"; b = "beta")
		items = __FOR_EACH_RENDER("{index}:{value}", "|", a, b, "gamma");
		itemsFiltered = __FOR_EACH_RENDER("{index}:{value|upper|replace(A, X)}", "|", a, b);
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	got, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", sourcePath, err)
	}

	if !strings.Contains(got.Text, `items = 0:alpha|1:beta|2:gamma;`) {
		t.Fatalf("expected __FOR_EACH_RENDER output, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `itemsFiltered = 0:XLPHX|1:BETX;`) {
		t.Fatalf("expected filtered __FOR_EACH_RENDER output, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_ForEachRender_InvalidArguments(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForEach
{
	class Root
	{
		items = __FOR_EACH_RENDER("{value}", "|");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode: PreprocessModeExtended,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for invalid __FOR_EACH_RENDER args, got %v", err)
	}
}

func TestPreprocessFile_ForEachRender_MaxItemsExceeded(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "root.cpp")
	source := `
class CfgForEach
{
	class Root
	{
		items = __FOR_EACH_RENDER("{value}", "|", "a", "b", "c");
	};
};
`

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("write source: %v", err)
	}

	_, err := PreprocessFile(sourcePath, PreprocessOptions{
		Mode:                 PreprocessModeExtended,
		ExtendedLoopMaxItems: 2,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for __FOR_EACH_RENDER loop limit, got %v", err)
	}
}

func TestParseIntrinsicArgs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "simple_three",
			input: `"a", "b", "c"`,
			want:  []string{`"a"`, `"b"`, `"c"`},
		},
		{
			name:  "nested_parentheses",
			input: `base + (1 + 2), "b"`,
			want:  []string{`base + (1 + 2)`, `"b"`},
		},
		{
			name:  "comma_inside_string",
			input: `"a,b", "c"`,
			want:  []string{`"a,b"`, `"c"`},
		},
		{
			name:    "unterminated_string",
			input:   `"a, "b"`,
			wantErr: true,
		},
		{
			name:    "unexpected_close_paren",
			input:   `a), b`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseIntrinsicArgs(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected parse error for input %q", tc.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("parseIntrinsicArgs(%q): %v", tc.input, err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("parseIntrinsicArgs(%q) len mismatch: got=%d want=%d", tc.input, len(got), len(tc.want))
			}

			for idx := range tc.want {
				if got[idx] != tc.want[idx] {
					t.Fatalf("parseIntrinsicArgs(%q)[%d]=%q want %q", tc.input, idx, got[idx], tc.want[idx])
				}
			}
		})
	}
}

func TestPreprocessFile_ExecEval_UnknownIdentifierFallback(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := `
class CfgExecEval
{
	class Root
	{
		value = __EVAL(unknown + 1);
	};
};
`

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableExecEvalIntrinsics: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, `value = "scalar";`) {
		t.Fatalf("expected BI-like fallback \"scalar\" for unknown eval identifier, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_BuiltInIntrinsics_Basic(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
class CfgBuiltins
{
	lineValue = __LINE__;
	fileValue = __FILE__;
	fileNameValue = __FILE_NAME__;
	fileShortValue = __FILE_SHORT__;
	stringLiteral = "__LINE__";
	// __FILE__
};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "lineValue = 2;") {
		t.Fatalf("expected __LINE__ to resolve with DayZ-style 0-based index, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "fileValue = "+quoteIntrinsicString(file)+";") {
		t.Fatalf("expected __FILE__ to resolve to current file path, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileNameValue = __FILE_NAME__;`) {
		t.Fatalf("expected __FILE_NAME__ to stay untouched in strict mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileShortValue = __FILE_SHORT__;`) {
		t.Fatalf("expected __FILE_SHORT__ to stay untouched in strict mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `stringLiteral = "__LINE__";`) {
		t.Fatalf("expected builtin in string literal to stay untouched, got:\n%s", got.Text)
	}

	if strings.Contains(got.Text, "// __FILE__") {
		t.Fatalf("expected comments to be removed by preprocessor, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_BuiltInIntrinsics_IncludeContext(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "root.cpp")
	include := filepath.Join(dir, "inc.hpp")

	rootSource := strings.TrimSpace(`
#include "inc.hpp"
class Root {};
`) + "\n"

	includeSource := strings.TrimSpace(`
class Included
{
	fileValue = __FILE__;
	lineValue = __LINE__;
	fileNameValue = __FILE_NAME__;
	fileShortValue = __FILE_SHORT__;
};
`) + "\n"

	if err := os.WriteFile(root, []byte(rootSource), 0o600); err != nil {
		t.Fatalf("write root source: %v", err)
	}

	if err := os.WriteFile(include, []byte(includeSource), 0o600); err != nil {
		t.Fatalf("write include source: %v", err)
	}

	got, err := PreprocessFile(root, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", root, err)
	}

	if !strings.Contains(got.Text, "lineValue = 3;") {
		t.Fatalf("expected __LINE__ in include to resolve with DayZ-style 0-based index, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileValue = "inc.hpp";`) {
		t.Fatalf("expected __FILE__ in include to use include literal path, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileNameValue = __FILE_NAME__;`) {
		t.Fatalf("expected __FILE_NAME__ in include to stay untouched in strict mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileShortValue = __FILE_SHORT__;`) {
		t.Fatalf("expected __FILE_SHORT__ in include to stay untouched in strict mode, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_BuiltInIntrinsics_FileNameEnabled(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
class CfgBuiltins
{
	fileNameValue = __FILE_NAME__;
	fileShortValue = __FILE_SHORT__;
};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableFileNameIntrinsics: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, `fileNameValue = "root.cpp";`) {
		t.Fatalf("expected __FILE_NAME__ to expand when enabled, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileShortValue = "root";`) {
		t.Fatalf("expected __FILE_SHORT__ to expand when enabled, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_DynamicIntrinsics_DefaultDisabled(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
class CfgDynamic
{
	counterValue = __COUNTER__;
	dateValue = __DATE_STR__;
	randValue = __RAND_INT8__;
};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "__COUNTER__") ||
		!strings.Contains(got.Text, "__DATE_STR__") ||
		!strings.Contains(got.Text, "__RAND_INT8__") {
		t.Fatalf("expected dynamic intrinsics to stay untouched by default, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_DynamicIntrinsics_Enabled(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
class CfgDynamic
{
	dates[] = {__DATE_ARR__, __DATE_STR__, __DATE_STR_ISO8601__};
	times[] = {__TIME__, __TIME_UTC__};
	parts[] = {__DAY__, __MONTH__, __YEAR__, __TIMESTAMP_UTC__};

	counterA = __COUNTER__;
	counterB = __COUNTER__;
	__COUNTER_RESET__
	counterC = __COUNTER__;

	randI8 = __RAND_INT8__;
	randI16 = __RAND_INT16__;
	randI32 = __RAND_INT32__;
	randI64 = __RAND_INT64__;
	randU8 = __RAND_UINT8__;
	randU16 = __RAND_UINT16__;
	randU32 = __RAND_UINT32__;
	randU64 = __RAND_UINT64__;
};
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableDynamicIntrinsics: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if strings.Contains(got.Text, "__COUNTER__") ||
		strings.Contains(got.Text, "__DATE_STR__") ||
		strings.Contains(got.Text, "__RAND_INT8__") {
		t.Fatalf("expected dynamic intrinsics to be expanded in enabled mode, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "counterA = 0;") ||
		!strings.Contains(got.Text, "counterB = 1;") ||
		!strings.Contains(got.Text, "counterC = 0;") {
		t.Fatalf("expected counter/reset sequence 0,1,0, got:\n%s", got.Text)
	}

	if !regexp.MustCompile(`dates\[] = \{\d{4},\d{1,2},\d{1,2},\d{1,2},\d{1,2},\d{1,2}, "\d{4}/\d{2}/\d{2}, \d{2}:\d{2}:\d{2}", "\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z"\};`).
		MatchString(got.Text) {
		t.Fatalf("expected expanded date intrinsics, got:\n%s", got.Text)
	}

	if !regexp.MustCompile(`times\[] = \{\d{2}:\d{2}:\d{2}, \d{2}:\d{2}:\d{2}\};`).MatchString(got.Text) {
		t.Fatalf("expected expanded time intrinsics, got:\n%s", got.Text)
	}

	if !regexp.MustCompile(`parts\[] = \{\d{1,2}, \d{1,2}, \d{4}, \d+\};`).MatchString(got.Text) {
		t.Fatalf("expected expanded day/month/year/timestamp intrinsics, got:\n%s", got.Text)
	}

	assertDynamicIntRange(t, got.Text, "randI8", -128, 127)
	assertDynamicIntRange(t, got.Text, "randI16", -32768, 32767)
	assertDynamicIntRange(t, got.Text, "randI32", -2147483648, 2147483647)
	assertDynamicIntPresent(t, got.Text, "randI64")

	assertDynamicUintRange(t, got.Text, "randU8", 255)
	assertDynamicUintRange(t, got.Text, "randU16", 65535)
	assertDynamicUintRange(t, got.Text, "randU32", 4294967295)
	assertDynamicUintPresent(t, got.Text, "randU64")
}

func TestPreprocessFile_IfNumericComparisons(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#define __GAME_VER_MIN__ 16
#if __GAME_VER_MIN__ >= 16
class VersionOk {};
#endif
#if __GAME_VER_MIN__ < 16
class VersionFail {};
#endif
#if (__GAME_VER_MIN__ == 16) && defined(__GAME_VER_MIN__)
class VersionAndDefined {};
#endif
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableIfExpressions: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if !strings.Contains(got.Text, "class VersionOk {};") {
		t.Fatalf("expected >= comparison branch enabled, got:\n%s", got.Text)
	}

	if strings.Contains(got.Text, "class VersionFail {};") {
		t.Fatalf("expected < comparison branch disabled, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class VersionAndDefined {};") {
		t.Fatalf("expected combined comparison + defined branch enabled, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_IfNumericComparisons_WithoutNumericMacroBody(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#define GAME_NAME DayZ
#if GAME_NAME >= 1
class ShouldNotAppear {};
#endif
#if defined(GAME_NAME)
class ShouldAppear {};
#endif
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableIfExpressions: true,
	})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	if strings.Contains(got.Text, "class ShouldNotAppear {};") {
		t.Fatalf("expected non-numeric macro compare to evaluate false, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "class ShouldAppear {};") {
		t.Fatalf("expected defined() branch enabled, got:\n%s", got.Text)
	}
}

func TestPreprocessFile_UnsupportedHasIncludeInIf(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "root.cpp")
	source := strings.TrimSpace(`
#if __has_include("x.hpp")
class ShouldNotParse {};
#endif
`) + "\n"

	if err := os.WriteFile(file, []byte(source), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	got, err := PreprocessFile(file, PreprocessOptions{
		EnableIfExpressions: true,
	})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic for __has_include, got %v", err)
	}

	found := false
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code == CodePPUnsupportedHasInclude {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected %d diagnostic, got: %+v", CodePPUnsupportedHasInclude, got.Diagnostics)
	}
}

func assertDynamicIntRange(t *testing.T, text string, name string, min int64, max int64) {
	t.Helper()

	m := regexp.MustCompile(name + ` = (-?\d+);`).FindStringSubmatch(text)
	if len(m) != 2 {
		t.Fatalf("expected %s assignment, got:\n%s", name, text)
	}

	value, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		t.Fatalf("parse %s as int64: %v", name, err)
	}

	if value < min || value > max {
		t.Fatalf("%s out of range: got=%d range=[%d,%d]", name, value, min, max)
	}
}

func assertDynamicUintRange(t *testing.T, text string, name string, max uint64) {
	t.Helper()

	m := regexp.MustCompile(name + ` = (\d+);`).FindStringSubmatch(text)
	if len(m) != 2 {
		t.Fatalf("expected %s assignment, got:\n%s", name, text)
	}

	value, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		t.Fatalf("parse %s as uint64: %v", name, err)
	}

	if value > max {
		t.Fatalf("%s out of range: got=%d max=%d", name, value, max)
	}
}

func assertDynamicIntPresent(t *testing.T, text string, name string) {
	t.Helper()

	if !regexp.MustCompile(name + ` = -?\d+;`).MatchString(text) {
		t.Fatalf("expected %s numeric assignment, got:\n%s", name, text)
	}
}

func assertDynamicUintPresent(t *testing.T, text string, name string) {
	t.Helper()

	if !regexp.MustCompile(name + ` = \d+;`).MatchString(text) {
		t.Fatalf("expected %s numeric assignment, got:\n%s", name, text)
	}
}

func stripHorizontalWhitespaceTest(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")

	return text
}
