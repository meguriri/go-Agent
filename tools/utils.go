package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// safePath returns an absolute path under current workspace.
// It rejects any path that escapes workspace via ".." or absolute path tricks.
func safePath(p string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get cwd failed: %w", err)
	}
	workdir, err = filepath.Abs(workdir)
	if err != nil {
		return "", fmt.Errorf("abs cwd failed: %w", err)
	}
	workdir = filepath.Clean(workdir)

	var target string
	if filepath.IsAbs(p) {
		target = p
	} else {
		target = filepath.Join(workdir, p)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("resolve path failed: %w", err)
	}
	target = filepath.Clean(target)
	rel, err := filepath.Rel(workdir, target)
	if err != nil {
		return "", fmt.Errorf("rel path failed: %w", err)
	}

	// rel == "." means exactly workdir itself
	if rel == "." {
		return target, nil
	}
	// outside workspace if starts with ".." segment
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes workspace: %s", p)
	}

	return target, nil
}
