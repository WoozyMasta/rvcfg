// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

// ProcessAndParseResult stores preprocess and parse stage outputs.
type ProcessAndParseResult struct {
	// Preprocess is preprocessor stage output.
	Preprocess PreprocessResult `json:"preprocess,omitzero" yaml:"preprocess,omitempty"`

	// Parse is parser stage output.
	Parse ParseResult `json:"parse,omitzero" yaml:"parse,omitempty"`
}

// ProcessAndParseFile runs processed mode pipeline: preprocess + parse.
//
// This is a convenience wrapper around:
//  1. PreprocessFile(path, preprocessOptions)
//  2. ParseBytes(path, []byte(preprocessed.Text), parseOptions)
func ProcessAndParseFile(path string, preprocessOptions PreprocessOptions, parseOptions ParseOptions) (ProcessAndParseResult, error) {
	// Processed pipeline always requires source map for diagnostic/AST position remap.
	preprocessOptions.TrackSourceMap = true

	preprocessed, preprocessErr := PreprocessFile(path, preprocessOptions)
	result := ProcessAndParseResult{
		Preprocess: preprocessed,
	}
	if preprocessErr != nil {
		return result, preprocessErr
	}

	parsed, parseErr := ParseBytes(path, []byte(preprocessed.Text), parseOptions)
	resolver := newSourceMapResolver(preprocessed.SourceMap)
	resolver.remapDiagnostics(parsed.Diagnostics)
	resolver.remapFile(&parsed.File)
	result.Parse = parsed
	if parseErr != nil {
		return result, parseErr
	}

	return result, nil
}
