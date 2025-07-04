package cmd

import (
	core "github.com/williampsena/ebenezer-cli/internal/core"
	"github.com/williampsena/ebenezer-cli/internal/process"
	settings "github.com/williampsena/ebenezer-cli/internal/settings"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

type BaseCmd struct {
	Logger         core.Logger            `kong:"-"`
	ProcessManager process.ProcessManager `kong:"-"`
	Shell          shell.Runner           `kong:"-"`
}

func (h *BaseCmd) SetupContext(ctx *Context) {
	if settings.IsTestMode {
		knownProcess := []string{"hyprland", "hyprctl", "waybar"}

		h.Logger = core.BuildLogger(false)
		h.ProcessManager = process.NewProcessManagerMock(
			knownProcess,
			knownProcess,
		)
		h.Shell = shell.NewRunnerMock(
			h.Logger,
			knownProcess,
			knownProcess,
			knownProcess,
		)
	} else {
		h.Logger = h.buildLogger(ctx)
		h.ProcessManager = process.NewProcessManager(h.Logger)
		h.Shell = shell.NewRunner(h.Logger)
	}
}

func (h *BaseCmd) buildLogger(ctx *Context) core.Logger {
	if ctx.Silent {
		return core.BuildSilentLogger()
	}

	return core.BuildLogger(ctx.Debug)
}
