package core

import (
	"os"
	"path/filepath"
)

func ResolvePath(path string) string {
	if path == "" {
		return ""
	}

	if path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[1:])
	}

	return path
}
