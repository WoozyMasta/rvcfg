// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"sync"

	"github.com/woozymasta/lintkit/lint"
)

const (
	// StageLex marks lexer diagnostics.
	StageLex lint.Stage = "lex"

	// StageParse marks parser diagnostics.
	StageParse lint.Stage = "parse"

	// StagePreprocess marks preprocessor diagnostics.
	StagePreprocess lint.Stage = "preprocess"
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodeLexUnexpectedCharacter reports unknown character tokenization.
	CodeLexUnexpectedCharacter lint.Code = 1001

	// CodeLexUnterminatedString reports non-closed string literal.
	CodeLexUnterminatedString lint.Code = 1002

	// CodeLexUnterminatedBlockComment reports non-closed /* ... */ comment.
	CodeLexUnterminatedBlockComment lint.Code = 1003
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodeParUnexpectedToken reports statement-level token mismatch.
	CodeParUnexpectedToken lint.Code = 2001

	// CodeParExpectedClassName reports missing class name token.
	CodeParExpectedClassName lint.Code = 2002

	// CodeParExpectedClassBodyOrSemicolon reports invalid class header terminator.
	CodeParExpectedClassBodyOrSemicolon lint.Code = 2003

	// CodeParExpectedClassClosingBrace reports missing class body closing brace.
	CodeParExpectedClassClosingBrace lint.Code = 2004

	// CodeParMissingClassSemicolon reports missing semicolon after class declaration.
	CodeParMissingClassSemicolon lint.Code = 2005

	// CodeParExpectedDeleteName reports missing delete declaration name.
	CodeParExpectedDeleteName lint.Code = 2006

	// CodeParMissingDeleteSemicolon reports missing semicolon after delete declaration.
	CodeParMissingDeleteSemicolon lint.Code = 2007

	// CodeParExpectedExternName reports missing extern declaration name.
	CodeParExpectedExternName lint.Code = 2008

	// CodeParMissingExternSemicolon reports missing semicolon after extern declaration.
	CodeParMissingExternSemicolon lint.Code = 2009

	// CodeParExpectedArrayRightBracket reports missing right bracket in array target.
	CodeParExpectedArrayRightBracket lint.Code = 2010

	// CodeParExpectedArrayAssignOperator reports missing = or += after array target.
	CodeParExpectedArrayAssignOperator lint.Code = 2011

	// CodeParMissingArrayAssignSemicolon reports missing semicolon after array assignment.
	CodeParMissingArrayAssignSemicolon lint.Code = 2012

	// CodeParExpectedAssign reports missing = in scalar assignment.
	CodeParExpectedAssign lint.Code = 2013

	// CodeParMissingAssignSemicolon reports missing semicolon after scalar assignment.
	CodeParMissingAssignSemicolon lint.Code = 2014

	// CodeParExpectedValueBeforeEOF reports value parse at EOF.
	CodeParExpectedValueBeforeEOF lint.Code = 2015

	// CodeParExpectedValue reports missing value token.
	CodeParExpectedValue lint.Code = 2016

	// CodeParExpectedScalarValue reports empty scalar fragment.
	CodeParExpectedScalarValue lint.Code = 2017

	// CodeParUnterminatedArrayLiteral reports non-closed array literal.
	CodeParUnterminatedArrayLiteral lint.Code = 2018

	// CodeParExpectedArrayDelimiter reports missing comma or closing brace in array literal.
	CodeParExpectedArrayDelimiter lint.Code = 2019

	// CodeParAutofixClassSemicolon marks parser autofix warning.
	CodeParAutofixClassSemicolon lint.Code = 2020

	// CodeParExpectedEnumBody reports missing enum body opening brace.
	CodeParExpectedEnumBody lint.Code = 2021

	// CodeParExpectedEnumItemName reports missing enum item name.
	CodeParExpectedEnumItemName lint.Code = 2022

	// CodeParExpectedEnumDelimiter reports missing comma or right brace in enum body.
	CodeParExpectedEnumDelimiter lint.Code = 2023

	// CodeParMissingEnumSemicolon reports missing semicolon after enum declaration.
	CodeParMissingEnumSemicolon lint.Code = 2024

	// CodeParStrictDigitLeadingClassName reports strict mode violation for class-like names.
	CodeParStrictDigitLeadingClassName lint.Code = 2025

	// CodeParDerivedNestedClassWithoutBase reports nested class declaration without explicit inheritance in derived class.
	CodeParDerivedNestedClassWithoutBase lint.Code = 2026
)

