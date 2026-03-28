// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

var (
	keywordClass  = []byte("class")
	keywordDelete = []byte("delete")
	keywordEnum   = []byte("enum")
)

// LexOptions configures lexer output volume and allocation profile.
type LexOptions struct {
	// CaptureLexeme is reserved for compatibility and currently has no effect.
	CaptureLexeme bool `json:"capture_lexeme,omitempty" yaml:"capture_lexeme,omitempty"`

	// EmitComments emits comment tokens. Disabled by default.
	EmitComments bool `json:"emit_comments,omitempty" yaml:"emit_comments,omitempty"`

	// EmitNewlines emits newline tokens. Disabled by default.
	EmitNewlines bool `json:"emit_newlines,omitempty" yaml:"emit_newlines,omitempty"`
}

// LexFile scans file and returns token stream with diagnostics.
func LexFile(path string) ([]Token, []Diagnostic, error) {
	return LexFileWithOptions(path, LexOptions{})
}

// LexFileWithOptions scans file and returns token stream with diagnostics.
func LexFileWithOptions(path string, opts LexOptions) ([]Token, []Diagnostic, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read source file %q: %w", path, err)
	}

	return LexBytesWithOptions(path, data, opts)
}

// LexBytes scans raw source bytes and returns token stream with diagnostics.
func LexBytes(filename string, data []byte) ([]Token, []Diagnostic, error) {
	return LexBytesWithOptions(filename, data, LexOptions{})
}

// LexBytesWithOptions scans raw source bytes with configurable output options.
func LexBytesWithOptions(filename string, data []byte, opts LexOptions) ([]Token, []Diagnostic, error) {
	l := newLexer(filename, data, opts)
	tokens, diagnostics := l.scan()

	for _, d := range diagnostics {
		if d.Severity == SeverityError {
			return tokens, diagnostics, ErrLex
		}
	}

	return tokens, diagnostics, nil
}

// lexer holds incremental scan state.
type lexer struct {
	filename string
	data     []byte
	opts     LexOptions
	index    int
	line     int
	column   int
}

// newLexer creates scanner with initial source location.
func newLexer(filename string, data []byte, opts LexOptions) *lexer {
	return &lexer{
		filename: filename,
		data:     data,
		opts:     opts,
		line:     1,
		column:   1,
	}
}

// scan tokenizes all input and appends EOF token.
func (l *lexer) scan() ([]Token, []Diagnostic) {
	tokens := make([]Token, 0, estimateTokenCap(len(l.data), l.opts))
	diagnostics := make([]Diagnostic, 0)

	for !l.eof() {
		ch := l.peek()

		if ch == '\r' || ch == '\n' {
			if l.opts.EmitNewlines {
				token := l.scanNewline()
				tokens = append(tokens, token)
			} else {
				l.consumeNewline()
			}

			continue
		}

		if isHorizontalWhitespace(ch) {
			l.advance()

			continue
		}

		start := l.pos()

		if isIdentifierStart(ch) {
			token := l.scanIdentifier()
			tokens = append(tokens, token)

			continue
		}

		if isDigit(ch) {
			token := l.scanNumber()
			tokens = append(tokens, token)

			continue
		}

		if ch == '"' {
			token, diag := l.scanString()
			tokens = append(tokens, token)

			if diag != nil {
				diagnostics = append(diagnostics, *diag)
			}

			continue
		}

		if ch == '/' && l.matchNext('/') {
			if l.opts.EmitComments {
				token := l.scanLineComment()
				tokens = append(tokens, token)
			} else {
				l.skipLineComment()
			}

			continue
		}

		if ch == '/' && l.matchNext('*') {
			token, diag := l.scanBlockComment()
			if l.opts.EmitComments {
				tokens = append(tokens, token)
			}

			if diag != nil {
				diagnostics = append(diagnostics, *diag)
			}

			continue
		}

		switch ch {
		case '#':
			startIdx := l.index
			l.advance()

			kind := TokenHash

			if l.peek() == '#' {
				l.advance()
				kind = TokenTokenPaste
			}

			tokens = append(tokens, l.makeToken(kind, start, startIdx))

		case '=':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenAssign, start, startIdx))

		case '+':
			startIdx := l.index
			l.advance()

			kind := TokenPlus

			if l.peek() == '=' {
				l.advance()
				kind = TokenPlusAssign
			}

			tokens = append(tokens, l.makeToken(kind, start, startIdx))

		case '-':
			if l.matchSignedNumberStart() {
				token := l.scanSignedNumber()
				tokens = append(tokens, token)

				continue
			}

			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenMinus, start, startIdx))

		case ':':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenColon, start, startIdx))

		case ';':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenSemicolon, start, startIdx))

		case ',':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenComma, start, startIdx))

		case '(':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenLParen, start, startIdx))

		case ')':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenRParen, start, startIdx))

		case '{':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenLBrace, start, startIdx))

		case '}':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenRBrace, start, startIdx))

		case '[':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenLBracket, start, startIdx))

		case ']':
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenRBracket, start, startIdx))

		default:
			startIdx := l.index
			l.advance()
			tokens = append(tokens, l.makeToken(TokenUnknown, start, startIdx))
			diagnostics = append(diagnostics, Diagnostic{
				Code:     CodeLexUnexpectedCharacter,
				Message:  fmt.Sprintf("unexpected character %q", ch),
				Severity: SeverityWarning,
				Start:    l.publicPos(start),
				End:      l.publicPos(l.pos()),
			})
		}
	}

	eofPos := l.pos()
	tokens = append(tokens, Token{
		Kind:  TokenEOF,
		Start: eofPos,
		End:   eofPos,
	})

	return tokens, diagnostics
}

