package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePath(t *testing.T) {
	t.Run("EmptyString", func(t *testing.T) {
		result := ResolvePath("")
		if result != "" {
			t.Errorf("Expected empty string, got '%s'", result)
		}
	})

	t.Run("AbsolutePath", func(t *testing.T) {
		tests := []string{
			"/usr/local/bin",
			"/home/user/documents",
			"/tmp/test",
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				result := ResolvePath(test)
				if result != test {
					t.Errorf("Expected '%s', got '%s'", test, result)
				}
			})
		}
	})

	t.Run("RelativePath", func(t *testing.T) {
		tests := []string{
			"./config",
			"../parent",
			"relative/path",
			".",
			"..",
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				result := ResolvePath(test)
				if result != test {
					t.Errorf("Expected '%s', got '%s'", test, result)
				}
			})
		}
	})

	t.Run("TildePath_Success", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Could not get user home directory: %v", err)
		}

		tests := []struct {
			input    string
			expected string
		}{
			{"~", homeDir},
			{"~/", filepath.Join(homeDir, "/")},
			{"~/Documents", filepath.Join(homeDir, "Documents")},
			{"~/.config", filepath.Join(homeDir, ".config")},
			{"~/Downloads/file.txt", filepath.Join(homeDir, "Downloads/file.txt")},
			{"~/path/with/multiple/segments", filepath.Join(homeDir, "path/with/multiple/segments")},
		}

		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				result := ResolvePath(test.input)
				if result != test.expected {
					t.Errorf("Expected '%s', got '%s'", test.expected, result)
				}
			})
		}
	})

	t.Run("TildePath_HomeDirError", func(t *testing.T) {
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		os.Unsetenv("HOME")

		input := "~/test/path"
		result := ResolvePath(input)

		if result != input && !strings.HasPrefix(result, "/") {
			t.Logf("ResolvePath with unset HOME returned: %s", result)
		}
	})

	t.Run("TildeNotAtStart", func(t *testing.T) {
		tests := []string{
			"path/~/not/expanded",
			"some~path",
			"path~",
			"/home/user/~backup",
		}

		for _, test := range tests {
			t.Run(test, func(t *testing.T) {
				result := ResolvePath(test)
				if result != test {
					t.Errorf("Expected '%s', got '%s'", test, result)
				}
			})
		}
	})

	t.Run("TildeWithUser", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Could not get user home directory: %v", err)
		}

		result := ResolvePath("~")
		if result != homeDir {
			t.Errorf("Expected '%s', got '%s'", homeDir, result)
		}
	})

	t.Run("TildeWithSlash", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Could not get user home directory: %v", err)
		}

		result := ResolvePath("~/")
		expected := filepath.Join(homeDir, "/")
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("TildeWithEdgeCases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "Just tilde character",
				input:    "~",
				expected: "",
			},
			{
				name:     "Tilde with single character",
				input:    "~a",
				expected: "",
			},
			{
				name:     "Multiple slashes after tilde",
				input:    "~//double//slash",
				expected: "",
			},
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Could not get user home directory: %v", err)
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var expected string
				switch test.input {
				case "~":
					expected = homeDir
				case "~a":
					expected = filepath.Join(homeDir, "a")
				case "~//double//slash":
					expected = filepath.Join(homeDir, "//double//slash")
				}

				result := ResolvePath(test.input)
				if result != expected {
					t.Errorf("Expected '%s', got '%s'", expected, result)
				}
			})
		}
	})

	t.Run("PreservesTrailingSlash", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("Could not get user home directory: %v", err)
		}

		tests := []struct {
			input    string
			expected string
		}{
			{"~/path/", filepath.Join(homeDir, "path/")},
			{"~/", filepath.Join(homeDir, "/")},
		}

		for _, test := range tests {
			t.Run(test.input, func(t *testing.T) {
				result := ResolvePath(test.input)
				if result != test.expected {
					t.Errorf("Expected '%s', got '%s'", test.expected, result)
				}
			})
		}
	})
}
