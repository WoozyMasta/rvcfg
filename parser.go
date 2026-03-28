// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"os"

	"github.com/woozymasta/lintkit/lint"
)

var (
	externIdentifier = []byte("extern")
)

// ParseOptions configures parser behavior.
type ParseOptions struct {
	// DisableRecovery disables statement-level error recovery.
	DisableRecovery bool `json:"disable_recovery,omitempty" yaml:"disable_recovery,omitempty"`

	// CaptureScalarRaw stores scalar literal text into Value.Raw.
	// Disabled by default to reduce parser allocations.
	CaptureScalarRaw bool `json:"capture_scalar_raw,omitempty" yaml:"capture_scalar_raw,omitempty"`

	// Strict enables conservative grammar validation.
	Strict bool `json:"strict,omitempty" yaml:"strict,omitempty"`

	// AutoFixMissingClassSemicolon enables safe autofix for missing class `;`.
	// Autofix is applied only when parser can prove clear statement boundary.
	// In ambiguous cases parser still emits PAR005 and fails.
	AutoFixMissingClassSemicolon bool `json:"auto_fix_missing_class_semicolon,omitempty" yaml:"auto_fix_missing_class_semicolon,omitempty"`

	// PreserveComments keeps leading/trailing statement comments in AST.
	PreserveComments bool `json:"preserve_comments,omitempty" yaml:"preserve_comments,omitempty"`
}

// ParseResult stores parser output with diagnostics.
type ParseResult struct {
	// Diagnostics are parser and lexer diagnostics.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty" yaml:"diagnostics,omitempty"`

	// File is parsed AST root.
	File File `json:"file,omitzero" yaml:"file,omitempty"`
}

// valueStopMask stores stop-token set for scalar parsing.
type valueStopMask uint8

const (
	stopSemicolon valueStopMask = 1 << iota
	stopRBrace
	stopComma
)

const arenaChunkSize = 64

// arena stores parser node payloads in fixed-size chunks to reduce per-node allocations.
type arena[T any] struct {
	chunks [][]T
	next   int
}

// alloc returns pointer to a zero-value node slot in arena.
func (a *arena[T]) alloc() *T {
	if len(a.chunks) == 0 || a.next >= len(a.chunks[len(a.chunks)-1]) {
		a.chunks = append(a.chunks, make([]T, arenaChunkSize))
		a.next = 0
	}

	chunk := a.chunks[len(a.chunks)-1]
	slot := &chunk[a.next]
	var zero T
	*slot = zero
	a.next++

	return slot
}

// ParseFile parses source file into AST in raw mode (without preprocess stage).
func ParseFile(path string, opts ParseOptions) (ParseResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ParseResult{}, fmt.Errorf("read source file %q: %w", path, err)
	}

	return ParseBytes(path, data, opts)
}

// ParseBytes parses source bytes into AST in raw mode (without preprocess stage).
func ParseBytes(filename string, data []byte, opts ParseOptions) (ParseResult, error) {
	lexOpts := LexOptions{}
	if opts.PreserveComments {
		lexOpts.CaptureLexeme = true
		lexOpts.EmitComments = true
		lexOpts.EmitNewlines = true
	}

	tokens, lexDiagnostics, lexErr := LexBytesWithOptions(filename, data, lexOpts)
	parseResult, parseErr := ParseTokens(filename, data, tokens, opts)
	diagnostics := make([]Diagnostic, 0, len(lexDiagnostics)+len(parseResult.Diagnostics))
	diagnostics = append(diagnostics, lexDiagnostics...)
	diagnostics = append(diagnostics, parseResult.Diagnostics...)

	parseResult.Diagnostics = diagnostics

	if lexErr != nil {
		return parseResult, lexErr
	}

	if parseErr != nil {
		return parseResult, parseErr
	}

	return parseResult, nil
}

