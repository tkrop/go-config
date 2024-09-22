// Package filepath provides path utility functions for Unix systems.
package filepath

import (
	"os"
	"path/filepath"
)

// Normalize the given path by expanding environment variables, resolving the
// absolute path, and cleaning the path.
func Normalize(path string) string {
	path = os.ExpandEnv(path)

	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	if path, err := filepath.Abs(path); err == nil {
		return filepath.Clean(path)
	}

	return filepath.Clean(path)
}
