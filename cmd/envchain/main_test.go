package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_MissingConfigFile(t *testing.T) {
	err := run([]string{"-config", "nonexistent.yaml"})
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestRun_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envchain.yaml")

	content := `chain:
  - name: dev
    vars:
      - name: APP_ENV
        required: true
  - name: prod
    vars:
      - name: APP_ENV
        required: true
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	t.Setenv("APP_ENV", "production")

	err := run([]string{"-config", cfgPath})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_FailingValidation(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envchain.yaml")

	content := `chain:
  - name: dev
    vars:
      - name: REQUIRED_BUT_MISSING
        required: true
`
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}

	os.Unsetenv("REQUIRED_BUT_MISSING")

	err := run([]string{"-config", cfgPath})
	if err == nil {
		t.Fatal("expected error for failing validation")
	}
}
