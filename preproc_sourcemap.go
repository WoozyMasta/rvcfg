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

// buildSourceMap compresses mapped lines into source map entries.
func buildSourceMap(lines []mappedLine, enabled bool) []SourceMapEntry {
	if !enabled || len(lines) == 0 {
		return nil
	}

	out := make([]SourceMapEntry, 0, len(lines))
	for idx, line := range lines {
		outputLine := idx + 1
		next := SourceMapEntry{
			Kind:            line.kind,
			SourceFile:      line.sourceFile,
			SourceStartLine: line.sourceLine,
			SourceEndLine:   line.sourceLine,
			OutputStartLine: outputLine,
			OutputEndLine:   outputLine,
			IncludeFile:     line.include,
		}

		if len(out) == 0 {
			out = append(out, next)

			continue
		}

		last := &out[len(out)-1]
		if canMergeSourceMapEntry(*last, next) {
			last.SourceEndLine = next.SourceEndLine
			last.OutputEndLine = next.OutputEndLine

			continue
		}

		out = append(out, next)
	}

	return out
}

// canMergeSourceMapEntry checks whether two adjacent entries can be merged.
func canMergeSourceMapEntry(left SourceMapEntry, right SourceMapEntry) bool {
	if left.Kind != "source" || right.Kind != "source" {
		return false
	}

	if left.SourceFile != right.SourceFile || left.IncludeFile != right.IncludeFile {
		return false
	}

	if right.SourceStartLine != left.SourceEndLine+1 {
		return false
	}

	return right.OutputStartLine == left.OutputEndLine+1
}
