package widgets

import (
	"encoding/json"
	"testing"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
)

func TestMemoryCmd_Render(t *testing.T) {
	tests := []struct {
		name            string
		vmTotal         uint64
		vmAvailable     uint64
		thresholdHigh   float64
		thresholdMedium float64
		burnEnabled     bool
		iconColor       string
		expectedOutput  map[string]interface{}
	}{
		{
			name:            "Low Memory usage",
			vmTotal:         1000000,
			vmAvailable:     800000,
			thresholdHigh:   80.0,
			thresholdMedium: 50.0,
			burnEnabled:     true,
			iconColor:       "",
			expectedOutput: map[string]interface{}{
				"icon":    "󰄧",
				"text":    "<span foreground='#f8f8f2'>󰄧</span> <span foreground='#f8f8f2'>20%</span>",
				"tooltip": "Memory usage: 20.00%",
				"class":   "low",
				"color":   color_low,
			},
		},
		{
			name:            "Medium Memory usage",
			vmTotal:         1000000,
			vmAvailable:     400000,
			thresholdHigh:   80.0,
			thresholdMedium: 50.0,
			burnEnabled:     true,
			iconColor:       "",
			expectedOutput: map[string]interface{}{
				"icon":    "󰄧",
				"text":    "<span foreground='#ffff00'>󰄧</span> <span foreground='#ffff00'>60%</span>",
				"tooltip": "Memory usage: 60.00%",
				"class":   "medium",
				"color":   color_medium,
			},
		},
		{
			name:            "High Memory usage with burn enabled",
			vmTotal:         1000000,
			vmAvailable:     100000,
			thresholdHigh:   80.0,
			thresholdMedium: 50.0,
			burnEnabled:     true,
			iconColor:       "",
			expectedOutput: map[string]interface{}{
				"icon":    "󰄧",
				"text":    "<span foreground='#ff0000'>󰄧</span> <span foreground='#ff0000'>90% 🔥</span>",
				"tooltip": "Memory usage: 90.00%",
				"class":   "high",
				"color":   color_high,
			},
		},
		{
			name:            "High Memory usage with burn disabled",
			vmTotal:         1000000,
			vmAvailable:     100000,
			thresholdHigh:   80.0,
			thresholdMedium: 50.0,
			burnEnabled:     false,
			iconColor:       "",
			expectedOutput: map[string]interface{}{
				"icon":    "󰄧",
				"text":    "<span foreground='#ff0000'>󰄧</span> <span foreground='#ff0000'>90%</span>",
				"tooltip": "Memory usage: 90.00%",
				"class":   "high",
				"color":   color_high,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock MemoryCmd
			memoryCmd := MemoryCmd{
				WidgetCmd: WidgetCmd{
					Format:    "waybar",
					IconColor: tt.iconColor,
				},
				Threshold:       tt.thresholdHigh,
				ThresholdMedium: tt.thresholdMedium,
				Burn:            tt.burnEnabled,
			}

			output, err := memoryCmd.Render(&cmd.Context{}, tt.vmTotal, tt.vmAvailable)
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
}
