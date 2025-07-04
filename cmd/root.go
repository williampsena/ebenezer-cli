package cmd

import (
	"github.com/williampsena/ebenezer-cli/cmd/desktop"
	"github.com/williampsena/ebenezer-cli/cmd/hyprland"
	"github.com/williampsena/ebenezer-cli/cmd/widgets"
)

type CLI struct {
	Debug    bool                   `help:"Enable debug mode."`
	Silent   bool                   `help:"Enable silent mode."`
	Desktop  desktop.DesktopGroup   `cmd:"" help:"Desktop commands"`
	Widgets  widgets.WidgetGroup    `cmd:"" help:"Waybar commands (JSON mode)"`
	Hyprland hyprland.HyprlandGroup `cmd:"" help:"Hyprland commands"`
}
