package rvcfg

import (
	"path/filepath"
	"testing"
)

func TestFormatGolden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		input  string
		golden string
		opts   FormatOptions
	}{
		{
			name:   "basic",
			input:  filepath.Join("testdata", "format", "basic.input.cpp"),
			golden: filepath.Join("testdata", "format", "basic.golden.cpp"),
		},
		{
			name:   "wrap_by_element_limit",
			input:  filepath.Join("testdata", "format", "wrap.input.cpp"),
			golden: filepath.Join("testdata", "format", "wrap.golden.cpp"),
			opts: FormatOptions{
				MaxLineWidth:           120,
				MaxInlineArrayElements: 8,
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			input := readTestFile(t, tc.input)
			want := readTestFile(t, tc.golden)
			got, err := FormatWithOptions(input, tc.opts)
			if err != nil {
				t.Fatalf("FormatWithOptions(%s) error: %v", tc.input, err)
			}

			if string(got) != string(want) {
				t.Fatalf("golden mismatch for %s\nwant:\n%s\ngot:\n%s", tc.name, string(want), string(got))
			}
		})
	}
}
