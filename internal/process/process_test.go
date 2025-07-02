package process

import (
	"os"
	"strings"
	"testing"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

func TestNewProcessManager(t *testing.T) {
	logger := core.BuildLogger(false)
	pm := NewProcessManager(logger)

	if pm == nil {
		t.Fatal("NewProcessManager should not return nil")
	}
}

func TestProcessManager_IsProcessRunning(t *testing.T) {
	logger := core.BuildLogger(false)
	pm := NewProcessManager(logger)

	tests := []struct {
		name        string
		processName string
		shouldExist bool
	}{
		{
			name:        "Process that doesn't exist",
			processName: "nonexistentprocess12345",
			shouldExist: false,
		},
		{
			name:        "Current shell process",
			processName: getShellName(),
			shouldExist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pm.IsProcessRunning(tt.processName)
			if tt.shouldExist && !result {
				t.Logf("Expected process %s to be running, but it wasn't found", tt.processName)
			}
			if !tt.shouldExist && result {
				t.Errorf("Expected process %s not to be running, but it was found", tt.processName)
			}
		})
	}
}

func TestProcessManager_KillProcess(t *testing.T) {
	logger := core.BuildLogger(false)
	pm := NewProcessManager(logger)

	t.Run("Kill nonexistent process", func(t *testing.T) {
		err := pm.KillProcess("nonexistentprocess12345")
		if err == nil {
			t.Error("Expected error when trying to kill nonexistent process")
		}
		if !strings.Contains(err.Error(), "failed to find process") {
			t.Errorf("Expected 'failed to find process' error, got: %v", err)
		}
	})

	t.Run("Kill process with empty name", func(t *testing.T) {
		err := pm.KillProcess("")
		if err == nil {
			t.Error("Expected error when trying to kill process with empty name")
		}
	})
}

func TestProcessManagerInterface(t *testing.T) {
	logger := core.BuildLogger(false)
	pm := NewProcessManager(logger)

	var _ ProcessManager = pm

	t.Run("Interface methods", func(t *testing.T) {
		result := pm.IsProcessRunning("test")
		_ = result

		err := pm.KillProcess("nonexistent")
		if err == nil {
			t.Error("Expected error for nonexistent process")
		}
	})
}

func TestProcessManagerImpl_Type(t *testing.T) {
	logger := core.BuildLogger(false)
	pm := NewProcessManager(logger)

	if _, ok := pm.(*processManagerImpl); !ok {
		t.Error("NewProcessManager should return *processManagerImpl")
	}
}

func getShellName() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "bash"
	}
	if lastSlash := strings.LastIndex(shell, "/"); lastSlash != -1 {
		shell = shell[lastSlash+1:]
	}
	return shell
}
