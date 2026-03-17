package rvcfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const cfgConvertEnvKey = "CFGCONVERT_EXE"

// requireCfgConvert resolves CfgConvert path or skips integration test.
func requireCfgConvert(t *testing.T) string {
	t.Helper()

	exe, err := resolveCfgConvertExe(".")
	if err != nil {
		t.Skipf("CfgConvert unavailable: %v", err)
	}

	return exe
}

// resolveCfgConvertExe resolves tool path from env and local .env.
func resolveCfgConvertExe(baseDir string) (string, error) {
	if fromEnv := strings.TrimSpace(os.Getenv(cfgConvertEnvKey)); fromEnv != "" {
		if isRegularFile(fromEnv) {
			return fromEnv, nil
		}

		return "", fmt.Errorf("%s points to missing file: %s", cfgConvertEnvKey, fromEnv)
	}

	if fromDotEnv, ok := lookupCfgConvertFromDotEnv(baseDir); ok {
		if isRegularFile(fromDotEnv) {
			return fromDotEnv, nil
		}

		return "", fmt.Errorf(".env %s path is missing: %s", cfgConvertEnvKey, fromDotEnv)
	}

	return "", fmt.Errorf("%s is unset and .env has no valid value", cfgConvertEnvKey)
}

// lookupCfgConvertFromDotEnv reads CFGCONVERT_EXE from local .env file.
func lookupCfgConvertFromDotEnv(baseDir string) (string, bool) {
	path := filepath.Join(baseDir, ".env")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		key, value, found := strings.Cut(trimmed, "=")
		if !found {
			continue
		}

		if strings.TrimSpace(key) != cfgConvertEnvKey {
			continue
		}

		value = strings.TrimSpace(value)
		value = strings.Trim(value, `"'`)
		if value == "" {
			return "", false
		}

		if filepath.IsAbs(value) {
			return value, true
		}

		return filepath.Join(baseDir, value), true
	}

	return "", false
}

// isRegularFile reports whether path exists and is not a directory.
func isRegularFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
