<!-- Automatically generated file, do not modify! -->

# Lint Rules Registry

This document contains the current registry of lint rules.

Total rules: 52.

## rvcfg

Real Virtuality Configs

> Lint rules for Real Virtuality config lexer, parser and preprocessor flows.

Rule groups for `rvcfg`:

* [lex](#lex)
* [parse](#parse)
* [preprocess](#preprocess)

### lex

> Tokenization diagnostics for invalid characters and unfinished literals.

Codes:
[RVCFG1001](#rvcfg1001),
[RVCFG1002](#rvcfg1002),
[RVCFG1003](#rvcfg1003),

#### `RVCFG1001`

Unexpected character

> Source contains byte sequence that does not belong to config token syntax.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unexpected-character` |
| Scope | `lex` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG1002`

Unterminated string literal

> String starts with quote but does not have valid closing quote on the same
> logical line.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unterminated-string-literal` |
| Scope | `lex` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG1003`

Unterminated block comment

> Block comment opened with /* is not closed with */ before end of file.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unterminated-block-comment` |
| Scope | `lex` |
| Severity | `error` |
| Enabled | `true` (implicit) |

### parse

> Grammar and structure diagnostics for classes, assignments, arrays, and enums.

Codes:
[RVCFG2001](#rvcfg2001),
[RVCFG2002](#rvcfg2002),
[RVCFG2003](#rvcfg2003),
[RVCFG2004](#rvcfg2004),
[RVCFG2005](#rvcfg2005),
[RVCFG2006](#rvcfg2006),
[RVCFG2007](#rvcfg2007),
[RVCFG2008](#rvcfg2008),
[RVCFG2009](#rvcfg2009),
[RVCFG2010](#rvcfg2010),
[RVCFG2011](#rvcfg2011),
[RVCFG2012](#rvcfg2012),
[RVCFG2013](#rvcfg2013),
[RVCFG2014](#rvcfg2014),
[RVCFG2015](#rvcfg2015),
[RVCFG2016](#rvcfg2016),
[RVCFG2017](#rvcfg2017),
[RVCFG2018](#rvcfg2018),
[RVCFG2019](#rvcfg2019),
[RVCFG2020](#rvcfg2020),
[RVCFG2021](#rvcfg2021),
[RVCFG2022](#rvcfg2022),
[RVCFG2023](#rvcfg2023),
[RVCFG2024](#rvcfg2024),
[RVCFG2025](#rvcfg2025),
[RVCFG2026](#rvcfg2026),
[RVCFG2027](#rvcfg2027),
[RVCFG2028](#rvcfg2028),
[RVCFG2029](#rvcfg2029),
[RVCFG2030](#rvcfg2030),
[RVCFG2031](#rvcfg2031),
[RVCFG2032](#rvcfg2032),
[RVCFG2033](#rvcfg2033),

#### `RVCFG2001`

Unexpected token

> Token order does not match grammar for current parser context.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.unexpected-token` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2002`

Expected class name

> Class declaration must provide one identifier after class keyword.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-class-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2003`

Expected class body or semicolon

> Class declaration must be either forward declaration with semicolon or class
> body in braces.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-class-body-or-semicolon` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2004`

Expected closing brace for class body

> Class body was opened but closing } token was not found at the expected
> boundary.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-closing-brace-for-class-body` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2005`

Missing semicolon after class declaration

> Class declaration must end with semicolon after closing brace.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-class-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2006`

Expected name after `delete`

> Delete statement must target one identifier.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-name-after-delete` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2007`

Missing semicolon after delete declaration

> Delete statement must end with semicolon.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-delete-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2008`

Expected `extern` declaration name

> Extern declaration must provide one symbol name.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-extern-declaration-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2009`

Missing semicolon after extern declaration

> Extern declaration must end with semicolon.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-extern-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2010`

Expected right bracket in array assignment target

> Array assignment target must use [] suffix with closing right bracket.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-right-bracket-in-array-assignment-target` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2011`

Expected assignment operator in array assignment

> Array assignment expects = or += operator after target.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-assignment-operator-in-array-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2012`

Missing semicolon after array assignment

> Array assignment statement must end with semicolon.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-array-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2013`

Expected assignment operator

> Property assignment expects = operator after left-hand identifier.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-assignment-operator` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2014`

Missing semicolon after assignment

> Property assignment statement must end with semicolon.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2015`

Expected value before end of file

> Parser reached end of file where value expression was required.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-value-before-end-of-file` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2016`

Expected value

> Assignment or array element requires scalar or array value.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-value` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2017`

Expected scalar value

> Scalar expression is empty or cannot be extracted from token range.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-scalar-value` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2018`

Unterminated array literal

> Array literal opened with { but matching } was not found.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.unterminated-array-literal` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2019`

Expected comma or right brace in array literal

> Array elements must be separated by comma and array must end with right brace.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-comma-or-right-brace-in-array-literal` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2020`

Autofix inserted missing class semicolon

> Parser recovered by inserting a missing semicolon after class declaration. Keep
> semicolons explicit to avoid parser recovery differences between tools.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.autofix-inserted-missing-class-semicolon` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2021`

Expected enum body

> Enum declaration must contain body enclosed in braces.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-body` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2022`

Expected enum item name

> Enum item must start with one identifier name.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-item-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2023`

Expected enum item delimiter

> Enum items must be separated by comma or terminated by closing brace.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-item-delimiter` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2024`

Missing semicolon after enum declaration

> Enum declaration must end with semicolon after closing brace.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-enum-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2025`

Class-like names must not start with digit in strict mode

> Strict parser mode rejects declarations whose class-like names start with
> digits.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.class-like-names-must-not-start-with-digit-in-strict-mode` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG2026`

Nested class in derived class has no explicit base

> This can replace inherited subtree instead of extending it. Add explicit
> inheritance to make merge behavior predictable.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.nested-class-in-derived-class-has-no-explicit-base` |
| Scope | `parse` |
| Severity | `info` |
| Enabled | `true` (implicit) |

#### `RVCFG2027`

Scalar may be unsupported by RAP encoder

> Scalar token sequence cannot be classified into RAP scalar subtype safely. Use
> quoted string, integer, float, or identifier-like scalar form.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.scalar-may-be-unsupported-by-rap-encoder` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2028`

Float loses precision in RAP float32 conversion

> RAP stores float scalars as float32. Parsed value differs noticeably after
> float64->float32 conversion.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-loses-precision-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2029`

Float may collapse to zero in RAP float32 conversion

> Very small non-zero float becomes zero when encoded as RAP float32.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-may-collapse-to-zero-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2030`

String escape style may be incompatible with BI/CfgConvert

> Detected C-style backslash quote escape in string scalar. Prefer doubled quote
> escaping for BI/CfgConvert compatibility.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.string-escape-style-may-be-incompatible-with-bi-cfgconvert` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2031`

Extreme float magnitude may normalize unexpectedly in RAP round-trip

> Extreme exponent/magnitude float can look unstable in text output after RAP
> float32 conversion.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.extreme-float-magnitude-may-normalize-unexpectedly-in-rap-round-trip` |
| Scope | `parse` |
| Severity | `info` |
| Enabled | `true` (implicit) |

#### `RVCFG2032`

Float overflows to Inf in RAP float32 conversion

> Float literal exceeds float32 finite range and becomes +Inf/-Inf during RAP
> numeric encoding.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-overflows-to-inf-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG2033`

Non-finite float literal is unsafe for RAP numeric encoding

> Detected NaN/Inf-like scalar literal. RAP numeric scalar encoding expects finite
> values.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.non-finite-float-literal-is-unsafe-for-rap-numeric-encoding` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

### preprocess

> Directive/macro/include diagnostics for preprocessing and conditional
> evaluation.

Codes:
[RVCFG3001](#rvcfg3001),
[RVCFG3002](#rvcfg3002),
[RVCFG3003](#rvcfg3003),
[RVCFG3004](#rvcfg3004),
[RVCFG3005](#rvcfg3005),
[RVCFG3006](#rvcfg3006),
[RVCFG3007](#rvcfg3007),
[RVCFG3008](#rvcfg3008),
[RVCFG3009](#rvcfg3009),
[RVCFG3011](#rvcfg3011),
[RVCFG3012](#rvcfg3012),
[RVCFG3013](#rvcfg3013),
[RVCFG3014](#rvcfg3014),
[RVCFG3015](#rvcfg3015),
[RVCFG3016](#rvcfg3016),
[RVCFG3017](#rvcfg3017),

#### `RVCFG3001`

Include target not found or unreadable

> Check include path, include roots, and file access permissions.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.include-target-not-found-or-unreadable` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3002`

Unsupported config intrinsic

> Current preprocessor mode does not support this intrinsic (for example `__EXEC`
> or `__EVAL`). Remove it or preprocess with toolchain that supports this
> intrinsic.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unsupported-config-intrinsic` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3003`

Macro expansion failure

> Usually caused by invalid invocation, argument count mismatch, or recursive
> expansion limit. Check macro definition and call site.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.macro-expansion-failure` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3004`

Unterminated conditional block

> Add missing `#endif` for active `#if` or `#ifdef` frame.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unterminated-conditional-block` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3005`

Invalid include syntax

> Use valid include directive format and quote style expected by this parser mode.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.invalid-include-syntax` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3006`

Unexpected elif directive

> Ensure `#elif` is placed after matching `#if` and before `#endif`.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-elif-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3007`

Unexpected else directive

> Check conditional directive structure: `#else` must be inside active block and
> only once per frame.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-else-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3008`

Unexpected endif directive

> Remove stray `#endif` or restore missing opening directive.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-endif-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3009`

`#Error` directive triggered

> Source contains explicit `#error` directive and preprocessing stops. Resolve
> condition or remove debug guard.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.error-directive-triggered` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3011`

Unsupported directive

> Directive token is recognized as preprocessor command but not supported by this
> implementation mode.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unsupported-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3012`

Missing macro name in define

> Add valid identifier after `#define` keyword.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.missing-macro-name-in-define` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3013`

Invalid macro name in define

> Macro names must match parser identifier rules.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.invalid-macro-name-in-define` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3014`

Unterminated macro parameter list

> Add missing closing parenthesis in function-like `#define`.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unterminated-macro-parameter-list` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3015`

Macro redefinition

> Later definition overrides previous one in current preprocessor state.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.macro-redefinition` |
| Scope | `preprocess` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `RVCFG3016`

Unresolved macro-like invocation

> Invocation looks like macro call but no matching definition was found. Check
> macro name, scope, and include order.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unresolved-macro-like-invocation` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `RVCFG3017`

Unsupported `__has_include` in conditional directive

> Conditional expression uses `__has_include`, which is not supported in this
> preprocessor mode.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unsupported-has-include-in-conditional-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

---

> Generated with
> [lintkit](https://github.com/woozymasta/lintkit)
> version `dev`
> commit `unknown`

<!-- Automatically generated file, do not modify! -->
