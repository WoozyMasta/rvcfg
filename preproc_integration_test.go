package rvcfg

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// preprocParityCase describes one preprocess parity fixture.
type preprocParityCase struct {
	name                  string
	inputPath             string
	whitespaceInsensitive bool
}

func TestPreprocessParityWithCfgConvert_SmokeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	cases := []preprocParityCase{
		{
			name:      "define_basic",
			inputPath: testDataPath("preproc", "define_basic", "input.cpp"),
		},
		{
			name:      "include_basic",
			inputPath: testDataPath("preproc", "include_basic", "input.cpp"),
		},
		{
			name:      "include_angle",
			inputPath: testDataPath("preproc", "include_angle", "input.cpp"),
		},
		{
			name:      "include_quote_nospace",
			inputPath: testDataPath("preproc", "include_quote_nospace", "input.cpp"),
		},
		{
			name:      "include_angle_nospace",
			inputPath: testDataPath("preproc", "include_angle_nospace", "input.cpp"),
		},
		{
			name:      "include_tail_text",
			inputPath: testDataPath("preproc", "include_tail_text", "input.cpp"),
		},
		{
			name:      "include_with_comment",
			inputPath: testDataPath("preproc", "include_with_comment", "input.cpp"),
		},
		{
			name:      "exec_eval_passthru",
			inputPath: testDataPath("preproc", "exec_eval_passthru", "input.cpp"),
		},
		{
			name:      "file_name_passthru",
			inputPath: testDataPath("preproc", "file_name_passthru", "input.cpp"),
		},
		{
			name:      "unresolved_macro_passthru",
			inputPath: testDataPath("preproc", "unresolved_macro_passthru", "input.cpp"),
		},
		{
			name:      "fn_macro_in_string",
			inputPath: testDataPath("preproc", "fn_macro_in_string", "input.cpp"),
		},
		{
			name:      "fn_macro_arg_mismatch",
			inputPath: testDataPath("preproc", "fn_macro_arg_mismatch", "input.cpp"),
		},
		{
			name:      "fn_macro_bad_call",
			inputPath: testDataPath("preproc", "fn_macro_bad_call", "input.cpp"),
		},
		{
			name:      "fn_macro_bad_two_arg",
			inputPath: testDataPath("preproc", "fn_macro_bad_two_arg", "input.cpp"),
		},
		{
			name:      "fn_macro_bad_quoted",
			inputPath: testDataPath("preproc", "fn_macro_bad_quoted", "input.cpp"),
		},
		{
			name:      "fn_macro_string_comma",
			inputPath: testDataPath("preproc", "fn_macro_string_comma", "input.cpp"),
		},
		{
			name:      "fn_macro_quote_modes",
			inputPath: testDataPath("preproc", "fn_macro_quote_modes", "input.cpp"),
		},
		{
			name:      "single_quote_parse",
			inputPath: testDataPath("preproc", "single_quote_parse", "input.cpp"),
		},
		{
			name:      "tokenpaste_in_string",
			inputPath: testDataPath("preproc", "tokenpaste_in_string", "input.cpp"),
		},
		{
			name:      "fn_macro_in_single_quote",
			inputPath: testDataPath("preproc", "fn_macro_in_single_quote", "input.cpp"),
		},
		{
			name:      "include_builtin_file",
			inputPath: testDataPath("preproc", "include_builtin_file", "input.cpp"),
		},
		{
			name:      "define_nospace",
			inputPath: testDataPath("preproc", "define_nospace", "input.cpp"),
		},
		{
			name:      "dynamic_passthru",
			inputPath: testDataPath("preproc", "dynamic_passthru", "input.cpp"),
		},
		{
			name:      "ifdef_nested",
			inputPath: testDataPath("preproc", "ifdef_nested", "input.cpp"),
		},
		{
			name:                  "stringify_tokenpaste",
			inputPath:             testDataPath("preproc", "stringify_tokenpaste", "input.cpp"),
			whitespaceInsensitive: true,
		},
		{
			name:                  "macros_config_semantic",
			inputPath:             testDataPath("cases", "macros", "config.cpp"),
			whitespaceInsensitive: true,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := PreprocessFile(tc.inputPath, PreprocessOptions{})
			if err != nil {
				t.Fatalf("PreprocessFile(%s) error: %v", tc.inputPath, err)
			}

			wantText := runCfgConvertPreprocess(t, exe, tc.inputPath)
			gotText := normalizeParityText(got.Text)
			wantText = normalizeParityText(wantText)
			if tc.whitespaceInsensitive {
				gotText = stripHorizontalWhitespace(gotText)
				wantText = stripHorizontalWhitespace(wantText)
			}

			if gotText != wantText {
				t.Fatalf(
					"preprocess parity mismatch for %s\n--- rvcfg ---\n%s\n--- CfgConvert ---\n%s",
					tc.name,
					gotText,
					wantText,
				)
			}
		})
	}
}

