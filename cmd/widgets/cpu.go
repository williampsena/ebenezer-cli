package widgets

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/williampsena/ebenezer-cli/internal/cmd"
	formatters "github.com/williampsena/ebenezer-cli/internal/cmd/formatters"
)

type CpuCmd struct {
	WidgetCmd
	Loop            bool    `help:"Run the command in a loop." default:"false"`
	Interval        int     `help:"Interval (in seconds) between CPU usage checks." default:"3"`
	Burn            bool    `help:"Show fire emoji when memory usage is high." default:"true"`
	Threshold       float64 `help:"CPU usage threshold for high usage in percentage." default:"80"`
	ThresholdMedium float64 `help:"CPU usage threshold for medium usage in percentage." default:"50"`
}

func (w *CpuCmd) Run(ctx *cmd.Context) error {
	for {
		percentages, err := cpu.Percent(500*time.Millisecond, false)
		if err != nil {
			log.Fatalf("Error fetching CPU usage: %v\n", err)
		}

		usage := percentages[0]

		output, err := w.Render(ctx, usage)
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

func (w *CpuCmd) Render(ctx *cmd.Context, usage float64) (string, error) {
	class := "low"
	color := color_low
	emoji := ""

	iconColor := w.IconColor
	if iconColor == "" {
		iconColor = color
	}

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

	data := map[string]interface{}{
		"icon":       "",
		"icon-color": iconColor,
		"text":       fmt.Sprintf("%.0f%%%s", usage, emoji),
		"tooltip":    fmt.Sprintf("CPU usage: %.2f%%", usage),
		"class":      class,
		"color":      color,
	}

	return formatters.FormatWidgetOutput(w.Format, data)
}
