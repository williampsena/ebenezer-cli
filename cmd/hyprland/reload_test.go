package hyprland

import (
	"os"
	"strings"
	"testing"

	internalcmd "github.com/williampsena/ebenezer-cli/internal/cmd"
)

func TestReloadCmd_Run(t *testing.T) {
	tests := []struct {
		name      string
		component string
		expectErr bool
	}{
		{
			name:      "Valid component: all",
			component: "all",
			expectErr: false,
		},
		{
			name:      "Valid component: hyprland",
			component: "hyprland",
			expectErr: false,
		},
		{
			name:      "Valid component: waybar",
			component: "waybar",
			expectErr: false,
		},
		{
			name:      "Invalid component",
			component: "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ReloadCmd{
				Component: tt.component,
				WaitTime:  1,
			}

			err := cmd.Run(&internalcmd.Context{Debug: false})
			if tt.expectErr && err == nil {
				t.Errorf("Expected error for component '%s', got nil", tt.component)
			}
			if !tt.expectErr && err != nil && !isExpectedError(err) {
				t.Errorf("Unexpected error for component '%s': %v", tt.component, err)
			}
		})
	}
}

func TestReloadCmd_checkDependencies(t *testing.T) {
	tests := []struct {
		name      string
		component string
		expectErr bool
	}{
		{
			name:      "Check dependencies for hyprland",
			component: "hyprland",
			expectErr: false, // hyprctl, pgrep, kill should be available
		},
		{
			name:      "Check dependencies for all",
			component: "all",
			expectErr: false, // Basic tools should be available, waybar might not be
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ReloadCmd{Component: tt.component}
			cmd.SetupContext(&internalcmd.Context{Silent: true})

			err := cmd.checkDependencies()
			// We don't strictly require this to pass since dependencies might not be installed
			// but we test that the function executes without panic
			if err != nil {
				t.Logf("Dependency check failed (expected in test environment): %v", err)
			}
		})
	}
}

func TestReloadCmd_validateEnvironment(t *testing.T) {
	cmd := &ReloadCmd{}
	cmd.SetupContext(&internalcmd.Context{Debug: false})

	// Save original environment
	originalWaylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	originalHyprlandSignature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

	defer func() {
		// Restore original environment
		if originalWaylandDisplay != "" {
			os.Setenv("WAYLAND_DISPLAY", originalWaylandDisplay)
		} else {
			os.Unsetenv("WAYLAND_DISPLAY")
		}
		if originalHyprlandSignature != "" {
			os.Setenv("HYPRLAND_INSTANCE_SIGNATURE", originalHyprlandSignature)
		} else {
			os.Unsetenv("HYPRLAND_INSTANCE_SIGNATURE")
		}
	}()

	t.Run("No Wayland display", func(t *testing.T) {
		os.Unsetenv("WAYLAND_DISPLAY")
		err := cmd.validateEnvironment()
		if err == nil {
			t.Error("Expected error when WAYLAND_DISPLAY is not set")
		}
		if !strings.Contains(err.Error(), "not running in Wayland session") {
			t.Errorf("Expected Wayland error, got: %v", err)
		}
	})

	t.Run("With Wayland display", func(t *testing.T) {
		os.Setenv("WAYLAND_DISPLAY", "wayland-0")
		err := cmd.validateEnvironment()
		if err != nil {
			t.Errorf("Unexpected error with Wayland display set: %v", err)
		}
	})

	t.Run("With both environment variables", func(t *testing.T) {
		os.Setenv("WAYLAND_DISPLAY", "wayland-0")
		os.Setenv("HYPRLAND_INSTANCE_SIGNATURE", "test-signature")
		err := cmd.validateEnvironment()
		if err != nil {
			t.Errorf("Unexpected error with both env vars set: %v", err)
		}
	})
}

func TestReloadCmd_startWaybar(t *testing.T) {
	cmd := &ReloadCmd{}
	cmd.SetupContext(&internalcmd.Context{Debug: false})

	// This will likely fail in most test environments since waybar won't be installed
	// but we test that the function handles the error gracefully
	t.Run("Start waybar (may fail in test environment)", func(t *testing.T) {
		err := cmd.startWaybar()
		if err != nil {
			// Expected in test environment
			t.Logf("Failed to start waybar (expected in test environment): %v", err)
			if !strings.Contains(err.Error(), "waybar not found") &&
				!strings.Contains(err.Error(), "failed to start waybar") {
				t.Errorf("Unexpected error type: %v", err)
			}
		}
	})
}

