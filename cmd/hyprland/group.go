package hyprland

type HyprlandGroup struct {
	Hyprlock  HyprlockCmd  `cmd:"" help:"Hyprland lock screen command"`
	Hyprpaper HyprpaperCmd `cmd:"" help:"Hyprland wallpaper management command"`
	Cron      CronCmd      `cmd:"" help:"Hyprland cron jobs command"`
	Reload    ReloadCmd    `cmd:"" help:"Reload Hyprland components (waybar, config, etc.)"`
}
