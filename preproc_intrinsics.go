// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type intrinsicValue struct {
	text     string
	number   float64
	isString bool
}

func (p *preprocessor) expandExecEvalIntrinsics(line string) (string, error) {
	out := line

	for pass := 0; pass < p.maxExpandDepth; pass++ {
		changed := false

		next, changedExec, err := p.expandExecIntrinsics(out)
		if err != nil {
			return out, err
		}

		next, changedEval, err := p.expandEvalIntrinsics(next)
		if err != nil {
			return out, err
		}

		out = next
		changed = changedExec || changedEval
		if !changed {
			return strings.TrimSpace(out), nil
		}
	}

	return out, errors.New("__EXEC/__EVAL expansion depth overflow")
}

func (p *preprocessor) expandExecIntrinsics(input string) (string, bool, error) {
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, body, end, ok, err := findIntrinsicCall(out, "__EXEC", searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		if err := p.executeExecBody(body); err != nil {
			return out, changedAny, err
		}

		out = out[:start] + out[end:]
		searchFrom = start
		changedAny = true
	}
}

func (p *preprocessor) expandEvalIntrinsics(input string) (string, bool, error) {
	changedAny := false
	out := input
	searchFrom := 0

	for {
		start, expr, end, ok, err := findIntrinsicCall(out, "__EVAL", searchFrom)
		if err != nil {
			return out, changedAny, err
		}

		if !ok {
			return out, changedAny, nil
		}

		evalOut, evalErr := p.evaluateEvalExpression(expr)
		if evalErr != nil {
			// DayZ CfgConvert fallback observed in compatibility probe.
			evalOut = `"scalar"`
		}

		out = out[:start] + evalOut + out[end:]
		searchFrom = start + len(evalOut)
		changedAny = true
	}
}

func (p *preprocessor) executeExecBody(body string) error {
	statements := splitExecStatements(body)
	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		eq := strings.Index(statement, "=")
		if eq <= 0 {
			return fmt.Errorf("__EXEC statement must be assignment, got %q", statement)
		}

		name := strings.TrimSpace(statement[:eq])
		if !isIntrinsicIdent(name) {
			return fmt.Errorf("__EXEC invalid variable name %q", name)
		}

		expr := strings.TrimSpace(statement[eq+1:])
		if expr == "" {
			return fmt.Errorf("__EXEC assignment missing expression for %q", name)
		}

		value, err := p.evalIntrinsicExpr(expr)
		if err != nil {
			return err
		}

		p.execEvalVars[name] = value
	}

	return nil
}

func (p *preprocessor) evaluateEvalExpression(expr string) (string, error) {
	value, err := p.evalIntrinsicExpr(expr)
	if err != nil {
		return "", err
	}

	if value.isString {
		return quoteIntrinsicString(value.text), nil
	}

	if isIntegralFloat(value.number) {
		return strconv.FormatInt(int64(value.number), 10), nil
	}

	return strconv.FormatFloat(value.number, 'f', -1, 64), nil
}

func (p *preprocessor) evalIntrinsicExpr(expr string) (intrinsicValue, error) {
	parser := newIntrinsicExprParser(expr, p.execEvalVars)
	value, err := parser.parseExpr()
	if err != nil {
		return intrinsicValue{}, err
	}

	parser.skipSpaces()
	if !parser.eof() {
		return intrinsicValue{}, fmt.Errorf("unexpected token at pos=%d", parser.pos)
	}

	return value, nil
}

func splitExecStatements(input string) []string {
	out := make([]string, 0, 4)
	start := 0
	depth := 0
	inString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '"' {
				inString = false
			}

			continue
		}

		if ch == '"' {
			inString = true

			continue
		}

		if ch == '(' {
			depth++

			continue
		}

		if ch == ')' {
			if depth > 0 {
				depth--
			}

			continue
		}

		if ch == ';' && depth == 0 {
			out = append(out, strings.TrimSpace(input[start:i]))
			start = i + 1
		}
	}

	if start < len(input) {
		out = append(out, strings.TrimSpace(input[start:]))
	}

	return out
}

