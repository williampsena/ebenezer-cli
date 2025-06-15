package cmd

import (
	"fmt"
)

type RawTextFormatter struct{}

func (p RawTextFormatter) Format(data map[string]interface{}) (string, error) {
	text := data["text"].(string)
	tooltip := data["tooltip"].(string)

	return fmt.Sprintf("%v\t%v", text, tooltip), nil
}
