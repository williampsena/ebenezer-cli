package widgets

import (
	"github.com/williampsena/ebenezer-cli/internal/core"
)

type WidgetCmd struct {
	Format    string `help:"Output format (e.g., waybar, polybar)" default:"waybar"`
	IconColor string `help:"Icon color for the widget." default:""`
	logger    core.Logger
}

func (h *WidgetCmd) SetupContext(debug bool) {
	h.logger = core.BuildLogger(debug)
}
