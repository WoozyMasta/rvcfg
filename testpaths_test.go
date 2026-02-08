package rvcfg

import "path/filepath"

// testDataPath builds repository-relative testdata fixture path.
func testDataPath(parts ...string) string {
	all := append([]string{"testdata"}, parts...)

	return filepath.Join(all...)
}