func findIntrinsicCall(input string, name string, from int) (int, string, int, bool, error) {
	if from < 0 {
		from = 0
	}

	for i := from; i < len(input); i++ {
		if !hasIdentifierAt(input, i, name) {
			continue
		}

		open := i + len(name)
		if open >= len(input) || input[open] != '(' {
			continue
		}

		body, end, err := parseIntrinsicCallBody(input, open)
		if err != nil {
			return 0, "", 0, false, err
		}

		return i, body, end, true, nil
	}

	return 0, "", 0, false, nil
}

func parseIntrinsicCallBody(input string, open int) (string, int, error) {
	if open >= len(input) || input[open] != '(' {
		return "", 0, errors.New("intrinsic call parse without opening parenthesis")
	}

	start := open + 1
	depth := 1
	inString := false

	for i := start; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '"' {
				inString = false
			}

			continue
		}

		if ch == '"' {
			inString = true

			continue
		}

		if ch == '(' {
			depth++

			continue
		}

		if ch == ')' {
			depth--
			if depth == 0 {
				return strings.TrimSpace(input[start:i]), i + 1, nil
			}
		}
	}

	return "", 0, errors.New("unterminated intrinsic call")
}

func quoteIntrinsicString(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
}

func isIntegralFloat(v float64) bool {
	const eps = 1e-9
	iv := float64(int64(v))
	if v >= 0 {
		return v-iv < eps
	}

	return iv-v < eps
}

func isIntrinsicIdent(name string) bool {
	if name == "" {
		return false
	}

	for i := 0; i < len(name); i++ {
		ch := rune(name[i])
		if i == 0 {
			if ch != '_' && ch != '$' && !unicode.IsLetter(ch) {
				return false
			}

			continue
		}

		if ch != '_' && ch != '$' && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}

	return true
}

type intrinsicExprParser struct {
	vars  map[string]intrinsicValue
	input string
	pos   int
}

func newIntrinsicExprParser(input string, vars map[string]intrinsicValue) *intrinsicExprParser {
	return &intrinsicExprParser{
		input: input,
		pos:   0,
		vars:  vars,
	}
}

func (p *intrinsicExprParser) parseExpr() (intrinsicValue, error) {
	return p.parseAddSub()
}

func (p *intrinsicExprParser) parseAddSub() (intrinsicValue, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return intrinsicValue{}, err
	}

	for {
		p.skipSpaces()
		if p.match('+') {
			right, rightErr := p.parseMulDiv()
			if rightErr != nil {
				return intrinsicValue{}, rightErr
			}

			left = addIntrinsicValues(left, right)
			continue
		}

		if p.match('-') {
			right, rightErr := p.parseMulDiv()
			if rightErr != nil {
				return intrinsicValue{}, rightErr
			}

			if left.isString || right.isString {
				return intrinsicValue{}, errors.New("operator '-' is not valid for string values")
			}

			left = intrinsicValue{number: left.number - right.number}
			continue
		}

		return left, nil
	}
}

func (p *intrinsicExprParser) parseMulDiv() (intrinsicValue, error) {
	left, err := p.parseUnary()
	if err != nil {
		return intrinsicValue{}, err
	}

	for {
		p.skipSpaces()
		if p.match('*') {
			right, rightErr := p.parseUnary()
			if rightErr != nil {
				return intrinsicValue{}, rightErr
			}

			if left.isString || right.isString {
				return intrinsicValue{}, errors.New("operator '*' is not valid for string values")
			}

			left = intrinsicValue{number: left.number * right.number}
			continue
		}

		if p.match('/') {
			right, rightErr := p.parseUnary()
			if rightErr != nil {
				return intrinsicValue{}, rightErr
			}

			if left.isString || right.isString {
				return intrinsicValue{}, errors.New("operator '/' is not valid for string values")
			}

			left = intrinsicValue{number: left.number / right.number}
			continue
		}

		return left, nil
	}
}

func (p *intrinsicExprParser) parseUnary() (intrinsicValue, error) {
	p.skipSpaces()
	if p.match('+') {
		return p.parseUnary()
	}

	if p.match('-') {
		value, err := p.parseUnary()
		if err != nil {
			return intrinsicValue{}, err
		}

		if value.isString {
			return intrinsicValue{}, errors.New("unary '-' is not valid for string values")
		}

		return intrinsicValue{number: -value.number}, nil
	}

	return p.parsePrimary()
}

