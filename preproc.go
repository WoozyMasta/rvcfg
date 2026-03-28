// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultMaxIncludeDepth = 32
	defaultMaxExpandDepth  = 32
)

// PreprocessMode defines high-level preprocessor behavior profile.
type PreprocessMode string

const (
	// PreprocessModeStrict mirrors DayZ CfgConvert behavior as default profile.
	PreprocessModeStrict PreprocessMode = "strict"
	// PreprocessModeCompat enables compatibility helpers on top of strict mode.
	PreprocessModeCompat PreprocessMode = "compat"
	// PreprocessModeExtended enables extended feature profile for custom workflows.
	PreprocessModeExtended PreprocessMode = "extended"
)

// PreprocessOptions configures preprocess behavior.
type PreprocessOptions struct {

	// IncludeResolver resolves include paths. Default resolver uses local filesystem.
	IncludeResolver IncludeResolver `json:"-" yaml:"-"`

	// Defines are additional external symbols.
	Defines map[string]string `json:"defines,omitempty" yaml:"defines,omitempty"`

	// Mode selects high-level preprocessor profile.
	// Empty value defaults to strict behavior.
	Mode PreprocessMode `json:"mode,omitempty" yaml:"mode,omitempty"`

	// ExtendedFSRoot restricts filesystem intrinsics to this root directory.
	// Empty value uses current source file directory.
	ExtendedFSRoot string `json:"extended_fs_root,omitempty" yaml:"extended_fs_root,omitempty"`

	// IncludeDirs are extra include search directories.
	IncludeDirs []string `json:"include_dirs,omitempty" yaml:"include_dirs,omitempty"`

	// MaxIncludeDepth limits nested include recursion.
	MaxIncludeDepth int `json:"max_include_depth,omitempty" yaml:"max_include_depth,omitempty"`

	// MaxExpandDepth limits macro expansion recursion.
	MaxExpandDepth int `json:"max_expand_depth,omitempty" yaml:"max_expand_depth,omitempty"`

	// ExtendedFSMaxItems limits matched file entries for filesystem intrinsics.
	// Non-positive value uses default internal limit.
	ExtendedFSMaxItems int `json:"extended_fs_max_items,omitempty" yaml:"extended_fs_max_items,omitempty"`

	// ExtendedLoopMaxItems limits iteration count for extended loop intrinsics.
	// Non-positive value uses default internal limit.
	ExtendedLoopMaxItems int `json:"extended_loop_max_items,omitempty" yaml:"extended_loop_max_items,omitempty"`

	// EmitIncludeMarkers inserts include boundary markers into output text.
	// Markers are emitted as line comments:
	//   // <include-start "path">
	//   // <include-end "path">
	EmitIncludeMarkers bool `json:"emit_include_markers,omitempty" yaml:"emit_include_markers,omitempty"`

	// TrackSourceMap enables source map generation for output lines.
	TrackSourceMap bool `json:"track_source_map,omitempty" yaml:"track_source_map,omitempty"`

	// EnableExecEvalIntrinsics enables compatibility mode for `__EXEC`/`__EVAL`.
	// Disabled by default to keep preprocess deterministic and side-effect free.
	EnableExecEvalIntrinsics bool `json:"enable_exec_eval_intrinsics,omitempty" yaml:"enable_exec_eval_intrinsics,omitempty"`

	// EnableDynamicIntrinsics enables non-deterministic/date/counter intrinsics:
	// `__DATE_ARR__`, `__DATE_STR__`, `__DATE_STR_ISO8601__`, `__TIME__`, `__TIME_UTC__`,
	// `__DAY__`, `__MONTH__`, `__YEAR__`, `__TIMESTAMP_UTC__`,
	// `__COUNTER__`, `__COUNTER_RESET__`, `__RAND_INT*N*__`, `__RAND_UINT*N*__`.
	//
	// Disabled by default.
	EnableDynamicIntrinsics bool `json:"enable_dynamic_intrinsics,omitempty" yaml:"enable_dynamic_intrinsics,omitempty"`

	// EnableFileNameIntrinsics enables `__FILE_NAME__` and `__FILE_SHORT__`.
	// Disabled by default because DayZ CfgConvert strict mode does not expand them.
	EnableFileNameIntrinsics bool `json:"enable_file_name_intrinsics,omitempty" yaml:"enable_file_name_intrinsics,omitempty"`

	// EnableMacroRedefinitionWarnings emits warnings when macro is redefined.
	// Disabled by default in strict mode for closer DayZ CfgConvert parity.
	EnableMacroRedefinitionWarnings bool `json:"enable_macro_redefinition_warnings,omitempty" yaml:"enable_macro_redefinition_warnings,omitempty"`

	// EnableIfExpressions enables `#if`/`#elif` expression evaluation.
	// Disabled by default in strict mode because DayZ CfgConvert rejects
	// `#if`/`#elif` directives.
	EnableIfExpressions bool `json:"enable_if_expressions,omitempty" yaml:"enable_if_expressions,omitempty"`
}

