package hyprland

import (
	core "github.com/williampsena/ebenezer-cli/internal/core"
)

type HyprlandCmd struct {
	logger *core.Logger
}

func (h *HyprlandCmd) BuildLogger(debug bool) {
	h.logger = core.BuildLogger(debug)
}
