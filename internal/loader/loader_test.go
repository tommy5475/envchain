package loader_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/loader"
)

var validYAML = []byte(`
chain: my-app
stages:
  - name: dev
    required:
      - APP_ENV
      - DB_URL
    optional:
      - LOG_LEVEL
  - name: prod
    required:
      - APP_ENV
      - DB_URL
      - SECRET_KEY
`)

func TestParse_Valid(t *testing.T) {
	cfg, err := loader.Parse(validYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Chain != "my-app" {
		t.Errorf("expected chain 'my-app', got %q", cfg.Chain)
	}
	if len(cfg.Stages) != 2 {
		t.Fatalf("expected 2 stages, got %d", len(cfg.Stages))
	}
	if cfg.Stages[0].Name != "dev" {
		t.Errorf("expected first stage 'dev', got %q", cfg.Stages[0].Name)
	}
}

func TestParse_MissingChain(t *testing.T) {
	data := []byte(`stages:\n  - name: dev\n    required: [FOO]\n`)
	_, err := loader.Parse(data)
	if err == nil {
		t.Fatal("expected error for missing 'chain' field, got nil")
	}
}

func TestParse_NoStages(t *testing.T) {
	data := []byte("chain: empty-app\nstages: []\n")
	_, err := loader.Parse(data)
	if err == nil {
		t.Fatal("expected error for empty stages, got nil")
	}
}

func TestToEnvSets_CorrectMapping(t *testing.T) {
	cfg, err := loader.Parse(validYAML)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	sets := loader.ToEnvSets(cfg)
	if len(sets) != len(cfg.Stages) {
		t.Fatalf("expected %d env sets, got %d", len(cfg.Stages), len(sets))
	}
	for i, s := range cfg.Stages {
		if sets[i].Name != s.Name {
			t.Errorf("stage %d: expected name %q, got %q", i, s.Name, sets[i].Name)
		}
		if len(sets[i].Required) != len(s.Required) {
			t.Errorf("stage %d: required length mismatch", i)
		}
	}
}
