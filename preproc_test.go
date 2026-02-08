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
		if d.Severity == SeverityError {
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

		if entry.Kind == "source" && entry.SourceFile == absInclude {
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
		expectedHash  = "ff8faade2d79d485e7faaa2599f5aeb2192839210ef469e52f2f5528c08c0511"
		expectedLines = 598
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

	got, err := PreprocessFile(file, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", file, err)
	}

	needles := []string{
		"nested[] = {func(1, (2 + 3)), finalValue};",
		`strPairs[] = {"left,right", "x,y"};`,
		"emptyLead[] = {, 2};",
		"emptyTail[] = {1, };",
	}

	for _, needle := range needles {
		if !strings.Contains(got.Text, needle) {
			t.Fatalf("expected expanded output to contain %q, got:\n%s", needle, got.Text)
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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

func TestPreprocessFile_UnresolvedMacroLikeInvocation(t *testing.T) {
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

	got, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnresolvedMacroInvocation) {
		t.Fatalf("expected ErrUnresolvedMacroInvocation, got %v", err)
	}

	found := false
	for _, diagnostic := range got.Diagnostics {
		if diagnostic.Code != CodePPUnresolvedMacroInvocation {
			continue
		}

		found = true

		break
	}

	if !found {
		t.Fatalf("expected %s diagnostic, got: %+v", CodePPUnresolvedMacroInvocation, got.Diagnostics)
	}
}

func TestPreprocessFile_UnresolvedMacroDetectionIgnoresStringsAndComments(t *testing.T) {
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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
		t.Fatalf("expected one %s warning after dedup, got=%d diagnostics=%+v", CodePPMacroRedefined, count, got.Diagnostics)
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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
		t.Fatalf("expected %s warning, diagnostics=%+v", CodePPMacroRedefined, got.Diagnostics)
	}

	if diagLine != 4 {
		t.Fatalf("expected %s warning on source line 4, got=%d diagnostics=%+v", CodePPMacroRedefined, diagLine, got.Diagnostics)
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

func TestPreprocessFile_ExecEval_DefaultUnsupported(t *testing.T) {
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

	_, err := PreprocessFile(file, PreprocessOptions{})
	if !errors.Is(err, ErrUnsupportedIntrinsic) {
		t.Fatalf("expected ErrUnsupportedIntrinsic, got %v", err)
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

	if !strings.Contains(got.Text, "lineValue = 3;") {
		t.Fatalf("expected __LINE__ to resolve to source line 3, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "fileValue = "+quoteIntrinsicString(file)+";") {
		t.Fatalf("expected __FILE__ to resolve to current file path, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileNameValue = "root.cpp";`) {
		t.Fatalf("expected __FILE_NAME__ to be root.cpp, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileShortValue = "root";`) {
		t.Fatalf("expected __FILE_SHORT__ to be root, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `stringLiteral = "__LINE__";`) {
		t.Fatalf("expected builtin in string literal to stay untouched, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, "// __FILE__") {
		t.Fatalf("expected builtin in comment to stay untouched, got:\n%s", got.Text)
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
		t.Fatalf("expected __LINE__ in include to resolve in include context, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileNameValue = "inc.hpp";`) {
		t.Fatalf("expected __FILE_NAME__ in include to be inc.hpp, got:\n%s", got.Text)
	}

	if !strings.Contains(got.Text, `fileShortValue = "inc";`) {
		t.Fatalf("expected __FILE_SHORT__ in include to be inc, got:\n%s", got.Text)
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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

	got, err := PreprocessFile(file, PreprocessOptions{})
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
		t.Fatalf("expected %s diagnostic, got: %+v", CodePPUnsupportedHasInclude, got.Diagnostics)
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
