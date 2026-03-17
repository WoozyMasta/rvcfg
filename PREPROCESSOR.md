# rvcfg preprocessor macros

`rvcfg` runs a preprocessor before parsing config text.
This document explains what each directive and intrinsic does, and when it
is available.

## What the preprocessor does

The preprocessor resolves `#include` files, expands `#define` macros, and
applies conditional directives so some blocks are kept and others are
skipped. In extended profiles it can also generate repetitive config
fragments with helper intrinsics.

## Modes

* `strict`: default profile with DayZ-oriented behavior.
* `compat`: `strict` plus compatibility extensions.
* `extended`: `compat` plus `rvcfg` utilities for strings, paths, and file
  list generation.

## Strict mode

Strict mode targets semantic compatibility with DayZ `CfgConvert`.
Known differences are mostly text-level `-pcpp` formatting quirks in some
malformed macro edge cases. The parsed structure is expected to remain
equivalent.

### `#include`

Includes another file and inserts its preprocessed content at this point.
Both quote form and angle form are supported.

```cpp
#include "common.hpp"
#include <defs/items.hpp>
```

### `#define` (object-like)

Defines a plain textual macro without arguments.
Every matching token is replaced with the macro body.

```cpp
#define MOD_TAG "MyMod"
displayName = MOD_TAG;
```

### `#define` (function-like)

Defines a macro with arguments.
Arguments are substituted into the macro body when invoked.

```cpp
#define MAT(BASE) "data/" BASE ".rvmat"
material = MAT("crate");
```

### `#undef`

Removes an existing macro definition.
After `#undef`, that macro name is no longer expanded.

```cpp
#undef MOD_TAG
```

### Stringify `#`

Converts a macro argument into a string literal.

```cpp
#define STR(X) #X
text = STR(HELLO_WORLD);
```

### Token paste `##`

Concatenates two macro tokens into one token.

```cpp
#define CAT(A,B) A##B
class CAT(My,Class) {};
```

### `#ifdef`

Keeps a block only when the macro name is defined.

```cpp
#ifdef DEBUG
scope = 2;
#endif
```

### `#ifndef`

Keeps a block only when the macro name is not defined.

```cpp
#ifndef RELEASE
class DevOnly {};
#endif
```

### `#else`

Switches to the alternative branch of the current conditional block.

```cpp
#ifdef DEBUG
scope = 2;
#else
scope = 1;
#endif
```

### `#endif`

Closes a conditional block started by `#ifdef` or `#ifndef`.

```cpp
#ifdef SOME_FLAG
class A {};
#endif
```

### `__LINE__`

Expands to the current line number in the currently processed file.
This is useful for diagnostics and generated IDs.

```cpp
lineValue = __LINE__; // example: 12
```

### `__FILE__`

Expands to the current source file path as a string.
For the root file this is usually the resolved path used for preprocess.
For include files this is the include display path used in that include
context.

```cpp
fileValue = __FILE__; // example: "config.cpp" or "defs/items.hpp"
```

## Compat mode

`compat` adds the following features on top of `strict`.

### `#if`, `#elif`, `defined(...)`

Enables expression-based conditionals.
Use this when branches depend on numeric expressions or macro presence.
Supported expression features:

* logical operators: `||`, `&&`, unary `!`;
* comparisons: `==`, `!=`, `>=`, `<=`, `>`, `<`;
* parentheses for grouping;
* integer literals and integer macro bodies;
* `defined(NAME)` and `defined NAME`.

Boundary behavior:

* non-numeric macro bodies are treated as non-numeric in comparisons;
* unsupported `__has_include(...)` is rejected.

```cpp
#if GAME_VER > 20
v = 3;
#elif GAME_VER > 16
v = 2;
#elif defined(DEBUG)
v = 1;
#endif
```

### `__EXEC(...)`

Enables execution of assignment-like setup statements for intrinsic
expression context.
Accepted body format is a semicolon-separated assignment list, for example:
`name = expr; other = expr`.

Boundary behavior:

* variable names must start with letter, `_`, or `$`;
* remaining name characters can be letters, digits, `_`, or `$`;
* each statement must be an assignment.

```cpp
__EXEC(base = "my/mod"; id = 7; item = "crate")
// variables are stored for later __EVAL usage
```

### `__EVAL(...)`

Enables expression evaluation and substitutes the result.
Supports numeric operations and values prepared via `__EXEC`.
Supported expression features:

* operators: unary `+`/`-`, binary `+`, `-`, `*`, `/`;
* parentheses;
* numeric literals (integer and float);
* quoted string literals;
* variables defined in `__EXEC`.

Boundary behavior:

* `+` works for both numeric addition and string concatenation;
* `-`, `*`, `/` are numeric-only;
* expression parse/eval errors fall back to `"scalar"` (string).