// PreprocessResult is preprocessor output bundle.
type PreprocessResult struct {
	// Text is preprocessed source text.
	Text string `json:"text,omitempty" yaml:"text,omitempty"`

	// Diagnostics are emitted warnings and errors.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty" yaml:"diagnostics,omitempty"`

	// Includes are resolved include files.
	Includes []string `json:"includes,omitempty" yaml:"includes,omitempty"`

	// SourceMap maps output line ranges to origin source lines.
	SourceMap []SourceMapEntry `json:"source_map,omitempty" yaml:"source_map,omitempty"`
}

// SourceMapEntry maps output line range to source location range.
type SourceMapEntry struct {
	// Kind describes segment type: "source", "include-start", "include-end".
	Kind string `json:"kind,omitempty" yaml:"kind,omitempty"`

	// SourceFile is origin source file for this line segment.
	SourceFile string `json:"source_file,omitempty" yaml:"source_file,omitempty"`

	// IncludeFile stores include target for include boundary markers.
	IncludeFile string `json:"include_file,omitempty" yaml:"include_file,omitempty"`

	// SourceStartLine is origin start line (1-based).
	SourceStartLine int `json:"source_start_line,omitempty" yaml:"source_start_line,omitempty"`

	// SourceEndLine is origin end line (1-based).
	SourceEndLine int `json:"source_end_line,omitempty" yaml:"source_end_line,omitempty"`

	// OutputStartLine is output start line (1-based).
	OutputStartLine int `json:"output_start_line,omitempty" yaml:"output_start_line,omitempty"`

	// OutputEndLine is output end line (1-based).
	OutputEndLine int `json:"output_end_line,omitempty" yaml:"output_end_line,omitempty"`

	// SourceStartColumn is origin start column (1-based, approximate for expanded content).
	SourceStartColumn int `json:"source_start_column,omitempty" yaml:"source_start_column,omitempty"`

	// SourceEndColumn is origin end column (1-based, approximate for expanded content).
	SourceEndColumn int `json:"source_end_column,omitempty" yaml:"source_end_column,omitempty"`

	// OutputStartColumn is output start column (1-based).
	OutputStartColumn int `json:"output_start_column,omitempty" yaml:"output_start_column,omitempty"`

	// OutputEndColumn is output end column (1-based).
	OutputEndColumn int `json:"output_end_column,omitempty" yaml:"output_end_column,omitempty"`
}

// mappedLine stores output line with origin metadata.
type mappedLine struct {
	kind       string
	text       string
	sourceFile string
	include    string
	sourceLine int
}

// logicalLine stores preprocessed input line with original source line number.
type logicalLine struct {
	// Text is logical source line after line-continuation merge.
	Text string

	// SourceLine is original physical source line number (1-based).
	SourceLine int
}

// PreprocessFile resolves includes and expands macros for file input.
func PreprocessFile(path string, opts PreprocessOptions) (PreprocessResult, error) {
	pp := newPreprocessor(opts)
	lines, err := pp.processFile(path, path, 0)
	text := joinMappedLines(lines)
	sourceMap := buildSourceMap(lines, pp.trackSourceMap)
	if err != nil {
		return PreprocessResult{
			Text:        text,
			Diagnostics: pp.diagnostics,
			Includes:    pp.includes,
			SourceMap:   sourceMap,
		}, err
	}

	return PreprocessResult{
		Text:        text,
		Diagnostics: pp.diagnostics,
		Includes:    pp.includes,
		SourceMap:   sourceMap,
	}, nil
}

// macroDefinition stores parsed macro metadata.
type macroDefinition struct {
	// Name is macro name.
	Name string

	// Body is macro replacement text.
	Body string

	// Params is function-like macro parameter list.
	Params []string

	// FunctionLike indicates macro expects argument list.
	FunctionLike bool
}

// conditionalFrame stores state for nested condition blocks.
type conditionalFrame struct {
	// ParentActive stores effective state of upper frames.
	ParentActive bool

	// Active marks whether current branch is active.
	Active bool

	// BranchTaken tracks if any previous branch was selected.
	BranchTaken bool
}

