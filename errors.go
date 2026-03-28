// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"

	"github.com/woozymasta/lintkit/lint"
)

var (
	// ErrIncludeNotFound indicates that #include target file was not resolved.
	ErrIncludeNotFound = errors.New("include file not found")

	// ErrUnsupportedDirective indicates unsupported preprocessor directive in v0.
	ErrUnsupportedDirective = errors.New("unsupported preprocessor directive")

	// ErrUnsupportedIntrinsic indicates unsupported config intrinsic like __EXEC/__EVAL.
	ErrUnsupportedIntrinsic = errors.New("unsupported config intrinsic")

	// ErrInvalidDirective indicates malformed preprocessor directive syntax.
	ErrInvalidDirective = errors.New("invalid preprocessor directive")

	// ErrMacroExpand indicates macro expansion failure or recursion overflow.
	ErrMacroExpand = errors.New("macro expansion failed")

	// ErrUnresolvedMacroInvocation indicates unresolved macro-like invocation after preprocess.
	ErrUnresolvedMacroInvocation = errors.New("unresolved macro-like invocation")

	// ErrLex indicates lexical scan failure.
	ErrLex = errors.New("lexical analysis failed")

	// ErrParse indicates syntax parse failure.
	ErrParse = errors.New("syntax parse failed")

	// ErrNilLintRuleRegistrar indicates nil lint rule registrar.
	ErrNilLintRuleRegistrar = lint.ErrNilRuleRegistrar
)