//nolint:gosec // Stable lint code tokens can match generic credential heuristics.
const (
	// CodePPIncludeNotFound reports unresolved include target.
	CodePPIncludeNotFound lint.Code = 3001

	// CodePPUnsupportedIntrinsic reports unsupported __EXEC/__EVAL usage.
	CodePPUnsupportedIntrinsic lint.Code = 3002

	// CodePPMacroExpand reports macro expansion failure at line level.
	CodePPMacroExpand lint.Code = 3003

	// CodePPUnterminatedConditional reports missing #endif.
	CodePPUnterminatedConditional lint.Code = 3004

	// CodePPInvalidIncludeSyntax reports malformed #include argument.
	CodePPInvalidIncludeSyntax lint.Code = 3005

	// CodePPUnexpectedElif reports #elif without matching active conditional frame.
	CodePPUnexpectedElif lint.Code = 3006

	// CodePPUnexpectedElse reports #else without matching active conditional frame.
	CodePPUnexpectedElse lint.Code = 3007

	// CodePPUnexpectedEndif reports #endif without matching active conditional frame.
	CodePPUnexpectedEndif lint.Code = 3008

	// CodePPDirectiveError reports #error directive hit.
	CodePPDirectiveError lint.Code = 3009

	// CodePPUnsupportedDirective reports unsupported unknown preprocessor directive.
	CodePPUnsupportedDirective lint.Code = 3011

	// CodePPMissingMacroName reports missing name in #define.
	CodePPMissingMacroName lint.Code = 3012

	// CodePPInvalidMacroName reports invalid macro name in #define.
	CodePPInvalidMacroName lint.Code = 3013

	// CodePPUnterminatedMacroParams reports unterminated parameter list in function-like #define.
	CodePPUnterminatedMacroParams lint.Code = 3014

	// CodePPMacroRedefined reports macro redefinition warning.
	CodePPMacroRedefined lint.Code = 3015

	// CodePPUnresolvedMacroInvocation reports unresolved macro-like invocation left in output.
	CodePPUnresolvedMacroInvocation lint.Code = 3016

	// CodePPUnsupportedHasInclude reports unsupported __has_include usage in #if/#elif.
	CodePPUnsupportedHasInclude lint.Code = 3017
)

var (
	// diagnosticCodeCatalogState stores lazy-initialized code catalog state.
	diagnosticCodeCatalogState struct {
		// once guards one-time catalog construction.
		once sync.Once

		// catalog stores constructed helper.
		catalog lint.CodeCatalog

		// err stores catalog construction error.
		err error
	}
)

// getDiagnosticCodeCatalog returns lazy-initialized code catalog helper.
func getDiagnosticCodeCatalog() (lint.CodeCatalog, error) {
	diagnosticCodeCatalogState.once.Do(func() {
		diagnosticCodeCatalogState.catalog, diagnosticCodeCatalogState.err = lint.NewCodeCatalog(
			lint.CodeCatalogConfig{
				Module:            LintModule,
				CodePrefix:        "CFG",
				ModuleName:        "Real Virtuality Configs",
				ModuleDescription: "Lint rules for Real Virtuality config lexer, parser and preprocessor flows.",
				ScopeDescriptions: map[lint.Stage]string{
					StageLex:        "Lexer diagnostics.",
					StageParse:      "Parser diagnostics.",
					StagePreprocess: "Preprocessor diagnostics.",
				},
			},
			diagnosticCatalog,
		)
	})

	return diagnosticCodeCatalogState.catalog, diagnosticCodeCatalogState.err
}

// DiagnosticRuleSpec converts one diagnostic spec into lint rule metadata.
func DiagnosticRuleSpec(spec lint.CodeSpec) lint.RuleSpec {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return lint.RuleSpec{}
	}

	return catalog.RuleSpec(spec)
}

// LintRuleID returns lint rule ID mapped from stable rvcfg diagnostic code.
func LintRuleID(code lint.Code) string {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return LintModule + ".unknown"
	}

	ruleID, err := catalog.RuleID(code)
	if err != nil {
		return LintModule + ".unknown"
	}

	return ruleID
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []lint.CodeSpec {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return nil
	}

	return catalog.CodeSpecs()
}

// DiagnosticByCode returns diagnostic metadata for code.
func DiagnosticByCode(code lint.Code) (lint.CodeSpec, bool) {
	catalog, err := getDiagnosticCodeCatalog()
	if err != nil {
		return lint.CodeSpec{}, false
	}

	return catalog.ByCode(code)
}