// estimateTokenCap returns initial token slice capacity for lex workload.
func estimateTokenCap(dataLen int, opts LexOptions) int {
	if dataLen <= 0 {
		return 0
	}

	// Default parse mode skips comments/newlines and still produces a dense token
	// stream for config grammar. Use a higher baseline to avoid costly growslice
	// on large files.
	capHint := dataLen / 2
	if opts.EmitComments || opts.EmitNewlines {
		capHint = (dataLen * 3) / 4
	}

	const minCap = 64
	if capHint < minCap {
		return minCap
	}

	return capHint
}

// matchSignedNumberStart checks whether current '-' starts numeric literal.
func (l *lexer) matchSignedNumberStart() bool {
	if l.peek() != '-' {
		return false
	}

	if l.index+1 >= len(l.data) {
		return false
	}

	next := l.data[l.index+1]
	if isDigit(next) {
		return true
	}

	if next == '.' && l.index+2 < len(l.data) && isDigit(l.data[l.index+2]) {
		return true
	}

	return false
}

// scanSignedNumber parses numbers with leading '-' sign.
func (l *lexer) scanSignedNumber() Token {
	start := l.pos()
	startIdx := l.index

	// consume sign
	l.advance()

	if !l.eof() && l.peek() == '.' {
		l.advance()
	}

	for !l.eof() && isDigit(l.peek()) {
		l.advance()
	}

	if !l.eof() && l.peek() == '.' {
		l.advance()

		for !l.eof() && isDigit(l.peek()) {
			l.advance()
		}
	}

	return l.makeToken(TokenNumber, start, startIdx)
}

// scanIdentifier parses keyword/identifier token.
func (l *lexer) scanIdentifier() Token {
	start := l.pos()
	startIdx := l.index

	for !l.eof() && isIdentifierPart(l.peek()) {
		l.advance()
	}

	kind, ok := l.keywordKindByRange(startIdx, l.index)
	if !ok {
		kind = TokenIdentifier
	}

	return l.makeToken(kind, start, startIdx)
}

// scanNumber parses integer and float literals.
func (l *lexer) scanNumber() Token {
	start := l.pos()
	startIdx := l.index

	for !l.eof() && isDigit(l.peek()) {
		l.advance()
	}

	if !l.eof() && l.peek() == '.' {
		l.advance()

		for !l.eof() && isDigit(l.peek()) {
			l.advance()
		}
	}

	return l.makeToken(TokenNumber, start, startIdx)
}

// scanString reads quoted string without C-style unescape semantics.
func (l *lexer) scanString() (Token, *Diagnostic) {
	start := l.pos()
	startIdx := l.index

	// consume opening quote
	l.advance()

	for !l.eof() {
		if l.peek() == '"' {
			l.advance()

			return l.makeToken(TokenString, start, startIdx), nil
		}

		if l.peek() == '\r' || l.peek() == '\n' {
			diag := Diagnostic{
				Code:     CodeLexUnterminatedString,
				Message:  "unterminated string literal",
				Severity: SeverityError,
				Start:    l.publicPos(start),
				End:      l.publicPos(l.pos()),
			}

			return l.makeToken(TokenString, start, startIdx), &diag
		}

		l.advance()
	}

	diag := Diagnostic{
		Code:     CodeLexUnterminatedString,
		Message:  "unterminated string literal",
		Severity: SeverityError,
		Start:    l.publicPos(start),
		End:      l.publicPos(l.pos()),
	}

	return l.makeToken(TokenString, start, startIdx), &diag
}

// scanLineComment consumes // comment until line end.
func (l *lexer) scanLineComment() Token {
	start := l.pos()
	startIdx := l.index

	l.advance()
	l.advance()

	for !l.eof() && l.peek() != '\n' && l.peek() != '\r' {
		l.advance()
	}

	return l.makeToken(TokenComment, start, startIdx)
}

