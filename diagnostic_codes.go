// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

const (
	// LintModule is stable lint module namespace for rvcfg diagnostics.
	LintModule = "rvcfg"

	// StageLex marks lexer diagnostics.
	StageLex Stage = "lex"

	// StageParse marks parser diagnostics.
	StageParse Stage = "parse"

	// StagePreprocess marks preprocessor diagnostics.
	StagePreprocess Stage = "preprocess"
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodeLexUnexpectedCharacter reports unknown character tokenization.
	CodeLexUnexpectedCharacter Code = 1001

	// CodeLexUnterminatedString reports non-closed string literal.
	CodeLexUnterminatedString Code = 1002

	// CodeLexUnterminatedBlockComment reports non-closed /* ... */ comment.
	CodeLexUnterminatedBlockComment Code = 1003
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodeParUnexpectedToken reports statement-level token mismatch.
	CodeParUnexpectedToken Code = 2001

	// CodeParExpectedClassName reports missing class name token.
	CodeParExpectedClassName Code = 2002

	// CodeParExpectedClassBodyOrSemicolon reports invalid class header terminator.
	CodeParExpectedClassBodyOrSemicolon Code = 2003

	// CodeParExpectedClassClosingBrace reports missing class body closing brace.
	CodeParExpectedClassClosingBrace Code = 2004

	// CodeParMissingClassSemicolon reports missing semicolon after class declaration.
	CodeParMissingClassSemicolon Code = 2005

	// CodeParExpectedDeleteName reports missing delete declaration name.
	CodeParExpectedDeleteName Code = 2006

	// CodeParMissingDeleteSemicolon reports missing semicolon after delete declaration.
	CodeParMissingDeleteSemicolon Code = 2007

	// CodeParExpectedExternName reports missing extern declaration name.
	CodeParExpectedExternName Code = 2008

	// CodeParMissingExternSemicolon reports missing semicolon after extern declaration.
	CodeParMissingExternSemicolon Code = 2009

	// CodeParExpectedArrayRightBracket reports missing right bracket in array target.
	CodeParExpectedArrayRightBracket Code = 2010

	// CodeParExpectedArrayAssignOperator reports missing = or += after array target.
	CodeParExpectedArrayAssignOperator Code = 2011

	// CodeParMissingArrayAssignSemicolon reports missing semicolon after array assignment.
	CodeParMissingArrayAssignSemicolon Code = 2012

	// CodeParExpectedAssign reports missing = in scalar assignment.
	CodeParExpectedAssign Code = 2013

	// CodeParMissingAssignSemicolon reports missing semicolon after scalar assignment.
	CodeParMissingAssignSemicolon Code = 2014

	// CodeParExpectedValueBeforeEOF reports value parse at EOF.
	CodeParExpectedValueBeforeEOF Code = 2015

	// CodeParExpectedValue reports missing value token.
	CodeParExpectedValue Code = 2016

	// CodeParExpectedScalarValue reports empty scalar fragment.
	CodeParExpectedScalarValue Code = 2017

	// CodeParUnterminatedArrayLiteral reports non-closed array literal.
	CodeParUnterminatedArrayLiteral Code = 2018

	// CodeParExpectedArrayDelimiter reports missing comma or closing brace in array literal.
	CodeParExpectedArrayDelimiter Code = 2019

	// CodeParAutofixClassSemicolon marks parser autofix warning.
	CodeParAutofixClassSemicolon Code = 2020

	// CodeParExpectedEnumBody reports missing enum body opening brace.
	CodeParExpectedEnumBody Code = 2021

	// CodeParExpectedEnumItemName reports missing enum item name.
	CodeParExpectedEnumItemName Code = 2022

	// CodeParExpectedEnumDelimiter reports missing comma or right brace in enum body.
	CodeParExpectedEnumDelimiter Code = 2023

	// CodeParMissingEnumSemicolon reports missing semicolon after enum declaration.
	CodeParMissingEnumSemicolon Code = 2024

	// CodeParStrictDigitLeadingClassName reports strict mode violation for class-like names.
	CodeParStrictDigitLeadingClassName Code = 2025

	// CodeParDerivedNestedClassWithoutBase reports nested class declaration without explicit inheritance in derived class.
	CodeParDerivedNestedClassWithoutBase Code = 2026

	// CodeParRAPUnsupportedScalar reports scalar raw that RAP encoder cannot classify safely.
	CodeParRAPUnsupportedScalar Code = 2027

	// CodeParRAPFloatPrecisionLoss reports notable precision loss during float64->float32 conversion.
	CodeParRAPFloatPrecisionLoss Code = 2028

	// CodeParRAPFloatUnderflowToZero reports non-zero float scalar collapsing to zero in float32.
	CodeParRAPFloatUnderflowToZero Code = 2029

	// CodeParRAPUnsafeStringEscape reports string escape style that is unsafe for BI/CfgConvert.
	CodeParRAPUnsafeStringEscape Code = 2030

	// CodeParRAPExtremeFloatMagnitude reports extreme float value likely to normalize unexpectedly in RAP text round-trip.
	CodeParRAPExtremeFloatMagnitude Code = 2031

	// CodeParRAPFloatOverflowToInf reports float scalar overflowing to Inf during float32 conversion.
	CodeParRAPFloatOverflowToInf Code = 2032

	// CodeParRAPNonFiniteFloat reports non-finite float scalar (NaN/Inf) unsupported for RAP numeric encoding.
	CodeParRAPNonFiniteFloat Code = 2033
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodePPIncludeNotFound reports unresolved include target.
	CodePPIncludeNotFound Code = 3001

	// CodePPUnsupportedIntrinsic reports unsupported __EXEC/__EVAL usage.
	CodePPUnsupportedIntrinsic Code = 3002

	// CodePPMacroExpand reports macro expansion failure at line level.
	CodePPMacroExpand Code = 3003

	// CodePPUnterminatedConditional reports missing #endif.
	CodePPUnterminatedConditional Code = 3004

	// CodePPInvalidIncludeSyntax reports malformed #include argument.
	CodePPInvalidIncludeSyntax Code = 3005

	// CodePPUnexpectedElif reports #elif without matching active conditional frame.
	CodePPUnexpectedElif Code = 3006

	// CodePPUnexpectedElse reports #else without matching active conditional frame.
	CodePPUnexpectedElse Code = 3007

	// CodePPUnexpectedEndif reports #endif without matching active conditional frame.
	CodePPUnexpectedEndif Code = 3008

	// CodePPDirectiveError reports #error directive hit.
	CodePPDirectiveError Code = 3009

	// CodePPUnsupportedDirective reports unsupported unknown preprocessor directive.
	CodePPUnsupportedDirective Code = 3011

	// CodePPMissingMacroName reports missing name in #define.
	CodePPMissingMacroName Code = 3012

	// CodePPInvalidMacroName reports invalid macro name in #define.
	CodePPInvalidMacroName Code = 3013

	// CodePPUnterminatedMacroParams reports unterminated parameter list in function-like #define.
	CodePPUnterminatedMacroParams Code = 3014

	// CodePPMacroRedefined reports macro redefinition warning.
	CodePPMacroRedefined Code = 3015

	// CodePPUnresolvedMacroInvocation reports unresolved macro-like invocation left in output.
	CodePPUnresolvedMacroInvocation Code = 3016

	// CodePPUnsupportedHasInclude reports unsupported __has_include usage in #if/#elif.
	CodePPUnsupportedHasInclude Code = 3017
)
