// formatter.go
package golog

import "log/slog"

// Formatter defines how to convert a slog.Record into styled output.
type Formatter interface {
	// Format returns a formatted byte slice for a given slog record.
	Format(r slog.Record) ([]byte, error)
}