// preprocessor is mutable execution state.
type preprocessor struct {
	dynamicNow           time.Time
	includeResolver      IncludeResolver
	macros               map[string]macroDefinition
	fileStack            map[string]bool
	macroRedefWarnedV0   map[string]struct{}
	execEvalVars         map[string]intrinsicValue
	extendedFSRoot       string
	diagnostics          []Diagnostic
	includes             []string
	includeDirs          []string
	objectMacroNamesV0   []string
	functionMacroNamesV0 []string
	maxIncludeDepth      int
	maxExpandDepth       int
	counter              uint64
	extendedFSMaxItems   int
	extendedLoopMaxItems int
	macroNamesDirty      bool
	emitIncludeMarkers   bool
	trackSourceMap       bool
	enableExecEval       bool
	enableDynamic        bool
	enableFileName       bool
	enableMacroRedefWarn bool
	enableIfExpr         bool
	enableExtended       bool
}

// newPreprocessor initializes mutable preprocess state.
func newPreprocessor(opts PreprocessOptions) *preprocessor {
	if opts.MaxIncludeDepth <= 0 {
		opts.MaxIncludeDepth = defaultMaxIncludeDepth
	}

	if opts.MaxExpandDepth <= 0 {
		opts.MaxExpandDepth = defaultMaxExpandDepth
	}

	resolver := opts.IncludeResolver
	if resolver == nil {
		resolver = defaultIncludeResolver{}
	}

	mode := opts.Mode
	if mode == "" {
		mode = PreprocessModeStrict
	}

	switch mode {
	case PreprocessModeCompat, PreprocessModeExtended:
		opts.EnableExecEvalIntrinsics = true
		opts.EnableDynamicIntrinsics = true
		opts.EnableFileNameIntrinsics = true
		opts.EnableMacroRedefinitionWarnings = true
		opts.EnableIfExpressions = true
	}

	if opts.ExtendedFSMaxItems <= 0 {
		opts.ExtendedFSMaxItems = 512
	}

	if opts.ExtendedLoopMaxItems <= 0 {
		opts.ExtendedLoopMaxItems = 2048
	}

	pp := &preprocessor{
		macros:               make(map[string]macroDefinition),
		fileStack:            make(map[string]bool),
		includeDirs:          opts.IncludeDirs,
		includeResolver:      resolver,
		macroRedefWarnedV0:   make(map[string]struct{}),
		macroNamesDirty:      true,
		emitIncludeMarkers:   opts.EmitIncludeMarkers,
		trackSourceMap:       opts.TrackSourceMap,
		enableExecEval:       opts.EnableExecEvalIntrinsics,
		enableDynamic:        opts.EnableDynamicIntrinsics,
		enableFileName:       opts.EnableFileNameIntrinsics,
		enableMacroRedefWarn: opts.EnableMacroRedefinitionWarnings,
		enableIfExpr:         opts.EnableIfExpressions,
		enableExtended:       mode == PreprocessModeExtended,
		extendedFSRoot:       opts.ExtendedFSRoot,
		extendedFSMaxItems:   opts.ExtendedFSMaxItems,
		extendedLoopMaxItems: opts.ExtendedLoopMaxItems,
		dynamicNow:           time.Now(),
		execEvalVars:         make(map[string]intrinsicValue, 32),
		maxIncludeDepth:      opts.MaxIncludeDepth,
		maxExpandDepth:       opts.MaxExpandDepth,
	}

	for key, value := range opts.Defines {
		pp.macros[key] = macroDefinition{
			Name: key,
			Body: value,
		}
	}

	return pp
}

// processFile preprocesses file with include recursion handling.
func (p *preprocessor) processFile(path string, displayPath string, depth int) ([]mappedLine, error) {
	if depth > p.maxIncludeDepth {
		return nil, fmt.Errorf("%w: include depth exceeded at %q", ErrIncludeNotFound, path)
	}

	if strings.TrimSpace(displayPath) == "" {
		displayPath = path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve absolute path for %q: %w", path, err)
	}

	if p.fileStack[absPath] {
		return nil, fmt.Errorf("%w: include cycle detected for %q", ErrIncludeNotFound, absPath)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		p.emitError(CodePPIncludeNotFound, absPath, 1, "include target not found: "+path)

		return nil, fmt.Errorf("%w: %q", ErrIncludeNotFound, path)
	}

	p.fileStack[absPath] = true
	defer delete(p.fileStack, absPath)

	text := normalizeLineEndings(string(data))
	out, err := p.processText(absPath, displayPath, text, depth)
	if err != nil {
		return out, err
	}

	return out, nil
}

