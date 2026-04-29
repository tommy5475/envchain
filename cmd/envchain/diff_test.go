package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "envchain.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

const diffSampleConfig = `
chain:
  - dev
  - staging
  - production
stages:
  dev:
    required:
      - DB_HOST
      - APP_SECRET
  staging:
    required:
      - DB_HOST
      - APP_SECRET
      - SENTRY_DSN
  production:
    required:
      - DB_HOST
      - APP_SECRET
      - SENTRY_DSN
    optional:
      - FEATURE_FLAG
`

func TestRunDiff_ShowsAdded(t *testing.T) {
	p := writeTempConfig(t, diffSampleConfig)
	var buf bytes.Buffer
	if err := runDiff(p, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ SENTRY_DSN") {
		t.Errorf("expected SENTRY_DSN to appear as added, got:\n%s", out)
	}
}

func TestRunDiff_NoChanges(t *testing.T) {
	const cfg = `
chain:
  - a
  - b
stages:
  a:
    required:
      - FOO
  b:
    required:
      - FOO
`
	p := writeTempConfig(t, cfg)
	var buf bytes.Buffer
	if err := runDiff(p, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "(no changes)") {
		t.Errorf("expected no-changes message, got:\n%s", buf.String())
	}
}

func TestRunDiff_MissingFile(t *testing.T) {
	err := runDiff("/nonexistent/path.yaml", &bytes.Buffer{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
