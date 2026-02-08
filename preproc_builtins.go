package rvcfg

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// expandBuiltInIntrinsics expands stable built-ins and optional dynamic intrinsics.
// Replacement uses identifier-token boundaries and skips strings/comments.
func (p *preprocessor) expandBuiltInIntrinsics(line string, filename string, lineNo int) string {
	out := line

	fileName := filepathBase(filename)
	fileShort := trimExt(fileName)

	out = replaceIdentifierTokens(out, "__LINE__", strconv.Itoa(lineNo))
	out = replaceIdentifierTokens(out, "__FILE__", quoteIntrinsicString(filename))
	out = replaceIdentifierTokens(out, "__FILE_NAME__", quoteIntrinsicString(fileName))
	out = replaceIdentifierTokens(out, "__FILE_SHORT__", quoteIntrinsicString(fileShort))

	if !p.enableDynamic {
		return out
	}

	return p.expandDynamicIntrinsics(out)
}

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

func filepathBase(path string) string {
	idx := strings.LastIndexAny(path, `/\`)
	if idx < 0 || idx+1 >= len(path) {
		return path
	}

	return path[idx+1:]
}

func trimExt(name string) string {
	idx := strings.LastIndexByte(name, '.')
	if idx <= 0 {
		return name
	}

	return name[:idx]
}
