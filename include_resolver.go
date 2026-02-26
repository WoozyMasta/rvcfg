// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/rvcfg

package rvcfg

import (
	"fmt"
	"path/filepath"
)

// IncludeResolver resolves include targets for preprocess stage.
type IncludeResolver interface {
	// Resolve maps include path from current source file into concrete readable path.
	Resolve(currentFile string, includePath string, includeDirs []string) (string, error)
}

// defaultIncludeResolver resolves include paths using local filesystem.
type defaultIncludeResolver struct{}

// Resolve maps include path to absolute local filesystem path.
func (r defaultIncludeResolver) Resolve(currentFile string, includePath string, includeDirs []string) (string, error) {
	candidates := make([]string, 0, 1+len(includeDirs))
	candidates = append(candidates, filepath.Join(filepath.Dir(currentFile), includePath))

	for _, dir := range includeDirs {
		candidates = append(candidates, filepath.Join(dir, includePath))
	}

	for _, candidate := range candidates {
		if !fileExists(candidate) {
			continue
		}

		absPath, err := filepath.Abs(candidate)
		if err != nil {
			return "", fmt.Errorf("resolve include path %q: %w", candidate, err)
		}

		return absPath, nil
	}

	return "", fmt.Errorf("%w: %s", ErrIncludeNotFound, includePath)
}
