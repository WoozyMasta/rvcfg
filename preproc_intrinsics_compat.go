// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// expandExecIntrinsics executes __EXEC(...) assignments.
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

// expandEvalIntrinsics evaluates __EVAL(...) expressions.
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

// executeExecBody parses and executes semicolon-separated __EXEC statements.
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

// evaluateEvalExpression evaluates expression and formats scalar/string result.
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

// evalIntrinsicExpr evaluates intrinsic arithmetic/string expression.
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

// splitExecStatements splits __EXEC body by top-level semicolons.
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

// isIntrinsicIdent checks whether name is valid __EXEC variable identifier.
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

// expandDynamicIntrinsics expands non-deterministic optional intrinsics.
func (p *preprocessor) expandDynamicIntrinsics(input string) string {
	if input == "" {
		return input
	}

	var out strings.Builder
	lastWrite := 0
	replaced := false

	inString := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(input); {
		if inLineComment {
			if input[i] == '\n' {
				inLineComment = false
			}

			i++

			continue
		}

		if inBlockComment {
			if input[i] == '*' && i+1 < len(input) && input[i+1] == '/' {
				inBlockComment = false
				i += 2

				continue
			}

			i++

			continue
		}

		if inString {
			if input[i] == '"' {
				inString = false
			}

			i++

			continue
		}

		if input[i] == '"' {
			inString = true
			i++

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '/' {
			inLineComment = true
			i += 2

			continue
		}

		if input[i] == '/' && i+1 < len(input) && input[i+1] == '*' {
			inBlockComment = true
			i += 2

			continue
		}

		if !isIdentifierPart(input[i]) {
			i++

			continue
		}

		start := i
		for i < len(input) && isIdentifierPart(input[i]) {
			i++
		}

		token := input[start:i]
		replacement, ok := p.dynamicIntrinsicReplacement(token)
		if !ok {
			continue
		}

		if !replaced {
			out.Grow(len(input))
			replaced = true
		}

		out.WriteString(input[lastWrite:start])
		out.WriteString(replacement)
		lastWrite = i
	}

	if !replaced {
		return input
	}

	out.WriteString(input[lastWrite:])

	return out.String()
}

// dynamicIntrinsicReplacement evaluates dynamic intrinsic token replacement.
func (p *preprocessor) dynamicIntrinsicReplacement(token string) (string, bool) {
	localNow := p.dynamicNow
	utcNow := p.dynamicNow.UTC()

	switch token {
	case "__DATE_ARR__":
		return strconv.Itoa(localNow.Year()) +
			"," + strconv.Itoa(int(localNow.Month())) +
			"," + strconv.Itoa(localNow.Day()) +
			"," + strconv.Itoa(localNow.Hour()) +
			"," + strconv.Itoa(localNow.Minute()) +
			"," + strconv.Itoa(localNow.Second()), true
	case "__DATE_STR__":
		return quoteIntrinsicString(localNow.Format("2006/01/02, 15:04:05")), true
	case "__DATE_STR_ISO8601__":
		return quoteIntrinsicString(utcNow.Format(time.RFC3339)), true
	case "__TIME__":
		return localNow.Format("15:04:05"), true
	case "__TIME_UTC__":
		return utcNow.Format("15:04:05"), true
	case "__DAY__":
		return strconv.Itoa(localNow.Day()), true
	case "__MONTH__":
		return strconv.Itoa(int(localNow.Month())), true
	case "__YEAR__":
		return strconv.Itoa(localNow.Year()), true
	case "__TIMESTAMP_UTC__":
		return strconv.FormatInt(utcNow.Unix(), 10), true
	case "__COUNTER__":
		value := p.counter
		p.counter++

		return strconv.FormatUint(value, 10), true
	case "__COUNTER_RESET__":
		p.counter = 0

		return "", true
	}

	if strings.HasPrefix(token, "__RAND_INT") && strings.HasSuffix(token, "__") {
		bits, ok := parseRandBits(token, "__RAND_INT", "__")
		if !ok {
			return "", false
		}

		value, ok := p.randomInt(bits)
		if !ok {
			return "", false
		}

		return strconv.FormatInt(value, 10), true
	}

	if strings.HasPrefix(token, "__RAND_UINT") && strings.HasSuffix(token, "__") {
		bits, ok := parseRandBits(token, "__RAND_UINT", "__")
		if !ok {
			return "", false
		}

		value, ok := p.randomUint(bits)
		if !ok {
			return "", false
		}

		return strconv.FormatUint(value, 10), true
	}

	return "", false
}

// parseRandBits parses bit-width suffix for random intrinsics.
func parseRandBits(token string, prefix string, suffix string) (int, bool) {
	if !strings.HasPrefix(token, prefix) || !strings.HasSuffix(token, suffix) {
		return 0, false
	}

	raw := strings.TrimSuffix(strings.TrimPrefix(token, prefix), suffix)
	if raw == "" {
		return 0, false
	}

	bits, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}

	switch bits {
	case 8, 16, 32, 64:
		return bits, true
	default:
		return 0, false
	}
}

