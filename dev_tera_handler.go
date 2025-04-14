// dev_tera_handler.go
package golog

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// TeraOptions holds style options for the Tera handler.
type TeraOptions struct {
	HeaderTitle      string         // e.g. "Traceback (most recent call last)"
	BorderColor      lipgloss.Color // For header/footer borders.
	PanelBorderColor lipgloss.Color // For attribute panel borders.
	DateFormat       string         // e.g. time.RFC3339; if empty a default is used.
}

// TeraHandler is a custom slog.Handler that formats log records
// with borders, colored text, and code snippets.
type TeraHandler struct {
	opts TeraOptions
}

// NewTeraHandler creates a new TeraHandler using the provided options.
func NewTeraHandler(opts TeraOptions) *TeraHandler {
	if opts.DateFormat == "" {
		opts.DateFormat = time.RFC822
	}
	return &TeraHandler{opts: opts}
}

// Enabled is part of the slog.Handler interface.
func (h *TeraHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Log every record (modify if needed)
	return true
}

// Handle formats the slog record in Tera style and writes it to stdout.
func (h *TeraHandler) Handle(ctx context.Context, r slog.Record) error {
	var contentBuf bytes.Buffer

	// Format the timestamp.
	timestamp := r.Time.Format(h.opts.DateFormat)

	// Create a header panel styled with lipgloss.
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(h.opts.BorderColor).
		Padding(0, 1).
		Bold(true)
	headerTitle := fmt.Sprintf("%s • %s", h.opts.HeaderTitle, timestamp)
	headerPanel := headerStyle.Render(headerTitle)
	contentBuf.WriteString(headerPanel + "\n")

	// Style the message based on log level
	var messageStyle lipgloss.Style
	var levelColor lipgloss.Color

	switch r.Level {
	case slog.LevelDebug:
		messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F87FF"))
		levelColor = lipgloss.Color("#5F87FF")
	case slog.LevelInfo:
		messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87FF5F"))
		levelColor = lipgloss.Color("#87FF5F")
	case slog.LevelWarn:
		messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF5F")).
			Bold(true)
		levelColor = lipgloss.Color("#FFAF5F")
	case slog.LevelError:
		messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)
		levelColor = lipgloss.Color("#FF5F87")
	default:
		messageStyle = lipgloss.NewStyle()
		levelColor = lipgloss.Color("#FFFFFF")
	}

	// Get log level as a nicely formatted string
	levelText := getLevelText(r.Level)
	levelDisplay := messageStyle.Copy().Bold(true).Render(levelText)

	// Format the message content
	message := messageStyle.Render(r.Message)

	// Combine the level indicator with the message
	contentBuf.WriteString(fmt.Sprintf("  %s  %s\n", levelDisplay, message))

	// Add code snippet if source info is available
	if r.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		if f.File != "" {
			config := DefaultCodeSnippetConfig()
			snippet := FormatCodeSnippet(f.File, f.Line, r.Level, config)

			// Create terminal-clickable path for tmux/neovim integration
			filePathLink := createClickableFilePath(f.File, f.Line)

			// Style for filename (bold, distinct color)
			filePathStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(levelColor).
				Padding(0, 1)

			// Use raw link without additional styling for the line number
			// to preserve the clickable link
			fileInfoBox := filePathStyle.Render(filePathLink)

			// Create a code panel with the code snippet only (no title text)
			codeStyle := lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(h.opts.PanelBorderColor).
				Padding(0, 1)

			// First add the file info box, then the code snippet
			processedSnippet := processCodeSnippet(snippet)
			codePanel := codeStyle.Render(fileInfoBox + "\n\n" + processedSnippet)
			contentBuf.WriteString("\n" + codePanel + "\n")
		}
	}

	// If the record contains attributes, render them in a separate panel.
	if r.NumAttrs() > 0 {
		var attrsBuf bytes.Buffer
		r.Attrs(func(a slog.Attr) bool {
			attrsBuf.WriteString(fmt.Sprintf("• %s: %v\n", a.Key, a.Value.Any()))
			return true
		})
		localsStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(h.opts.PanelBorderColor).
			Padding(0, 1)
		localsPanel := localsStyle.Render(attrsBuf.String())
		contentBuf.WriteString("\n" + localsPanel)
	}

	// Now wrap the entire log entry in a container
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(levelColor). // Use the level color for the border
		Padding(0, 1).
		MarginBottom(1) // Add space between log entries

	// Render the final container
	finalOutput := containerStyle.Render(contentBuf.String())

	// Output the formatted text.
	fmt.Print(finalOutput + "\n")

	return nil
}

// WithAttrs returns a new handler with additional attributes.
func (h *TeraHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we return a new handler with the same options.
	return &TeraHandler{opts: h.opts}
}

// WithGroup returns a new handler with an additional group.
func (h *TeraHandler) WithGroup(name string) slog.Handler {
	// For this handler, groups are not additionally processed.
	return h
}

// Helper function to get nicely formatted level text
func getLevelText(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO "
	case slog.LevelWarn:
		return "WARN "
	case slog.LevelError:
		return "ERROR"
	default:
		return fmt.Sprintf("%-5s", level.String())
	}
}

// Helper function to process code snippets for consistent formatting
func processCodeSnippet(snippet string) string {
	// Split into lines and process each line
	lines := strings.Split(snippet, "\n")

	// For now, we just join them back together
	// This function can be expanded later if needed
	return strings.Join(lines, "\n")
}

// Helper function to create terminal-clickable file paths
// Returns an OSC 8 hyperlink that works in modern terminals
func createClickableFilePath(path string, line int) string {
	// Get absolute path to ensure it works correctly
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path // Fallback to original path
	}

	// Create OSC 8 hyperlink format - terminal standard for clickable links
	// Format: \033]8;;file://{full_path}\033\\{display_text}\033]8;;\033\\
	fullPath := fmt.Sprintf("file://%s", absPath)

	// Extract just the filename and parent directory for display
	filename := filepath.Base(path)
	dir := filepath.Base(filepath.Dir(path))
	displayPath := dir + "/" + filename

	// Return terminal hyperlink with short display but clickable full path
	return fmt.Sprintf("\033]8;;%s#%d\033\\%s\033]8;;\033\\",
		fullPath, line, displayPath)
}
