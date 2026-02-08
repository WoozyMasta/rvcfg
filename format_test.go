package rvcfg

import (
	"os"
	"strings"
	"testing"
)

func intPtr(v int) *int {
	return &v
}

func TestFormatCanonicalConfig(t *testing.T) {
	t.Parallel()

	input := []byte(`class CfgVehicles{class Car: Vehicle{wheels[]={1,2,{3,4}};speed=-1;};};`)
	got, err := Format(input)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	want := "" +
		"class CfgVehicles\n" +
		"{\n" +
		"  class Car: Vehicle\n" +
		"  {\n" +
		"    wheels[] = {1, 2, {3, 4}};\n" +
		"    speed = -1;\n" +
		"  };\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected formatted output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWithOptionsTabIndent(t *testing.T) {
	t.Parallel()

	input := []byte(`class Cfg{class A{value=1;};};`)
	got, err := FormatWithOptions(input, FormatOptions{
		IndentChar: "\t",
		IndentSize: 1,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"class Cfg\n" +
		"{\n" +
		"\tclass A\n" +
		"\t{\n" +
		"\t\tvalue = 1;\n" +
		"\t};\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected tab-indented output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatFallbackForInvalidInput(t *testing.T) {
	t.Parallel()

	input := []byte("class Broken {\r\nvalue = 1\r\n}\r\n")
	got, err := Format(input)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	want := "class Broken {\nvalue = 1\n}\n"
	if string(got) != want {
		t.Fatalf("unexpected fallback output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWrapsLongArrayByMaxLineWidth(t *testing.T) {
	t.Parallel()

	input := []byte(`healthLevels[] = {{1, {}}, {0.7, {}}, {0.5, {}}, {0.3, {}}, {0, {}}};`)
	got, err := FormatWithOptions(input, FormatOptions{
		MaxLineWidth: 40,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"healthLevels[] =\n" +
		"{\n" +
		"  {1, {}},\n" +
		"  {0.7, {}},\n" +
		"  {0.5, {}},\n" +
		"  {0.3, {}},\n" +
		"  {0, {}},\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected wrapped output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatKeepsShortArrayInlineByMaxLineWidth(t *testing.T) {
	t.Parallel()

	input := []byte(`healthLevels[] = {{1, {}}, {0.7, {}}, {0.5, {}}, {0.3, {}}, {0, {}}};`)
	got, err := FormatWithOptions(input, FormatOptions{
		MaxLineWidth: 120,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "healthLevels[] = {{1, {}}, {0.7, {}}, {0.5, {}}, {0.3, {}}, {0, {}}};\n"
	if string(got) != want {
		t.Fatalf("unexpected inline output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWrapsArrayByElementLimit(t *testing.T) {
	t.Parallel()

	input := []byte(`pressureBySpeed[] = {0, 0.5, 10, 0.46, 30, 0.43, 40, 0.4, 60, 0.5, 80, 0.6};`)
	got, err := FormatWithOptions(input, FormatOptions{
		MaxLineWidth:           200,
		MaxInlineArrayElements: 8,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"pressureBySpeed[] =\n" +
		"{\n" +
		"  0,\n" +
		"  0.5,\n" +
		"  10,\n" +
		"  0.46,\n" +
		"  30,\n" +
		"  0.43,\n" +
		"  40,\n" +
		"  0.4,\n" +
		"  60,\n" +
		"  0.5,\n" +
		"  80,\n" +
		"  0.6,\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected element-limit output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatCompactsEmptyClassByDefault(t *testing.T) {
	t.Parallel()

	input := []byte("class Health: Health{};")
	got, err := Format(input)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	want := "class Health: Health {};\n"
	if string(got) != want {
		t.Fatalf("unexpected compact empty class output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatCanDisableCompactEmptyClass(t *testing.T) {
	t.Parallel()

	input := []byte("class Health: Health{};")
	got, err := FormatWithOptions(input, FormatOptions{
		DisableCompactEmptyClass: true,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"class Health: Health\n" +
		"{\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected expanded empty class output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWithAutoFixMissingClassSemicolon(t *testing.T) {
	t.Parallel()

	input := []byte(`class Root{class Child{value=1;}};`)

	gotDefault, err := Format(input)
	if err != nil {
		t.Fatalf("Format default options error: %v", err)
	}

	if string(gotDefault) != "class Root{class Child{value=1;}};" {
		t.Fatalf("unexpected default fallback output: %q", string(gotDefault))
	}

	gotCompat, err := FormatWithOptions(input, FormatOptions{
		AutoFixMissingClassSemicolon: true,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions compat error: %v", err)
	}

	wantCompat := "" +
		"class Root\n" +
		"{\n" +
		"  class Child\n" +
		"  {\n" +
		"    value = 1;\n" +
		"  };\n" +
		"};\n"

	if string(gotCompat) != wantCompat {
		t.Fatalf("unexpected compat formatted output\nwant:\n%s\ngot:\n%s", wantCompat, string(gotCompat))
	}
}

func TestFormatEnumDeclaration(t *testing.T) {
	t.Parallel()

	input := []byte(`enum EType{A,B=2,C=B+3,};`)
	got, err := Format(input)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	want := "" +
		"enum EType\n" +
		"{\n" +
		"  A,\n" +
		"  B = 2,\n" +
		"  C = B+3,\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected enum formatted output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWithPreserveComments(t *testing.T) {
	t.Parallel()

	input := []byte("" +
		"// top comment\n" +
		"class Cfg{value=1; // inline value\n" +
		"}; // inline class\n")

	got, err := FormatWithOptions(input, FormatOptions{
		PreserveComments: true,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"// top comment\n" +
		"class Cfg\n" +
		"{\n" +
		"  value = 1; // inline value\n" +
		"}; // inline class\n"

	if string(got) != want {
		t.Fatalf("unexpected preserve-comments output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatWithPreserveComments_Fixture(t *testing.T) {
	t.Parallel()

	path := testDataPath("cases", "comments", "mixed.cpp")
	input, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}

	got, err := FormatWithOptions(input, FormatOptions{
		PreserveComments: true,
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	text := string(got)
	needles := []string{
		"/* top block comment */",
		"// top line comment",
		"class CfgComments",
		"// property leading comment",
		"value = 1; // property trailing comment",
		"/* array leading block comment",
		"keeps inner indentation",
		"and irregular spacing */",
		"values[] = {1, 2, 3}; // array trailing comment",
		"}; // class trailing comment",
		"/* footer block comment */",
	}

	for _, needle := range needles {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected preserve-comments output to contain %q, got:\n%s", needle, text)
		}
	}
}

func TestFormatPreserveBlankLinesDefaultOne(t *testing.T) {
	t.Parallel()

	input := []byte("class A{};\n\n \t \nclass B{};\n")
	got, err := Format(input)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	want := "" +
		"class A {};\n" +
		"\n" +
		"class B {};\n"

	if string(got) != want {
		t.Fatalf("unexpected blank-line preservation output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestFormatPreserveBlankLinesDisabled(t *testing.T) {
	t.Parallel()

	input := []byte("class A{};\n\n\nclass B{};\n")
	got, err := FormatWithOptions(input, FormatOptions{
		PreserveBlankLines: intPtr(0),
	})
	if err != nil {
		t.Fatalf("FormatWithOptions error: %v", err)
	}

	want := "" +
		"class A {};\n" +
		"class B {};\n"

	if string(got) != want {
		t.Fatalf("unexpected blank-line disabled output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}
