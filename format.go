package rvcfg

import (
	"fmt"
	"strings"
)

// FormatOptions configures canonical config formatter.
type FormatOptions struct {
	// PreserveBlankLines limits how many empty lines to keep between neighbor statements.
	// nil means default 1. Explicit 0 disables blank-line preservation.
	PreserveBlankLines *int `json:"preserve_blank_lines,omitempty" yaml:"preserve_blank_lines,omitempty"`

	// IndentChar is indentation symbol. Supported values: " " or "\t".
	IndentChar string `json:"indent_char,omitempty" yaml:"indent_char,omitempty"`

	// IndentSize is count of IndentChar symbols per nesting level.
	IndentSize int `json:"indent_size,omitempty" yaml:"indent_size,omitempty"`

	// MaxLineWidth is soft target line width for wrapping long array assignments.
	// Values <= 0 disable width-based wrapping.
	MaxLineWidth int `json:"max_line_width,omitempty" yaml:"max_line_width,omitempty"`

	// MaxInlineArrayElements is soft limit for keeping arrays on a single line.
	// Values <= 0 disable count-based wrapping.
	MaxInlineArrayElements int `json:"max_inline_array_elements,omitempty" yaml:"max_inline_array_elements,omitempty"`

	// DisableCompactEmptyClass disables one-line empty class form `class X: Y {};`.
	DisableCompactEmptyClass bool `json:"disable_compact_empty_class,omitempty" yaml:"disable_compact_empty_class,omitempty"`

	// AutoFixMissingClassSemicolon enables safe autofix for missing class `;`
	// in formatter parse stage.
	AutoFixMissingClassSemicolon bool `json:"auto_fix_missing_class_semicolon,omitempty" yaml:"auto_fix_missing_class_semicolon,omitempty"`

	// PreserveComments keeps leading/trailing statement comments in formatted output.
	// This mode preserves comments as standalone lines near statement boundaries.
	PreserveComments bool `json:"preserve_comments,omitempty" yaml:"preserve_comments,omitempty"`
}

// Format applies canonical structural formatting for config syntax.
//
// When input cannot be parsed as config syntax, formatter falls back to deterministic
// line ending normalization to keep behavior useful for mixed corpuses.
func Format(input []byte) ([]byte, error) {
	return FormatWithOptions(input, FormatOptions{})
}

// FormatWithOptions applies canonical structural formatting with writer options.
func FormatWithOptions(input []byte, opts FormatOptions) ([]byte, error) {
	normalized := normalizeLineEndings(string(input))
	result, _ := ParseBytes("format-input", []byte(normalized), ParseOptions{
		CaptureScalarRaw:             true,
		AutoFixMissingClassSemicolon: opts.AutoFixMissingClassSemicolon,
		PreserveComments:             opts.PreserveComments,
	})
	if hasErrorDiagnostics(result.Diagnostics) {
		return []byte(normalized), nil
	}

	writer := newFormatter(opts)
	formatted, err := writer.formatFile(result.File)
	if err != nil {
		return nil, err
	}

	return []byte(formatted), nil
}

// formatter emits canonical text from parsed AST.
type formatter struct {
	builder                strings.Builder
	level                  int
	indentSize             int
	maxLineWidth           int
	maxInlineArrayElements int
	indentChar             byte
	compactEmptyClass      bool
	preserveComments       bool
	preserveBlankLines     int
}

// newFormatter creates AST formatter from options.
func newFormatter(opts FormatOptions) *formatter {
	indentChar := byte(' ')
	if opts.IndentChar == "\t" {
		indentChar = '\t'
	}

	indentSize := opts.IndentSize
	if indentSize <= 0 {
		indentSize = 2
	}

	compactEmptyClass := !opts.DisableCompactEmptyClass
	preserveBlankLines := 1
	if opts.PreserveBlankLines != nil {
		preserveBlankLines = *opts.PreserveBlankLines
		if preserveBlankLines < 0 {
			preserveBlankLines = 0
		}
	}

	return &formatter{
		maxLineWidth:           opts.MaxLineWidth,
		maxInlineArrayElements: opts.MaxInlineArrayElements,
		compactEmptyClass:      compactEmptyClass,
		indentSize:             indentSize,
		indentChar:             indentChar,
		preserveComments:       opts.PreserveComments,
		preserveBlankLines:     preserveBlankLines,
	}
}

// formatFile renders complete file statements.
func (f *formatter) formatFile(file File) (string, error) {
	for i, statement := range file.Statements {
		if err := f.writeStatement(statement); err != nil {
			return "", err
		}

		if i+1 < len(file.Statements) {
			f.writeInterStatementBlankLines(statement, file.Statements[i+1])
		}
	}

	return f.builder.String(), nil
}

