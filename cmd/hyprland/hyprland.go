package hyprland

import (
	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	"github.com/williampsena/ebenezer-cli/internal/process"
)

type HyprlandCmd struct {
	cmd.BaseCmd
}

// Useful for testing purposes
func (h *HyprlandCmd) injectProcessManager(processManager process.ProcessManager) {
	h.BaseCmd.ProcessManager = processManager
}
