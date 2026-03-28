// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import "github.com/woozymasta/lintkit/lint"

// withDescription attaches optional documentation text to one catalog spec.
func withDescription(spec lint.CodeSpec, description string) lint.CodeSpec {
	spec.Description = description
	return spec
}

// diagnosticCatalog stores stable diagnostics metadata table.
var diagnosticCatalog = []lint.CodeSpec{
	lint.WarningCodeSpec(
		CodeLexUnexpectedCharacter,
		StageLex,
		"unexpected character",
	),
	lint.ErrorCodeSpec(
		CodeLexUnterminatedString,
		StageLex,
		"unterminated string literal",
	),
	lint.ErrorCodeSpec(
		CodeLexUnterminatedBlockComment,
		StageLex,
		"unterminated block comment",
	),
	lint.ErrorCodeSpec(
		CodeParUnexpectedToken,
		StageParse,
		"unexpected token",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedClassName,
		StageParse,
		"expected class name",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedClassBodyOrSemicolon,
		StageParse,
		"expected class body or semicolon",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedClassClosingBrace,
		StageParse,
		"expected closing brace for class body",
	),
	lint.ErrorCodeSpec(
		CodeParMissingClassSemicolon,
		StageParse,
		"missing semicolon after class declaration",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedDeleteName,
		StageParse,
		"expected name after `delete`",
	),
	lint.ErrorCodeSpec(
		CodeParMissingDeleteSemicolon,
		StageParse,
		"missing semicolon after delete declaration",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedExternName,
		StageParse,
		"expected `extern` declaration name",
	),
	lint.ErrorCodeSpec(
		CodeParMissingExternSemicolon,
		StageParse,
		"missing semicolon after extern declaration",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedArrayRightBracket,
		StageParse,
		"expected right bracket in array assignment target",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedArrayAssignOperator,
		StageParse,
		"expected assignment operator in array assignment",
	),
	lint.ErrorCodeSpec(
		CodeParMissingArrayAssignSemicolon,
		StageParse,
		"missing semicolon after array assignment",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedAssign,
		StageParse,
		"expected assignment operator",
	),
	lint.ErrorCodeSpec(
		CodeParMissingAssignSemicolon,
		StageParse,
		"missing semicolon after assignment",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedValueBeforeEOF,
		StageParse,
		"expected value before end of file",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedValue,
		StageParse,
		"expected value",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedScalarValue,
		StageParse,
		"expected scalar value",
	),
	lint.ErrorCodeSpec(
		CodeParUnterminatedArrayLiteral,
		StageParse,
		"unterminated array literal",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedArrayDelimiter,
		StageParse,
		"expected comma or right brace in array literal",
	),
	lint.ErrorCodeSpec(
		CodeParStrictDigitLeadingClassName,
		StageParse,
		"class-like names must not start with digit in strict mode",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParAutofixClassSemicolon,
			StageParse,
			"autofix inserted missing class semicolon",
		),
		"Parser recovered by inserting a missing semicolon after class "+
			"declaration. Keep semicolons explicit to avoid parser recovery "+
			"differences between tools.",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedEnumBody,
		StageParse,
		"expected enum body",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedEnumItemName,
		StageParse,
		"expected enum item name",
	),
	lint.ErrorCodeSpec(
		CodeParExpectedEnumDelimiter,
		StageParse,
		"expected enum item delimiter",
	),
	lint.ErrorCodeSpec(
		CodeParMissingEnumSemicolon,
		StageParse,
		"missing semicolon after enum declaration",
	),
	withDescription(
		lint.InfoCodeSpec(
			CodeParDerivedNestedClassWithoutBase,
			StageParse,
			"nested class in derived class has no explicit base",
		),
		"This can replace inherited subtree instead of extending it. Add explicit "+
			"inheritance to make merge behavior predictable.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPIncludeNotFound,
			StagePreprocess,
			"include target not found or unreadable",
		),
		"Check include path, include roots, and file access permissions.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnsupportedIntrinsic,
			StagePreprocess,
			"unsupported config intrinsic",
		),
		"Current preprocessor mode does not support this intrinsic "+
			"(for example `__EXEC` or `__EVAL`). Remove it or preprocess with "+
			"toolchain "+
			"that supports this intrinsic.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPMacroExpand,
			StagePreprocess,
			"macro expansion failure",
		),
		"Usually caused by invalid invocation, argument count "+
			"mismatch, or recursive expansion limit. Check macro definition and "+
			"call site.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnterminatedConditional,
			StagePreprocess,
			"unterminated conditional block",
		),
		"Add missing `#endif` for active `#if` or `#ifdef` frame.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPInvalidIncludeSyntax,
			StagePreprocess,
			"invalid include syntax",
		),
		"Use valid include directive format and "+
			"quote style expected by this parser mode.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnexpectedElif,
			StagePreprocess,
			"unexpected elif directive",
		),
		"Ensure `#elif` is placed "+
			"after matching `#if` and before `#endif`.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnexpectedElse,
			StagePreprocess,
			"unexpected else directive",
		),
		"Check conditional directive structure: `#else` must be inside active block "+
			"and only once per frame.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnexpectedEndif,
			StagePreprocess,
			"unexpected endif directive",
		),
		"Remove stray `#endif` or restore missing opening directive.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPDirectiveError,
			StagePreprocess,
			"`#error` directive triggered",
		),
		"Source contains explicit `#error` directive and preprocessing stops. "+
			"Resolve condition or remove debug guard.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnsupportedDirective,
			StagePreprocess,
			"unsupported directive",
		),
		"Directive token is recognized as preprocessor command but not supported "+
			"by this implementation mode.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPMissingMacroName,
			StagePreprocess,
			"missing macro name in define",
		),
		"Add valid identifier after `#define` keyword.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPInvalidMacroName,
			StagePreprocess,
			"invalid macro name in define",
		),
		"Macro names must match parser identifier rules.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnterminatedMacroParams,
			StagePreprocess,
			"unterminated macro parameter list",
		),
		"Add missing closing parenthesis in function-like `#define`.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodePPMacroRedefined,
			StagePreprocess,
			"macro redefinition",
		),
		"Later definition overrides previous one in current preprocessor state.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnresolvedMacroInvocation,
			StagePreprocess,
			"unresolved macro-like invocation",
		),
		"Invocation looks like macro call but no matching definition was found. "+
			"Check macro name, scope, and include order.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodePPUnsupportedHasInclude,
			StagePreprocess,
			"unsupported `__has_include` in conditional directive",
		),
		"Conditional expression uses `__has_include`, which is not supported in "+
			"this preprocessor mode.",
	),
}
