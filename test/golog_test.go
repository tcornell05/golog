// golog_test.go
package test

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tcornell05/golog"
)

// captureOutput redirects stdout for the duration of f and returns what was printed.
func captureOutput(f func()) string {
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stdout = oldOut
	return buf.String()
}

// captureStderr redirects stderr for the duration of f and returns what was printed.
func captureStderr(f func()) string {
	oldErr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	os.Stderr = oldErr
	return buf.String()
}

// TestDevelopmentLogger tests the development logger's output.
func TestDevelopmentLogger(t *testing.T) {
	logger := golog.NewDevelopment()
	output := captureStderr(func() {
		// Log an INFO message with an attribute.
		logger.Info("Test development log", slog.String("key", "value"))
	})
	if !strings.Contains(output, "Test development log") {
		t.Errorf("Development logger output missing message. Got: %s", output)
	}
	if !strings.Contains(output, "Attrib") && !strings.Contains(output, "key:") {
		t.Errorf("Development logger output missing attribute formatting. Got: %s", output)
	}
}

// TestDiscardLogger tests that the discard logger produces no output.
func TestDiscardLogger(t *testing.T) {
	logger := golog.NewDiscard()
	output := captureStderr(func() {
		logger.Info("This should be discarded")
	})
	if len(strings.TrimSpace(output)) != 0 {
		t.Errorf("Discard logger produced output, got: %s", output)
	}
}

// TestProductionLogger tests that the production logger outputs JSON-formatted logs.
func TestProductionLogger(t *testing.T) {
	logger := golog.NewProduction()
	output := captureStderr(func() {
		logger.Info("Test production log", slog.String("file", "main.go"))
	})
	// Check for JSON structure elements (keys "msg" and "file").
	if !strings.Contains(output, "\"msg\":\"Test production log\"") {
		t.Errorf("Production logger output doesn't seem to be JSON formatted. Got: %s", output)
	}
	if !strings.Contains(output, "\"file\":\"main.go\"") {
		t.Errorf("Production logger output missing expected attribute, got: %s", output)
	}
}

// TestTeraHandler tests the custom TeraHandler (from dev_tera_handler.go).
func TestTeraHandler(t *testing.T) {
	// Create a TeraHandler with specific options.
	teraOpts := golog.TeraOptions{
		HeaderTitle:      "Traceback (most recent call last)",
		BorderColor:      "#FF5F87", // lipgloss.Color is a string alias.
		PanelBorderColor: "#5F87FF",
		DateFormat:       time.RFC822,
	}
	teraHandler := golog.NewTeraHandler(teraOpts)
	logger := slog.New(teraHandler)

	// The TeraHandler writes its output to stdout.
	output := captureOutput(func() {
		logger.Error("Multi-line error\nSecond line",
			slog.String("file", "/home/user/app.go"),
			slog.Int("line", 42),
		)
	})

	// Check that output includes the header, line numbers and attribute details.
	if !strings.Contains(output, "Traceback") {
		t.Errorf("TeraHandler header missing, got: %s", output)
	}
	if !strings.Contains(output, "│") {
		t.Errorf("TeraHandler output missing formatted border/line numbers, got: %s", output)
	}
	if !strings.Contains(output, "file") {
		t.Errorf("TeraHandler output missing attribute details, got: %s", output)
	}
}
