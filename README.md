# rvcfg

Go library for preprocessing, parsing, querying, and formatting  
DayZ / Real Virtuality config-like text sources
(`config.cpp`, `model.cfg`, `.rvmat`-like syntax blocks).

* Full preprocess pipeline:
  `#include`, `#define`, function-like macros, conditionals.
* Lexer + parser with diagnostics, stable codes, and source positions.
* AST with classes, properties, arrays, extern/delete, and enum declarations.
* Deterministic formatter with configurable indentation and soft array wrapping.
* Query helpers for walking statements/classes and building external lint rules.
* Built-in lint diagnostics and rule catalog [RULES.md](RULES.md).

## Install

```bash
go get github.com/woozymasta/rvcfg
```

## Usage

Raw parse (no preprocess):

```go
result, err := rvcfg.ParseFile("config.cpp", rvcfg.ParseOptions{
  CaptureScalarRaw: true,
})
if err != nil {
  // handle
}

_ = result.File
_ = result.Diagnostics
```

Processed parse (preprocess + parse):

```go
processed, err := rvcfg.ProcessAndParseFile(
  "config.cpp",
  rvcfg.PreprocessOptions{
    Defines: map[string]string{
      "SOME_FLAG": "1",
    },
    IncludeDirs: []string{"./include"},
  },
  rvcfg.ParseOptions{
    CaptureScalarRaw: true,
  },
)
if err != nil {
  // handle
}

_ = processed.Preprocess.Text
_ = processed.Preprocess.SourceMap
_ = processed.Parse.File
```

Format source:

```go
out, err := rvcfg.FormatWithOptions(input, rvcfg.FormatOptions{
  IndentChar:             " ",
  IndentSize:             2,
  MaxLineWidth:           120,
  MaxInlineArrayElements: 8,
})
```

Render AST:

```go
rendered, err := rvcfg.RenderFileWithOptions(file, rvcfg.FormatOptions{
  IndentChar:             "\t",
  IndentSize:             1,
  MaxLineWidth:           120,
  MaxInlineArrayElements: 0,
  PreserveComments:       true,
})
```

Walk AST for custom lint rules:

```go
parsed, _ := rvcfg.ParseFile("config.cpp", rvcfg.ParseOptions{})
parsed.File.WalkStatements(func(ref rvcfg.StatementRef) bool {
  _ = ref.PathString()
  _ = ref.Start
  _ = ref.Statement

  return true
})
```

## Options overview

Preprocess:

```go
rvcfg.PreprocessOptions{
  Mode:            rvcfg.PreprocessModeStrict, // strict|compat|extended
  Defines:         map[string]string{"DEBUG": "1"},
  IncludeDirs:     []string{"./include"},
  IncludeResolver: nil, // optional custom resolver
  // strict default keeps DayZ-like behavior:
  EnableIfExpressions:      false, // #if/#elif disabled in strict
  EnableExecEvalIntrinsics: false, // __EXEC/__EVAL opt-in
  EnableDynamicIntrinsics:  false, // date/time/counter/rand opt-in
  EnableFileNameIntrinsics: false, // __FILE_NAME__/__FILE_SHORT__ opt-in
  ExtendedFSRoot:           "",    // optional root for __FILES_* in extended
  ExtendedFSMaxItems:       512,   // safety cap for __FILES_* matches
  ExtendedLoopMaxItems:     2048,  // safety cap for __FOR_RANGE_RENDER
}
```

`PreprocessModeExtended` enables additional deterministic helper intrinsics
for path normalization, string transforms, file-list rendering, and
range/list templating.

See [PREPROCESSOR.md](PREPROCESSOR.md) for the complete intrinsic list,
syntax, limits, and compatibility details.

Parity note:

* strict mode targets semantic parity by default;
* exact `CfgConvert -pcpp` text formatting is not a strict-mode target.

Parser:

```go
rvcfg.ParseOptions{
  CaptureScalarRaw:             true,
  Strict:                       false,
  DisableRecovery:              false,
  PreserveComments:             false, // keep statement-level comments in AST
  AutoFixMissingClassSemicolon: false, // opt-in compatibility mode
}
```

Formatter:

```go
rvcfg.FormatOptions{
  IndentChar:               " ",  // or "\t"
  IndentSize:               2,
  MaxLineWidth:             120,  // <=0 disables width wrapping
  MaxInlineArrayElements:   0,    // <=0 disables count wrapping
  ArrayWrapByName: map[string]int{
    "SkeletonBones": 2, // group elements in wrapped rows
  },
  PreserveBlankLines:       nil,  // nil -> keep up to 1, &0 -> disable
  PreserveComments:         false, // keep statement-level comments
  DisableCompactEmptyClass: false,
}
```

AST renderer:

```go
rvcfg.FormatOptions{
  IndentChar:               " ",  // or "\t"
  IndentSize:               2,
  MaxLineWidth:             120,  // <=0 disables width wrapping
  MaxInlineArrayElements:   0,    // <=0 disables count wrapping
  ArrayWrapByName: map[string]int{
    "SkeletonBones": 2, // group elements in wrapped rows
  },
  PreserveBlankLines:       nil,  // nil -> keep up to 1, &0 -> disable
  PreserveComments:         true, // keep statement-level comments
  DisableCompactEmptyClass: false,
}
```

## Diagnostics

Diagnostics are stable and machine-readable:

* code (`RVCFG2001`, `RVCFG3015`, ...)
* severity (`error`/`warning`)
* exact source position

Catalog API:

```go
all := rvcfg.DiagnosticCatalog()
spec, ok := rvcfg.DiagnosticByCode(rvcfg.CodeParUnexpectedToken)
_, _, _ = all, spec, ok
```

Lint rules are available as a machine-readable snapshot in
[rules.yaml](rules.yaml), and detailed rule documentation is provided in
[RULES.md](RULES.md).

## lintkit integration

`rvcfg` exposes parser and preprocess diagnostics as `lintkit` rules.

```go
sourceBytes, err := os.ReadFile("config.cpp")
if err != nil {
  return err
}

parsed, err := rvcfg.ParseFile("config.cpp", rvcfg.ParseOptions{})
if err != nil {
  return err
}

optional := rvcfg.AnalyzeFile(
  parsed.File,
  sourceBytes, // original source bytes (or preprocessed text bytes)
  rvcfg.AnalyzeOptions{},
)

allDiagnostics := append(parsed.Diagnostics, optional...)

engine := linting.NewEngine()
if err := rvcfg.RegisterLintRules(engine); err != nil {
  return err
}

// same provider form:
// _ = lint.RegisterRuleProviders(engine, rvcfg.LintRulesProvider{})

runCtx := lint.RunContext{
  TargetPath: "config.cpp",
  TargetKind: "rvcfg.config",
}
rvcfg.AttachLintDiagnostics(&runCtx, allDiagnostics)

result, err := engine.Run(context.Background(), runCtx, nil)
if err != nil {
  return err
}
_ = result
```

Exported short code format uses catalog prefix `RVCFG` + numeric code.
Rule IDs use semantic form `rvcfg.<stage>.<description-slug>`.
