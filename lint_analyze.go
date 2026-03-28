// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	// rapFloatPrecisionAbsThreshold is absolute drift threshold for float32 loss warning.
	rapFloatPrecisionAbsThreshold = 1e-6

	// rapFloatPrecisionRelThreshold is relative drift threshold for float32 loss warning.
	rapFloatPrecisionRelThreshold = 1e-8

	// rapExtremeFloatUpperThreshold is magnitude threshold for extreme float info.
	rapExtremeFloatUpperThreshold = 1e20

	// rapExtremeFloatLowerThreshold is minimum non-zero magnitude threshold for extreme float info.
	rapExtremeFloatLowerThreshold = 1e-20
)

const (
	// rapScalarKindUnknown marks unsupported scalar shape for RAP subtype mapping.
	rapScalarKindUnknown rapScalarKind = iota

	// rapScalarKindString marks quoted scalar string.
	rapScalarKindString

	// rapScalarKindFloat marks float scalar supported by RAP float32 subtype.
	rapScalarKindFloat

	// rapScalarKindInt32 marks integer scalar supported by RAP int32 subtype.
	rapScalarKindInt32

	// rapScalarKindInt64 marks integer scalar supported by RAP int64 subtype.
	rapScalarKindInt64

	// rapScalarKindIdentifier marks identifier-like scalar encoded as RAP string subtype.
	rapScalarKindIdentifier
)

// AnalyzeOptions configures optional lint passes over parsed AST.
type AnalyzeOptions struct {
	// DisableInheritanceHints disables PAR026 inheritance risk hints.
	DisableInheritanceHints bool `json:"disable_inheritance_hints,omitempty" yaml:"disable_inheritance_hints,omitempty"`

	// DisableRAPScalarHints disables RAP scalar compatibility lint pass.
	DisableRAPScalarHints bool `json:"disable_rap_scalar_hints,omitempty" yaml:"disable_rap_scalar_hints,omitempty"`
}

// rapScalarKind stores scalar classification group used by RAP lint pass.
type rapScalarKind uint8

// rapScalarClass stores RAP scalar classification result.
type rapScalarClass struct {
	// Kind stores scalar compatibility group.
	Kind rapScalarKind

	// Float64 stores parsed float value for float diagnostics.
	Float64 float64
}

// rapScalarLintContext stores mutable state for RAP scalar lint traversal.
type rapScalarLintContext struct {
	// source stores original source bytes for scalar raw extraction.
	source []byte

	// diagnostics accumulates emitted RAP parse diagnostics.
	diagnostics []Diagnostic
}

// AnalyzeFile runs optional lint passes on parsed file and returns diagnostics.
func AnalyzeFile(file File, source []byte, opts AnalyzeOptions) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, 16)

	if !opts.DisableInheritanceHints {
		diagnostics = append(diagnostics, collectInheritanceHints(file)...)
	}

	if !opts.DisableRAPScalarHints {
		diagnostics = append(diagnostics, collectRAPScalarHints(file, source)...)
	}

	return diagnostics
}

// collectInheritanceHints emits info diagnostics for potentially risky inheritance overrides.
func collectInheritanceHints(file File) []Diagnostic {
	diagnostics := make([]Diagnostic, 0)
	walkDerivedNestedClassWithoutBase(file.Statements, nil, &diagnostics)

	return diagnostics
}

// walkDerivedNestedClassWithoutBase recursively walks class statements.
func walkDerivedNestedClassWithoutBase(
	statements []Statement,
	parentClass *ClassDecl,
	diagnostics *[]Diagnostic,
) {
	for index := range statements {
		statement := &statements[index]
		if statement.Kind != NodeClass ||
			statement.Class == nil ||
			statement.Class.Forward {
			continue
		}

		if parentClass != nil &&
			parentClass.Base != "" &&
			statement.Class.Base == "" {
			*diagnostics = append(*diagnostics, Diagnostic{
				Code: CodeParDerivedNestedClassWithoutBase,
				Message: fmt.Sprintf(
					"nested class %q in derived class %q has no explicit inheritance and may replace parent subtree",
					statement.Class.Name,
					parentClass.Name,
				),
				Severity: SeverityInfo,
				Start:    statement.Start,
				End:      statement.End,
			})
		}

		walkDerivedNestedClassWithoutBase(
			statement.Class.Body,
			statement.Class,
			diagnostics,
		)
	}
}

// collectRAPScalarHints emits RAP-focused diagnostics for parsed scalar values.
func collectRAPScalarHints(file File, source []byte) []Diagnostic {
	context := rapScalarLintContext{
		source:      source,
		diagnostics: make([]Diagnostic, 0),
	}

	context.walkStatements(file.Statements)

	return context.diagnostics
}

