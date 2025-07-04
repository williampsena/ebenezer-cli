package desktop

import (
	cmd "github.com/williampsena/ebenezer-cli/internal/cmd"
	"github.com/williampsena/ebenezer-cli/internal/shell"
)

type NotificationsCmd struct {
	cmd.BaseCmd
	Clear    bool   `help:"Clear all notifications" default:"false"`
	Provider string `help:"Notification provider to use" default:"swaync"`
}

func (d *NotificationsCmd) Run(ctx *cmd.Context) error {
	d.SetupContext(ctx)

	if d.Clear {
		return d.ClearNotifications()
	}

	d.Logger.Error("No action specified.")

	return nil
}

func (d *NotificationsCmd) ClearNotifications() error {
	switch d.Provider {
	case "swaync":
		return d.ClearSwaync()
	case "dunst":
		return d.ClearDunst()
	default:
		d.Logger.Error("Unsupported notification provider", "provider", d.Provider)
	}

	return nil
}

func (d *NotificationsCmd) ClearSwaync() error {
	_, err := d.Shell.Run(shell.RunnerExecutionArgs{
		Command: "swaync-client",
		Args:    []string{"--close-all"},
	})

	if err != nil {
		d.Logger.Error("Failed to clear notifications with swaync", "error", err)
		return err
	}

	d.Logger.Info("All notifications cleared with swaync")

	return nil
}

func (d *NotificationsCmd) ClearDunst() error {
	_, err := d.Shell.Run(shell.RunnerExecutionArgs{
		Command: "dunstctl",
		Args:    []string{"history-clear"},
	})

	if err != nil {
		d.Logger.Error("Failed to clear notifications with dunst", "error", err)
		return err
	}

	d.Logger.Info("All notifications cleared with dunst")

	return nil
}
