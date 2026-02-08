package rvcfg

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestFormatRoundTripStabilityOnCorpus(t *testing.T) {
	t.Parallel()

	files := []string{
		testDataPath("cases", "parse", "vehicle", "config.cpp"),
		testDataPath("cases", "parse", "backpacks", "config.cpp"),
		testDataPath("cases", "parse", "model", "model.cfg"),
		testDataPath("cases", "macros", "config-explained.cpp"),
	}

	for _, path := range files {
		path := path
		t.Run(path, func(t *testing.T) {
			t.Parallel()

			input := readTestFile(t, path)
			opts := FormatOptions{}

			formattedA, err := FormatWithOptions(input, opts)
			if err != nil {
				t.Fatalf("FormatWithOptions(%s) error: %v", path, err)
			}

			parsedA, err := ParseBytes(path, formattedA, ParseOptions{
				CaptureScalarRaw: true,
			})
			if err != nil {
				t.Fatalf("ParseBytes(formattedA) error: %v", err)
			}

			formattedB, err := FormatWithOptions(formattedA, opts)
			if err != nil {
				t.Fatalf("FormatWithOptions(formattedA) error: %v", err)
			}

			if string(formattedA) != string(formattedB) {
				t.Fatal("formatter output is not idempotent on second pass")
			}

			parsedB, err := ParseBytes(path, formattedB, ParseOptions{
				CaptureScalarRaw: true,
			})
			if err != nil {
				t.Fatalf("ParseBytes(formattedB) error: %v", err)
			}

			digestA := astDigest(parsedA.File)
			digestB := astDigest(parsedB.File)
			if digestA != digestB {
				t.Fatal("AST digest mismatch after round-trip")
			}
		})
	}
}

// astDigest creates deterministic semantic digest for parser output.
// Positions and source labels are intentionally excluded.
func astDigest(file File) string {
	var b strings.Builder
	appendStatementDigest(&b, file.Statements)

	return b.String()
}

// appendStatementDigest appends stable statement/value shape recursively.
func appendStatementDigest(b *strings.Builder, statements []Statement) {
	for _, statement := range statements {
		b.WriteString("K=")
		b.WriteString(string(statement.Kind))
		b.WriteByte(';')

		switch statement.Kind {
		case NodeClass:
			if statement.Class == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.Class.Name)
			b.WriteString(";B=")
			b.WriteString(statement.Class.Base)
			b.WriteString(";F=")
			b.WriteString(strconv.FormatBool(statement.Class.Forward))
			b.WriteByte(';')
			appendStatementDigest(b, statement.Class.Body)
		case NodeDelete:
			if statement.Delete == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.Delete.Name)
			b.WriteByte(';')
		case NodeExtern:
			if statement.Extern == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.Extern.Name)
			b.WriteString(";C=")
			b.WriteString(strconv.FormatBool(statement.Extern.Class))
			b.WriteByte(';')
		case NodeProperty:
			if statement.Property == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.Property.Name)
			b.WriteByte(';')
			appendValueDigest(b, statement.Property.Value)
		case NodeArrayAssign:
			if statement.ArrayAssign == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.ArrayAssign.Name)
			b.WriteString(";A=")
			b.WriteString(strconv.FormatBool(statement.ArrayAssign.Append))
			b.WriteByte(';')
			appendValueDigest(b, statement.ArrayAssign.Value)
		case NodeEnum:
			if statement.Enum == nil {
				b.WriteString("nil;")
				continue
			}

			b.WriteString("N=")
			b.WriteString(statement.Enum.Name)
			b.WriteString(";L=")
			b.WriteString(strconv.Itoa(len(statement.Enum.Items)))
			b.WriteByte(';')

			for _, item := range statement.Enum.Items {
				b.WriteString("I=")
				b.WriteString(item.Name)
				b.WriteString(";V=")
				b.WriteString(item.ValueRaw)
				b.WriteByte(';')
			}
		}
	}
}

// appendValueDigest appends stable value shape recursively.
func appendValueDigest(b *strings.Builder, value Value) {
	b.WriteString("V=")
	b.WriteString(string(value.Kind))
	b.WriteByte(';')

	if value.Kind == ValueScalar {
		b.WriteString("R=")
		b.WriteString(value.Raw)
		b.WriteByte(';')

		return
	}

	if value.Kind != ValueArray {
		return
	}

	b.WriteString("L=")
	b.WriteString(strconv.Itoa(len(value.Elements)))
	b.WriteByte(';')

	for _, element := range value.Elements {
		appendValueDigest(b, element)
	}
}

// readTestFile loads fixture from repository with fatal-on-error behavior.
func readTestFile(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}

	return data
}