// skipLineComment consumes // comment without emitting token.
func (l *lexer) skipLineComment() {
	l.advance()
	l.advance()

	for !l.eof() && l.peek() != '\n' && l.peek() != '\r' {
		l.advance()
	}
}

// scanBlockComment consumes /* ... */ and reports unclosed comment.
func (l *lexer) scanBlockComment() (Token, *Diagnostic) {
	start := l.pos()
	startIdx := l.index

	l.advance()
	l.advance()

	for !l.eof() {
		if l.peek() == '*' && l.matchNext('/') {
			l.advance()
			l.advance()

			return l.makeToken(TokenComment, start, startIdx), nil
		}

		l.advance()
	}

	diag := Diagnostic{
		Code:     CodeLexUnterminatedBlockComment,
		Message:  "unterminated block comment",
		Severity: SeverityError,
		Start:    l.publicPos(start),
		End:      l.publicPos(l.pos()),
	}

	return l.makeToken(TokenComment, start, startIdx), &diag
}

// scanNewline consumes \n, \r, and \r\n as single newline token.
func (l *lexer) scanNewline() Token {
	start := l.pos()
	startIdx := l.index

	if l.peek() == '\r' {
		l.advance()

		if !l.eof() && l.peek() == '\n' {
			l.advance()

			return l.makeToken(TokenNewline, start, startIdx)
		}

		return l.makeToken(TokenNewline, start, startIdx)
	}

	l.advance()

	return l.makeToken(TokenNewline, start, startIdx)
}

// consumeNewline consumes \n, \r, and \r\n without token emission.
func (l *lexer) consumeNewline() {
	if l.peek() == '\r' {
		l.advance()
		if !l.eof() && l.peek() == '\n' {
			l.advance()
		}

		return
	}

	l.advance()
}

// makeToken builds token from start and current scanner position.
func (l *lexer) makeToken(kind TokenKind, start TokenPosition, startIdx int) Token {
	_ = startIdx

	end := l.pos()
	if end.Offset > start.Offset {
		end.Offset--

		if end.Column > 1 {
			end.Column--
		}
	}

	return Token{
		Kind:  kind,
		Start: start,
		End:   end,
	}
}

// keywordKindByRange resolves keyword kind without allocating identifier string.
func (l *lexer) keywordKindByRange(start int, end int) (TokenKind, bool) {
	slice := l.data[start:end]

	switch len(slice) {
	case 4:
		if bytes.Equal(slice, keywordEnum) {
			return TokenKeywordEnum, true
		}
	case 5:
		if bytes.Equal(slice, keywordClass) {
			return TokenKeywordClass, true
		}
	case 6:
		if bytes.Equal(slice, keywordDelete) {
			return TokenKeywordDelete, true
		}
	}

	return TokenUnknown, false
}

// pos returns current scanner source position.
func (l *lexer) pos() TokenPosition {
	return TokenPosition{
		Line:   compactPosValue(l.line),
		Column: compactPosValue(l.column),
		Offset: compactPosValue(l.index),
	}
}

// compactPosValue converts int position component to uint32 with saturation.
func compactPosValue(v int) uint32 {
	if v <= 0 {
		return 0
	}

	u := uint64(v)
	maxUint32 := uint64(^uint32(0))
	if u > maxUint32 {
		return ^uint32(0)
	}

	return uint32(u)
}

// publicPos converts compact token position to public position with filename.
func (l *lexer) publicPos(pos TokenPosition) Position {
	return Position{
		File:   l.filename,
		Line:   int(pos.Line),
		Column: int(pos.Column),
		Offset: int(pos.Offset),
	}
}

// peek returns current byte without consuming it.
func (l *lexer) peek() byte {
	if l.eof() {
		return 0
	}

	return l.data[l.index]
}

// matchNext checks next byte, without consuming current byte.
func (l *lexer) matchNext(expected byte) bool {
	if l.index+1 >= len(l.data) {
		return false
	}

	return l.data[l.index+1] == expected
}

// advance consumes current byte and updates source location.
func (l *lexer) advance() {
	if l.eof() {
		return
	}

	ch := l.data[l.index]
	l.index++

	if ch == '\n' {
		l.line++
		l.column = 1

		return
	}

	l.column++
}

// eof indicates scanner end of input.
func (l *lexer) eof() bool {
	return l.index >= len(l.data)
}

// isIdentifierStart checks identifier first byte.
func isIdentifierStart(ch byte) bool {
	return ch == '$' || ch == '_' || unicode.IsLetter(rune(ch))
}

// isIdentifierPart checks identifier continuation byte.
func isIdentifierPart(ch byte) bool {
	return ch == '$' || ch == '_' || unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch))
}

// isDigit checks decimal digit byte.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isHorizontalWhitespace checks non-newline whitespace.
func isHorizontalWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\f' || ch == '\v'
}