// processText applies directives and macro expansion to text.
func (p *preprocessor) processText(filename string, displayFilename string, text string, depth int) ([]mappedLine, error) {
	lines := mergeLineContinuationsWithSourceLines(text)
	out := make([]mappedLine, 0, len(lines))
	frames := make([]conditionalFrame, 0, 8)
	inBlockComment := false

	for _, sourceLine := range lines {
		lineNo := sourceLine.SourceLine
		line := sourceLine.Text

		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			includeLines, err := p.handleDirective(filename, lineNo, trimmed, depth, &frames)
			if err != nil {
				return out, err
			}

			if len(includeLines) > 0 {
				out = append(out, includeLines...)
			}

			continue
		}

		if !framesActive(frames) {
			continue
		}

		expanded, err := p.expandLine(line)
		if err != nil {
			p.emitError(CodePPMacroExpand, filename, lineNo, err.Error())

			return out, err
		}

		expanded = p.expandBuiltInIntrinsics(expanded, displayFilename, lineNo)

		if p.shouldExpandCompatIntrinsics(expanded) {
			expanded, err = p.expandExecEvalIntrinsics(expanded, filename)
			if err != nil {
				p.emitError(CodePPUnsupportedIntrinsic, filename, lineNo, err.Error())

				return out, ErrUnsupportedIntrinsic
			}
		}

		expanded, inBlockComment = stripComments(expanded, inBlockComment)
		if strings.TrimSpace(expanded) == "" {
			continue
		}

		out = append(out, mappedLine{
			kind:       "source",
			text:       expanded,
			sourceFile: displayFilename,
			sourceLine: lineNo,
		})
	}

	if len(frames) > 0 {
		lastLine := 1
		if len(lines) > 0 {
			lastLine = lines[len(lines)-1].SourceLine
		}

		p.emitError(CodePPUnterminatedConditional, filename, lastLine, "unterminated conditional block")

		return out, ErrInvalidDirective
	}

	return out, nil
}

