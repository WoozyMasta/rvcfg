// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

// DiagnosticStage groups diagnostics by pipeline stage.
type DiagnosticStage string

const (
	// StageLex marks lexer diagnostics.
	StageLex DiagnosticStage = "lex"

	// StageParse marks parser diagnostics.
	StageParse DiagnosticStage = "parse"

	// StagePreprocess marks preprocessor diagnostics.
	StagePreprocess DiagnosticStage = "preprocess"
)

const (
	// CodeLexUnexpectedCharacter reports unknown character tokenization.
	CodeLexUnexpectedCharacter DiagnosticCode = "LEX001"

	// CodeLexUnterminatedString reports non-closed string literal.
	CodeLexUnterminatedString DiagnosticCode = "LEX002"

	// CodeLexUnterminatedBlockComment reports non-closed /* ... */ comment.
	CodeLexUnterminatedBlockComment DiagnosticCode = "LEX003"
)

const (
	// CodeParUnexpectedToken reports statement-level token mismatch.
	CodeParUnexpectedToken DiagnosticCode = "PAR001"

	// CodeParExpectedClassName reports missing class name token.
	CodeParExpectedClassName DiagnosticCode = "PAR002"

	// CodeParExpectedClassBodyOrSemicolon reports invalid class header terminator.
	CodeParExpectedClassBodyOrSemicolon DiagnosticCode = "PAR003"

	// CodeParExpectedClassClosingBrace reports missing class body closing brace.
	CodeParExpectedClassClosingBrace DiagnosticCode = "PAR004"

	// CodeParMissingClassSemicolon reports missing semicolon after class declaration.
	CodeParMissingClassSemicolon DiagnosticCode = "PAR005"

	// CodeParExpectedDeleteName reports missing delete declaration name.
	CodeParExpectedDeleteName DiagnosticCode = "PAR006"

	// CodeParMissingDeleteSemicolon reports missing semicolon after delete declaration.
	CodeParMissingDeleteSemicolon DiagnosticCode = "PAR007"

	// CodeParExpectedExternName reports missing extern declaration name.
	CodeParExpectedExternName DiagnosticCode = "PAR008"

	// CodeParMissingExternSemicolon reports missing semicolon after extern declaration.
	CodeParMissingExternSemicolon DiagnosticCode = "PAR009"

	// CodeParExpectedArrayRightBracket reports missing right bracket in array target.
	CodeParExpectedArrayRightBracket DiagnosticCode = "PAR010"

	// CodeParExpectedArrayAssignOperator reports missing = or += after array target.
	CodeParExpectedArrayAssignOperator DiagnosticCode = "PAR011"

	// CodeParMissingArrayAssignSemicolon reports missing semicolon after array assignment.
	CodeParMissingArrayAssignSemicolon DiagnosticCode = "PAR012"

	// CodeParExpectedAssign reports missing = in scalar assignment.
	CodeParExpectedAssign DiagnosticCode = "PAR013"

	// CodeParMissingAssignSemicolon reports missing semicolon after scalar assignment.
	CodeParMissingAssignSemicolon DiagnosticCode = "PAR014"

	// CodeParExpectedValueBeforeEOF reports value parse at EOF.
	CodeParExpectedValueBeforeEOF DiagnosticCode = "PAR015"

	// CodeParExpectedValue reports missing value token.
	CodeParExpectedValue DiagnosticCode = "PAR016"

	// CodeParExpectedScalarValue reports empty scalar fragment.
	CodeParExpectedScalarValue DiagnosticCode = "PAR017"

	// CodeParUnterminatedArrayLiteral reports non-closed array literal.
	CodeParUnterminatedArrayLiteral DiagnosticCode = "PAR018"

	// CodeParExpectedArrayDelimiter reports missing comma or closing brace in array literal.
	CodeParExpectedArrayDelimiter DiagnosticCode = "PAR019"

	// CodeParStrictDigitLeadingClassName reports strict mode violation for class-like names.
	CodeParStrictDigitLeadingClassName DiagnosticCode = "PAR025"

	// CodeParAutofixClassSemicolon marks parser autofix warning.
	CodeParAutofixClassSemicolon DiagnosticCode = "PAR020"

	// CodeParExpectedEnumBody reports missing enum body opening brace.
	CodeParExpectedEnumBody DiagnosticCode = "PAR021"

	// CodeParExpectedEnumItemName reports missing enum item name.
	CodeParExpectedEnumItemName DiagnosticCode = "PAR022"

	// CodeParExpectedEnumDelimiter reports missing comma or right brace in enum body.
	CodeParExpectedEnumDelimiter DiagnosticCode = "PAR023"

	// CodeParMissingEnumSemicolon reports missing semicolon after enum declaration.
	CodeParMissingEnumSemicolon DiagnosticCode = "PAR024"
)

