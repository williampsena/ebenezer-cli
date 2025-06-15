package widgets

import "testing"

func TestColors(t *testing.T) {
	tests := []struct {
		name     string
		colorVar string
		expected string
	}{
		{"Low Color", color_low, "#f8f8f2"},
		{"Medium Color", color_medium, "#ffff00"},
		{"High Color", color_high, "#ff0000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.colorVar != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.colorVar)
			}
		})
	}
}