func TestPreprocessParityWithCfgConvert_HasIncludeFailsInStrictMode(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	path := filepath.Join(t.TempDir(), "has_include.cpp")
	source := strings.TrimSpace(`
#if __has_include("x.hpp")
class A {};
#else
class B {};
#endif
`) + "\n"
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, oursErr := PreprocessFile(path, PreprocessOptions{})
	if oursErr == nil {
		t.Fatalf("expected PreprocessFile to fail for __has_include in strict mode")
	}

	cmd := exec.Command(exe, "-pcpp", "-dst", filepath.Join(t.TempDir(), "out.cpp"), path)
	_, cfgErr := cmd.CombinedOutput()
	if cfgErr == nil {
		t.Fatalf("expected CfgConvert -pcpp to fail for __has_include")
	}
}

func TestPreprocessParityWithCfgConvert_UnsupportedDirectivesFailInStrictMode(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	cases := []struct {
		name   string
		source string
	}{
		{
			name: "if",
			source: strings.TrimSpace(`
#if 1
class A {};
#endif
`) + "\n",
		},
		{
			name: "elif",
			source: strings.TrimSpace(`
#ifdef A
#elif 1
class A {};
#endif
`) + "\n",
		},
		{
			name: "if_inside_inactive_ifdef",
			source: strings.TrimSpace(`
#ifdef MISSING
#if 1
class A {};
#endif
#endif
`) + "\n",
		},
		{
			name: "warning",
			source: strings.TrimSpace(`
#warning this is unsupported in DayZ strict mode
class A {};
`) + "\n",
		},
		{
			name: "line",
			source: strings.TrimSpace(`
#line 42
class A {};
`) + "\n",
		},
		{
			name: "error",
			source: strings.TrimSpace(`
#error hard stop
class A {};
`) + "\n",
		},
		{
			name: "undef_no_arg",
			source: strings.TrimSpace(`
#undef
class A {};
`) + "\n",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			path := filepath.Join(t.TempDir(), tc.name+".cpp")
			if err := os.WriteFile(path, []byte(tc.source), 0o600); err != nil {
				t.Fatalf("write fixture: %v", err)
			}

			_, oursErr := PreprocessFile(path, PreprocessOptions{})
			if oursErr == nil {
				t.Fatalf("expected PreprocessFile to fail for #%s", tc.name)
			}

			cmd := exec.Command(exe, "-pcpp", "-dst", filepath.Join(t.TempDir(), "out.cpp"), path)
			_, cfgErr := cmd.CombinedOutput()
			if cfgErr == nil {
				t.Fatalf("expected CfgConvert -pcpp to fail for #%s", tc.name)
			}
		})
	}
}

// runCfgConvertPreprocess executes CfgConvert -pcpp and returns generated text.
func runCfgConvertPreprocess(t *testing.T, exe string, sourcePath string) string {
	t.Helper()

	outPath := filepath.Join(t.TempDir(), "cfgconvert.out.cpp")
	cmd := exec.Command(exe, "-pcpp", "-dst", outPath, sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf(
			"CfgConvert -pcpp failed for %s: %v\noutput=%s",
			sourcePath,
			err,
			string(output),
		)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read CfgConvert output %s: %v", outPath, err)
	}

	return string(data)
}

// normalizeParityText removes non-semantic formatting differences.
func normalizeParityText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		out = append(out, trimmed)
	}

	return strings.TrimSpace(strings.Join(out, "\n"))
}

// stripHorizontalWhitespace removes spaces and tabs for semantic text comparison.
func stripHorizontalWhitespace(text string) string {
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")

	return text
}