```cpp
numberValue = __EVAL(id + 5);                 // 12
pathValue = __EVAL(base + "/data/" + item);   // "my/mod/data/crate"
```

### `__FILE_NAME__`

Expands to the file base name with extension.

```cpp
nameValue = __FILE_NAME__;
```

### `__FILE_SHORT__`

Expands to the file base name without extension.

```cpp
shortValue = __FILE_SHORT__;
```

### Dynamic date/time intrinsics

These values change between runs and are not deterministic.

* `__DATE_ARR__`: date/time tuple, for example `2026,3,11,15,4,9`.
* `__DATE_STR__`: local datetime string, for example
  `"2026/03/11, 15:04:09"`.
* `__DATE_STR_ISO8601__`: UTC ISO-8601 datetime, for example
  `"2026-03-11T12:04:09Z"`.
* `__TIME__`: local time, for example `15:04:09`.
* `__TIME_UTC__`: UTC time, for example `12:04:09`.
* `__DAY__`: current day number, for example `11`.
* `__MONTH__`: current month number, for example `3`.
* `__YEAR__`: current year number, for example `2026`.
* `__TIMESTAMP_UTC__`: unix timestamp, for example `1773221049`.

### Counter intrinsics

* `__COUNTER__`: returns current counter and increments it.
* `__COUNTER_RESET__`: resets the counter to zero.

### Random number intrinsics

Generates pseudo-random values of specific bit width.
Only 8/16/32/64-bit suffixes are supported.

Ranges:

* `__RAND_INT8__`: `-128..127`;
* `__RAND_INT16__`: `-32768..32767`;
* `__RAND_INT32__`: `-2147483648..2147483647`;
* `__RAND_INT64__`: full `int64` range;
* `__RAND_UINT8__`: `0..255`;
* `__RAND_UINT16__`: `0..65535`;
* `__RAND_UINT32__`: `0..4294967295`;
* `__RAND_UINT64__`: full `uint64` range.

Implementation note: random values are generated via `crypto/rand`.

## Extended mode

`extended` adds `rvcfg` utilities on top of `compat`.

### `__PATH_NORM(expr)`

Normalizes path separators to a stable backslash form.
Useful when path fragments come from mixed slash styles.

```cpp
__EXEC(base = "mods/demo")
path = __PATH_NORM(base + "/assets\\data//x.paa");
// result: "mods\demo\assets\data\x.paa"
```

### `__STR_TRIM(value)`

Removes leading and trailing whitespace.

```cpp
v = __STR_TRIM("  demo  "); // "demo"
```

### `__STR_LOWER(value)`

Converts text to lowercase.

```cpp
v = __STR_LOWER("Demo Item"); // "demo item"
```

### `__STR_UPPER(value)`

Converts text to uppercase.

```cpp
v = __STR_UPPER("Demo Item"); // "DEMO ITEM"
```

### `__STR_REPLACE(text, old, new)`

Replaces all `old` occurrences with `new`.

```cpp
v = __STR_REPLACE("a/b/c", "/", "\\"); // "a\b\c"
```

### `__STR_JOIN(delimiter, value1, ...)`

Joins values into one string with a delimiter.
Requires at least one value after delimiter.

```cpp
v = __STR_JOIN("/", "mods", "demo", "data"); // "mods/demo/data"
```

### `__STR_SPLIT(text, delimiter, index)`

Splits by delimiter and returns item at `index`.
Out-of-range index returns an empty string.
`delimiter` must be non-empty.
`index` must be an integer and must be `>= 0`.

```cpp
v1 = __STR_SPLIT("a_b_c", "_", 1); // "b"
v2 = __STR_SPLIT("a_b_c", "_", 9); // ""
```

### `__STR_QUOTE(value)`

Evaluates the expression and returns it as a quoted string literal.

```cpp
v = __STR_QUOTE("mods/demo"); // "\"mods/demo\""
```

### `__STR_PASCAL(value)`

Converts text to `UpperCamelCase`.
Separators (`space`, `_`, `-`, punctuation) are treated as word boundaries.

```cpp
v = __STR_PASCAL("hello-world_demo x"); // "HelloWorldDemoX"
```

### `__STR_CAMEL(value)`

Converts text to `lowerCamelCase`.
Word boundary rules are the same as `__STR_PASCAL`.

```cpp
v = __STR_CAMEL("hello-world_demo x"); // "helloWorldDemoX"
```

### `__STR_SNAKE(value)`

Converts text to `lower_snake_case`.

```cpp
v = __STR_SNAKE("HelloWorld Demo"); // "hello_world_demo"
```

### `__STR_CONST(value)`

Converts text to `UPPER_SNAKE_CASE`.

```cpp
v = __STR_CONST("HelloWorld Demo"); // "HELLO_WORLD_DEMO"
```

### `__FILES_JOIN(pattern, delimiter)`

