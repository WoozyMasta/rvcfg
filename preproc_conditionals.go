package rvcfg

import (
	"strconv"
	"strings"
)

func containsHasIncludeIntrinsic(expr string) bool {
	return strings.Contains(expr, "__has_include")
}

// evalIfExpression evaluates minimal #if expressions for v0.
func (p *preprocessor) evalIfExpression(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return false
	}

	if parts := splitTopLevelBy(expr, "||"); len(parts) > 1 {
		for _, part := range parts {
			if p.evalIfExpression(part) {
				return true
			}
		}

		return false
	}

	if parts := splitTopLevelBy(expr, "&&"); len(parts) > 1 {
		for _, part := range parts {
			if !p.evalIfExpression(part) {
				return false
			}
		}

		return true
	}

	expr = trimOuterParens(expr)
	if strings.HasPrefix(expr, "!") {
		return !p.evalIfExpression(strings.TrimSpace(expr[1:]))
	}

	if idx, op, ok := findTopLevelComparisonOp(expr); ok {
		leftRaw := strings.TrimSpace(expr[:idx])
		rightRaw := strings.TrimSpace(expr[idx+len(op):])
		left, leftOK := p.evalIfNumericValue(leftRaw)
		right, rightOK := p.evalIfNumericValue(rightRaw)
		if !leftOK || !rightOK {
			return false
		}

		switch op {
		case "==":
			return left == right
		case "!=":
			return left != right
		case ">=":
			return left >= right
		case "<=":
			return left <= right
		case ">":
			return left > right
		case "<":
			return left < right
		default:
			return false
		}
	}

	if value, ok := p.evalIfNumericValue(expr); ok {
		return value != 0
	}

	if p.macroExists(expr) {
		return true
	}

	return false
}

func (p *preprocessor) evalIfNumericValue(expr string) (int64, bool) {
	expr = trimOuterParens(strings.TrimSpace(expr))
	if expr == "" {
		return 0, false
	}

	if strings.HasPrefix(expr, "defined") {
		rest := strings.TrimSpace(strings.TrimPrefix(expr, "defined"))
		if strings.HasPrefix(rest, "(") && strings.HasSuffix(rest, ")") {
			name := strings.TrimSpace(rest[1 : len(rest)-1])
			if p.macroExists(name) {
				return 1, true
			}

			return 0, true
		}

		if p.macroExists(rest) {
			return 1, true
		}

		return 0, true
	}

	if iv, err := strconv.ParseInt(expr, 0, 64); err == nil {
		return iv, true
	}

	if !p.macroExists(expr) {
		return 0, false
	}

	body := strings.TrimSpace(p.macros[expr].Body)
	if body == "" {
		return 1, true
	}

	iv, err := strconv.ParseInt(body, 0, 64)
	if err != nil {
		return 0, false
	}

	return iv, true
}

func splitTopLevelBy(expr string, sep string) []string {
	if expr == "" || sep == "" {
		return []string{expr}
	}

	parts := make([]string, 0, 2)
	start := 0
	depth := 0
	inString := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
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

		if depth == 0 && i+len(sep) <= len(expr) && expr[i:i+len(sep)] == sep {
			parts = append(parts, strings.TrimSpace(expr[start:i]))
			i += len(sep) - 1
			start = i + 1
		}
	}

	if len(parts) == 0 {
		return []string{expr}
	}

	parts = append(parts, strings.TrimSpace(expr[start:]))

	return parts
}

func findTopLevelComparisonOp(expr string) (int, string, bool) {
	ops := []string{"==", "!=", ">=", "<=", ">", "<"}
	depth := 0
	inString := false

	for i := 0; i < len(expr); i++ {
		ch := expr[i]
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

		if depth != 0 {
			continue
		}

		for _, op := range ops {
			if i+len(op) <= len(expr) && expr[i:i+len(op)] == op {
				return i, op, true
			}
		}
	}

	return 0, "", false
}

func trimOuterParens(expr string) string {
	out := strings.TrimSpace(expr)
	for strings.HasPrefix(out, "(") && strings.HasSuffix(out, ")") {
		depth := 0
		valid := true

		for i := 0; i < len(out); i++ {
			ch := out[i]
			if ch == '(' {
				depth++

				continue
			}

			if ch == ')' {
				depth--
				if depth == 0 && i != len(out)-1 {
					valid = false

					break
				}
			}
		}

		if !valid || depth != 0 {
			break
		}

		out = strings.TrimSpace(out[1 : len(out)-1])
	}

	return out
}

// macroExists checks macro table by name.
func (p *preprocessor) macroExists(name string) bool {
	_, ok := p.macros[name]

	return ok
}

// framesActive checks whether all conditional frames allow output.
func framesActive(frames []conditionalFrame) bool {
	for _, frame := range frames {
		if !frame.Active {
			return false
		}
	}

	return true
}

// pushConditional appends new conditional frame.
func pushConditional(frames *[]conditionalFrame, condition bool) {
	parent := framesActive(*frames)
	frame := conditionalFrame{
		ParentActive: parent,
		Active:       parent && condition,
		BranchTaken:  condition,
	}

	*frames = append(*frames, frame)
}

// updateElif switches frame branch by elif condition.
func updateElif(frames *[]conditionalFrame, condition bool) bool {
	if len(*frames) == 0 {
		return false
	}

	last := &(*frames)[len(*frames)-1]
	if !last.ParentActive {
		last.Active = false

		return true
	}

	if last.BranchTaken {
		last.Active = false

		return true
	}

	last.Active = condition
	if condition {
		last.BranchTaken = true
	}

	return true
}

// updateElse switches frame into else branch.
func updateElse(frames *[]conditionalFrame) bool {
	if len(*frames) == 0 {
		return false
	}

	last := &(*frames)[len(*frames)-1]
	if !last.ParentActive {
		last.Active = false

		return true
	}

	last.Active = !last.BranchTaken
	last.BranchTaken = true

	return true
}

// popConditional closes top conditional frame.
func popConditional(frames *[]conditionalFrame) bool {
	if len(*frames) == 0 {
		return false
	}

	*frames = (*frames)[:len(*frames)-1]

	return true
}
