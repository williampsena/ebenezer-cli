package core

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSetupContext(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		logger := BuildLogger(true)
		if logger == nil {
			t.Fatal("SetupContext should return a Logger instance")
		}
		if !logger.IsDebugEnabled() {
			t.Error("Logger should have debug enabled")
		}

		logger = BuildLogger(false)
		if logger == nil {
			t.Fatal("SetupContext should return a Logger instance")
		}
		if logger.IsDebugEnabled() {
			t.Error("Logger should have debug disabled")
		}
	})

	t.Run("DebugEnabled", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(true)
		logger.Debug("test debug message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "test debug message") {
			t.Error("Debug message should be printed when debug is enabled")
		}
	})

	t.Run("DebugDisabled", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Debug("test debug message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if strings.Contains(output, "test debug message") {
			t.Error("Debug message should not be printed when debug is disabled")
		}
	})

	t.Run("Info", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Info("test info message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "test info message") {
			t.Error("Info message should be printed")
		}
	})

	t.Run("Warning", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Warning("test warning message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "test warning message") {
			t.Error("Warning message should be printed")
		}
	})

	t.Run("Error", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Error("test error message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "test error message") {
			t.Error("Error message should be printed")
		}
	})

	t.Run("WithArgs", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Info("test %s with %d args", "message", 2)

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "test message with 2 args") {
			t.Error("Message with args should be formatted correctly")
		}
	})

	t.Run("WithoutArgs", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := BuildLogger(false)
		logger.Info("simple message")

		w.Close()
		os.Stdout = oldStdout

		var buf bytes.Buffer
		buf.ReadFrom(r)
		output := buf.String()

		if !strings.Contains(output, "simple message") {
			t.Error("Simple message should be printed as-is")
		}
	})
}
