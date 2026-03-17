// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"strconv"
	"strings"
)

// expandBuiltInIntrinsics expands stable strict built-ins and optional dynamic intrinsics.
// Replacement uses identifier-token boundaries and skips strings/comments.
func (p *preprocessor) expandBuiltInIntrinsics(line string, filename string, lineNo int) string {
	if !strings.Contains(line, "__") {
		return line
	}

	out := line
	if strings.Contains(out, "__LINE__") {
		out = replaceIdentifierTokens(out, "__LINE__", strconv.Itoa(maxInt(0, lineNo-1)))
	}

	if strings.Contains(out, "__FILE__") {
		out = replaceIdentifierTokens(out, "__FILE__", quoteIntrinsicString(filename))
	}

	if p.enableFileName {
		needFileName := strings.Contains(out, "__FILE_NAME__")
		needFileShort := strings.Contains(out, "__FILE_SHORT__")
		if needFileName || needFileShort {
			fileName := filepathBase(filename)
			if needFileName {
				out = replaceIdentifierTokens(out, "__FILE_NAME__", quoteIntrinsicString(fileName))
			}

			if needFileShort {
				out = replaceIdentifierTokens(out, "__FILE_SHORT__", quoteIntrinsicString(trimExt(fileName)))
			}
		}
	}

	if !p.enableDynamic {
		return out
	}

	if !strings.Contains(out, "__") {
		return out
	}

	return p.expandDynamicIntrinsics(out)
}

// filepathBase returns last path segment for slash/backslash separators.
func filepathBase(path string) string {
	idx := strings.LastIndexAny(path, `/\`)
	if idx < 0 || idx+1 >= len(path) {
		return path
	}

	return path[idx+1:]
}

// trimExt trims file extension suffix from name.
func trimExt(name string) string {
	idx := strings.LastIndexByte(name, '.')
	if idx <= 0 {
		return name
	}

	return name[:idx]
}

// maxInt returns maximum from two ints.
func maxInt(a int, b int) int {
	if a > b {
		return a
	}

	return b
}
