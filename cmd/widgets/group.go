package widgets

type WidgetGroup struct {
	Logo          LogoCmd          `cmd:"" help:"Widget Logo"`
	Cpu           CpuCmd           `cmd:"" help:"Widget CPU"`
	Memory        MemoryCmd        `cmd:"" help:"Widget Memory"`
	Temperature   TemperatureCmd   `cmd:"" help:"Widget Temperature"`
	Notifications NotificationsCmd `cmd:"" help:"Widget Notifications"`
}
