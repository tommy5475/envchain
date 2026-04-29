package differ_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/differ"
	"github.com/yourorg/envchain/internal/envset"
)

func makeSet(required, optional []string, defaults map[string]string) envset.EnvSet {
	if defaults == nil {
		defaults = map[string]string{}
	}
	return envset.EnvSet{
		Required: required,
		Optional: optional,
		Defaults: defaults,
	}
}

func TestDiff_NoChanges(t *testing.T) {
	a := makeSet([]string{"FOO", "BAR"}, nil, nil)
	b := makeSet([]string{"FOO", "BAR"}, nil, nil)

	result := differ.Diff("staging", a, "production", b)

	if result.HasChanges() {
		t.Errorf("expected no changes, got %d", len(result.Changes))
	}
}

func TestDiff_Added(t *testing.T) {
	a := makeSet([]string{"FOO"}, nil, nil)
	b := makeSet([]string{"FOO", "BAR"}, nil, nil)

	result := differ.Diff("staging", a, "production", b)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "BAR" && c.Kind == differ.Added {
			found = true
		}
	}
	if !found {
		t.Error("expected BAR to be marked as added")
	}
}

func TestDiff_Removed(t *testing.T) {
	a := makeSet([]string{"FOO", "LEGACY"}, nil, nil)
	b := makeSet([]string{"FOO"}, nil, nil)

	result := differ.Diff("staging", a, "production", b)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, c := range result.Changes {
		if c.Key == "LEGACY" && c.Kind == differ.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected LEGACY to be marked as removed")
	}
}

func TestDiff_Changed(t *testing.T) {
	a := makeSet([]string{"FOO"}, nil, map[string]string{"FOO": "old"})
	b := makeSet([]string{"FOO"}, nil, map[string]string{"FOO": "new"})

	result := differ.Diff("staging", a, "production", b)

	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	for _, c := range result.Changes {
		if c.Key == "FOO" {
			if c.Kind != differ.Changed {
				t.Errorf("expected Changed, got %s", c.Kind)
			}
			return
		}
	}
	t.Error("expected FOO in changes")
}

func TestDiff_StageNames(t *testing.T) {
	a := makeSet(nil, nil, nil)
	b := makeSet(nil, nil, nil)

	result := differ.Diff("alpha", a, "beta", b)

	if result.FromStage != "alpha" || result.ToStage != "beta" {
		t.Errorf("unexpected stage names: %s -> %s", result.FromStage, result.ToStage)
	}
}
