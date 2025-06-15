package widgets

import (
	"encoding/json"
	"testing"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
)

func TestCpuCmd_Render(t *testing.T) {
	t.Run("waybar format", func(t *testing.T) {
		tests := []struct {
			name            string
			usage           float64
			thresholdHigh   float64
			thresholdMedium float64
			burnEnabled     bool
			iconColor       string
			expectedOutput  map[string]interface{}
		}{
			{
				name:            "Low CPU usage",
				usage:           20.0,
				thresholdHigh:   80.0,
				thresholdMedium: 50.0,
				burnEnabled:     true,
				iconColor:       "",
				expectedOutput: map[string]interface{}{
					"icon":    "",
					"text":    "<span foreground='tt.iconColor'></span> <span foreground='#f8f8f2'>20%</span>",
					"tooltip": "CPU usage: 20.00%",
					"class":   "low",
					"color":   color_low,
				},
			},
			{
				name:            "Medium CPU usage",
				usage:           60.0,
				thresholdHigh:   80.0,
				thresholdMedium: 50.0,
				burnEnabled:     true,
				iconColor:       "",
				expectedOutput: map[string]interface{}{
					"icon":    "",
					"text":    "<span foreground='tt.iconColor'></span> <span foreground='#ffff00'>60%</span>",
					"tooltip": "CPU usage: 60.00%",
					"class":   "medium",
					"color":   color_medium,
				},
			},
			{
				name:            "High CPU usage with burn enabled",
				usage:           90.0,
				thresholdHigh:   80.0,
				thresholdMedium: 50.0,
				burnEnabled:     true,
				iconColor:       "",
				expectedOutput: map[string]interface{}{
					"icon":    "",
					"text":    "<span foreground='tt.iconColor'></span> <span foreground='#ff0000'>90% 🔥</span>",
					"tooltip": "CPU usage: 90.00%",
					"class":   "high",
					"color":   color_high,
				},
			},
			{
				name:            "High CPU usage with burn disabled",
				usage:           90.0,
				thresholdHigh:   80.0,
				thresholdMedium: 50.0,
				burnEnabled:     false,
				iconColor:       "",
				expectedOutput: map[string]interface{}{
					"icon":    "",
					"text":    "<span foreground='tt.iconColor'></span> <span foreground='#ff0000'>90%</span>",
					"tooltip": "CPU usage: 90.00%",
					"class":   "high",
					"color":   color_high,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cpuCmd := CpuCmd{
					WidgetCmd: WidgetCmd{
						Format:    "waybar",
						IconColor: "tt.iconColor",
					},
					Threshold:       tt.thresholdHigh,
					ThresholdMedium: tt.thresholdMedium,
					Burn:            tt.burnEnabled,
				}

				output, err := cpuCmd.Render(&cmd.Context{}, tt.usage)
				if err != nil {
					t.Fatalf("Render failed: %v", err)
				}

				var result map[string]interface{}
				err = json.Unmarshal([]byte(output), &result)
				if err != nil {
					t.Fatalf("Failed to parse JSON output: %v", err)
				}

				for key, expectedValue := range tt.expectedOutput {
					if result[key] != expectedValue {
						t.Errorf("Expected %s: %v, got: %v", key, expectedValue, result[key])
					}
				}
			})
		}
	})
}
