package envset_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/envset"
)

func TestValidate_AllPresent(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")

	s := &envset.EnvSet{
		Name:     "database",
		Required: []string{"DB_HOST", "DB_PORT"},
	}

	errs := s.Validate()
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	// DB_PORT intentionally not set

	s := &envset.EnvSet{
		Name:     "database",
		Required: []string{"DB_HOST", "DB_PORT"},
	}

	errs := s.Validate()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
}

func TestValidate_EmptyRequired(t *testing.T) {
	t.Setenv("API_KEY", "")

	s := &envset.EnvSet{
		Name:     "api",
		Required: []string{"API_KEY"},
	}

	errs := s.Validate()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for empty var, got %d", len(errs))
	}
}

func TestValidate_MultipleMissing(t *testing.T) {
	// Neither DB_HOST nor DB_PORT are set
	s := &envset.EnvSet{
		Name:     "database",
		Required: []string{"DB_HOST", "DB_PORT", "DB_NAME"},
	}

	errs := s.Validate()
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors for all missing vars, got %d: %v", len(errs), errs)
	}
}

func TestResolve_IncludesOptional(t *testing.T) {
	t.Setenv("APP_HOST", "0.0.0.0")
	t.Setenv("APP_DEBUG", "true")
	// APP_METRICS not set — optional, should be absent from result

	s := &envset.EnvSet{
		Name:     "app",
		Required: []string{"APP_HOST"},
		Optional: []string{"APP_DEBUG", "APP_METRICS"},
	}

	resolved := s.Resolve()

	if resolved["APP_HOST"] != "0.0.0.0" {
		t.Errorf("expected APP_HOST=0.0.0.0, got %q", resolved["APP_HOST"])
	}
	if resolved["APP_DEBUG"] != "true" {
		t.Errorf("expected APP_DEBUG=true, got %q", resolved["APP_DEBUG"])
	}
	if _, ok := resolved["APP_METRICS"]; ok {
		t.Error("APP_METRICS should not be in resolved map when unset")
	}
}
