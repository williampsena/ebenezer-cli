package widgets

type WidgetCmd struct {
	Format    string `help:"Output format (e.g., waybar, polybar)" default:"waybar"`
	IconColor string `help:"Icon color for the widget." default:""`
}