// walkStatements walks statements recursively and inspects scalar values.
func (context *rapScalarLintContext) walkStatements(statements []Statement) {
	for index := range statements {
		statement := statements[index]

		switch statement.Kind {
		case NodeClass:
			if statement.Class == nil {
				continue
			}

			context.walkStatements(statement.Class.Body)

		case NodeProperty:
			if statement.Property == nil {
				continue
			}

			context.inspectValue(statement.Property.Value)

		case NodeArrayAssign:
			if statement.ArrayAssign == nil {
				continue
			}

			context.inspectValue(statement.ArrayAssign.Value)
		}
	}
}

// inspectValue checks one value and recurses into nested arrays.
func (context *rapScalarLintContext) inspectValue(value Value) {
	switch value.Kind {
	case ValueScalar:
		context.inspectScalar(value)

	case ValueArray:
		for index := range value.Elements {
			context.inspectValue(value.Elements[index])
		}
	}
}

// inspectScalar emits RAP scalar diagnostics for one scalar value.
func (context *rapScalarLintContext) inspectScalar(value Value) {
	raw := context.scalarRaw(value)
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return
	}

	if isNonFiniteFloatLiteral(trimmed) {
		context.emitDiagnostic(
			CodeParRAPNonFiniteFloat,
			SeverityWarning,
			value,
			fmt.Sprintf(
				"non-finite float literal %q is unsafe for RAP numeric encoding",
				trimmed,
			),
		)

		return
	}

	if hasUnsafeBackslashQuoteEscape(trimmed) {
		context.emitDiagnostic(
			CodeParRAPUnsafeStringEscape,
			SeverityWarning,
			value,
			"string contains legacy \\\" escape that BI/CfgConvert may reject",
		)
	}

	if isFloat32OverflowLiteral(trimmed) {
		context.emitDiagnostic(
			CodeParRAPFloatOverflowToInf,
			SeverityWarning,
			value,
			fmt.Sprintf(
				"float %q overflows to Inf after RAP float32 conversion",
				trimmed,
			),
		)

		return
	}

	class := classifyRAPScalarTrimmed(trimmed)
	if class.Kind == rapScalarKindUnknown {
		context.emitDiagnostic(
			CodeParRAPUnsupportedScalar,
			SeverityWarning,
			value,
			fmt.Sprintf("scalar %q cannot be mapped to RAP scalar subtype safely", trimmed),
		)

		return
	}

	if class.Kind != rapScalarKindFloat {
		return
	}

	if math.IsNaN(class.Float64) || math.IsInf(class.Float64, 0) {
		return
	}

	float32Value := float32(class.Float64)
	roundTrip := float64(float32Value)

	if class.Float64 != 0 && float32Value == 0 {
		context.emitDiagnostic(
			CodeParRAPFloatUnderflowToZero,
			SeverityWarning,
			value,
			fmt.Sprintf("float %q collapses to 0 after RAP float32 conversion", trimmed),
		)
	}

	loss := math.Abs(class.Float64 - roundTrip)
	if loss > rapFloatPrecisionAbsThreshold &&
		loss > math.Abs(class.Float64)*rapFloatPrecisionRelThreshold {
		context.emitDiagnostic(
			CodeParRAPFloatPrecisionLoss,
			SeverityWarning,
			value,
			fmt.Sprintf(
				"float %q becomes %g after RAP float32 conversion (abs drift=%g)",
				trimmed,
				roundTrip,
				loss,
			),
		)
	}

	if isExtremeFloatLiteral(trimmed, class.Float64) {
		context.emitDiagnostic(
			CodeParRAPExtremeFloatMagnitude,
			SeverityInfo,
			value,
			fmt.Sprintf("float %q uses extreme exponent/magnitude and may normalize heavily", trimmed),
		)
	}
}

// scalarRaw resolves scalar raw text from captured parser value or source offsets.
func (context *rapScalarLintContext) scalarRaw(value Value) string {
	if value.Raw != "" {
		return value.Raw
	}

	start := value.Start.Offset
	end := value.End.Offset
	if start < 0 || end < start || end >= len(context.source) {
		return ""
	}

	return string(context.source[start : end+1])
}

// emitDiagnostic appends one diagnostic with severity from code catalog.
func (context *rapScalarLintContext) emitDiagnostic(
	code Code,
	severity Severity,
	value Value,
	message string,
) {
	context.diagnostics = append(context.diagnostics, Diagnostic{
		Code:     code,
		Message:  message,
		Severity: severity,
		Start:    value.Start,
		End:      value.End,
	})
}