const (
	// CodePPIncludeNotFound reports unresolved include target.
	CodePPIncludeNotFound DiagnosticCode = "PP001"

	// CodePPUnsupportedIntrinsic reports unsupported __EXEC/__EVAL usage.
	CodePPUnsupportedIntrinsic DiagnosticCode = "PP002"

	// CodePPMacroExpand reports macro expansion failure at line level.
	CodePPMacroExpand DiagnosticCode = "PP003"

	// CodePPUnterminatedConditional reports missing #endif.
	CodePPUnterminatedConditional DiagnosticCode = "PP004"

	// CodePPInvalidIncludeSyntax reports malformed #include argument.
	CodePPInvalidIncludeSyntax DiagnosticCode = "PP005"

	// CodePPUnexpectedElif reports #elif without matching active conditional frame.
	CodePPUnexpectedElif DiagnosticCode = "PP006"

	// CodePPUnexpectedElse reports #else without matching active conditional frame.
	CodePPUnexpectedElse DiagnosticCode = "PP007"

	// CodePPUnexpectedEndif reports #endif without matching active conditional frame.
	CodePPUnexpectedEndif DiagnosticCode = "PP008"

	// CodePPDirectiveError reports #error directive hit.
	CodePPDirectiveError DiagnosticCode = "PP009"

	// CodePPUnsupportedDirective reports unsupported unknown preprocessor directive.
	CodePPUnsupportedDirective DiagnosticCode = "PP011"

	// CodePPMissingMacroName reports missing name in #define.
	CodePPMissingMacroName DiagnosticCode = "PP012"

	// CodePPInvalidMacroName reports invalid macro name in #define.
	CodePPInvalidMacroName DiagnosticCode = "PP013"

	// CodePPUnterminatedMacroParams reports unterminated parameter list in function-like #define.
	CodePPUnterminatedMacroParams DiagnosticCode = "PP014"

	// CodePPMacroRedefined reports macro redefinition warning.
	CodePPMacroRedefined DiagnosticCode = "PP015"

	// CodePPUnresolvedMacroInvocation reports unresolved macro-like invocation left in output.
	CodePPUnresolvedMacroInvocation DiagnosticCode = "PP016"

	// CodePPUnsupportedHasInclude reports unsupported __has_include usage in #if/#elif.
	CodePPUnsupportedHasInclude DiagnosticCode = "PP017"
)

// DiagnosticSpec describes one stable diagnostic code.
type DiagnosticSpec struct {
	// Code is stable machine-readable diagnostic identifier.
	Code DiagnosticCode `json:"code,omitempty" yaml:"code,omitempty"`

	// Stage is pipeline stage where diagnostic can be produced.
	Stage DiagnosticStage `json:"stage,omitempty" yaml:"stage,omitempty"`

	// Severity is default diagnostic severity.
	Severity Severity `json:"severity,omitempty" yaml:"severity,omitempty"`

	// Summary is short human-readable description.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
}

