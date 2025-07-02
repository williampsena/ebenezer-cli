package widgets

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/williampsena/ebenezer-cli/internal/cmd"
	formatters "github.com/williampsena/ebenezer-cli/internal/cmd/formatters"
)

type TemperatureCmd struct {
	WidgetCmd
	Loop            bool    `help:"Run the command in a loop." default:"false"`
	Interval        int     `help:"Interval (in seconds) between temperature checks." default:"2"`
	Burn            bool    `help:"Show fire emoji when memory usage is high." default:"true"`
	Sensor          string  `help:"Specific sensor key to monitor. If empty, all sensors will be monitored." default:""`
	Threshold       float64 `help:"Temperature threshold for high usage in degrees Celsius." default:"70"`
	ThresholdMedium float64 `help:"Temperature threshold for high usage in degrees Celsius." default:"60"`
}

func (w *TemperatureCmd) Run(ctx *cmd.Context) error {
	w.SetupContext(ctx.Debug)

	for {
		temps, err := host.SensorsTemperatures()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching temperature data: %v\n", err)
			return err
		}

		output, err := w.Render(temps)
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

func (w *TemperatureCmd) Render(temps []host.TemperatureStat) (string, error) {
	var avgTemp float64

	if w.Sensor == "" {
		avgTemp = w.averageAllSensors(temps)
	} else {
		avgTemp = w.averageFilteredSensors(temps)
	}

	class := "normal"
	color := color_low
	emoji := ""

	if avgTemp > w.Threshold {
		class = "high"
		color = color_high
		if w.Burn {
			emoji = " 🔥"
		}
	} else if avgTemp > w.ThresholdMedium {
		class = "medium"
		color = color_medium
	}

	iconColor := w.IconColor
	if iconColor == "" {
		iconColor = color
	}

	output := map[string]interface{}{
		"icon":       "",
		"icon-color": iconColor,
		"text":       fmt.Sprintf("%.1f°C%s", avgTemp, emoji),
		"tooltip":    fmt.Sprintf("Average temperature: %.1f°C", avgTemp),
		"class":      class,
		"color":      color,
	}

	return formatters.FormatWidgetOutput(w.Format, output)
}

func (t *TemperatureCmd) averageAllSensors(temps []host.TemperatureStat) float64 {
	var avgTemp float64
	if len(temps) > 0 {
		for _, temp := range temps {
			avgTemp += temp.Temperature
		}
		avgTemp /= float64(len(temps))
	}

	return avgTemp
}

func (t *TemperatureCmd) averageFilteredSensors(temps []host.TemperatureStat) float64 {
	var avgTemp float64
	var sensorCount int

	for _, temp := range temps {
		if temp.SensorKey == t.Sensor {
			sensorCount++
			avgTemp += temp.Temperature
		}
	}

	if sensorCount == 0 {
		return 0.0
	}

	avgTemp /= float64(sensorCount)

	return avgTemp
}
