# Macros And Preprocessor: User Guide

This document explains macros in general:
what they are, why they are useful, and how to start using them
in real `config.cpp` and similar C-like config files.

A short technical section about the current `rvcfg` behavior is kept at the end.

## What It Is In Simple Terms

Macros and the preprocessor help you build final config text automatically.

Instead of manual copy-paste, you:

* extract repeated fragments into `#define`;
* include shared files with `#include`;
* enable or disable blocks with `#if/#ifdef`.

Result: less routine work, fewer mistakes, faster edits.

## Why It Is Useful

* write a template once and reuse it many times;
* update values across many classes faster;
* keep configs modular instead of one huge file;
* build different variants (for example dev/release, server/client).

## Minimal Syntax You Should Know

### 1. `#define` (constant)

`#define` creates a named text substitution.  
When the preprocessor sees that name, it replaces it with the defined value.

```cpp
#define MOD_TAG "MyMod"
#define CARGO_SLOTS 24
```

### 2. `#define` (macro with arguments)

This form works like a small text template with parameters.

```cpp
#define HL_ROW(PCT, MAT) {PCT, {MAT}}
```

### 3. `#include`

`#include` inserts another file content at this exact location.

```cpp
#include "my_block.hpp"
```

### 4. Conditional blocks

Conditional directives decide which text blocks are included in final output.

```cpp
#ifdef EXPERIMENTAL
class ExperimentalFeature {};
#endif
```

## Simple Example

```cpp
#define MOD_TAG "MyMod"
#define CARGO_SLOTS 24

class CfgVehicles
{
  class Backpack_Base;
  class MyBag: Backpack_Base
  {
    scope = 2;
    displayName = MOD_TAG;
    itemsCargoSize[] = {CARGO_SLOTS, CARGO_SLOTS};
  };
};
```

After preprocessing, values are substituted as plain text:

```cpp
displayName = "MyMod";
itemsCargoSize[] = {24, 24};
```

## Example With `#include`

`config.cpp`:

```cpp
class CfgPatches
{
  #include "patches.hpp"
};
```

`patches.hpp`:

```cpp
class MyTag_Main
{
  requiredVersion = 0.1;
  requiredAddons[] = {"DZ_Data"};
};
```

`#include` inserts `patches.hpp` content at the directive location.

## Routine Automation With Double-Nested Macro

Goal: pass only one base material name and generate
the full `healthLevels[]` block automatically.

```c
#define CAMPING_DATA_PATH "DZ/gear/camping/data/"
#define MAT_OK(BASE) CAMPING_DATA_PATH BASE ".rvmat"
#define MAT_DMG(BASE) CAMPING_DATA_PATH BASE "_damage.rvmat"
#define MAT_DEST(BASE) CAMPING_DATA_PATH BASE "_destruct.rvmat"
#define HL_ROW(PCT, MAT) {PCT, {MAT}}
#define HL_STD(MAT_OK, MAT_DMG, MAT_DEST) \
{ \
  HL_ROW(1.0, MAT_OK), \
  HL_ROW(0.5, MAT_DMG), \
  HL_ROW(0.0, MAT_DEST) \
}
#define HL_FROM_BASE(BASE) HL_STD(MAT_OK(BASE), MAT_DMG(BASE), MAT_DEST(BASE))
```

Now we just call it and get a full `healthLevels[]` definition:

```cpp
healthLevels[] = HL_FROM_BASE("bagpack");
```

Expanded result (after preprocessing):

```cpp
healthLevels[] = {
  {1.0, {"DZ/gear/camping/data/bagpack.rvmat"}},
  {0.5, {"DZ/gear/camping/data/bagpack_damage.rvmat"}},
  {0.0, {"DZ/gear/camping/data/bagpack_destruct.rvmat"}}
};
```

What happens:

* `MAT_OK/MAT_DMG/MAT_DEST` build full material paths from one base name;
* `HL_STD(...)` builds the standard damage-level array shape;
* `HL_FROM_BASE(...)` combines both levels into one call site;
* you pass only `"bagpack"` and get all three material variants automatically.

## Practical Rules

* keep macro names uppercase (`MOD_TAG`, `HL_STD`);
* avoid over-complicated 20+ line macros;
* prefer multiple simple macros over one magic macro;
* move shared macros into dedicated include files;
* add short comments for non-obvious macros.

## Common Mistakes

* missing closing `#endif`;
* typo in macro name;
* wrong path in `#include`;
* macro defined but not used where expected.

## Technical Notes For Current `rvcfg` Behavior

The following points describe the current implementation behavior.

### Processing Order

1. Preprocessor runs first (`#include`, `#define`, conditional branches).
2. Final produced text is processed next (parse, format, and other stages).

### Currently Supported

