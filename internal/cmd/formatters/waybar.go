package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

type WaybarOutput struct {
	Icon    string `json:"icon,omitempty"`
	Text    string `json:"text,omitempty"`
	Tooltip string `json:"tooltip,omitempty"`
	Class   string `json:"class,omitempty"`
	Color   string `json:"color,omitempty"`
}

type WaybarFormatter struct{}

func (w WaybarFormatter) Format(data map[string]interface{}) (string, error) {
	output := WaybarOutput{
		Icon:    core.GetMapValue(data, "icon", "").(string),
		Text:    w.buildWidgetText(data),
		Tooltip: data["tooltip"].(string),
		Class:   data["class"].(string),
		Color:   data["color"].(string),
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		return "", err
	}
	return string(jsonOutput), nil
}

func (w WaybarFormatter) buildWidgetText(data map[string]interface{}) string {
	text := data["text"].(string)
	color := data["color"].(string)
	noIcon := core.GetMapValue(data, "no-icon", false).(bool)

	icon := core.GetMapValue(data, "icon", "")
	iconColor := core.GetMapValue(data, "icon-color", "")

	if !noIcon {
		if iconColor != "" && icon != "" {
			return fmt.Sprintf("<span foreground='%v'>%v</span> <span foreground='%v'>%v</span>", iconColor, icon, color, text)
		}

		return fmt.Sprintf("%v %v", icon, text)
	}

	return fmt.Sprintf("<span foreground='%v'>%v</span>", color, text)
}
