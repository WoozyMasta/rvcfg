<!-- Automatically generated file, do not modify! -->

# Lint Rules Registry

This document contains the current registry of lint rules.

Total rules: 52.

## rvcfg

Real Virtuality Configs

> Lint rules for Real Virtuality config lexer, parser and preprocessor flows.

### lex

> Lexer diagnostics.

#### `CFG1001`

Unexpected character

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unexpected-character` |
| Scope | `lex` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG1002`

Unterminated string literal

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unterminated-string-literal` |
| Scope | `lex` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG1003`

Unterminated block comment

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.lex.unterminated-block-comment` |
| Scope | `lex` |
| Severity | `error` |
| Enabled | `true` (implicit) |

### parse

> Parser diagnostics.

#### `CFG2001`

Unexpected token

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.unexpected-token` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2002`

Expected class name

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-class-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2003`

Expected class body or semicolon

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-class-body-or-semicolon` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2004`

Expected closing brace for class body

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-closing-brace-for-class-body` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2005`

Missing semicolon after class declaration

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-class-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2006`

Expected name after `delete`

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-name-after-delete` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2007`

Missing semicolon after delete declaration

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-delete-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2008`

Expected `extern` declaration name

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-extern-declaration-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2009`

Missing semicolon after extern declaration

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-extern-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2010`

Expected right bracket in array assignment target

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-right-bracket-in-array-assignment-target` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2011`

Expected assignment operator in array assignment

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-assignment-operator-in-array-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2012`

Missing semicolon after array assignment

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-array-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2013`

Expected assignment operator

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-assignment-operator` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2014`

Missing semicolon after assignment

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-assignment` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2015`

Expected value before end of file

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-value-before-end-of-file` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2016`

Expected value

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-value` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2017`

Expected scalar value

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-scalar-value` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2018`

Unterminated array literal

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.unterminated-array-literal` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2019`

Expected comma or right brace in array literal

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-comma-or-right-brace-in-array-literal` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2020`

Autofix inserted missing class semicolon

> Parser recovered by inserting a missing semicolon after class declaration. Keep
> semicolons explicit to avoid parser recovery differences between tools.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.autofix-inserted-missing-class-semicolon` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2021`

Expected enum body

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-body` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2022`

Expected enum item name

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-item-name` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2023`

Expected enum item delimiter

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.expected-enum-item-delimiter` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2024`

Missing semicolon after enum declaration

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.missing-semicolon-after-enum-declaration` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2025`

Class-like names must not start with digit in strict mode

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.class-like-names-must-not-start-with-digit-in-strict-mode` |
| Scope | `parse` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG2026`

Nested class in derived class has no explicit base

> This can replace inherited subtree instead of extending it. Add explicit
> inheritance to make merge behavior predictable.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.nested-class-in-derived-class-has-no-explicit-base` |
| Scope | `parse` |
| Severity | `info` |
| Enabled | `true` (implicit) |

#### `CFG2027`

Scalar may be unsupported by RAP encoder

> Scalar token sequence cannot be classified into RAP scalar subtype safely. Use
> quoted string, integer, float, or identifier-like scalar form.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.scalar-may-be-unsupported-by-rap-encoder` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2028`

Float loses precision in RAP float32 conversion

> RAP stores float scalars as float32. Parsed value differs noticeably after
> float64->float32 conversion.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-loses-precision-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2029`

Float may collapse to zero in RAP float32 conversion

> Very small non-zero float becomes zero when encoded as RAP float32.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-may-collapse-to-zero-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2030`

String escape style may be incompatible with BI/CfgConvert

> Detected C-style backslash quote escape in string scalar. Prefer doubled quote
> escaping for BI/CfgConvert compatibility.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.string-escape-style-may-be-incompatible-with-bi-cfgconvert` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2031`

Extreme float magnitude may normalize unexpectedly in RAP round-trip

> Extreme exponent/magnitude float can look unstable in text output after RAP
> float32 conversion.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.extreme-float-magnitude-may-normalize-unexpectedly-in-rap-round-trip` |
| Scope | `parse` |
| Severity | `info` |
| Enabled | `true` (implicit) |

#### `CFG2032`

Float overflows to Inf in RAP float32 conversion

> Float literal exceeds float32 finite range and becomes +Inf/-Inf during RAP
> numeric encoding.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.parse.float-overflows-to-inf-in-rap-float32-conversion` |
| Scope | `parse` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG2033`

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

> Preprocessor diagnostics.

#### `CFG3001`

Include target not found or unreadable

> Check include path, include roots, and file access permissions.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.include-target-not-found-or-unreadable` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3002`

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

#### `CFG3003`

Macro expansion failure

> Usually caused by invalid invocation, argument count mismatch, or recursive
> expansion limit. Check macro definition and call site.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.macro-expansion-failure` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3004`

Unterminated conditional block

> Add missing `#endif` for active `#if` or `#ifdef` frame.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unterminated-conditional-block` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3005`

Invalid include syntax

> Use valid include directive format and quote style expected by this parser mode.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.invalid-include-syntax` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3006`

Unexpected elif directive

> Ensure `#elif` is placed after matching `#if` and before `#endif`.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-elif-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3007`

Unexpected else directive

> Check conditional directive structure: `#else` must be inside active block and
> only once per frame.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-else-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3008`

Unexpected endif directive

> Remove stray `#endif` or restore missing opening directive.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unexpected-endif-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3009`

`#Error` directive triggered

> Source contains explicit `#error` directive and preprocessing stops. Resolve
> condition or remove debug guard.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.error-directive-triggered` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3011`

Unsupported directive

> Directive token is recognized as preprocessor command but not supported by this
> implementation mode.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unsupported-directive` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3012`

Missing macro name in define

> Add valid identifier after `#define` keyword.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.missing-macro-name-in-define` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3013`

Invalid macro name in define

> Macro names must match parser identifier rules.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.invalid-macro-name-in-define` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3014`

Unterminated macro parameter list

> Add missing closing parenthesis in function-like `#define`.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unterminated-macro-parameter-list` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3015`

Macro redefinition

> Later definition overrides previous one in current preprocessor state.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.macro-redefinition` |
| Scope | `preprocess` |
| Severity | `warning` |
| Enabled | `true` (implicit) |

#### `CFG3016`

Unresolved macro-like invocation

> Invocation looks like macro call but no matching definition was found. Check
> macro name, scope, and include order.

| Field | Value |
| --- | --- |
| Rule ID | `rvcfg.preprocess.unresolved-macro-like-invocation` |
| Scope | `preprocess` |
| Severity | `error` |
| Enabled | `true` (implicit) |

#### `CFG3017`

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
