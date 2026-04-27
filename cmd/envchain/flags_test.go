package main

import (
	"testing"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.configFile != "envchain.yaml" {
		t.Errorf("expected default config file, got %q", cfg.configFile)
	}
	if cfg.format != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.format)
	}
	if cfg.upTo != "" {
		t.Errorf("expected empty up-to, got %q", cfg.upTo)
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	cfg, err := parseFlags([]string{
		"-config", "custom.yaml",
		"-up-to", "staging",
		"-format", "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.configFile != "custom.yaml" {
		t.Errorf("expected custom.yaml, got %q", cfg.configFile)
	}
	if cfg.upTo != "staging" {
		t.Errorf("expected staging, got %q", cfg.upTo)
	}
	if cfg.format != "json" {
		t.Errorf("expected json, got %q", cfg.format)
	}
}

func TestParseFlags_InvalidFormat(t *testing.T) {
	_, err := parseFlags([]string{"-format", "xml"})
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "value"})
	if err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}