// writeStatement serializes one AST statement node.
func (f *formatter) writeStatement(statement Statement) error {
	if f.preserveComments {
		f.writeLeadingComments(statement.LeadingComments)
	}

	switch statement.Kind {
	case NodeClass:
		if err := f.writeClass(statement.Class); err != nil {
			return err
		}
	case NodeDelete:
		if err := f.writeDelete(statement.Delete); err != nil {
			return err
		}
	case NodeProperty:
		if err := f.writeProperty(statement.Property); err != nil {
			return err
		}
	case NodeArrayAssign:
		if err := f.writeArrayAssign(statement.ArrayAssign); err != nil {
			return err
		}
	case NodeExtern:
		if err := f.writeExtern(statement.Extern); err != nil {
			return err
		}
	case NodeEnum:
		if err := f.writeEnum(statement.Enum); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported statement kind: %s", statement.Kind)
	}

	if f.preserveComments && statement.TrailingComment != nil {
		f.writeTrailingComment(statement.TrailingComment.Text)
	}

	if f.preserveComments && len(statement.TrailingComments) > 0 {
		f.writeLeadingComments(statement.TrailingComments)
	}

	return nil
}

// writeLeadingComments writes statement-level leading comments.
func (f *formatter) writeLeadingComments(comments []Comment) {
	for _, comment := range comments {
		text := strings.TrimSpace(comment.Text)
		if text == "" {
			continue
		}

		f.writeLine(text)
	}
}

// writeTrailingComment writes trailing comment inline when line width allows.
func (f *formatter) writeTrailingComment(commentText string) {
	comment := strings.TrimSpace(commentText)
	if comment == "" {
		return
	}

	if f.tryInlineTrailingComment(comment) {
		return
	}

	f.writeLine(comment)
}

// tryInlineTrailingComment appends trailing comment to the last emitted line.
func (f *formatter) tryInlineTrailingComment(comment string) bool {
	all := f.builder.String()
	if len(all) == 0 || all[len(all)-1] != '\n' {
		return false
	}

	end := len(all) - 1
	start := strings.LastIndex(all[:end], "\n") + 1
	lastLine := all[start:end]
	if strings.TrimSpace(lastLine) == "" {
		return false
	}

	inline := lastLine + " " + comment
	if f.maxLineWidth > 0 && len(inline) > f.maxLineWidth {
		return false
	}

	f.builder.Reset()
	f.builder.Grow(len(all) + 1 + len(comment))
	f.builder.WriteString(all[:start])
	f.builder.WriteString(inline)
	f.builder.WriteByte('\n')

	return true
}

// writeClass serializes class declaration.
func (f *formatter) writeClass(classDecl *ClassDecl) error {
	if classDecl == nil {
		return fmt.Errorf("class payload is nil")
	}

	header := "class " + classDecl.Name
	if classDecl.Base != "" {
		header += ": " + classDecl.Base
	}

	if classDecl.Forward {
		f.writeLine(header + ";")

		return nil
	}

	if f.compactEmptyClass && len(classDecl.Body) == 0 {
		f.writeLine(header + " {};")

		return nil
	}

	f.writeLine(header)
	f.writeLine("{")
	f.level++

	for i, statement := range classDecl.Body {
		if err := f.writeStatement(statement); err != nil {
			return err
		}

		if i+1 < len(classDecl.Body) {
			f.writeInterStatementBlankLines(statement, classDecl.Body[i+1])
		}
	}

	f.level--
	f.writeLine("};")

	return nil
}

// writeInterStatementBlankLines preserves original blank-line spacing up to configured limit.
func (f *formatter) writeInterStatementBlankLines(prev Statement, next Statement) {
	if f.preserveBlankLines <= 0 {
		return
	}

	if prev.End.Line <= 0 || next.Start.Line <= 0 {
		return
	}

	gap := next.Start.Line - prev.End.Line - 1
	if gap <= 0 {
		return
	}

	if gap > f.preserveBlankLines {
		gap = f.preserveBlankLines
	}

	for i := 0; i < gap; i++ {
		f.builder.WriteByte('\n')
	}
}

// writeDelete serializes delete statement.
func (f *formatter) writeDelete(deleteStmt *DeleteStmt) error {
	if deleteStmt == nil {
		return fmt.Errorf("delete payload is nil")
	}

	f.writeLine("delete " + deleteStmt.Name + ";")

	return nil
}

// writeExtern serializes extern declaration.
func (f *formatter) writeExtern(extern *ExternDecl) error {
	if extern == nil {
		return fmt.Errorf("extern payload is nil")
	}

	prefix := "extern "
	if extern.Class {
		prefix += "class "
	}

	f.writeLine(prefix + extern.Name + ";")

	return nil
}

// writeEnum serializes enum declaration.
func (f *formatter) writeEnum(enumDecl *EnumDecl) error {
	if enumDecl == nil {
		return fmt.Errorf("enum payload is nil")
	}

	header := "enum"
	if enumDecl.Name != "" {
		header += " " + enumDecl.Name
	}

	f.writeLine(header)
	f.writeLine("{")
	f.level++

	for _, item := range enumDecl.Items {
		line := item.Name
		if item.ValueRaw != "" {
			line += " = " + item.ValueRaw
		}

		f.writeLine(line + ",")
	}

	f.level--
	f.writeLine("};")

	return nil
}

