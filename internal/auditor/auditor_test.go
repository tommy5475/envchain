package auditor_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/envchain/internal/auditor"
)

func tempLog(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestRecord_And_ReadAll(t *testing.T) {
	path := tempLog(t)
	a := auditor.New(path)

	entry := auditor.Entry{
		Timestamp:  time.Now().UTC(),
		ConfigFile: "envchain.yaml",
		Passed:     true,
		Results: []auditor.StageResult{
			{Stage: "dev", Passed: true},
		},
	}

	if err := a.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := a.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].ConfigFile != "envchain.yaml" {
		t.Errorf("unexpected config file: %s", entries[0].ConfigFile)
	}
	if !entries[0].Passed {
		t.Error("expected passed=true")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	path := tempLog(t)
	a := auditor.New(path)

	for i := 0; i < 3; i++ {
		err := a.Record(auditor.Entry{
			Timestamp:  time.Now().UTC(),
			ConfigFile: "envchain.yaml",
			Passed:     i%2 == 0,
		})
		if err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	entries, err := a.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestReadAll_MissingFile(t *testing.T) {
	a := auditor.New("/nonexistent/path/audit.log")
	entries, err := a.ReadAll()
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file")
	}
}

func TestRecord_WithMissingVars(t *testing.T) {
	path := tempLog(t)
	a := auditor.New(path)

	entry := auditor.Entry{
		Timestamp:  time.Now().UTC(),
		ConfigFile: "envchain.yaml",
		Passed:     false,
		Results: []auditor.StageResult{
			{Stage: "prod", Passed: false, Missing: []string{"DB_PASSWORD", "API_KEY"}},
		},
	}

	if err := a.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := a.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	got := entries[0].Results[0].Missing
	if len(got) != 2 || got[0] != "DB_PASSWORD" {
		t.Errorf("unexpected missing vars: %v", got)
	}
}

func TestRecord_BadPath(t *testing.T) {
	a := auditor.New("/no/such/dir/audit.log")
	err := a.Record(auditor.Entry{Timestamp: time.Now().UTC()})
	if err == nil {
		t.Error("expected error writing to bad path")
	}
	_ = os.Remove("/no/such/dir/audit.log")
}
