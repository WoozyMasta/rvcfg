// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

// TokenKind is lexical token type.
type TokenKind uint8

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
	// Start is compact token start location.
	Start TokenPosition `json:"start,omitzero" yaml:"start,omitempty"`

	// End is compact token end location.
	End TokenPosition `json:"end,omitzero" yaml:"end,omitempty"`

	// Kind is token type.
	Kind TokenKind `json:"kind,omitzero" yaml:"kind,omitempty"`
}

// TokenPosition is compact token location without duplicated filename.
type TokenPosition struct {
	// Line is 1-based source line.
	Line uint32 `json:"line,omitzero" yaml:"line,omitempty"`

	// Column is 1-based source column.
	Column uint32 `json:"column,omitzero" yaml:"column,omitempty"`

	// Offset is 0-based byte offset.
	Offset uint32 `json:"offset,omitzero" yaml:"offset,omitempty"`
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
