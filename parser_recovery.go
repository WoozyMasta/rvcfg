package rvcfg

// recoverStatement skips tokens to next statement boundary.
func (p *parser) recoverStatement(stopAtBrace bool) {
	if !p.recovery {
		return
	}

	for !p.isEOF() {
		switch p.peek().Kind {
		case TokenSemicolon:
			p.advance()

			return
		case TokenRBrace:
			if stopAtBrace {
				return
			}
		}

		p.advance()
	}
}

// recoverArrayItem skips invalid array item until comma or closing brace.
func (p *parser) recoverArrayItem() {
	if !p.recovery {
		return
	}

	for !p.isEOF() {
		switch p.peek().Kind {
		case TokenComma, TokenRBrace:
			return
		}

		p.advance()
	}
}

// recoverEnumItem skips invalid enum item until comma or closing brace.
func (p *parser) recoverEnumItem() {
	if !p.recovery {
		return
	}

	for !p.isEOF() {
		switch p.peek().Kind {
		case TokenComma, TokenRBrace:
			return
		}

		p.advance()
	}
}

// skipTrivia skips optional newline/comment tokens.
func (p *parser) skipTrivia() {
	for !p.isEOF() {
		switch p.peek().Kind {
		case TokenComment, TokenNewline:
			p.advance()
		default:
			return
		}
	}
}

// expect consumes required token kind or emits parse diagnostic.
func (p *parser) expect(kind TokenKind, code DiagnosticCode, message string) bool {
	p.skipTrivia()
	if p.peek().Kind != kind {
		p.emitError(code, p.peek().Start, message)

		return false
	}

	p.advance()

	return true
}

// canRecoverImplicitClassSemicolon checks whether current token starts safe boundary.
func (p *parser) canRecoverImplicitClassSemicolon() bool {
	p.skipTrivia()

	switch p.peek().Kind {
	case TokenEOF, TokenRBrace, TokenKeywordClass, TokenKeywordDelete, TokenKeywordEnum, TokenIdentifier:
		return true
	default:
		return false
	}
}
