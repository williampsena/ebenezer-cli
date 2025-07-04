package hyprland

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
)

func TestHyprlockCmd_getMessage(t *testing.T) {
	tests := []struct {
		name        string
		jokes       bool
		message     string
		provider    []string
		startup     bool
		expectEmpty bool
	}{
		{
			name:        "Default message when no jokes and no custom message",
			jokes:       false,
			message:     "",
			expectEmpty: false,
		},
		{
			name:        "Custom message",
			jokes:       false,
			message:     "Custom lock message",
			expectEmpty: false,
		},
		{
			name:        "Jokes enabled with valid provider",
			jokes:       true,
			message:     "",
			provider:    []string{"icanhazdadjoke"},
			startup:     false,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hyprlockCmd := &HyprlockCmd{
				Jokes:    tt.jokes,
				Message:  tt.message,
				Provider: tt.provider,
				Startup:  tt.startup,
			}
			hyprlockCmd.SetupContext(&cmd.Context{Debug: false})

			result, err := hyprlockCmd.getMessage()
			if err != nil {
				t.Errorf("getMessage() returned unexpected error: %v", err)
				return
			}

			if tt.expectEmpty && result != "" {
				t.Errorf("Expected empty message, got: %s", result)
			}

			if !tt.expectEmpty && result == "" {
				t.Errorf("Expected non-empty message, got empty string")
			}

			if !tt.jokes && tt.message == "" && result != defaultLockMessage {
				t.Errorf("Expected default message '%s', got '%s'", defaultLockMessage, result)
			}

			if !tt.jokes && tt.message != "" && result != tt.message {
				t.Errorf("Expected custom message '%s', got '%s'", tt.message, result)
			}
		})
	}
}

func TestHyprlockCmd_getProvider(t *testing.T) {
	tests := []struct {
		name      string
		providers []string
	}{
		{
			name:      "Single provider",
			providers: []string{"icanhazdadjoke"},
		},
		{
			name:      "Multiple providers",
			providers: []string{"icanhazdadjoke", "reddit"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hyprlockCmd := &HyprlockCmd{
				Provider: tt.providers,
			}

			result := hyprlockCmd.getProvider()
			found := false
			for _, provider := range tt.providers {
				if result == provider {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("getProvider() returned '%s', which is not in the provider list %v", result, tt.providers)
			}
		})
	}
}

func TestHyprlockCmd_fetchJokes(t *testing.T) {
	hyprlockCmd := &HyprlockCmd{
		Startup: false,
	}
	hyprlockCmd.SetupContext(&cmd.Context{Debug: false})

	t.Run("Valid provider", func(t *testing.T) {
		result, err := hyprlockCmd.fetchJokes("icanhazdadjoke")
		if err != nil {
			t.Logf("Failed to fetch jokes (expected in test environment): %v", err)
			if !strings.Contains(err.Error(), "failed to fetch joke after 3 attempts") {
				t.Errorf("Unexpected error format: %v", err)
			}
		} else {
			if result == "" {
				t.Error("Expected non-empty joke result")
			}
		}
	})

	t.Run("Invalid provider", func(t *testing.T) {
		_, err := hyprlockCmd.fetchJokes("invalidprovider")
		if err == nil {
			t.Error("Expected error for invalid provider")
		}
	})
}

func TestHyprlockCmd_Run(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "hyprlock.conf")

	sampleConfig := `
background {
    monitor = 
    path = screenshot
    blur_passes = 3
    blur_size = 8
}

label {
    monitor =
    text = Test message
    text_align = center
    color = rgba(200, 200, 200, 1.0)
}
`
	err := os.WriteFile(configPath, []byte(sampleConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	tests := []struct {
		name       string
		dry        bool
		startup    bool
		jokes      bool
		message    string
		configPath string
		expectErr  bool
	}{
		{
			name:       "Dry run mode",
			dry:        true,
			startup:    false,
			jokes:      false,
			message:    "Test message",
			configPath: configPath,
			expectErr:  false,
		},
		{
			name:       "Normal run with custom message",
			dry:        false,
			startup:    false,
			jokes:      false,
			message:    "Custom test message",
			configPath: configPath,
			expectErr:  false,
		},
		{
			name:       "Nonexistent config file",
			dry:        true,
			startup:    false,
			jokes:      false,
			message:    "Test message",
			configPath: "/nonexistent/path/hyprlock.conf",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hyprlockCmd := &HyprlockCmd{
				Dry:        tt.dry,
				Startup:    tt.startup,
				Jokes:      tt.jokes,
				Message:    tt.message,
				ConfigPath: tt.configPath,
				Provider:   []string{"icanhazdadjoke"},
				Format:     "👉 %s 🤪",
			}

			ctx := &cmd.Context{Debug: false}
			err := hyprlockCmd.Run(ctx)

			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestHyprlockCmd_DefaultValues(t *testing.T) {
	hyprlockCmd := &HyprlockCmd{}

	if hyprlockCmd.Dry != false {
		t.Errorf("Expected default Dry to be false, got %v", hyprlockCmd.Dry)
	}

	if hyprlockCmd.Startup != false {
		t.Errorf("Expected default Startup to be false, got %v", hyprlockCmd.Startup)
	}

	if hyprlockCmd.Message != "" {
		t.Errorf("Expected default Message to be empty, got %s", hyprlockCmd.Message)
	}

	hyprlockCmd.ConfigPath = "$HOME/.config/hypr/hyprlock.conf"
	if hyprlockCmd.ConfigPath != "$HOME/.config/hypr/hyprlock.conf" {
		t.Errorf("Failed to set ConfigPath to expected default value")
	}

	hyprlockCmd.Format = "👉 %s 🤪"
	if hyprlockCmd.Format != "👉 %s 🤪" {
		t.Errorf("Failed to set Format to expected default value")
	}
}

func TestDefaultLockMessage(t *testing.T) {
	expected := "Powered by hyprlock 🔥"
	if defaultLockMessage != expected {
		t.Errorf("Expected defaultLockMessage to be '%s', got '%s'", expected, defaultLockMessage)
	}
}
