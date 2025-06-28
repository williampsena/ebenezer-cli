package core

import (
	"fmt"
	"os"
)

type Logger struct{ debug bool }

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func (l *Logger) Debug(msg string, args ...any) {
	if !l.debug {
		return
	}

	logWithColor(msg, args, colorBlue)
}

func (l *Logger) Info(msg string, args ...any) {
	logWithColor(msg, args, colorReset)
}

func (l *Logger) Warning(msg string, args ...any) {
	logWithColor(msg, args, colorYellow)
}

func (l *Logger) Error(msg string, args ...any) {
	logWithColor(msg, args, colorRed)
}

func logWithColor(msg string, args []any, color string) {
	var formatted string
	if len(args) > 0 {
		formatted = fmt.Sprintf(msg, args...)
	} else {
		formatted = msg
	}
	fmt.Fprintf(os.Stdout, "%s%s%s\n", color, formatted, colorReset)
}

func BuildLogger(debug bool) *Logger {
	return &Logger{debug: debug}
}
