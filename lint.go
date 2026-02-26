// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"strings"
)

// LintOptions configures lint engine behavior.
type LintOptions struct {
	// RootDir is optional filesystem/project root for external rules.
	RootDir string `json:"root_dir,omitempty" yaml:"root_dir,omitempty"`
}

// LintContext stores shared mutable state for lint rules.
type LintContext struct {
	// Values is optional shared key-value storage for cross-rule data.
	Values map[string]any `json:"values,omitempty" yaml:"values,omitempty"`
	// RootDir is copied from LintOptions.
	RootDir string `json:"root_dir,omitempty" yaml:"root_dir,omitempty"`

	// Source is current file source label during CheckFile call.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`
}

// LintIssue is rule-local issue payload converted into Diagnostic.
type LintIssue struct {
	// Code is rule-local issue code. When empty, RuleID fallback is used.
	Code string `json:"code,omitempty" yaml:"code,omitempty"`

	// Message is issue description.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Severity is issue severity. Defaults to warning when empty.
	Severity Severity `json:"severity,omitzero" yaml:"severity,omitempty"`

	// Start is issue start position.
	Start Position `json:"start,omitzero" yaml:"start,omitempty"`

	// End is issue end position.
	End Position `json:"end,omitzero" yaml:"end,omitempty"`
}

// LintEmit appends one rule issue.
type LintEmit func(issue LintIssue)

// LintRule is external rule contract for file-level checks.
type LintRule interface {
	// RuleID returns stable rule code suffix, for example "PATH001".
	RuleID() string

	// CheckFile checks one parsed file and emits zero or more issues.
	CheckFile(ctx *LintContext, file File, emit LintEmit)
}

// LintFinalizeRule is optional hook called once after all files are checked.
type LintFinalizeRule interface {
	// Finalize runs after all CheckFile calls and can emit aggregate issues.
	Finalize(ctx *LintContext, emit LintEmit)
}

// LintResult stores unified diagnostics collected from registered rules.
type LintResult struct {
	// Diagnostics are emitted lint diagnostics converted to shared model.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty" yaml:"diagnostics,omitempty"`
}

// LintEngine is rule registry and execution engine for external lint checks.
type LintEngine struct {
	options LintOptions
	rules   []lintRegisteredRule
}

// lintRegisteredRule stores one registered rule with namespace prefix.
type lintRegisteredRule struct {
	rule   LintRule
	prefix string
}

// NewLintEngine builds empty lint engine with provided options.
func NewLintEngine(options LintOptions) *LintEngine {
	return &LintEngine{
		rules:   make([]lintRegisteredRule, 0),
		options: options,
	}
}

// Register adds one or more rules under external namespace prefix.
//
// Prefix format is strict uppercase letters/numbers/underscore, for example:
//
//	LINT_PATH
//	PACK
//
// Resulting diagnostic code format:
//
//	<prefix>_<ruleCode>
func (e *LintEngine) Register(prefix string, rules ...LintRule) error {
	prefix = strings.TrimSpace(prefix)
	if !isValidDiagnosticToken(prefix) {
		return fmt.Errorf("invalid lint prefix %q", prefix)
	}

	for idx := range rules {
		rule := rules[idx]
		if rule == nil {
			return fmt.Errorf("nil lint rule at index %d", idx)
		}

		ruleID := strings.TrimSpace(rule.RuleID())
		if !isValidDiagnosticToken(ruleID) {
			return fmt.Errorf("invalid lint rule id %q", ruleID)
		}

		e.rules = append(e.rules, lintRegisteredRule{
			prefix: prefix,
			rule:   rule,
		})
	}

	return nil
}

// RunFiles executes all registered rules against provided files.
func (e *LintEngine) RunFiles(files ...File) LintResult {
	if len(e.rules) == 0 || len(files) == 0 {
		return LintResult{}
	}

	context := &LintContext{
		RootDir: e.options.RootDir,
		Values:  make(map[string]any),
	}
	diagnostics := make([]Diagnostic, 0, 32)

	for ridx := range e.rules {
		registry := e.rules[ridx]
		emit := e.emitForRule(&diagnostics, context, registry)

		for fidx := range files {
			file := files[fidx]
			context.Source = file.Source
			registry.rule.CheckFile(context, file, emit)
		}

		finalizer, ok := registry.rule.(LintFinalizeRule)
		if !ok {
			continue
		}

		context.Source = ""
		finalizer.Finalize(context, emit)
	}

	return LintResult{
		Diagnostics: diagnostics,
	}
}

// emitForRule builds diagnostic-emitter closure for one registered rule.
func (e *LintEngine) emitForRule(target *[]Diagnostic, context *LintContext, registry lintRegisteredRule) LintEmit {
	return func(issue LintIssue) {
		localCode := strings.TrimSpace(issue.Code)
		if localCode == "" {
			localCode = strings.TrimSpace(registry.rule.RuleID())
		}

		if !isValidDiagnosticToken(localCode) {
			localCode = "INVALID"
		}

		fullCode := joinDiagnosticCode(registry.prefix, localCode)
		severity := issue.Severity
		if severity == "" {
			severity = SeverityWarning
		}

		start := issue.Start
		if start.File == "" {
			start.File = context.Source
		}

		end := issue.End
		if isZeroPosition(end) {
			end = start
		}

		*target = append(*target, Diagnostic{
			Code:     DiagnosticCode(fullCode),
			Message:  issue.Message,
			Severity: severity,
			Start:    start,
			End:      end,
		})
	}
}

// joinDiagnosticCode builds prefixed diagnostic code for external rules.
func joinDiagnosticCode(prefix string, code string) string {
	return prefix + "_" + code
}

// isValidDiagnosticToken validates strict uppercase token used in external diagnostic IDs.
func isValidDiagnosticToken(value string) bool {
	if value == "" {
		return false
	}

	hasLetter := false
	for i := 0; i < len(value); i++ {
		ch := value[i]

		if ch >= 'A' && ch <= 'Z' {
			hasLetter = true

			continue
		}

		if ch >= '0' && ch <= '9' {
			continue
		}

		if ch == '_' {
			continue
		}

		return false
	}

	return hasLetter
}

// isZeroPosition checks whether position has no source coordinates.
func isZeroPosition(position Position) bool {
	return position.File == "" && position.Line == 0 && position.Column == 0 && position.Offset == 0
}
