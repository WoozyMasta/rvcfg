// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

// parseValue parses scalar token sequence or nested array literal.
func (p *parser) parseValue(stopMask valueStopMask) (Value, bool) {
	p.skipTrivia()
	if p.isEOF() {
		p.emitError(CodeParExpectedValueBeforeEOF, p.prev().End, "expected value before end of file")

		return Value{}, false
	}

	if p.peek().Kind == TokenLBrace {
		return p.parseArrayValue()
	}

	if p.isStopToken(p.peek().Kind, stopMask) {
		p.emitError(CodeParExpectedValue, p.peek().Start, "expected value")

		return Value{}, false
	}

	start := p.peek().Start
	startOffset := p.peek().Start.Offset
	endOffset := p.peek().End.Offset
	endPos := p.peek().End

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
		endPos = token.End
		p.advance()
	}

	raw := ""
	if startOffset < 0 || endOffset < startOffset {
		p.emitError(CodeParExpectedScalarValue, start, "expected scalar value")

		return Value{}, false
	}

	if p.captureScalarRaw {
		raw = p.rawByOffsets(startOffset, endOffset)
		if raw == "" {
			p.emitError(CodeParExpectedScalarValue, start, "expected scalar value")

			return Value{}, false
		}
	}

	return Value{
		Kind:  ValueScalar,
		Raw:   raw,
		Start: start,
		End:   endPos,
	}, true
}

// parseArrayValue parses `{...}` with nested arrays and trailing commas.
func (p *parser) parseArrayValue() (Value, bool) {
	startToken := p.consume()
	capHint := p.estimateArrayElementCap()
	value := Value{
		Kind:  ValueArray,
		Start: startToken.Start,
	}
	if capHint > 0 {
		value.Elements = make([]Value, 0, capHint)
	}

	for {
		p.skipTrivia()
		if p.isEOF() {
			p.emitError(CodeParUnterminatedArrayLiteral, startToken.Start, "unterminated array literal")

			return Value{}, false
		}

		if p.match(TokenRBrace) {
			value.End = p.prev().End

			return value, true
		}

		item, ok := p.parseValue(stopComma | stopRBrace)
		if !ok {
			p.recoverArrayItem()
		} else {
			if value.Elements == nil {
				value.Elements = make([]Value, 0, 4)
			}

			value.Elements = append(value.Elements, item)
		}

		p.skipTrivia()
		if p.match(TokenComma) {
			p.skipTrivia()
			if p.peek().Kind == TokenRBrace {
				// trailing comma before closing brace is valid in game configs.
				continue
			}

			continue
		}

		if p.match(TokenRBrace) {
			value.End = p.prev().End

			return value, true
		}

		p.emitError(CodeParExpectedArrayDelimiter, p.peek().Start, "expected ',' or '}' in array literal")
		p.recoverArrayItem()
	}
}

// estimateArrayElementCap estimates element count until matching array closing brace.
func (p *parser) estimateArrayElementCap() int {
	depth := 1
	count := 0
	hasItem := false

	for idx := p.index; idx < len(p.tokens); idx++ {
		kind := p.tokens[idx].Kind
		if kind == TokenComment || kind == TokenNewline {
			continue
		}

		switch kind {
		case TokenLBrace:
			if depth == 1 {
				hasItem = true
			}

			depth++
		case TokenRBrace:
			depth--
			if depth == 0 {
				if hasItem {
					count++
				}

				return count
			}
		case TokenComma:
			if depth == 1 && hasItem {
				count++
				hasItem = false
			}
		default:
			if depth == 1 {
				hasItem = true
			}
		}
	}

	return 0
}

// parseBaseExpression captures tokens between `:` and class body/semicolon.
func (p *parser) parseBaseExpression() string {
	startOffset := -1
	endOffset := -1

	for !p.isEOF() {
		token := p.peek()
		if token.Kind == TokenLBrace || token.Kind == TokenSemicolon {
			break
		}

		if token.Kind == TokenComment || token.Kind == TokenNewline {
			p.advance()

			continue
		}

		if startOffset < 0 {
			startOffset = token.Start.Offset
		}

		endOffset = token.End.Offset
		p.advance()
	}

	if startOffset < 0 || endOffset < startOffset {
		return ""
	}

	return p.rawByOffsets(startOffset, endOffset)
}

// isStopToken checks whether token kind belongs to parser stop-mask.
func (p *parser) isStopToken(kind TokenKind, mask valueStopMask) bool {
	switch kind {
	case TokenSemicolon:
		return mask&stopSemicolon != 0
	case TokenRBrace:
		return mask&stopRBrace != 0
	case TokenComma:
		return mask&stopComma != 0
	default:
		return false
	}
}
