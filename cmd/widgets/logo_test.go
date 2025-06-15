package widgets

import (
	"encoding/json"
	"testing"

	"github.com/williampsena/ebenezer-cli/internal/cmd"
)

func TestLogoCmd_Render(t *testing.T) {
	tests := []struct {
		name           string
		distroName     string
		version        string
		kernelVersion  string
		typeOption     string
		iconColor      string
		expectedOutput map[string]interface{}
	}{
		{
			name:          "Icon only",
			distroName:    "ubuntu",
			version:       "Ubuntu 22.04.2 LTS",
			kernelVersion: "5.15.0-72-generic",
			typeOption:    "icon",
			iconColor:     "#ffffff",
			expectedOutput: map[string]interface{}{
				"icon":    "",
				"text":    "<span foreground='#ffffff'></span> <span foreground='#ffffff'></span>",
				"tooltip": "Ubuntu 22.04.2 LTS 5.15.0-72-generic",
				"class":   "normal",
				"color":   "#ffffff",
			},
		},
		{
			name:          "Icon + name",
			distroName:    "ubuntu",
			version:       "Ubuntu 22.04.2 LTS",
			kernelVersion: "5.15.0-72-generic",
			typeOption:    "icon+name",
			iconColor:     "#ffffff",
			expectedOutput: map[string]interface{}{
				"icon":    "",
				"text":    "<span foreground='#ffffff'></span> <span foreground='#ffffff'>ubuntu</span>",
				"tooltip": "Ubuntu 22.04.2 LTS 5.15.0-72-generic",
				"class":   "normal",
				"color":   "#ffffff",
			},
		},
		{
			name:          "Name only",
			distroName:    "ubuntu",
			version:       "Ubuntu 22.04.2 LTS",
			kernelVersion: "5.15.0-72-generic",
			typeOption:    "name",
			iconColor:     "#ffffff",
			expectedOutput: map[string]interface{}{
				"icon":    "",
				"text":    "<span foreground='#ffffff'>ubuntu</span>",
				"tooltip": "Ubuntu 22.04.2 LTS 5.15.0-72-generic",
				"class":   "normal",
				"color":   "#ffffff",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logoCmd := LogoCmd{
				WidgetCmd: WidgetCmd{
					Format:    "waybar",
					IconColor: tt.iconColor,
				},
				Name: tt.distroName,
				Type: tt.typeOption,
			}

			output, err := logoCmd.Render(&cmd.Context{}, tt.version, tt.kernelVersion)
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

func TestLogoCmd_getLogo(t *testing.T) {
	tests := []struct {
		name         string
		distroName   string
		expectedLogo string
	}{
		{"Ubuntu Logo", "ubuntu", ""},
		{"Fedora Logo", "fedora", ""},
		{"Unknown Distro", "unknown-distro", "󰌽"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logoCmd := LogoCmd{}
			logo := logoCmd.getLogo(tt.distroName)
			if logo != tt.expectedLogo {
				t.Errorf("Expected logo %s, got %s", tt.expectedLogo, logo)
			}
		})
	}
}

func TestLogoCmd_buildText(t *testing.T) {
	tests := []struct {
		name         string
		distroName   string
		typeOption   string
		expectedText string
	}{
		{"Icon Only", "ubuntu", "icon", ""},
		{"Icon + Name", "ubuntu", "icon+name", "ubuntu"},
		{"Name Only", "ubuntu", "name", "ubuntu"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logoCmd := LogoCmd{Type: tt.typeOption}
			text := logoCmd.buildText(tt.distroName)
			if text != tt.expectedText {
				t.Errorf("Expected text %s, got %s", tt.expectedText, text)
			}
		})
	}
}
