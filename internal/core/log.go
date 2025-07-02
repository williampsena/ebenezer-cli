package core

import (
	"fmt"
	"os"
)

type Logger interface {
	IsDebugEnabled() bool            // Check if debug logging is enabled
	Debug(msg string, args ...any)   // Debug messages, typically verbose and for development
	Info(msg string, args ...any)    // General information messages
	Warning(msg string, args ...any) // Warnings that may require attention but are not critical
	Error(msg string, args ...any)   // Errors that need to be reported, typically to stderr
}

type logger struct{ debug bool }

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func (l *logger) IsDebugEnabled() bool {
	return l.debug
}

func (l *logger) Debug(msg string, args ...any) {
	if !l.debug {
		return
	}

	logWithColor(msg, args, colorBlue)
}

func (l *logger) Info(msg string, args ...any) {
	logWithColor(msg, args, colorReset)
}

func (l *logger) Warning(msg string, args ...any) {
	logWithColor(msg, args, colorYellow)
}

func (l *logger) Error(msg string, args ...any) {
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

func BuildLogger(debug bool) Logger {
	return &logger{debug: debug}
}

type SilentLogger struct{}

func (l *SilentLogger) IsDebugEnabled() bool {
	return false
}
func (l *SilentLogger) Debug(msg string, args ...any)   {}
func (l *SilentLogger) Info(msg string, args ...any)    {}
func (l *SilentLogger) Warning(msg string, args ...any) {}
func (l *SilentLogger) Error(msg string, args ...any) {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, msg, args...)
	} else {
		fmt.Fprint(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr)
}

func BuildSilentLogger() Logger {
	return &SilentLogger{}
}