// ParseTokens parses pre-tokenized source into AST without lexer stage.
func ParseTokens(filename string, data []byte, tokens []Token, opts ParseOptions) (ParseResult, error) {
	p := newParser(tokens, data, opts)
	file := p.parseFile(filename)

	result := ParseResult{
		File:        file,
		Diagnostics: p.diagnostics,
	}

	hasError := false
	for _, d := range result.Diagnostics {
		if d.Severity != lint.SeverityError {
			continue
		}

		hasError = true

		break
	}

	if !hasError {
		result.Diagnostics = append(
			result.Diagnostics,
			collectInheritanceHints(result.File)...,
		)
	}

	for _, d := range result.Diagnostics {
		if d.Severity == lint.SeverityError {
			return result, ErrParse
		}
	}

	return result, nil
}

// parser stores mutable recursive-descent parse state.
type parser struct {
	classArena       arena[ClassDecl]
	deleteArena      arena[DeleteStmt]
	externArena      arena[ExternDecl]
	enumArena        arena[EnumDecl]
	propertyArena    arena[PropertyAssign]
	source           []byte
	tokens           []Token
	diagnostics      []Diagnostic
	arrayAssignArena arena[ArrayAssign]
	index            int
	captureScalarRaw bool
	autoFixClassSem  bool
	preserveComments bool
	strict           bool
	recovery         bool
}

// newParser builds parser with configured recovery mode.
func newParser(tokens []Token, source []byte, opts ParseOptions) *parser {
	return &parser{
		source:           source,
		tokens:           tokens,
		diagnostics:      make([]Diagnostic, 0),
		captureScalarRaw: opts.CaptureScalarRaw,
		autoFixClassSem:  opts.AutoFixMissingClassSemicolon,
		preserveComments: opts.PreserveComments,
		strict:           opts.Strict,
		recovery:         !opts.DisableRecovery,
	}
}

// parseFile parses complete file statements until EOF.
func (p *parser) parseFile(filename string) File {
	file := File{
		Source:     filename,
		Statements: p.parseStatements(false),
	}

	if len(p.tokens) == 0 {
		return file
	}

	file.Start = p.tokens[0].Start
	file.End = p.tokens[len(p.tokens)-1].End

	return file
}

// parseStatements parses statement list until EOF or optional closing brace.
func (p *parser) parseStatements(stopAtBrace bool) []Statement {
	capHint := 16
	if !stopAtBrace {
		capHint = 128
	}

	statements := make([]Statement, 0, capHint)

	for {
		leadingComments := []Comment(nil)
		if p.preserveComments {
			leadingComments = p.consumeLeadingComments()
		} else {
			p.skipTrivia()
		}

		if p.isEOF() {
			if p.preserveComments && len(leadingComments) > 0 && len(statements) > 0 {
				last := &statements[len(statements)-1]
				last.TrailingComments = append(last.TrailingComments, leadingComments...)
			}

			break
		}

		if stopAtBrace && p.peek().Kind == TokenRBrace {
			break
		}

		start := p.index
		stmt, ok := p.parseStatement(stopAtBrace)
		if ok {
			if p.preserveComments {
				if len(leadingComments) > 0 {
					stmt.LeadingComments = leadingComments
				}

				if trailing, ok := p.consumeTrailingComment(stmt.End.Line); ok {
					stmt.TrailingComment = &trailing
				}
			}

			statements = append(statements, stmt)

			continue
		}

		// Failsafe to avoid infinite loops in any unhandled branch.
		if p.index == start {
			p.advance()
		}

		if !p.recovery {
			break
		}
	}

	return statements
}

// parseStatement parses one declaration or assignment statement.
func (p *parser) parseStatement(stopAtBrace bool) (Statement, bool) {
	token := p.peek()

	switch token.Kind {
	case TokenSemicolon:
		p.advance()

		return Statement{}, false

	case TokenKeywordClass:
		return p.parseClass(stopAtBrace)

	case TokenKeywordDelete:
		return p.parseDelete(stopAtBrace)

	case TokenKeywordEnum:
		return p.parseEnum(stopAtBrace)

	case TokenIdentifier:
		if p.tokenEquals(token, externIdentifier) {
			return p.parseExtern(stopAtBrace)
		}

		return p.parseAssignment(stopAtBrace)
	}

	p.emitError(CodeParUnexpectedToken, token.Start, "unexpected token "+token.Kind.String())
	p.recoverStatement(stopAtBrace)

	return Statement{}, false
}