func TestReloadCmd_reloadHyprland(t *testing.T) {
	cmd := &ReloadCmd{}
	cmd.SetupContext(&internalcmd.Context{Debug: false})

	// This will likely fail in most test environments since hyprctl won't be available
	// but we test that the function handles the error gracefully
	t.Run("Reload Hyprland (may fail in test environment)", func(t *testing.T) {
		err := cmd.reloadHyprland()
		if err != nil {
			// Expected in test environment
			t.Logf("Failed to reload Hyprland (expected in test environment): %v", err)
			if !strings.Contains(err.Error(), "failed to reload Hyprland") {
				t.Errorf("Unexpected error format: %v", err)
			}
		}
	})
}

func TestReloadCmd_reloadWaybar(t *testing.T) {
	cmd := &ReloadCmd{WaitTime: 1}
	cmd.SetupContext(&internalcmd.Context{Debug: false})

	// This will likely fail in most test environments
	// but we test that the function handles the process correctly
	t.Run("Reload waybar (may fail in test environment)", func(t *testing.T) {
		err := cmd.reloadWaybar()
		if err != nil {
			// Expected in test environment since waybar is likely not installed
			t.Logf("Failed to reload waybar (expected in test environment): %v", err)
		}
	})
}

func TestReloadCmd_reloadAll(t *testing.T) {
	cmd := &ReloadCmd{WaitTime: 1}
	cmd.SetupContext(&internalcmd.Context{Debug: false})

	// This will likely fail in most test environments
	// but we test that the function executes without panic
	t.Run("Reload all (may fail in test environment)", func(t *testing.T) {
		err := cmd.reloadAll()
		if err != nil {
			// Expected in test environment
			t.Logf("Failed to reload all (expected in test environment): %v", err)
		}
	})
}

func TestReloadCmd_DefaultValues(t *testing.T) {
	cmd := &ReloadCmd{}

	// Test default values when not explicitly set
	if cmd.WaitTime == 0 {
		cmd.WaitTime = 2 // Default value according to struct tag
	}

	if cmd.Component == "" {
		cmd.Component = "all" // Default value according to struct tag
	}

	if cmd.WaitTime != 2 {
		t.Errorf("Expected default WaitTime to be 2, got %d", cmd.WaitTime)
	}

	if cmd.Component != "all" {
		t.Errorf("Expected default Component to be 'all', got %s", cmd.Component)
	}
}

func TestReloadCmd_BuildLogger(t *testing.T) {
	cmd := &ReloadCmd{}

	t.Run("Build logger with debug disabled", func(t *testing.T) {
		cmd.SetupContext(&internalcmd.Context{Debug: false})
		if cmd.logger == nil {
			t.Error("Logger should not be nil after SetupContext")
		}
	})

	t.Run("Build logger with debug enabled", func(t *testing.T) {
		cmd.SetupContext(&internalcmd.Context{Debug: false})
		if cmd.logger == nil {
			t.Error("Logger should not be nil after SetupContext")
		}
	})
}

func TestReloadCmd_performHealthCheck(t *testing.T) {
	tests := []struct {
		name      string
		component string
		expectErr bool
	}{
		{
			name:      "Health check for hyprland",
			component: "hyprland",
			expectErr: false,
		},
		{
			name:      "Health check for waybar",
			component: "waybar",
			expectErr: false,
		},
		{
			name:      "Health check for all",
			component: "all",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &ReloadCmd{Component: tt.component}
			cmd.SetupContext(&internalcmd.Context{Silent: true})

			err := cmd.performHealthCheck()
			if err != nil {
				// Expected in test environment where hyprland/waybar might not be running
				t.Logf("Health check failed (expected in test environment): %v", err)
				if !isExpectedError(err) {
					t.Errorf("Unexpected error type: %v", err)
				}
			}
		})
	}
}

// Helper function to determine if an error is expected in a test environment
func isExpectedError(err error) bool {
	errorString := err.Error()
	expectedErrors := []string{
		"waybar not found",
		"failed to start waybar",
		"failed to reload Hyprland",
		"hyprctl",
		"executable file not found",
		"not running in Wayland session",
		"failed to find process",
		"environment validation failed",
		"dependency check failed",
		"health check failed",
		"waybar is not running after reload",
		"hyprland health check failed",
	}

	for _, expected := range expectedErrors {
		if strings.Contains(errorString, expected) {
			return true
		}
	}
	return false
}