// handleDirective routes preprocessor directives.
func (p *preprocessor) handleDirective(
	filename string,
	lineNo int,
	line string,
	depth int,
	frames *[]conditionalFrame,
) ([]mappedLine, error) {
	name, arg := splitDirective(line)
	active := framesActive(*frames)
	conditional := isConditionalDirective(name)

	if !active && !conditional {
		return nil, nil
	}

	switch name {
	case "include":
		includePath, tail, err := parseIncludePathWithTail(arg)
		if err != nil {
			p.emitError(CodePPInvalidIncludeSyntax, filename, lineNo, err.Error())

			return nil, fmt.Errorf("%w: %s", ErrInvalidDirective, err)
		}

		resolved, err := p.resolveInclude(filename, includePath)
		if err != nil {
			p.emitError(CodePPIncludeNotFound, filename, lineNo, err.Error())

			return nil, err
		}

		p.includes = append(p.includes, resolved)
		lines, err := p.processFile(resolved, includePath, depth+1)
		if err != nil {
			return nil, err
		}

		tailText, err := p.expandDirectiveTailText(filename, lineNo, tail)
		if err != nil {
			return nil, err
		}

		if tailText != "" {
			lines = append(lines, mappedLine{
				kind:       "source",
				text:       tailText,
				sourceFile: filename,
				sourceLine: lineNo,
			})
		}

		if !p.emitIncludeMarkers {
			return lines, nil
		}

		out := make([]mappedLine, 0, len(lines)+2)
		out = append(out, mappedLine{
			kind:       "include-start",
			text:       fmt.Sprintf(`// <include-start "%s">`, resolved),
			sourceFile: filename,
			sourceLine: lineNo,
			include:    resolved,
		})
		out = append(out, lines...)
		out = append(out, mappedLine{
			kind:       "include-end",
			text:       fmt.Sprintf(`// <include-end "%s">`, resolved),
			sourceFile: filename,
			sourceLine: lineNo,
			include:    resolved,
		})

		return out, nil

	case "define":
		if err := p.defineMacro(arg, filename, lineNo); err != nil {
			return nil, err
		}

		return nil, nil

	case "undef":
		macroName, tail := splitDirectiveHeadTail(arg)
		if macroName == "" {
			p.emitError(CodePPMissingMacroName, filename, lineNo, "missing macro name in #undef")

			return nil, ErrInvalidDirective
		}

		delete(p.macros, macroName)
		p.macroNamesDirty = true

		return p.emitDirectiveTail(filename, lineNo, tail)

	case "ifdef":
		macroName, tail := splitDirectiveHeadTail(arg)
		if macroName == "" {
			p.emitError(CodePPMissingMacroName, filename, lineNo, "missing macro name in #ifdef")

			return nil, ErrInvalidDirective
		}

		exists := p.macroExists(macroName)
		pushConditional(frames, exists)

		if !framesActive(*frames) {
			return nil, nil
		}

		return p.emitDirectiveTail(filename, lineNo, tail)

	case "ifndef":
		macroName, tail := splitDirectiveHeadTail(arg)
		if macroName == "" {
			p.emitError(CodePPMissingMacroName, filename, lineNo, "missing macro name in #ifndef")

			return nil, ErrInvalidDirective
		}

		exists := p.macroExists(macroName)
		pushConditional(frames, !exists)

		if !framesActive(*frames) {
			return nil, nil
		}

		return p.emitDirectiveTail(filename, lineNo, tail)

	case "if":
		if !p.enableIfExpr {
			p.emitError(CodePPUnsupportedDirective, filename, lineNo, "unsupported directive #if")

			return nil, ErrUnsupportedDirective
		}

		if containsHasIncludeIntrinsic(arg) {
			p.emitError(CodePPUnsupportedHasInclude, filename, lineNo, "__has_include is unsupported in v0")

			return nil, ErrUnsupportedIntrinsic
		}

		cond := p.evalIfExpression(arg)
		pushConditional(frames, cond)

		return nil, nil

	case "elif":
		if !p.enableIfExpr {
			p.emitError(CodePPUnsupportedDirective, filename, lineNo, "unsupported directive #elif")

			return nil, ErrUnsupportedDirective
		}

		if containsHasIncludeIntrinsic(arg) {
			p.emitError(CodePPUnsupportedHasInclude, filename, lineNo, "__has_include is unsupported in v0")

			return nil, ErrUnsupportedIntrinsic
		}

		ok := updateElif(frames, p.evalIfExpression(arg))
		if !ok {
			p.emitError(CodePPUnexpectedElif, filename, lineNo, "unexpected #elif")

			return nil, ErrInvalidDirective
		}

		return nil, nil

	case "else":
		ok := updateElse(frames)
		if !ok {
			p.emitError(CodePPUnexpectedElse, filename, lineNo, "unexpected #else")

			return nil, ErrInvalidDirective
		}

		if !framesActive(*frames) {
			return nil, nil
		}

		return p.emitDirectiveTail(filename, lineNo, arg)

	case "endif":
		ok := popConditional(frames)
		if !ok {
			p.emitError(CodePPUnexpectedEndif, filename, lineNo, "unexpected #endif")

			return nil, ErrInvalidDirective
		}

		if !framesActive(*frames) {
			return nil, nil
		}

		return p.emitDirectiveTail(filename, lineNo, arg)

	case "error":
		msg := strings.TrimSpace(arg)
		if msg == "" {
			msg = "#error"
		}

		p.emitError(CodePPDirectiveError, filename, lineNo, msg)

		return nil, ErrInvalidDirective
	}

	p.emitError(CodePPUnsupportedDirective, filename, lineNo, "unsupported directive #"+name)

	return nil, ErrUnsupportedDirective
}

// resolveInclude resolves include path from current file and include dirs.
func (p *preprocessor) resolveInclude(currentFile string, includePath string) (string, error) {
	return p.includeResolver.Resolve(currentFile, includePath, p.includeDirs)
}

// emitDirectiveTail expands and emits extra text that follows directive arguments.
func (p *preprocessor) emitDirectiveTail(filename string, lineNo int, tail string) ([]mappedLine, error) {
	tailText, err := p.expandDirectiveTailText(filename, lineNo, tail)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(tailText) == "" {
		return nil, nil
	}

	return []mappedLine{
		{
			kind:       "source",
			text:       tailText,
			sourceFile: filename,
			sourceLine: lineNo,
		},
	}, nil
}

// expandDirectiveTailText expands tail text that follows directive arguments.
func (p *preprocessor) expandDirectiveTailText(filename string, lineNo int, tail string) (string, error) {
	if tail == "" {
		return "", nil
	}

	expanded, err := p.expandLine(tail)
	if err != nil {
		p.emitError(CodePPMacroExpand, filename, lineNo, err.Error())

		return "", err
	}

	return expanded, nil
}

// shouldExpandCompatIntrinsics checks if line contains enabled compat/extended intrinsics.
func (p *preprocessor) shouldExpandCompatIntrinsics(line string) bool {
	if p.enableExecEval && (strings.Contains(line, "__EXEC") || strings.Contains(line, "__EVAL")) {
		return true
	}

	if !p.enableExtended {
		return false
	}

	return strings.Contains(line, "__PATH_NORM") ||
		strings.Contains(line, "__STR_") ||
		strings.Contains(line, "__FILES_") ||
		strings.Contains(line, "__FOR_")
}
