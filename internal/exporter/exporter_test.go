package exporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/exporter"
)

func TestExporter_Shell(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(&buf, exporter.FormatShell)

	vars := map[string]string{"FOO": "bar", "BAZ": "hello world"}
	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "export FOO=bar") {
		t.Errorf("expected 'export FOO=bar' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export BAZ='hello world'") {
		t.Errorf("expected quoted BAZ in output, got:\n%s", out)
	}
}

func TestExporter_Dotenv(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(&buf, exporter.FormatDotenv)

	vars := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected 'DB_HOST=localhost' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected 'DB_PORT=5432' in output, got:\n%s", out)
	}
}

func TestExporter_JSON(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(&buf, exporter.FormatJSON)

	vars := map[string]string{"KEY": "value"}
	if err := e.Export(vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["KEY"] != "value" {
		t.Errorf("expected KEY=value in JSON output, got %v", result)
	}
}

func TestExporter_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	e := exporter.New(&buf, exporter.Format("xml"))

	err := e.Export(map[string]string{"X": "y"})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
