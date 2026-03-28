// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import "strings"

// joinMappedLines joins mapped output lines into one LF-normalized text.
func joinMappedLines(lines []mappedLine) string {
	if len(lines) == 0 {
		return ""
	}

	totalSize := len(lines) - 1
	for _, line := range lines {
		totalSize += len(line.text)
	}

	var out strings.Builder
	out.Grow(totalSize)

	for idx, line := range lines {
		if idx > 0 {
			out.WriteByte('\n')
		}

		out.WriteString(line.text)
	}

	return out.String()
}

// buildSourceMap converts mapped lines into output-line source map entries.
func buildSourceMap(lines []mappedLine, enabled bool) []SourceMapEntry {
	if !enabled || len(lines) == 0 {
		return nil
	}

	out := make([]SourceMapEntry, 0, len(lines))
	outputLine := 1
	for _, line := range lines {
		physicalLines := strings.Split(line.text, "\n")
		if len(physicalLines) == 0 {
			physicalLines = []string{""}
		}

		for _, physical := range physicalLines {
			width := len(physical)
			if width <= 0 {
				width = 1
			}

			out = append(out, SourceMapEntry{
				Kind:              line.kind,
				SourceFile:        line.sourceFile,
				SourceStartLine:   line.sourceLine,
				SourceEndLine:     line.sourceLine,
				OutputStartLine:   outputLine,
				OutputEndLine:     outputLine,
				SourceStartColumn: 1,
				SourceEndColumn:   width,
				OutputStartColumn: 1,
				OutputEndColumn:   width,
				IncludeFile:       line.include,
			})
			outputLine++
		}
	}

	return out
}