func (p *intrinsicExprParser) parsePrimary() (intrinsicValue, error) {
	p.skipSpaces()
	if p.eof() {
		return intrinsicValue{}, errors.New("expected expression value")
	}

	if p.match('(') {
		value, err := p.parseExpr()
		if err != nil {
			return intrinsicValue{}, err
		}

		p.skipSpaces()
		if !p.match(')') {
			return intrinsicValue{}, errors.New("expected ')'")
		}

		return value, nil
	}

	if p.peek() == '"' {
		return p.parseString()
	}

	if isDigitByte(p.peek()) || p.peek() == '.' {
		return p.parseNumber()
	}

	if isIntrinsicIdentStart(p.peek()) {
		name := p.parseIdent()
		value, ok := p.vars[name]
		if !ok {
			return intrinsicValue{}, fmt.Errorf("unknown identifier %q", name)
		}

		return value, nil
	}

	return intrinsicValue{}, fmt.Errorf("unexpected token at pos=%d", p.pos)
}

func (p *intrinsicExprParser) parseString() (intrinsicValue, error) {
	if !p.match('"') {
		return intrinsicValue{}, errors.New("expected string literal")
	}

	start := p.pos
	for !p.eof() {
		ch := p.peek()
		if ch == '"' {
			text := p.input[start:p.pos]
			p.pos++

			return intrinsicValue{isString: true, text: text}, nil
		}

		p.pos++
	}

	return intrinsicValue{}, errors.New("unterminated string literal")
}

func (p *intrinsicExprParser) parseNumber() (intrinsicValue, error) {
	start := p.pos
	dotSeen := false

	for !p.eof() {
		ch := p.peek()
		if ch == '.' {
			if dotSeen {
				break
			}

			dotSeen = true
			p.pos++
			continue
		}

		if !isDigitByte(ch) {
			break
		}

		p.pos++
	}

	raw := strings.TrimSpace(p.input[start:p.pos])
	if raw == "" || raw == "." {
		return intrinsicValue{}, errors.New("invalid number literal")
	}

	num, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return intrinsicValue{}, fmt.Errorf("invalid number literal %q", raw)
	}

	return intrinsicValue{number: num}, nil
}

func (p *intrinsicExprParser) parseIdent() string {
	start := p.pos
	p.pos++

	for !p.eof() {
		ch := p.peek()
		if !isIntrinsicIdentPart(ch) {
			break
		}

		p.pos++
	}

	return p.input[start:p.pos]
}

func (p *intrinsicExprParser) skipSpaces() {
	for !p.eof() {
		switch p.input[p.pos] {
		case ' ', '\t', '\r', '\n':
			p.pos++
		default:
			return
		}
	}
}

func (p *intrinsicExprParser) match(ch byte) bool {
	if p.eof() || p.input[p.pos] != ch {
		return false
	}

	p.pos++

	return true
}

func (p *intrinsicExprParser) peek() byte {
	return p.input[p.pos]
}

func (p *intrinsicExprParser) eof() bool {
	return p.pos >= len(p.input)
}

func addIntrinsicValues(left intrinsicValue, right intrinsicValue) intrinsicValue {
	if left.isString || right.isString {
		leftText := intrinsicValueToString(left)
		rightText := intrinsicValueToString(right)

		return intrinsicValue{
			isString: true,
			text:     leftText + rightText,
		}
	}

	return intrinsicValue{
		number: left.number + right.number,
	}
}

func intrinsicValueToString(value intrinsicValue) string {
	if value.isString {
		return value.text
	}

	if isIntegralFloat(value.number) {
		return strconv.FormatInt(int64(value.number), 10)
	}

	return strconv.FormatFloat(value.number, 'f', -1, 64)
}

func isDigitByte(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isIntrinsicIdentStart(ch byte) bool {
	return ch == '_' || ch == '$' || unicode.IsLetter(rune(ch))
}

func isIntrinsicIdentPart(ch byte) bool {
	return isIntrinsicIdentStart(ch) || unicode.IsDigit(rune(ch))
}
