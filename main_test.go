package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestEnvv(t *testing.T) {
	tests := []struct {
		environ []string
		files   map[string]string
		want    []string
	}{
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{},
			want:    []string{"foo=bar"},
		},
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{"baz": "baz"},
			want:    []string{"foo=bar", "baz=baz"},
		},
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{"baz": "baz\t\n"},
			want:    []string{"foo=bar", "baz=baz"},
		},
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{"baz": "baz   "},
			want:    []string{"foo=bar", "baz=baz"},
		},
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{"baz": "baz\nbaz"},
			want:    []string{"foo=bar", "baz=baz\nbaz"},
		},
		{
			environ: []string{"foo=bar"},
			files:   map[string]string{"foo": ""},
			want:    []string{},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			dir, err := os.MkdirTemp("", "envv")
			if err != nil {
				t.Fatalf("unexpected err creating tempdir: %v", err)
			}

			for name, contents := range test.files {
				os.WriteFile(filepath.Join(dir, name), []byte(contents), os.ModePerm)
			}

			got, err := envv(test.environ, dir)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			if !slices.Equal(test.want, got) {
				t.Fatalf("envv slices did not match want=%v got=%v", test.want, got)
			}
		})
	}
}
