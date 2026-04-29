package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain/internal/snapshot"
)

func TestNew_CopiesVars(t *testing.T) {
	original := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := snapshot.New("production", original)

	original["FOO"] = "mutated"
	if s.Vars["FOO"] != "bar" {
		t.Errorf("expected snapshot to be isolated from original map")
	}
	if s.Stage != "production" {
		t.Errorf("expected stage 'production', got %q", s.Stage)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	vars := map[string]string{"API_KEY": "secret", "PORT": "8080"}
	s := snapshot.New("staging", vars)

	if err := snapshot.Save(path, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Stage != "staging" {
		t.Errorf("expected stage 'staging', got %q", loaded.Stage)
	}
	if loaded.Vars["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY='secret', got %q", loaded.Vars["API_KEY"])
	}
	if loaded.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT='8080', got %q", loaded.Vars["PORT"])
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json{"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestLoad_MissingStage(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nostage.json")
	os.WriteFile(path, []byte(`{"vars":{"FOO":"bar"}}`), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for missing stage field, got nil")
	}
}
