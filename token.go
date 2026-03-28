// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import "github.com/woozymasta/lintkit/lint"

// TokenKind is lexical token type.
type TokenKind int

const (
	// TokenUnknown is fallback token for unknown lexeme.
	TokenUnknown TokenKind = iota

	// TokenEOF is end-of-file marker.
	TokenEOF

	// TokenNewline is line break marker.
	TokenNewline

	// TokenIdentifier is identifier token.
	TokenIdentifier

	// TokenNumber is numeric literal token.
	TokenNumber

	// TokenString is string literal token.
	TokenString

	// TokenComment is line or block comment token.
	TokenComment

	// TokenHash is # token.
	TokenHash

	// TokenTokenPaste is ## token-paste operator.
	TokenTokenPaste

	// TokenAssign is = token.
	TokenAssign

	// TokenPlus is + token.
	TokenPlus

	// TokenMinus is - token.
	TokenMinus

	// TokenPlusAssign is += token.
	TokenPlusAssign

	// TokenColon is : token.
	TokenColon

	// TokenSemicolon is ; token.
	TokenSemicolon

	// TokenComma is , token.
	TokenComma

	// TokenLParen is ( token.
	TokenLParen

	// TokenRParen is ) token.
	TokenRParen

	// TokenLBrace is { token.
	TokenLBrace

	// TokenRBrace is } token.
	TokenRBrace

	// TokenLBracket is [ token.
	TokenLBracket

	// TokenRBracket is ] token.
	TokenRBracket

	// TokenKeywordClass is class keyword token.
	TokenKeywordClass

	// TokenKeywordDelete is delete keyword token.
	TokenKeywordDelete

	// TokenKeywordEnum is enum keyword token.
	TokenKeywordEnum
)

// Token is lexical token with source positions.
type Token struct {
	// Lexeme is exact source fragment.
	Lexeme string `json:"lexeme,omitempty" yaml:"lexeme,omitempty"`

	// Start is token start location.
	Start lint.Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is token end location.
	End lint.Position `json:"end,omitzero" yaml:"end,omitempty"`

	// Kind is token type.
	Kind TokenKind `json:"kind,omitzero" yaml:"kind,omitempty"`
}

// String renders readable token kind name.
func (k TokenKind) String() string {
	switch k {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "Newline"
	case TokenIdentifier:
		return "Identifier"
	case TokenNumber:
		return "Number"
	case TokenString:
		return "String"
	case TokenComment:
		return "Comment"
	case TokenHash:
		return "Hash"
	case TokenTokenPaste:
		return "TokenPaste"
	case TokenAssign:
		return "Assign"
	case TokenPlus:
		return "Plus"
	case TokenMinus:
		return "Minus"
	case TokenPlusAssign:
		return "PlusAssign"
	case TokenColon:
		return "Colon"
	case TokenSemicolon:
		return "Semicolon"
	case TokenComma:
		return "Comma"
	case TokenLParen:
		return "LParen"
	case TokenRParen:
		return "RParen"
	case TokenLBrace:
		return "LBrace"
	case TokenRBrace:
		return "RBrace"
	case TokenLBracket:
		return "LBracket"
	case TokenRBracket:
		return "RBracket"
	case TokenKeywordClass:
		return "KeywordClass"
	case TokenKeywordDelete:
		return "KeywordDelete"
	case TokenKeywordEnum:
		return "KeywordEnum"
	default:
		return "Unknown"
	}
}