// classifyRAPScalarTrimmed maps scalar raw text to RAP-supported scalar families.
func classifyRAPScalarTrimmed(raw string) rapScalarClass {
	if raw == "" {
		return rapScalarClass{
			Kind: rapScalarKindUnknown,
		}
	}

	if isRVCfgString(raw) {
		return rapScalarClass{
			Kind: rapScalarKindString,
		}
	}

	if strings.HasPrefix(raw, `@"`) && strings.HasSuffix(raw, `"`) {
		if isRVCfgString(strings.TrimPrefix(raw, "@")) {
			return rapScalarClass{
				Kind: rapScalarKindString,
			}
		}
	}

	intValue, err := strconv.ParseInt(raw, 10, 64)
	if err == nil {
		if intValue >= math.MinInt32 && intValue <= math.MaxInt32 {
			return rapScalarClass{
				Kind: rapScalarKindInt32,
			}
		}

		return rapScalarClass{
			Kind: rapScalarKindInt64,
		}
	}

	if looksFloatRaw(raw) {
		_, err32 := strconv.ParseFloat(raw, 32)
		if err32 == nil {
			floatValue64, err64 := strconv.ParseFloat(raw, 64)
			if err64 == nil {
				return rapScalarClass{
					Kind:    rapScalarKindFloat,
					Float64: floatValue64,
				}
			}

			floatValue32, _ := strconv.ParseFloat(raw, 32)
			return rapScalarClass{
				Kind:    rapScalarKindFloat,
				Float64: floatValue32,
			}
		}
	}

	if isIdentifierLike(raw) {
		return rapScalarClass{
			Kind: rapScalarKindIdentifier,
		}
	}

	return rapScalarClass{
		Kind: rapScalarKindUnknown,
	}
}

// hasUnsafeBackslashQuoteEscape checks for C-style \" escapes unsafe for BI parser.
func hasUnsafeBackslashQuoteEscape(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, `@"`) {
		trimmed = strings.TrimPrefix(trimmed, "@")
	}

	if len(trimmed) < 2 || trimmed[0] != '"' || trimmed[len(trimmed)-1] != '"' {
		return false
	}

	body := trimmed[1 : len(trimmed)-1]
	for index := 0; index+1 < len(body); index++ {
		if body[index] != '\\' || body[index+1] != '"' {
			continue
		}

		if index+2 < len(body) && body[index+2] == '"' {
			continue
		}

		return true
	}

	return false
}

// isRVCfgString validates BI-style quoted string scalar syntax.
func isRVCfgString(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) < 2 || trimmed[0] != '"' || trimmed[len(trimmed)-1] != '"' {
		return false
	}

	body := trimmed[1 : len(trimmed)-1]
	if body == "" {
		return true
	}

	for index := 0; index < len(body); index++ {
		char := body[index]

		if char == '"' {
			if index+1 < len(body) && body[index+1] == '"' {
				index++
				continue
			}

			return false
		}

		if char == '\\' && index+1 < len(body) && body[index+1] == '"' &&
			(index+2 >= len(body) || body[index+2] != '"') {
			index++
			continue
		}
	}

	return true
}

// looksFloatRaw checks whether scalar text is likely intended as float syntax.
func looksFloatRaw(raw string) bool {
	return strings.Contains(raw, ".") || strings.Contains(raw, "e") || strings.Contains(raw, "E")
}

// isIdentifierLike reports whether scalar can be encoded as identifier-like string.
func isIdentifierLike(raw string) bool {
	if raw == "" {
		return false
	}

	for index := 0; index < len(raw); index++ {
		char := raw[index]
		if char >= 'a' && char <= 'z' {
			continue
		}

		if char >= 'A' && char <= 'Z' {
			continue
		}

		if char >= '0' && char <= '9' {
			continue
		}

		if char == '_' || char == '.' {
			continue
		}

		return false
	}

	return true
}

// isExtremeFloatLiteral checks for exponent form with extreme absolute magnitude.
func isExtremeFloatLiteral(raw string, value float64) bool {
	if !strings.Contains(raw, "e") && !strings.Contains(raw, "E") {
		return false
	}

	abs := math.Abs(value)
	if abs >= rapExtremeFloatUpperThreshold {
		return true
	}

	return abs > 0 && abs <= rapExtremeFloatLowerThreshold
}

// isFloat32OverflowLiteral checks whether float literal exceeds float32 range.
func isFloat32OverflowLiteral(raw string) bool {
	if !looksFloatRaw(raw) {
		return false
	}

	value32, err32 := strconv.ParseFloat(raw, 32)
	if err32 == nil {
		return false
	}

	var numErr *strconv.NumError
	if !errors.As(err32, &numErr) || numErr.Err != strconv.ErrRange {
		return false
	}

	return math.IsInf(value32, 0)
}

// isNonFiniteFloatLiteral checks whether scalar literal is NaN/Inf-like token.
func isNonFiniteFloatLiteral(raw string) bool {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return false
	}

	return math.IsNaN(value) || math.IsInf(value, 0)
}