var diagnosticCatalog = []DiagnosticSpec{
	{Code: CodeLexUnexpectedCharacter, Stage: StageLex, Severity: SeverityWarning, Summary: "unexpected character"},
	{Code: CodeLexUnterminatedString, Stage: StageLex, Severity: SeverityError, Summary: "unterminated string literal"},
	{Code: CodeLexUnterminatedBlockComment, Stage: StageLex, Severity: SeverityError, Summary: "unterminated block comment"},
	{Code: CodeParUnexpectedToken, Stage: StageParse, Severity: SeverityError, Summary: "unexpected token"},
	{Code: CodeParExpectedClassName, Stage: StageParse, Severity: SeverityError, Summary: "expected class name"},
	{
		Code:     CodeParExpectedClassBodyOrSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected class body or semicolon",
	},
	{
		Code:     CodeParExpectedClassClosingBrace,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected closing brace for class body",
	},
	{
		Code:     CodeParMissingClassSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after class declaration",
	},
	{Code: CodeParExpectedDeleteName, Stage: StageParse, Severity: SeverityError, Summary: "expected name after delete"},
	{
		Code:     CodeParMissingDeleteSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after delete declaration",
	},
	{
		Code:     CodeParExpectedExternName,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected extern declaration name",
	},
	{
		Code:     CodeParMissingExternSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after extern declaration",
	},
	{
		Code:     CodeParExpectedArrayRightBracket,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected right bracket in array assignment target",
	},
	{
		Code:     CodeParExpectedArrayAssignOperator,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected assignment operator for array assignment",
	},
	{
		Code:     CodeParMissingArrayAssignSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after array assignment",
	},
	{
		Code:     CodeParExpectedAssign,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected assignment operator",
	},
	{
		Code:     CodeParMissingAssignSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after assignment",
	},
	{
		Code:     CodeParExpectedValueBeforeEOF,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected value before end of file",
	},
	{Code: CodeParExpectedValue, Stage: StageParse, Severity: SeverityError, Summary: "expected value"},
	{Code: CodeParExpectedScalarValue, Stage: StageParse, Severity: SeverityError, Summary: "expected scalar value"},
	{
		Code:     CodeParUnterminatedArrayLiteral,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "unterminated array literal",
	},
	{
		Code:     CodeParExpectedArrayDelimiter,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected comma or right brace in array literal",
	},
	{
		Code:     CodeParStrictDigitLeadingClassName,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "strict mode rejects class-like names starting with digit",
	},
	{
		Code:     CodeParAutofixClassSemicolon,
		Stage:    StageParse,
		Severity: SeverityWarning,
		Summary:  "autofix inserted missing class semicolon",
	},
	{
		Code:     CodeParExpectedEnumBody,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected enum body",
	},
	{
		Code:     CodeParExpectedEnumItemName,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected enum item name",
	},
	{
		Code:     CodeParExpectedEnumDelimiter,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "expected enum item delimiter",
	},
	{
		Code:     CodeParMissingEnumSemicolon,
		Stage:    StageParse,
		Severity: SeverityError,
		Summary:  "missing semicolon after enum declaration",
	},
	{
		Code:     CodePPIncludeNotFound,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "include target not found or unreadable",
	},
	{
		Code:     CodePPUnsupportedIntrinsic,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unsupported config intrinsic",
	},
	{
		Code:     CodePPMacroExpand,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "macro expansion failure",
	},
	{
		Code:     CodePPUnterminatedConditional,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unterminated conditional block",
	},
	{
		Code:     CodePPInvalidIncludeSyntax,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "invalid include syntax",
	},
	{
		Code:     CodePPUnexpectedElif,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unexpected elif directive",
	},
	{
		Code:     CodePPUnexpectedElse,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unexpected else directive",
	},
	{
		Code:     CodePPUnexpectedEndif,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unexpected endif directive",
	},
	{
		Code:     CodePPDirectiveError,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "error directive triggered",
	},
	{
		Code:     CodePPUnsupportedDirective,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unsupported directive",
	},
	{
		Code:     CodePPMissingMacroName,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "missing macro name in define",
	},
	{
		Code:     CodePPInvalidMacroName,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "invalid macro name in define",
	},
	{
		Code:     CodePPUnterminatedMacroParams,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unterminated macro parameter list",
	},
	{
		Code:     CodePPMacroRedefined,
		Stage:    StagePreprocess,
		Severity: SeverityWarning,
		Summary:  "macro redefinition",
	},
	{
		Code:     CodePPUnresolvedMacroInvocation,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unresolved macro-like invocation",
	},
	{
		Code:     CodePPUnsupportedHasInclude,
		Stage:    StagePreprocess,
		Severity: SeverityError,
		Summary:  "unsupported __has_include in conditional directive",
	},
}

var diagnosticCatalogByCode = buildDiagnosticCatalogByCode()

// buildDiagnosticCatalogByCode builds O(1) lookup map from stable catalog slice.
func buildDiagnosticCatalogByCode() map[DiagnosticCode]DiagnosticSpec {
	out := make(map[DiagnosticCode]DiagnosticSpec, len(diagnosticCatalog))

	for _, spec := range diagnosticCatalog {
		out[spec.Code] = spec
	}

	return out
}

// DiagnosticCatalog returns stable diagnostics metadata list.
func DiagnosticCatalog() []DiagnosticSpec {
	out := make([]DiagnosticSpec, len(diagnosticCatalog))
	copy(out, diagnosticCatalog)

	return out
}

// DiagnosticByCode returns diagnostic metadata for code.
func DiagnosticByCode(code DiagnosticCode) (DiagnosticSpec, bool) {
	spec, ok := diagnosticCatalogByCode[code]

	return spec, ok
}
