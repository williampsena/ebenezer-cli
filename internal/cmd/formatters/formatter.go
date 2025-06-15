package cmd

import (
	"fmt"
	"os"
)

var formatters = map[string]WidgetFormatter{
	"waybar":  WaybarFormatter{},
	"polybar": PolybarFormatter{},
	"text":    RawTextFormatter{},
}

type WidgetFormatter interface {
	Format(data map[string]interface{}) (string, error)
}

func FormatWidgetOutput(format string, data map[string]interface{}) (string, error) {
	formatter, exists := formatters[format]
	if !exists {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	return formatAndWriteOutput(formatter, data)
}

func formatAndWriteOutput(formatter WidgetFormatter, data map[string]interface{}) (string, error) {
	output, err := formatter.Format(data)
	if err != nil {
		return "", fmt.Errorf("error formatting widget output: %w", err)
	}

	return output, nil
}

func WriteToStdout(output string) error {
	_, err := os.Stdout.Write([]byte(output))
	return err
}