* `#include "path/file.hpp"`;
* `#define NAME value`;
* `#define NAME(arg1, arg2) body`;
* `#undef NAME`;
* `#if`, `#ifdef`, `#ifndef`, `#elif`, `#else`, `#endif`;
* `defined(NAME)` and `defined NAME`;
* `!`, `&&`, `||` inside `#if`
  (**rvcfg extension**, not guaranteed to match BI tools exactly);
* numeric comparisons in `#if`: `==`, `!=`, `>=`, `<=`, `>`, `<`
  (also **rvcfg extension**);
* `##` (token paste);
* `#` stringify in function-like macros;
* stable built-ins:
  * `__LINE__`
  * `__FILE__`
  * `__FILE_NAME__`
  * `__FILE_SHORT__`.

### `#` Stringify And `##` Token-Paste Example

Both operators are supported in function-like macros:

```cpp
#define STRINGIFY(S) #S
#define GLUE(A,B) A##B

stringify_result = STRINGIFY(ABC);
glue_result = GLUE(12,34);
```

Expected expanded text in `rvcfg` preprocess output:

```cpp
stringify_result = "ABC";
glue_result = 1234;
```

DayZ CfgConvert note: in some probe-style files,
`-pcpp` output may drop such probe lines instead of showing textual expansion.
`rvcfg` keeps explicit expanded lines in preprocess text.

### Not Supported

* `__has_include(...)` in `#if/#elif` (reported as preprocess error);
* complex arithmetic and comparison expressions in `#if`.

### `__EXEC` And `__EVAL` (Opt-In)

`rvcfg` supports `__EXEC` / `__EVAL` in compatibility mode.
It is disabled by default.

Enable with:

```go
rvcfg.PreprocessOptions{
  EnableExecEvalIntrinsics: true,
}
```

Examples:

```cpp
__EXEC(testVar = 7)
value = __EVAL(testVar + 5); // -> value = 12;
```

```cpp
value = __EVAL(1 + 2); // -> value = 3;
```

`__EVAL` support scope in current `rvcfg` implementation:

* supported:
  * numbers (`1`, `1.25`, `-2`);
  * variables assigned via `__EXEC`;
  * arithmetic operators `+`, `-`, `*`, `/`;
  * parentheses and unary `+`/`-`;
  * string values from `__EXEC` (returned as quoted strings).
* not supported:
  * engine/runtime expressions like `getResolution select 2`;
  * script commands or external VM state.

When `__EVAL(...)` contains unsupported pieces,
`rvcfg` uses compatibility fallback:

```cpp
value = __EVAL(getResolution select 2); // -> value = "scalar";
```

Example that works (all variables defined by `__EXEC`):

```cpp
__EXEC(someVar = 4; safeZoneW = 1.0)
w = __EVAL(safeZoneW - (5 * ((1 / (someVar - 2)) * 1.25 * 4)));
```

### `EnableDynamicIntrinsics` (Opt-In)

`rvcfg` also supports date/time/counter/random intrinsics in a separate mode.
This mode is disabled by default.

Enable with:

```go
rvcfg.PreprocessOptions{
  EnableDynamicIntrinsics: true,
}
```

Supported dynamic intrinsics and examples:

* `__DATE_ARR__` -> `2026,2,8,18,07,11`
* `__DATE_STR__` -> `"2026/02/08, 18:07:11"`
* `__DATE_STR_ISO8601__` -> `"2026-02-08T17:07:11Z"`
* `__TIME__` -> `18:07:11`
* `__TIME_UTC__` -> `17:07:11`
* `__DAY__` -> `8`
* `__MONTH__` -> `2`
* `__YEAR__` -> `2026`
* `__TIMESTAMP_UTC__` -> `176...`
* `__COUNTER__` -> `0`, next call `1`, then `2`, ...
* `__COUNTER_RESET__` -> resets counter to `0`
* `__RAND_INT8__` -> random signed 8-bit integer (`-128..127`)
* `__RAND_INT16__` -> random signed 16-bit integer (`-32768..32767`)
* `__RAND_INT32__` -> random signed 32-bit integer (`-2147483648..2147483647`)
* `__RAND_INT64__` -> random signed 64-bit integer (`int64` range)
* `__RAND_UINT8__` -> random unsigned 8-bit integer (`0..255`)
* `__RAND_UINT16__` -> random unsigned 16-bit integer (`0..65535`)
* `__RAND_UINT32__` -> random unsigned 32-bit integer (`0..4294967295`)
* `__RAND_UINT64__` -> random unsigned 64-bit integer (`uint64` range)

Counter example:

```cpp
a = __COUNTER__;      // 0
b = __COUNTER__;      // 1
__COUNTER_RESET__
c = __COUNTER__;      // 0
```

## References

* <https://community.bistudio.com/wiki/PreProcessor_Commands>
