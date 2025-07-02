package hyprland

import (
	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	core "github.com/williampsena/ebenezer-cli/internal/core"
	"github.com/williampsena/ebenezer-cli/internal/process"
	settings "github.com/williampsena/ebenezer-cli/internal/settings"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

type HyprlandCmd struct {
	logger         core.Logger
	processManager process.ProcessManager
	shell          shell.Runner
}

func (h *HyprlandCmd) SetupContext(ctx *cmd.Context) {
	if settings.IsTestMode {
		knownProcess := []string{"hyprland", "hyprctl", "waybar"}

		h.logger = core.BuildLogger(false)
		h.processManager = process.NewProcessManagerMock(
			knownProcess,
			knownProcess,
		)
		h.shell = shell.NewRunnerMock(
			h.logger,
			knownProcess,
			knownProcess,
			knownProcess,
		)
	} else {
		h.logger = h.buildLogger(ctx)
		h.processManager = process.NewProcessManager(h.logger)
		h.shell = shell.NewRunner(h.logger)
	}
}

// Useful for testing purposes
func (h *HyprlandCmd) injectProcessManager(processManager process.ProcessManager) {
	h.processManager = processManager
}

func (h *HyprlandCmd) buildLogger(ctx *cmd.Context) core.Logger {
	if ctx.Silent {
		return core.BuildSilentLogger()
	}

	return core.BuildLogger(ctx.Debug)
}
