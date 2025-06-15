package cmd

import (
	"fmt"

	"github.com/williampsena/ebenezer-cli/internal/core"
)

type PolybarFormatter struct{}

func (p PolybarFormatter) Format(data map[string]interface{}) (string, error) {
	color := data["color"].(string)
	text := p.buildWidgetText(data)

	return fmt.Sprintf("%%{F%v} %v%%{F-}", color, text), nil
}

func (w PolybarFormatter) buildWidgetText(data map[string]interface{}) string {
	text := data["text"].(string)
	color := data["color"].(string)
	icon := core.GetMapValue(data, "icon", "")
	iconColor := core.GetMapValue(data, "icon-color", "")
	noIcon := core.GetMapValue(data, "no-icon", false).(bool)

	if !noIcon {
		if iconColor != "" && icon != "" {
			return fmt.Sprintf("%%{F%v}%v%%{F-}%%{F%v}%v%%{F-}", iconColor, icon, color, text)
		}

		return fmt.Sprintf("%v %v", icon, text)
	}

	return fmt.Sprintf("%%{F%v}%v%%{F-}", color, text)
}
