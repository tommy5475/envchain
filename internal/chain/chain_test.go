package chain_test

import (
	"os"
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/envset"
)

func makeEnvSet(required, optional []string) *envset.EnvSet {
	return envset.New(required, optional)
}

func TestChain_AllStagesPass(t *testing.T) {
	os.Setenv("APP_HOST", "localhost")
	os.Setenv("APP_PORT", "8080")
	defer os.Unsetenv("APP_HOST")
	defer os.Unsetenv("APP_PORT")

	c := chain.New()
	c.AddStage("dev", makeEnvSet([]string{"APP_HOST"}, nil))
	c.AddStage("staging", makeEnvSet([]string{"APP_PORT"}, nil))

	errs := c.Validate()
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}
}

func TestChain_OneStageFails(t *testing.T) {
	os.Setenv("APP_HOST", "localhost")
	defer os.Unsetenv("APP_HOST")
	os.Unsetenv("DB_URL")

	c := chain.New()
	c.AddStage("dev", makeEnvSet([]string{"APP_HOST"}, nil))
	c.AddStage("staging", makeEnvSet([]string{"DB_URL"}, nil))

	errs := c.Validate()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errs), errs)
	}
	if ve, ok := errs[0].(*chain.ValidationError); !ok || ve.Stage != "staging" {
		t.Errorf("expected staging stage error, got: %v", errs[0])
	}
}

func TestChain_ValidateUpTo_StopsAtFailure(t *testing.T) {
	os.Unsetenv("REQUIRED_VAR")

	c := chain.New()
	c.AddStage("dev", makeEnvSet([]string{"REQUIRED_VAR"}, nil))
	c.AddStage("staging", makeEnvSet([]string{"REQUIRED_VAR"}, nil))

	errs := c.ValidateUpTo("staging")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error (stopped at first), got %d", len(errs))
	}
}

func TestChain_ValidateUpTo_StopsAtNamedStage(t *testing.T) {
	os.Setenv("APP_HOST", "localhost")
	defer os.Unsetenv("APP_HOST")
	os.Unsetenv("DB_URL")

	c := chain.New()
	c.AddStage("dev", makeEnvSet([]string{"APP_HOST"}, nil))
	c.AddStage("staging", makeEnvSet([]string{"APP_HOST"}, nil))
	c.AddStage("prod", makeEnvSet([]string{"DB_URL"}, nil))

	errs := c.ValidateUpTo("staging")
	if len(errs) != 0 {
		t.Fatalf("expected no errors (prod not reached), got: %v", errs)
	}
}
