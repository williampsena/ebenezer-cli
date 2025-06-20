package cmd

import (
	"github.com/williampsena/ebenezer-cli/cmd/widgets"
)

type CLI struct {
	Debug   bool                `help:"Enable debug mode."`
	Widgets widgets.WidgetGroup `cmd:"" help:"Waybar commands (JSON mode)"`
}
