// code_snippet.go
package golog

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CodeSnippetConfig holds configuration options for code snippet display
type CodeSnippetConfig struct {
	// Number of context lines to show before and after the target line
	ContextLines int
	// Styling options
	ErrorLineStyle    lipgloss.Style
	DebugLineStyle    lipgloss.Style
	InfoLineStyle     lipgloss.Style
	WarnLineStyle     lipgloss.Style
	ContextLineStyle  lipgloss.Style
	LineNumberStyle   lipgloss.Style
	ErrorPointerStyle lipgloss.Style
	// Text to show under different log levels
	ErrorPointerText string
	DebugPointerText string
	InfoPointerText  string
	WarnPointerText  string
}

// DefaultCodeSnippetConfig returns the default configuration for code snippets
func DefaultCodeSnippetConfig() CodeSnippetConfig {
	return CodeSnippetConfig{
		ContextLines: 2,
		ErrorLineStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true),
		DebugLineStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5F87FF")).
			Bold(true),
		InfoLineStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87FF5F")).
			Bold(true),
		WarnLineStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF5F")).
			Bold(true),
		ContextLineStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")),
		LineNumberStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#777777")),
		ErrorPointerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true),
		ErrorPointerText: "^ Error occurs here",
		DebugPointerText: "^ Debug log emitted here",
		InfoPointerText:  "^ Info log emitted here",
		WarnPointerText:  "^ Warning log emitted here",
	}
}

// FormatCodeSnippet extracts and formats a code snippet from a file at the specified line
func FormatCodeSnippet(filePath string, lineNumber int, level slog.Level, config CodeSnippetConfig) string {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Could not open source file: %s", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error reading source file: %s", err)
	}

	// Check if lineNumber is within range
	if lineNumber <= 0 || lineNumber > len(lines) {
		return fmt.Sprintf("Line number %d is out of range (file has %d lines)", lineNumber, len(lines))
	}

	// Create the snippet with context
	startLine := max(0, lineNumber-config.ContextLines-1)
	endLine := min(len(lines)-1, lineNumber+config.ContextLines-1)

	// Build the snippet
	var sb strings.Builder

	for i := startLine; i <= endLine; i++ {
		lineNum := i + 1
		linePrefix := fmt.Sprintf("%4d | ", lineNum)

		// Add highlighting for the target line based on log level
		if lineNum == lineNumber {
			var style lipgloss.Style
			var pointerText string

			switch level {
			case slog.LevelError:
				style = config.ErrorLineStyle
				pointerText = config.ErrorPointerText
			case slog.LevelWarn:
				style = config.WarnLineStyle
				pointerText = config.WarnPointerText
			case slog.LevelInfo:
				style = config.InfoLineStyle
				pointerText = config.InfoPointerText
			case slog.LevelDebug:
				style = config.DebugLineStyle
				pointerText = config.DebugPointerText
			default:
				style = config.ContextLineStyle
				pointerText = "^ Code location"
			}

			sb.WriteString(style.Render(linePrefix+lines[i]) + "\n")

			// Add the pointer under the line
			arrowPos := len(linePrefix)
			padding := strings.Repeat(" ", arrowPos)
			sb.WriteString(style.Render(padding+pointerText) + "\n")
		} else {
			sb.WriteString(config.ContextLineStyle.Render(linePrefix+lines[i]) + "\n")
		}
	}

	return sb.String()
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
