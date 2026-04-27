package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/reporter"
)

func makeResults() []chain.StageResult {
	return []chain.StageResult{
		{Stage: "dev", OK: true, Missing: nil, Empty: nil},
		{Stage: "staging", OK: false, Missing: []string{"DB_URL"}, Empty: []string{"API_KEY"}},
	}
}

func TestReporter_Text_PassStage(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	r.Report(makeResults())
	out := buf.String()
	if !strings.Contains(out, "✓ PASS") {
		t.Errorf("expected PASS marker, got:\n%s", out)
	}
	if !strings.Contains(out, "stage: dev") {
		t.Errorf("expected stage name 'dev', got:\n%s", out)
	}
}

func TestReporter_Text_FailStage(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	r.Report(makeResults())
	out := buf.String()
	if !strings.Contains(out, "✗ FAIL") {
		t.Errorf("expected FAIL marker, got:\n%s", out)
	}
	if !strings.Contains(out, "missing: DB_URL") {
		t.Errorf("expected missing var listed, got:\n%s", out)
	}
	if !strings.Contains(out, "empty:   API_KEY") {
		t.Errorf("expected empty var listed, got:\n%s", out)
	}
}

func TestReporter_JSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	r.Report(makeResults())
	out := buf.String()
	if !strings.Contains(out, `"stage":"dev"`) {
		t.Errorf("expected JSON stage field, got:\n%s", out)
	}
	if !strings.Contains(out, `"ok":true`) {
		t.Errorf("expected ok:true for dev, got:\n%s", out)
	}
	if !strings.Contains(out, `"ok":false`) {
		t.Errorf("expected ok:false for staging, got:\n%s", out)
	}
	if !strings.Contains(out, `"DB_URL"`) {
		t.Errorf("expected missing var in JSON, got:\n%s", out)
	}
}

func TestReporter_JSON_EmptySlices(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	r.Report([]chain.StageResult{
		{Stage: "prod", OK: true},
	})
	out := buf.String()
	if !strings.Contains(out, `"missing":[]`) {
		t.Errorf("expected empty missing array, got:\n%s", out)
	}
}
