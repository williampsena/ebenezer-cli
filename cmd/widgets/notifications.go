package widgets

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
	formatters "github.com/williampsena/ebenezer-cli/internal/cmd/formatters"
)

var bells = []string{"󰂚", "󰂞"}

type NotificationsCmd struct {
	WidgetCmd
	Loop     bool   `help:"Run the command in a loop." default:"false"`
	Interval int    `help:"Interval (in seconds) between notification checks." default:"5"`
	Provider string `help:"Notification provider to use (dunst or swaync). If empty, both will be checked." default:"swaync"`
}

func (w *NotificationsCmd) Run(ctx *cmd.Context) error {
	w.BuildLogger(ctx.Debug)

	for {
		count, err := getUnseenNotificationsCount(w.Provider)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching notifications count: %v\n", err)
			return err
		}

		output, err := w.Render(ctx, count)
		if err != nil {
			return err
		}

		if err := formatters.WriteToStdout(output); err != nil {
			return fmt.Errorf("error writing to stdout: %w", err)
		}

		if !w.Loop {
			return nil
		}

		time.Sleep(time.Duration(w.Interval) * time.Second)
	}
}

func (w *NotificationsCmd) Render(ctx *cmd.Context, count int) (string, error) {
	data := map[string]interface{}{
		"icon":       w.getIcon(count),
		"icon-color": w.IconColor,
		"text":       fmt.Sprintf("%d", count),
		"tooltip":    fmt.Sprintf("Unseen notifications: %d", count),
		"class":      "normal",
		"color":      "#ffffff",
	}

	return formatters.FormatWidgetOutput(w.Format, data)
}

func (w *NotificationsCmd) getIcon(count int) string {
	if count == 0 {
		return ""
	} else {
		return w.buildAnimatedBellIcon()
	}
}

func (w *NotificationsCmd) buildAnimatedBellIcon() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := r.Intn(len(bells))

	return bells[randomNum]
}

func getUnseenNotificationsCount(provider string) (int, error) {
	var count int
	var err error

	if provider == "dunst" {
		count, err := dunstCount()
		if err == nil {
			return count, nil
		}
	} else if provider == "swaync" {
		count, err = swayncCount()
		if err == nil {
			return count, nil
		}
	}

	return 0, fmt.Errorf("failed to query both notifications: %v", err)
}

func dunstCount() (int, error) {
	output, err := exec.Command("dunstctl", "count").Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute dunstctl: %v", err)
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("failed to parse dunstctl output: %v", err)
	}

	return count, nil
}

func swayncCount() (int, error) {
	output, err := exec.Command("swaync-client", "--count").Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute swaync-client: %v", err)
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("failed to parse swaync-client output: %v", err)
	}

	return count, nil
}
