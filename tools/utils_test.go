package tools

import (
	"path/filepath"
	"testing"
)

func TestSafePath(t *testing.T) {
	cwd, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("failed to resolve cwd: %v", err)
	}

	insideAbs := filepath.Join(cwd, "tools", "bash.go")
	outsideAbs := filepath.Join(filepath.Dir(cwd), "outside.txt")

	tests := []struct {
		name      string
		input     string
		wantErr   bool
		wantExact string
	}{
		{
			name:      "dot path",
			input:     ".",
			wantErr:   false,
			wantExact: cwd,
		},
		{
			name:    "relative inside",
			input:   "tools/bash.go",
			wantErr: false,
		},
		{
			name:      "absolute inside",
			input:     insideAbs,
			wantErr:   false,
			wantExact: insideAbs,
		},
		{
			name:    "relative escape",
			input:   "../go.mod",
			wantErr: true,
		},
		{
			name:    "absolute outside",
			input:   outsideAbs,
			wantErr: true,
		},
		{
			name:    "normalized escape",
			input:   "tools/../../outside.txt",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := safePath(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got path: %s", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantExact != "" && got != tc.wantExact {
				t.Fatalf("expected exact path %q, got %q", tc.wantExact, got)
			}

			rel, relErr := filepath.Rel(cwd, got)
			if relErr != nil {
				t.Fatalf("failed to compute rel path: %v", relErr)
			}
			if rel == ".." || len(rel) >= 3 && rel[:3] == "../" {
				t.Fatalf("resolved path escaped cwd: %q (rel=%q)", got, rel)
			}
		})
	}
}