// randomUint returns cryptographically secure unsigned random value for bits.
func (p *preprocessor) randomUint(bits int) (uint64, bool) {
	shift, ok := randBitShift(bits)
	if !ok {
		return 0, false
	}

	upperBound := new(big.Int).Lsh(big.NewInt(1), shift)
	n, err := rand.Int(rand.Reader, upperBound)
	if err != nil {
		return 0, false
	}

	return n.Uint64(), true
}

// randomInt returns cryptographically secure signed random value for bits.
func (p *preprocessor) randomInt(bits int) (int64, bool) {
	shift, ok := randBitShift(bits)
	if !ok {
		return 0, false
	}

	span := new(big.Int).Lsh(big.NewInt(1), shift)
	n, err := rand.Int(rand.Reader, span)
	if err != nil {
		return 0, false
	}

	offset := new(big.Int).Lsh(big.NewInt(1), shift-1)
	n.Sub(n, offset)
	if !n.IsInt64() {
		return 0, false
	}

	return n.Int64(), true
}

// randBitShift validates bit-width for random intrinsics.
func randBitShift(bits int) (uint, bool) {
	switch bits {
	case 8:
		return 8, true
	case 16:
		return 16, true
	case 32:
		return 32, true
	case 64:
		return 64, true
	default:
		return 0, false
	}
}

// intrinsicExprParser parses __EVAL expression into intrinsicValue.
type intrinsicExprParser struct {
	vars  map[string]intrinsicValue
	input string
	pos   int
}

// newIntrinsicExprParser creates parser for one expression.
func newIntrinsicExprParser(input string, vars map[string]intrinsicValue) *intrinsicExprParser {
	return &intrinsicExprParser{
		input: input,
		pos:   0,
		vars:  vars,
	}
}

// parseExpr parses expression root.
func (p *intrinsicExprParser) parseExpr() (intrinsicValue, error) {
	return p.parseAddSub()
}

// parseAddSub parses + and - precedence layer.
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

// parseMulDiv parses * and / precedence layer.
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

// parseUnary parses unary +/- operators.
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

// parsePrimary parses literals, identifiers, and parenthesized expressions.
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

// parseString parses quoted string literal without unescaping semantics.
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

// parseNumber parses integer/float numeric literal.
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

// parseIdent parses identifier token.
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

// skipSpaces advances parser over ASCII whitespace.
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

// match consumes ch if present at current parser position.
func (p *intrinsicExprParser) match(ch byte) bool {
	if p.eof() || p.input[p.pos] != ch {
		return false
	}

	p.pos++

	return true
}

// peek returns current byte without consuming.
func (p *intrinsicExprParser) peek() byte {
	return p.input[p.pos]
}

// eof reports parser end of input.
func (p *intrinsicExprParser) eof() bool {
	return p.pos >= len(p.input)
}

// addIntrinsicValues applies + semantics for numeric/string values.
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

// intrinsicValueToString formats intrinsic value as string.
func intrinsicValueToString(value intrinsicValue) string {
	if value.isString {
		return value.text
	}

	if isIntegralFloat(value.number) {
		return strconv.FormatInt(int64(value.number), 10)
	}

	return strconv.FormatFloat(value.number, 'f', -1, 64)
}

// isDigitByte reports whether byte is ASCII decimal digit.
func isDigitByte(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// isIntrinsicIdentStart reports whether byte can start intrinsic identifier.
func isIntrinsicIdentStart(ch byte) bool {
	return ch == '_' || ch == '$' || unicode.IsLetter(rune(ch))
}

// isIntrinsicIdentPart reports whether byte can continue intrinsic identifier.
func isIntrinsicIdentPart(ch byte) bool {
	return isIntrinsicIdentStart(ch) || unicode.IsDigit(rune(ch))
}
