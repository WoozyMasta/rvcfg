package rvcfg

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkLexVehicleConfig(b *testing.B) {
	data := benchLargeConfigData(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := LexBytes("bench-large.cpp", data)
		if err != nil {
			b.Fatalf("LexBytes benchmark error: %v", err)
		}
	}
}

func BenchmarkPreprocessConfig(b *testing.B) {
	root := benchIncludeFixtureRoot(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := PreprocessFile(root, PreprocessOptions{})
		if err != nil {
			b.Fatalf("PreprocessFile benchmark error: %v", err)
		}
	}
}

func BenchmarkParseVehicleConfig(b *testing.B) {
	data := benchLargeConfigData(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseBytes("bench-large.cpp", data, ParseOptions{})
		if err != nil {
			b.Fatalf("ParseBytes benchmark error: %v", err)
		}
	}
}

func BenchmarkParseVehicleConfigWithRaw(b *testing.B) {
	data := benchLargeConfigData(b)
	opts := ParseOptions{
		CaptureScalarRaw: true,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseBytes("bench-large.cpp", data, opts)
		if err != nil {
			b.Fatalf("ParseBytes benchmark error: %v", err)
		}
	}
}

func BenchmarkParseMalformedConfigRecoveryModes(b *testing.B) {
	data := benchLargeConfigData(b)

	// Replace semicolons to force parser errors and measure recovery overhead.
	malformed := make([]byte, len(data))
	copy(malformed, data)

	for i := 0; i < len(malformed); i++ {
		if malformed[i] == ';' {
			malformed[i] = ' '
		}
	}

	b.Run("RecoveryOn", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = ParseBytes("malformed.cpp", malformed, ParseOptions{})
		}
	})

	b.Run("RecoveryOff", func(b *testing.B) {
		opts := ParseOptions{
			DisableRecovery: true,
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, _ = ParseBytes("malformed.cpp", malformed, opts)
		}
	})
}

func BenchmarkFormatVehicleConfig(b *testing.B) {
	data := benchLargeConfigData(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Format(data)
		if err != nil {
			b.Fatalf("Format benchmark error: %v", err)
		}
	}
}

// benchLargeConfigData builds large deterministic config body for parser/formatter benchmarks.
func benchLargeConfigData(b *testing.B) []byte {
	b.Helper()

	const classCount = 1728
	var out strings.Builder
	out.Grow(classCount * 320)
	out.WriteString("class CfgPatches { class BenchPkg { requiredAddons[] = {}; }; };\n")
	out.WriteString("class CfgVehicles\n{\n")

	for id := 1; id <= classCount; id++ {
		appendBenchVehicleClass(&out, id, "  ")
	}

	out.WriteString("};\n")

	return []byte(out.String())
}

// benchIncludeFixtureRoot creates include-heavy preprocess fixture in benchmark temp directory.
func benchIncludeFixtureRoot(b *testing.B) string {
	b.Helper()

	const (
		partCount      = 48
		classesPerPart = 36
	)

	rootDir := b.TempDir()
	partsDir := filepath.Join(rootDir, "parts")
	if err := os.MkdirAll(partsDir, 0o700); err != nil {
		b.Fatalf("create parts dir: %v", err)
	}

	var root strings.Builder
	root.Grow(partCount * 48)
	root.WriteString("class CfgPatches { class BenchPkg { requiredAddons[] = {}; }; };\n")
	root.WriteString("class CfgVehicles\n{\n")

	for p := 1; p <= partCount; p++ {
		partName := "part_" + pad2(p) + ".hpp"
		partPath := filepath.Join(partsDir, partName)

		var part strings.Builder
		part.Grow(classesPerPart * 320)
		for i := 1; i <= classesPerPart; i++ {
			id := (p-1)*classesPerPart + i
			appendBenchVehicleClass(&part, id, "  ")
		}

		if err := os.WriteFile(partPath, []byte(part.String()), 0o600); err != nil {
			b.Fatalf("write include part %s: %v", partPath, err)
		}

		root.WriteString("#include \"parts/" + partName + "\"\n")
	}

	root.WriteString("};\n")
	rootPath := filepath.Join(rootDir, "root.cpp")
	if err := os.WriteFile(rootPath, []byte(root.String()), 0o600); err != nil {
		b.Fatalf("write include root: %v", err)
	}

	return rootPath
}

// appendBenchVehicleClass writes one deterministic benchmark vehicle class.
func appendBenchVehicleClass(out *strings.Builder, id int, indent string) {
	out.WriteString(indent)
	out.WriteString("class BenchVehicle_")
	out.WriteString(strconv.Itoa(id))
	out.WriteByte('\n')
	out.WriteString(indent)
	out.WriteString("{\n")
	out.WriteString(indent)
	out.WriteString("  scope = 2;\n")
	out.WriteString(indent)
	out.WriteString("  displayName = \"BenchVehicle_")
	out.WriteString(strconv.Itoa(id))
	out.WriteString("\";\n")
	out.WriteString(indent)
	out.WriteString("  wheelCount = 4;\n")
	out.WriteString(indent)
	out.WriteString("  pressureBySpeed[] = {0, 0.5, 10, 0.46, 30, 0.43, 40, 0.4, 60, 0.5, 80, 0.6};\n")
	out.WriteString(indent)
	out.WriteString("  class DamageSystem\n")
	out.WriteString(indent)
	out.WriteString("  {\n")
	out.WriteString(indent)
	out.WriteString("    class GlobalHealth\n")
	out.WriteString(indent)
	out.WriteString("    {\n")
	out.WriteString(indent)
	out.WriteString("      class Health\n")
	out.WriteString(indent)
	out.WriteString("      {\n")
	out.WriteString(indent)
	out.WriteString("        hitpoints = 500;\n")
	out.WriteString(indent)
	out.WriteString("        healthLevels[] = {{1, {}}, {0.7, {}}, {0.5, {}}, {0.3, {}}, {0, {}}};\n")
	out.WriteString(indent)
	out.WriteString("      };\n")
	out.WriteString(indent)
	out.WriteString("    };\n")
	out.WriteString(indent)
	out.WriteString("  };\n")
	out.WriteString(indent)
	out.WriteString("};\n")
}

// pad2 formats positive integer with zero padding to 2 digits.
func pad2(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}

	return strconv.Itoa(v)
}
