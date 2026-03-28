// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

// sourceMapResolver resolves output parser positions back to source positions.
type sourceMapResolver struct {
	entries []SourceMapEntry
}

// newSourceMapResolver builds immutable resolver over preprocessor source map.
func newSourceMapResolver(entries []SourceMapEntry) sourceMapResolver {
	return sourceMapResolver{entries: entries}
}

// remapDiagnostics remaps diagnostic positions when source map is present.
func (r sourceMapResolver) remapDiagnostics(diagnostics []Diagnostic) {
	for index := range diagnostics {
		diagnostics[index].Start = r.remapPosition(diagnostics[index].Start)
		diagnostics[index].End = r.remapPosition(diagnostics[index].End)
	}
}

// remapFile remaps all AST positions recursively when source map is present.
func (r sourceMapResolver) remapFile(file *File) {
	if file == nil {
		return
	}

	file.Start = r.remapPosition(file.Start)
	file.End = r.remapPosition(file.End)
	r.remapStatements(file.Statements)
}

// remapStatements remaps statement positions recursively.
func (r sourceMapResolver) remapStatements(statements []Statement) {
	for index := range statements {
		statement := &statements[index]
		statement.Start = r.remapPosition(statement.Start)
		statement.End = r.remapPosition(statement.End)

		if statement.TrailingComment != nil {
			statement.TrailingComment.Start = r.remapPosition(statement.TrailingComment.Start)
			statement.TrailingComment.End = r.remapPosition(statement.TrailingComment.End)
		}

		for commentIndex := range statement.LeadingComments {
			statement.LeadingComments[commentIndex].Start = r.remapPosition(statement.LeadingComments[commentIndex].Start)
			statement.LeadingComments[commentIndex].End = r.remapPosition(statement.LeadingComments[commentIndex].End)
		}

		for commentIndex := range statement.TrailingComments {
			statement.TrailingComments[commentIndex].Start = r.remapPosition(statement.TrailingComments[commentIndex].Start)
			statement.TrailingComments[commentIndex].End = r.remapPosition(statement.TrailingComments[commentIndex].End)
		}

		switch statement.Kind {
		case NodeClass:
			if statement.Class == nil {
				continue
			}

			r.remapStatements(statement.Class.Body)

		case NodeProperty:
			if statement.Property == nil {
				continue
			}

			r.remapValue(&statement.Property.Value)

		case NodeArrayAssign:
			if statement.ArrayAssign == nil {
				continue
			}

			r.remapValue(&statement.ArrayAssign.Value)
		}
	}
}

// remapValue remaps nested value positions recursively.
func (r sourceMapResolver) remapValue(value *Value) {
	if value == nil {
		return
	}

	value.Start = r.remapPosition(value.Start)
	value.End = r.remapPosition(value.End)
	for index := range value.Elements {
		r.remapValue(&value.Elements[index])
	}
}

// remapPosition remaps one output position into source map coordinates.
func (r sourceMapResolver) remapPosition(position Position) Position {
	if len(r.entries) == 0 || position.Line <= 0 {
		return position
	}

	entry, ok := r.findEntryByOutputLine(position.Line)
	if !ok {
		return position
	}

	remapped := position
	if entry.SourceFile != "" {
		remapped.File = entry.SourceFile
	}

	lineDelta := max(position.Line-entry.OutputStartLine, 0)

	remapped.Line = entry.SourceStartLine + lineDelta
	remapped.Column = remapColumn(position.Column, entry)

	return remapped
}

// findEntryByOutputLine finds source map entry containing output line.
func (r sourceMapResolver) findEntryByOutputLine(outputLine int) (SourceMapEntry, bool) {
	low := 0
	high := len(r.entries) - 1

	for low <= high {
		mid := low + (high-low)/2
		entry := r.entries[mid]
		if outputLine < entry.OutputStartLine {
			high = mid - 1

			continue
		}

		if outputLine > entry.OutputEndLine {
			low = mid + 1

			continue
		}

		return entry, true
	}

	return SourceMapEntry{}, false
}

// remapColumn maps output column into source column range, best effort.
func remapColumn(outputColumn int, entry SourceMapEntry) int {
	if outputColumn <= 0 {
		if entry.SourceStartColumn > 0 {
			return entry.SourceStartColumn
		}

		return 1
	}

	outputStart := entry.OutputStartColumn
	if outputStart <= 0 {
		outputStart = 1
	}

	sourceStart := entry.SourceStartColumn
	if sourceStart <= 0 {
		sourceStart = 1
	}

	column := max(sourceStart+(outputColumn-outputStart), 1)

	return column
}