// parseClass parses class declaration with body or forward form.
func (p *parser) parseClass(stopAtBrace bool) (Statement, bool) {
	startToken := p.consume()
	name, ok := p.parseClassLikeName(CodeParExpectedClassName, "expected class name")
	if !ok {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	base := ""
	if p.match(TokenColon) {
		base = p.parseBaseExpression()
	}

	decl := ClassDecl{
		Name: name,
		Base: base,
	}

	if p.match(TokenSemicolon) {
		classDecl := p.classArena.alloc()
		*classDecl = decl
		classDecl.Forward = true

		stmt := Statement{
			Kind:  NodeClass,
			Class: classDecl,
			Start: startToken.Start,
			End:   p.prev().End,
		}

		return stmt, true
	}

	if !p.match(TokenLBrace) {
		p.emitError(CodeParExpectedClassBodyOrSemicolon, p.peek().Start, "expected '{' or ';' after class declaration")
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	decl.Body = p.parseStatements(true)
	if !p.expect(TokenRBrace, CodeParExpectedClassClosingBrace, "expected '}' to close class body") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	if !p.match(TokenSemicolon) {
		if p.autoFixClassSem && p.canRecoverImplicitClassSemicolon() {
			p.emitWarning(CodeParAutofixClassSemicolon, p.prev().End, "autofix: inserted missing ';' after class declaration")
		} else {
			p.emitError(CodeParMissingClassSemicolon, p.peek().Start, "missing ';' after class declaration")
			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}
	}

	classDecl := p.classArena.alloc()
	*classDecl = decl

	stmt := Statement{
		Kind:  NodeClass,
		Class: classDecl,
		Start: startToken.Start,
		End:   p.prev().End,
	}

	return stmt, true
}

// parseDelete parses delete declaration.
func (p *parser) parseDelete(stopAtBrace bool) (Statement, bool) {
	startToken := p.consume()
	name, ok := p.parseClassLikeName(CodeParExpectedDeleteName, "expected name after delete")
	if !ok {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	if !p.expect(TokenSemicolon, CodeParMissingDeleteSemicolon, "missing ';' after delete statement") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	deleteNode := p.deleteArena.alloc()
	deleteNode.Name = name

	stmt := Statement{
		Kind:   NodeDelete,
		Delete: deleteNode,
		Start:  startToken.Start,
		End:    p.prev().End,
	}

	return stmt, true
}

// parseExtern parses extern declaration, with optional `class` keyword.
func (p *parser) parseExtern(stopAtBrace bool) (Statement, bool) {
	startToken := p.consume()
	extern := p.externArena.alloc()

	if p.match(TokenKeywordClass) {
		extern.Class = true
	}

	name, ok := p.parseClassLikeName(CodeParExpectedExternName, "expected name in extern declaration")
	if !ok {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	extern.Name = name
	if !p.expect(TokenSemicolon, CodeParMissingExternSemicolon, "missing ';' after extern declaration") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	stmt := Statement{
		Kind:   NodeExtern,
		Extern: extern,
		Start:  startToken.Start,
		End:    p.prev().End,
	}

	return stmt, true
}

// parseEnum parses enum declaration.
func (p *parser) parseEnum(stopAtBrace bool) (Statement, bool) {
	startToken := p.consume()
	enumDecl := EnumDecl{
		Items: make([]EnumItem, 0, 8),
	}

	p.skipTrivia()
	if p.peek().Kind == TokenIdentifier {
		enumDecl.Name = p.tokenText(p.consume())
	}

	if !p.expect(TokenLBrace, CodeParExpectedEnumBody, "expected '{' to open enum body") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	for {
		p.skipTrivia()
		if p.isEOF() {
			p.emitError(CodeParExpectedEnumDelimiter, p.prev().End, "unterminated enum declaration")

			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}

		if p.match(TokenRBrace) {
			break
		}

		if p.peek().Kind != TokenIdentifier {
			p.emitError(CodeParExpectedEnumItemName, p.peek().Start, "expected enum item name")
			p.recoverEnumItem()

			if p.match(TokenComma) {
				continue
			}

			if p.match(TokenRBrace) {
				break
			}

			continue
		}

		nameToken := p.consume()
		item := EnumItem{
			Name: p.tokenText(nameToken),
		}

		if p.match(TokenAssign) {
			valueRaw, ok := p.parseEnumValueRaw(stopComma | stopRBrace)
			if !ok {
				p.recoverEnumItem()
			} else {
				item.ValueRaw = valueRaw
			}
		}

		enumDecl.Items = append(enumDecl.Items, item)

		p.skipTrivia()
		if p.match(TokenComma) {
			p.skipTrivia()
			if p.peek().Kind == TokenRBrace {
				continue
			}

			continue
		}

		if p.match(TokenRBrace) {
			break
		}

		p.emitError(CodeParExpectedEnumDelimiter, p.peek().Start, "expected ',' or '}' in enum declaration")
		p.recoverEnumItem()
	}

	if !p.expect(TokenSemicolon, CodeParMissingEnumSemicolon, "missing ';' after enum declaration") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	node := p.enumArena.alloc()
	*node = enumDecl

	stmt := Statement{
		Kind:  NodeEnum,
		Enum:  node,
		Start: startToken.Start,
		End:   p.prev().End,
	}

	return stmt, true
}

// parseEnumValueRaw parses enum member value up to comma or right brace.
func (p *parser) parseEnumValueRaw(stopMask valueStopMask) (string, bool) {
	p.skipTrivia()
	if p.isEOF() {
		p.emitError(CodeParExpectedValueBeforeEOF, p.prev().End, "expected enum item value before end of file")

		return "", false
	}

	if p.isStopToken(p.peek().Kind, stopMask) {
		p.emitError(CodeParExpectedValue, p.peek().Start, "expected enum item value")

		return "", false
	}

	startOffset := p.peek().Start.Offset
	endOffset := p.peek().End.Offset

	for !p.isEOF() {
		token := p.peek()
		if token.Kind == TokenComment || token.Kind == TokenNewline {
			p.advance()

			continue
		}

		if p.isStopToken(token.Kind, stopMask) {
			break
		}

		endOffset = token.End.Offset
		p.advance()
	}

	if startOffset < 0 || endOffset < startOffset {
		p.emitError(CodeParExpectedValue, p.peek().Start, "expected enum item value")

		return "", false
	}

	return p.rawByOffsets(startOffset, endOffset), true
}

// parseAssignment parses scalar and array assignment statements.
func (p *parser) parseAssignment(stopAtBrace bool) (Statement, bool) {
	nameToken := p.consume()
	stopMask := stopSemicolon
	if stopAtBrace {
		stopMask |= stopRBrace
	}

	if p.match(TokenLBracket) {
		if !p.expect(TokenRBracket, CodeParExpectedArrayRightBracket, "expected ']' in array assignment target") {
			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}

		appendMode := false
		switch p.peek().Kind {
		case TokenAssign:
			p.advance()
		case TokenPlusAssign:
			appendMode = true
			p.advance()
		default:
			p.emitError(CodeParExpectedArrayAssignOperator, p.peek().Start, "expected '=' or '+=' after array target")
			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}

		value, ok := p.parseValue(stopMask)
		if !ok {
			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}

		if !p.expect(TokenSemicolon, CodeParMissingArrayAssignSemicolon, "missing ';' after array assignment") {
			p.recoverStatement(stopAtBrace)

			return Statement{}, false
		}

		stmt := Statement{
			Kind:        NodeArrayAssign,
			ArrayAssign: p.arrayAssignNode(nameToken, appendMode, value),
			Start:       nameToken.Start,
			End:         p.prev().End,
		}

		return stmt, true
	}

	if !p.expect(TokenAssign, CodeParExpectedAssign, "expected '=' after assignment target") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	value, ok := p.parseValue(stopMask)
	if !ok {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	if !p.expect(TokenSemicolon, CodeParMissingAssignSemicolon, "missing ';' after assignment") {
		p.recoverStatement(stopAtBrace)

		return Statement{}, false
	}

	stmt := Statement{
		Kind:     NodeProperty,
		Property: p.propertyNode(nameToken, value),
		Start:    nameToken.Start,
		End:      p.prev().End,
	}

	return stmt, true
}