// writeProperty serializes scalar property assignment.
func (f *formatter) writeProperty(property *PropertyAssign) error {
	if property == nil {
		return fmt.Errorf("property payload is nil")
	}

	return f.writeAssignment(property.Name, "=", property.Value, false)
}

// writeArrayAssign serializes array assignment or append statement.
func (f *formatter) writeArrayAssign(arrayAssign *ArrayAssign) error {
	if arrayAssign == nil {
		return fmt.Errorf("array assignment payload is nil")
	}

	operator := "="
	if arrayAssign.Append {
		operator = "+="
	}

	return f.writeAssignment(arrayAssign.Name+"[]", operator, arrayAssign.Value, true)
}

// writeAssignment serializes property and array assignment statements with soft-wrap.
func (f *formatter) writeAssignment(name string, operator string, value Value, wrapArray bool) error {
	inlineValue, err := f.valueString(value)
	if err != nil {
		return err
	}

	inlineLine := name + " " + operator + " " + inlineValue + ";"
	if !wrapArray || value.Kind != ValueArray || !f.shouldWrapArray(value, inlineLine) {
		f.writeLine(inlineLine)

		return nil
	}

	f.writeLine(name + " " + operator)
	if err := f.writeWrappedArray(value); err != nil {
		return err
	}

	f.writeLine("};")

	return nil
}

// writeWrappedArray serializes array value in multiline block form.
func (f *formatter) writeWrappedArray(value Value) error {
	if value.Kind != ValueArray {
		return fmt.Errorf("wrapped array expected ValueArray, got %s", value.Kind)
	}

	f.writeLine("{")
	f.level++

	for _, element := range value.Elements {
		inlineElement, err := f.valueString(element)
		if err != nil {
			return err
		}

		inlineElementLine := inlineElement + ","
		if element.Kind == ValueArray && f.shouldWrapArray(element, inlineElementLine) {
			if err := f.writeWrappedArrayElement(element); err != nil {
				return err
			}

			continue
		}

		f.writeLine(inlineElementLine)
	}

	f.level--

	return nil
}

// writeWrappedArrayElement serializes nested long array as multiline element.
func (f *formatter) writeWrappedArrayElement(value Value) error {
	if value.Kind != ValueArray {
		return fmt.Errorf("wrapped array element expected ValueArray, got %s", value.Kind)
	}

	f.writeLine("{")
	f.level++

	for _, element := range value.Elements {
		inlineElement, err := f.valueString(element)
		if err != nil {
			return err
		}

		inlineElementLine := inlineElement + ","
		if element.Kind == ValueArray && f.shouldWrapArray(element, inlineElementLine) {
			if err := f.writeWrappedArrayElement(element); err != nil {
				return err
			}

			continue
		}

		f.writeLine(inlineElementLine)
	}

	f.level--
	f.writeLine("},")

	return nil
}

// shouldWrap checks whether current line exceeds configured soft width.
func (f *formatter) shouldWrap(line string) bool {
	if f.maxLineWidth <= 0 {
		return false
	}

	return f.currentIndentWidth()+len(line) > f.maxLineWidth
}

// shouldWrapArray combines width-based and element-count based wrapping conditions.
func (f *formatter) shouldWrapArray(value Value, line string) bool {
	if f.shouldWrap(line) {
		return true
	}

	if value.Kind != ValueArray {
		return false
	}

	// Count-based wrap: useful for dense arrays that still fit width but become hard to read.
	if f.maxInlineArrayElements > 0 && len(value.Elements) > f.maxInlineArrayElements {
		return true
	}

	return false
}

// currentIndentWidth returns current indentation width in bytes.
func (f *formatter) currentIndentWidth() int {
	return f.level * f.indentSize
}

// valueString serializes scalar or nested array value.
func (f *formatter) valueString(value Value) (string, error) {
	switch value.Kind {
	case ValueScalar:
		return value.Raw, nil
	case ValueArray:
		parts := make([]string, 0, len(value.Elements))

		for _, element := range value.Elements {
			text, err := f.valueString(element)
			if err != nil {
				return "", err
			}

			parts = append(parts, text)
		}

		return "{" + strings.Join(parts, ", ") + "}", nil
	default:
		return "", fmt.Errorf("unsupported value kind: %s", value.Kind)
	}
}

// writeLine appends one formatted line with current indentation.
func (f *formatter) writeLine(line string) {
	n := f.level * f.indentSize
	for i := 0; i < n; i++ {
		f.builder.WriteByte(f.indentChar)
	}

	f.builder.WriteString(line)
	f.builder.WriteByte('\n')
}

// hasErrorDiagnostics checks whether diagnostics contain any error-level issue.
func hasErrorDiagnostics(diagnostics []Diagnostic) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}

	return false
}