// firstDiffLine returns first differing line and values for diagnostics.
func firstDiffLine(left string, right string) (int, string, string) {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
	max := len(leftLines)
	if len(rightLines) > max {
		max = len(rightLines)
	}

	for idx := 0; idx < max; idx++ {
		lv := "<EOF>"
		if idx < len(leftLines) {
			lv = leftLines[idx]
		}

		rv := "<EOF>"
		if idx < len(rightLines) {
			rv = rightLines[idx]
		}

		if lv != rv {
			return idx + 1, lv, rv
		}
	}

	return 0, "", ""
}

func TestPreprocessParityWithCfgConvert_MissingFileReturnsErrorCode(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	missingPath := filepath.Join(t.TempDir(), "missing_input.cpp")
	cmd := exec.Command(exe, "-pcpp", "-dst", filepath.Join(t.TempDir(), "out.cpp"), missingPath)
	_, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected CfgConvert to fail for missing file: %s", missingPath)
	}
}

func TestPreprocessParityWithCfgConvert_DirectiveTailBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	path := filepath.Join(t.TempDir(), "directive_tail.cpp")
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

	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	got, err := PreprocessFile(path, PreprocessOptions{})
	if err != nil {
		t.Fatalf("PreprocessFile(%s) error: %v", path, err)
	}

	gotText := normalizeParityText(got.Text)
	wantText := normalizeParityText(runCfgConvertPreprocess(t, exe, path))
	if gotText != wantText {
		t.Fatalf(
			"directive tail parity mismatch\n--- rvcfg ---\n%s\n--- CfgConvert ---\n%s",
			gotText,
			wantText,
		)
	}
}

// realParityCase describes one real-world file parity check case.
type realParityCase struct {
	name string
	path string
}

func TestPreprocessParityWithCfgConvert_RealProjectFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skip CfgConvert integration test in -short mode")
	}

	exe := requireCfgConvert(t)
	cases := []realParityCase{
		{name: "warpbox_assets_config", path: `P:\warpbox\assets\config.cpp`},
		{name: "warpbox_assets_data", path: `P:\warpbox\assets\data.hpp`},
		{name: "warpbox_sounds_config", path: `P:\warpbox\sounds\config.cpp`},
		{name: "warpbox_sounds_defs", path: `P:\warpbox\sounds\sounds.hpp`},
		{name: "warpbox_utils", path: `P:\warpbox\utils.hpp`},
		{name: "utesplus_surfaces_config", path: `P:\utesplus\data\surfaces\config.cpp`},
		{name: "utesplus_animals", path: `P:\utesplus\data\surfaces\sounds_animals.hpp`},
		{name: "utesplus_character", path: `P:\utesplus\data\surfaces\sounds_character.hpp`},
		{name: "utesplus_infected", path: `P:\utesplus\data\surfaces\sounds_infected.hpp`},
	}

	checked := 0
	exactMatches := 0
	var statsMu sync.Mutex

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if !isRegularFile(tc.path) {
				t.Skipf("real-world fixture is missing: %s", tc.path)
			}

			got, err := PreprocessFile(tc.path, PreprocessOptions{})
			if err != nil {
				t.Fatalf("PreprocessFile(%s) error: %v", tc.path, err)
			}

			wantText := runCfgConvertPreprocess(t, exe, tc.path)
			gotExact := normalizeParityText(got.Text)
			wantExact := normalizeParityText(wantText)
			gotSemantic := stripHorizontalWhitespace(gotExact)
			wantSemantic := stripHorizontalWhitespace(wantExact)
			if gotSemantic != wantSemantic {
				t.Fatalf(
					"semantic preprocess parity mismatch for %s\n--- rvcfg ---\n%s\n--- CfgConvert ---\n%s",
					tc.path,
					gotExact,
					wantExact,
				)
			}

			statsMu.Lock()
			checked++
			if gotExact == wantExact {
				exactMatches++
			}
			statsMu.Unlock()

			if gotExact != wantExact {
				line, oursLine, theirsLine := firstDiffLine(gotExact, wantExact)
				t.Logf(
					"exact mismatch (semantic ok) at line %d\nours: %s\ntheirs: %s",
					line,
					oursLine,
					theirsLine,
				)
			}
		})
	}

	t.Cleanup(func() {
		if checked == 0 {
			t.Logf("no real-world files found; all cases skipped")

			return
		}

		t.Logf("real-world semantic parity: %d/%d, exact parity: %d/%d", checked, checked, exactMatches, checked)
	})
}