Finds files by glob pattern, sorts matches deterministically, and joins
them into one string. If delimiter is omitted, `;` is used.
Requires `pattern` argument. Optional second argument is delimiter.

### `__FILES_COUNT(pattern)`

Returns the number of files matched by glob pattern.
Requires exactly one argument.

### `__FILES_GET(pattern, index)`

Returns one matched file by index from the sorted list.
Out-of-range index returns an empty string.
`index` must be an integer and must be `>= 0`.

### `__FILES_RENDER(pattern, template, delimiter)`

Renders each matched file with a template and joins rendered items.
This emits raw text (not an auto-quoted string).
Requires `pattern` and `template`.
Optional third argument is delimiter (default newline).

> Template placeholders support filters. See **Template filters** below.

Available placeholders:

* `{path}`: relative matched file path.
* `{name}`: file name with extension.
* `{stem}`: file name without extension.
* `{ext}`: extension.
* `{index}`: zero-based position in sorted match list.

Examples below assume this sorted match set:
`sounds/a.ogg`, `sounds/b.ogg`.

```cpp
joined = __FILES_JOIN("sounds/*.ogg", "|");
// "sounds\a.ogg|sounds\b.ogg"

count = __FILES_COUNT("sounds/*.ogg");
// 2

first = __FILES_GET("sounds/*.ogg", 0);
// "sounds\a.ogg"
```

```cpp
rows = __FILES_RENDER(
  "sounds/*.ogg",
  "{index}:{stem}={path}",
  "|"
);
// result: 0:a=sounds\a.ogg|1:b=sounds\b.ogg
```

### `__FOR_RANGE_RENDER(start, end, template, delimiter)`

Renders a numeric inclusive range with a template and joins results.
Works in both directions (`1..3` and `3..1`).
This emits raw text (not an auto-quoted string).
`start` and `end` must be integers.
Optional fourth argument is delimiter (default newline).

> Template placeholders support filters. See **Template filters** below.

Placeholders:

* `{index}`: zero-based iteration index.
* `{value}`: current numeric value.

```cpp
rowsAsc = __FOR_RANGE_RENDER(1, 3, "{index}:{value}", "|");
rowsDesc = __FOR_RANGE_RENDER(3, 1, "{value}", ",");
// rowsAsc:  0:1|1:2|2:3
// rowsDesc: 3,2,1
```

### `__FOR_EACH_RENDER(template, delimiter, value1, ...)`

Renders an explicit list of values with a template and joins results.
Values are processed in argument order.
This emits raw text (not an auto-quoted string).
Requires at least one value argument.

> Template placeholders support filters. See **Template filters** below.

Placeholders:

* `{index}`: zero-based iteration index.
* `{value}`: current list value.

```cpp
items = __FOR_EACH_RENDER(
  "{index}:{value}",
  "|",
  "alpha",
  "beta",
  "gamma"
);
// result: 0:alpha|1:beta|2:gamma
```

### Template filters

Template filters are inline transformations for placeholder values inside
render templates. Filters are applied left-to-right and can be chained.

Template filters are available in:

* `__FILES_RENDER(...)`
* `__FOR_RANGE_RENDER(...)`
* `__FOR_EACH_RENDER(...)`

Supported filters:

* `trim` - see `__STR_TRIM()`.
* `lower` - see `__STR_LOWER()`.
* `upper` - see `__STR_UPPER()`.
* `replace(old, new)` - see `__STR_REPLACE()`.
* `split(delimiter, index)` - see `__STR_SPLIT()`.
* `quote` - see `__STR_QUOTE()`.
* `pascal` - see `__STR_PASCAL()`.
* `camel` - see `__STR_CAMEL()`.
* `snake` - see `__STR_SNAKE()`.
* `const` - see `__STR_CONST()`.
* `path_norm` (alias: `slash_norm`) - see `__PATH_NORM()`.

Examples:

```cpp
rows = __FILES_RENDER(
  "sounds/*.ogg",
  "class Snd_{stem|pascal} { file = {name|const}; };",
  "\n"
);
```

```cpp
rows[] = {
__FILES_RENDER(
  "sounds/*.ogg",
  "{index}:{path|lower|replace(sounds, sfx)};"
)};
```

## Extended limits

Extended mode has safety limits:

* `ExtendedFSRoot`:
  limits filesystem root for `__FILES_*` lookups.
* `ExtendedFSMaxItems`:
  caps glob match count for `__FILES_*` (default `512`).
* `ExtendedLoopMaxItems`:
  caps iteration count for `__FOR_RANGE_RENDER` and
  `__FOR_EACH_RENDER` (default `2048`).

General preprocess limits:

* include depth limit (`MaxIncludeDepth`) is `32` by default;
* macro/intrinsic expansion depth (`MaxExpandDepth`) is `32` by default.
