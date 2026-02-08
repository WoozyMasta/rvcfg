// Package rvcfg provides frontend primitives for DayZ/Real Virtuality config text.
//
// Main capabilities:
//   - lexical analysis for config-like sources
//   - preprocess pipeline (#include, #define, conditionals, token-paste, stringify)
//   - parser/AST for class/property/array/extern/delete/enum declarations
//   - deterministic formatter with configurable indentation/wrapping
//
// Pipeline modes:
//   - Raw mode: ParseFile/ParseBytes parse source without preprocess stage.
//   - Processed mode: PreprocessFile resolves includes and expands macros.
//   - Orchestrated mode: ProcessAndParseFile runs preprocess + parse helper flow.
//
// Raw mode example:
//
//	parsed, err := ParseFile("config.cpp", ParseOptions{
//		CaptureScalarRaw: true,
//	})
//	if err != nil {
//		// handle parse error
//	}
//
//	_ = parsed.File
//
// Processed mode example:
//
//	pre, err := PreprocessFile("config.cpp", PreprocessOptions{
//		IncludeDirs: []string{"./include"},
//		Defines:     map[string]string{"DEBUG": "1"},
//		// opt-in compatibility modes:
//		EnableExecEvalIntrinsics: false,
//		EnableDynamicIntrinsics:  false,
//	})
//	if err != nil {
//		// handle preprocess error
//	}
//
//	parsed, err := ParseBytes("config.cpp", []byte(pre.Text), ParseOptions{
//		CaptureScalarRaw: true,
//		PreserveComments: false,
//	})
//	if err != nil {
//		// handle parse error on expanded text
//	}
//
//	_ = pre.SourceMap
//	_ = parsed.File
//
// Formatter example:
//
//	out, err := FormatWithOptions([]byte("class A{v=1;};"), FormatOptions{
//		IndentChar:             " ",
//		IndentSize:             2,
//		MaxLineWidth:           120,
//		MaxInlineArrayElements: 8,
//		// nil => keep up to one blank line; explicit 0 => disable.
//		PreserveBlankLines: nil,
//		PreserveComments:   false,
//	})
//	if err != nil {
//		// handle format error
//	}
//
//	_ = out
//
// Compatibility notes:
//   - ParseOptions.AutoFixMissingClassSemicolon is opt-in and disabled by default.
//   - __EXEC/__EVAL and dynamic intrinsics are opt-in preprocess modes.
//   - __has_include is explicitly unsupported in v0 and reported as preprocess error.
//
// AST/query example:
//
//	parsed.File.WalkStatements(func(ref StatementRef) bool {
//		// ref.Start/ref.End and ref.PathString() are ready for external lint diagnostics.
//		return true
//	})
package rvcfg
