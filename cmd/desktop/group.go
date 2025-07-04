package desktop

type DesktopGroup struct {
	Notifications NotificationsCmd `cmd:"notifications" help:"Desktop notifications management command"`
}
