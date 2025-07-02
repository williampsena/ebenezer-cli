package cmd

import (
	"encoding/json"
	"testing"
)

func TestWaybarFormatter(t *testing.T) {
	t.Run("FormatSuccess", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"icon":    "🎵",
			"text":    "Test Text",
			"tooltip": "Test Tooltip",
			"class":   "test-class",
			"color":   "#ff0000",
		}

		result, err := formatter.Format(data)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var output WaybarOutput
		err = json.Unmarshal([]byte(result), &output)
		if err != nil {
			t.Errorf("Failed to unmarshal result: %v", err)
		}

		if output.Icon != "🎵" {
			t.Errorf("Expected icon '🎵', got '%s'", output.Icon)
		}
		if output.Tooltip != "Test Tooltip" {
			t.Errorf("Expected tooltip 'Test Tooltip', got '%s'", output.Tooltip)
		}
		if output.Class != "test-class" {
			t.Errorf("Expected class 'test-class', got '%s'", output.Class)
		}
		if output.Color != "#ff0000" {
			t.Errorf("Expected color '#ff0000', got '%s'", output.Color)
		}
	})
}

func TestBuildWidgetText(t *testing.T) {
	t.Run("withIcon", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"text":  "Test",
			"color": "#ff0000",
			"icon":  "🎵",
		}

		result := formatter.buildWidgetText(data)
		expected := "🎵 Test"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("withIconColor", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"text":       "Test",
			"color":      "#ff0000",
			"icon":       "🎵",
			"icon-color": "#00ff00",
		}

		result := formatter.buildWidgetText(data)
		expected := "<span foreground='#00ff00'>🎵</span> <span foreground='#ff0000'>Test</span>"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("noIcon", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"text":    "Test",
			"color":   "#ff0000",
			"no-icon": true,
		}

		result := formatter.buildWidgetText(data)
		expected := "<span foreground='#ff0000'>Test</span>"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("emptyIcon", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"text":       "Test",
			"color":      "#ff0000",
			"icon":       "",
			"icon-color": "#00ff00",
		}

		result := formatter.buildWidgetText(data)
		expected := " Test"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("emptyIconColor", func(t *testing.T) {
		formatter := WaybarFormatter{}
		data := map[string]interface{}{
			"text":       "Test",
			"color":      "#ff0000",
			"icon":       "🎵",
			"icon-color": "",
		}

		result := formatter.buildWidgetText(data)
		expected := "🎵 Test"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})
}

func TestWaybarOutput(t *testing.T) {
	t.Run("JSONMarshaling", func(t *testing.T) {
		output := WaybarOutput{
			Icon:    "🎵",
			Text:    "Test Text",
			Tooltip: "Test Tooltip",
			Class:   "test-class",
			Color:   "#ff0000",
		}

		data, err := json.Marshal(output)
		if err != nil {
			t.Errorf("Failed to marshal WaybarOutput: %v", err)
		}

		var unmarshaled WaybarOutput
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			t.Errorf("Failed to unmarshal WaybarOutput: %v", err)
		}

		if unmarshaled.Icon != output.Icon {
			t.Errorf("Icon mismatch: expected '%s', got '%s'", output.Icon, unmarshaled.Icon)
		}
		if unmarshaled.Text != output.Text {
			t.Errorf("Text mismatch: expected '%s', got '%s'", output.Text, unmarshaled.Text)
		}
	})
}
