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
  Defines:         map[string]string{"DEBUG": "1"},
  IncludeDirs:     []string{"./include"},
  IncludeResolver: nil, // optional custom resolver
  EnableExecEvalIntrinsics: false, // opt-in compatibility mode
  EnableDynamicIntrinsics:  false, // opt-in date/time/counter/rand intrinsics
}
```

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
  PreserveBlankLines:       nil,  // nil -> keep up to 1, &0 -> disable
  PreserveComments:         false, // keep statement-level comments
  DisableCompactEmptyClass: false,
}
```

## Diagnostics

Diagnostics are stable and machine-readable:

* code (`PAR001`, `PP015`, ...)
* severity (`error`/`warning`)
* exact source position

Catalog API:

```go
all := rvcfg.DiagnosticCatalog()
spec, ok := rvcfg.DiagnosticByCode("PAR001")
```

## Docs

* [MACROS.md](MACROS.md) - macro and preprocessor guide.
* [DIAGNOSTICS.md](DIAGNOSTICS.md) - diagnostic code registry.
