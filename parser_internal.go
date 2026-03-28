// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"unsafe"
)

// parseClassLikeName parses class-like names and accepts digit-leading forms like `1kHz_*`.
func (p *parser) parseClassLikeName(code Code, message string) (string, bool) {
	p.skipTrivia()

	token := p.peek()
	if token.Kind == TokenIdentifier {
		return p.tokenText(p.consume()), true
	}

	if token.Kind != TokenNumber {
		p.emitError(code, p.tokStart(token), message)

		return "", false
	}

	numberToken := p.consume()
	name := p.tokenText(numberToken)

	// DayZ config corpus contains names like `1kHz_mono_1s_SoundSet`.
	if p.peek().Kind == TokenIdentifier {
		if p.strict {
			p.emitError(
				CodeParStrictDigitLeadingClassName,
				p.tokStart(numberToken),
				"strict mode: class-like name must not start with digit",
			)

			return "", false
		}

		name += p.tokenText(p.consume())

		return name, true
	}

	p.emitError(code, p.tokStart(token), message)

	return "", false
}

// emitError appends parser error diagnostic.
func (p *parser) emitError(code Code, at Position, message string) {
	p.emit(code, SeverityError, at, message)
}

// emitWarning appends parser warning diagnostic.
func (p *parser) emitWarning(code Code, at Position, message string) {
	p.emit(code, SeverityWarning, at, message)
}

// emit appends parser diagnostic with explicit severity.
func (p *parser) emit(code Code, severity Severity, at Position, message string) {
	if at.File == "" {
		at.File = p.filename
	}

	p.diagnostics = append(p.diagnostics, Diagnostic{
		Code:     code,
		Message:  message,
		Severity: severity,
		Start:    at,
		End:      at,
	})
}

// propertyNode allocates property assignment payload in arena.
func (p *parser) propertyNode(nameToken Token, value Value) *PropertyAssign {
	node := p.propertyArena.alloc()
	node.Name = p.tokenText(nameToken)
	node.Value = value

	return node
}

// arrayAssignNode allocates array assignment payload in arena.
func (p *parser) arrayAssignNode(nameToken Token, appendMode bool, value Value) *ArrayAssign {
	node := p.arrayAssignArena.alloc()
	node.Name = p.tokenText(nameToken)
	node.Append = appendMode
	node.Value = value

	return node
}

// tokenText returns token text using captured lexeme or source offsets.
func (p *parser) tokenText(token Token) string {
	if token.End.Offset >= token.Start.Offset &&
		int(token.End.Offset) < len(p.sourceText) {
		return p.sourceText[int(token.Start.Offset) : int(token.End.Offset)+1]
	}

	return token.Kind.String()
}

// tokStart converts compact token start position to public position.
func (p *parser) tokStart(token Token) Position {
	return Position{
		Line:   int(token.Start.Line),
		Column: int(token.Start.Column),
		Offset: int(token.Start.Offset),
	}
}

// tokEnd converts compact token end position to public position.
func (p *parser) tokEnd(token Token) Position {
	return Position{
		Line:   int(token.End.Line),
		Column: int(token.End.Column),
		Offset: int(token.End.Offset),
	}
}

// tokenEquals checks token text without allocating when possible.
func (p *parser) tokenEquals(token Token, literal string) bool {
	if token.End.Offset < token.Start.Offset ||
		int(token.End.Offset) >= len(p.sourceText) {
		return false
	}

	return p.sourceText[int(token.Start.Offset):int(token.End.Offset)+1] == literal
}

// rawByOffsets returns source fragment by inclusive byte offsets.
func (p *parser) rawByOffsets(start int, end int) string {
	if start < 0 || end < start || end >= len(p.sourceText) {
		return ""
	}

	return p.sourceText[start : end+1]
}

// bytesToStringView returns zero-copy string view over source bytes.
// Input bytes must not be mutated while parser uses this view.
func bytesToStringView(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	return unsafe.String(unsafe.SliceData(data), len(data))
}

// prev returns previously consumed token or zero token at parser start.
func (p *parser) prev() Token {
	if p.index == 0 || len(p.tokens) == 0 {
		return Token{}
	}

	return p.tokens[p.index-1]
}

// peek returns current token, or EOF token at stream end.
func (p *parser) peek() Token {
	if p.index >= len(p.tokens) {
		return Token{
			Kind: TokenEOF,
		}
	}

	return p.tokens[p.index]
}

// consume returns current token and advances parser.
func (p *parser) consume() Token {
	token := p.peek()
	p.advance()

	return token
}

// match consumes token kind and returns true on success.
func (p *parser) match(kind TokenKind) bool {
	p.skipTrivia()
	if p.peek().Kind != kind {
		return false
	}

	p.advance()

	return true
}

// advance moves parser index forward by one token.
func (p *parser) advance() {
	if p.index < len(p.tokens) {
		p.index++
	}
}

// isEOF checks current parser token kind.
func (p *parser) isEOF() bool {
	if p.index >= len(p.tokens) {
		return true
	}

	return p.tokens[p.index].Kind == TokenEOF
}

// consumeLeadingComments reads leading comments and skips blank lines before statement.
func (p *parser) consumeLeadingComments() []Comment {
	comments := make([]Comment, 0, 2)
	for !p.isEOF() {
		token := p.peek()
		switch token.Kind {
		case TokenComment:
			comments = append(comments, Comment{
				Text:  p.tokenText(token),
				Start: p.tokStart(token),
				End:   p.tokEnd(token),
			})
			p.advance()
		case TokenNewline:
			p.advance()
		default:
			if len(comments) == 0 {
				return nil
			}

			return comments
		}
	}

	if len(comments) == 0 {
		return nil
	}

	return comments
}

// consumeTrailingComment captures optional inline comment immediately after statement.
// It only attaches comment when the next trivia token is comment (not newline).
func (p *parser) consumeTrailingComment(line int) (Comment, bool) {
	_ = line

	if p.isEOF() {
		return Comment{}, false
	}

	token := p.peek()
	if token.Kind != TokenComment {
		return Comment{}, false
	}

	p.advance()

	return Comment{
		Text:  p.tokenText(token),
		Start: p.tokStart(token),
		End:   p.tokEnd(token),
	}, true
}
