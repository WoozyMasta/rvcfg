// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"github.com/woozymasta/lintkit/lint"
)

// withDescription attaches optional documentation text to one catalog spec.
func withDescription(spec lint.CodeSpec, description string) lint.CodeSpec {
	spec.Description = description
	return spec
}

// diagnosticCatalog stores stable diagnostics metadata table.
var diagnosticCatalog = []lint.CodeSpec{
	withDescription(
		lint.WarningCodeSpec(
			CodeLexUnexpectedCharacter,
			StageLex,
			"unexpected character",
		),
		"Source contains byte sequence that does not belong to config token syntax.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeLexUnterminatedString,
			StageLex,
			"unterminated string literal",
		),
		"String starts with quote but does not have valid closing quote on the same logical line.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeLexUnterminatedBlockComment,
			StageLex,
			"unterminated block comment",
		),
		"Block comment opened with /* is not closed with */ before end of file.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParUnexpectedToken,
			StageParse,
			"unexpected token",
		),
		"Token order does not match grammar for current parser context.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedClassName,
			StageParse,
			"expected class name",
		),
		"Class declaration must provide one identifier after class keyword.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedClassBodyOrSemicolon,
			StageParse,
			"expected class body or semicolon",
		),
		"Class declaration must be either forward declaration with semicolon or class body in braces.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedClassClosingBrace,
			StageParse,
			"expected closing brace for class body",
		),
		"Class body was opened but closing } token was not found at the expected boundary.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingClassSemicolon,
			StageParse,
			"missing semicolon after class declaration",
		),
		"Class declaration must end with semicolon after closing brace.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedDeleteName,
			StageParse,
			"expected name after `delete`",
		),
		"Delete statement must target one identifier.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingDeleteSemicolon,
			StageParse,
			"missing semicolon after delete declaration",
		),
		"Delete statement must end with semicolon.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedExternName,
			StageParse,
			"expected `extern` declaration name",
		),
		"Extern declaration must provide one symbol name.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingExternSemicolon,
			StageParse,
			"missing semicolon after extern declaration",
		),
		"Extern declaration must end with semicolon.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedArrayRightBracket,
			StageParse,
			"expected right bracket in array assignment target",
		),
		"Array assignment target must use [] suffix with closing right bracket.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedArrayAssignOperator,
			StageParse,
			"expected assignment operator in array assignment",
		),
		"Array assignment expects = or += operator after target.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingArrayAssignSemicolon,
			StageParse,
			"missing semicolon after array assignment",
		),
		"Array assignment statement must end with semicolon.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedAssign,
			StageParse,
			"expected assignment operator",
		),
		"Property assignment expects = operator after left-hand identifier.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingAssignSemicolon,
			StageParse,
			"missing semicolon after assignment",
		),
		"Property assignment statement must end with semicolon.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedValueBeforeEOF,
			StageParse,
			"expected value before end of file",
		),
		"Parser reached end of file where value expression was required.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedValue,
			StageParse,
			"expected value",
		),
		"Assignment or array element requires scalar or array value.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedScalarValue,
			StageParse,
			"expected scalar value",
		),
		"Scalar expression is empty or cannot be extracted from token range.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParUnterminatedArrayLiteral,
			StageParse,
			"unterminated array literal",
		),
		"Array literal opened with { but matching } was not found.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedArrayDelimiter,
			StageParse,
			"expected comma or right brace in array literal",
		),
		"Array elements must be separated by comma and array must end with right brace.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParStrictDigitLeadingClassName,
			StageParse,
			"class-like names must not start with digit in strict mode",
		),
		"Strict parser mode rejects declarations whose class-like names start with digits.",
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
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedEnumBody,
			StageParse,
			"expected enum body",
		),
		"Enum declaration must contain body enclosed in braces.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedEnumItemName,
			StageParse,
			"expected enum item name",
		),
		"Enum item must start with one identifier name.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParExpectedEnumDelimiter,
			StageParse,
			"expected enum item delimiter",
		),
		"Enum items must be separated by comma or terminated by closing brace.",
	),
	withDescription(
		lint.ErrorCodeSpec(
			CodeParMissingEnumSemicolon,
			StageParse,
			"missing semicolon after enum declaration",
		),
		"Enum declaration must end with semicolon after closing brace.",
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
		lint.WarningCodeSpec(
			CodeParRAPUnsupportedScalar,
			StageParse,
			"scalar may be unsupported by RAP encoder",
		),
		"Scalar token sequence cannot be classified into RAP scalar subtype safely. "+
			"Use quoted string, integer, float, or identifier-like scalar form.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParRAPFloatPrecisionLoss,
			StageParse,
			"float loses precision in RAP float32 conversion",
		),
		"RAP stores float scalars as float32. Parsed value differs noticeably after "+
			"float64->float32 conversion.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParRAPFloatUnderflowToZero,
			StageParse,
			"float may collapse to zero in RAP float32 conversion",
		),
		"Very small non-zero float becomes zero when encoded as RAP float32.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParRAPUnsafeStringEscape,
			StageParse,
			"string escape style may be incompatible with BI/CfgConvert",
		),
		"Detected C-style backslash quote escape in string scalar. Prefer doubled "+
			"quote escaping for BI/CfgConvert compatibility.",
	),
	withDescription(
		lint.InfoCodeSpec(
			CodeParRAPExtremeFloatMagnitude,
			StageParse,
			"extreme float magnitude may normalize unexpectedly in RAP round-trip",
		),
		"Extreme exponent/magnitude float can look unstable in text output after "+
			"RAP float32 conversion.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParRAPFloatOverflowToInf,
			StageParse,
			"float overflows to Inf in RAP float32 conversion",
		),
		"Float literal exceeds float32 finite range and becomes +Inf/-Inf during "+
			"RAP numeric encoding.",
	),
	withDescription(
		lint.WarningCodeSpec(
			CodeParRAPNonFiniteFloat,
			StageParse,
			"non-finite float literal is unsafe for RAP numeric encoding",
		),
		"Detected NaN/Inf-like scalar literal. RAP numeric scalar encoding expects "+
			"finite values.",
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
