// Package exporter provides functionality for exporting validated
// environment variable sets to various output formats (shell, dotenv, JSON).
package exporter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// Exporter writes environment variables to an output stream.
type Exporter struct {
	w      io.Writer
	format Format
}

// New creates a new Exporter writing to w in the given format.
func New(w io.Writer, format Format) *Exporter {
	return &Exporter{w: w, format: format}
}

// Export writes the provided key-value environment variables to the output.
func (e *Exporter) Export(vars map[string]string) error {
	switch e.format {
	case FormatShell:
		return e.writeShell(vars)
	case FormatDotenv:
		return e.writeDotenv(vars)
	case FormatJSON:
		return e.writeJSON(vars)
	default:
		return fmt.Errorf("unsupported export format: %q", e.format)
	}
}

func (e *Exporter) writeShell(vars map[string]string) error {
	for k, v := range vars {
		_, err := fmt.Fprintf(e.w, "export %s=%s\n", k, shellQuote(v))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeDotenv(vars map[string]string) error {
	for k, v := range vars {
		_, err := fmt.Fprintf(e.w, "%s=%s\n", k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(vars map[string]string) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(vars)
}

// shellQuote wraps a value in single quotes if it contains special characters.
func shellQuote(v string) string {
	if strings.ContainsAny(v, " \t\n$\"'\\`") {
		return "'" + strings.ReplaceAll(v, "'", "'\\'''") + "'"
	}
	return v
}
