package widgets

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/williampsena/ebenezer-cli/internal/cmd"
	formatters "github.com/williampsena/ebenezer-cli/internal/cmd/formatters"
)

type MemoryCmd struct {
	WidgetCmd
	Loop            bool    `help:"Run the command in a loop." default:"false"`
	Interval        int     `help:"Interval (in seconds) between Memory usage checks." default:"3"`
	Burn            bool    `help:"Show fire emoji when memory usage is high." default:"true"`
	Threshold       float64 `help:"Memory usage threshold for high usage in percentage." default:"80"`
	ThresholdMedium float64 `help:"Memory usage threshold for medium usage in percentage." default:"50"`
}

func (w *MemoryCmd) Run(ctx *cmd.Context) error {
	for {
		vm, err := mem.VirtualMemory()
		if err != nil {
			return err
		}

		output, err := w.Render(ctx, vm.Total, vm.Available)
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

func (w *MemoryCmd) Render(ctx *cmd.Context, vmTotal uint64, vmAvailable uint64) (string, error) {
	used := vmTotal - vmAvailable
	usage := (float64(used) / float64(vmTotal)) * 100

	class := "low"
	color := color_low
	emoji := ""

	switch {
	case usage > w.Threshold:
		class = "high"
		color = color_high
		if w.Burn {
			emoji = " 🔥"
		}
	case usage > w.ThresholdMedium:
		class = "medium"
		color = color_medium
	}

	iconColor := w.IconColor
	if iconColor == "" {
		iconColor = color
	}

	data := map[string]interface{}{
		"icon":       "󰄧",
		"icon-color": iconColor,
		"text":       fmt.Sprintf("%.0f%%%s", usage, emoji),
		"tooltip":    fmt.Sprintf("Memory usage: %.2f%%", usage),
		"class":      class,
		"color":      color,
	}

	return formatters.FormatWidgetOutput(w.Format, data)
}
