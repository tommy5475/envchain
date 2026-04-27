package chain

import (
	"fmt"

	"github.com/envchain/envchain/internal/envset"
)

// Stage represents a named deployment stage with its required environment set.
type Stage struct {
	Name   string
	EnvSet *envset.EnvSet
}

// Chain holds an ordered sequence of stages to validate in order.
type Chain struct {
	Stages []*Stage
}

// New creates a new Chain with no stages.
func New() *Chain {
	return &Chain{}
}

// AddStage appends a stage to the chain.
func (c *Chain) AddStage(name string, es *envset.EnvSet) {
	c.Stages = append(c.Stages, &Stage{Name: name, EnvSet: es})
}

// ValidationError captures which stage failed and why.
type ValidationError struct {
	Stage string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("stage %q failed validation: %v", e.Stage, e.Err)
}

// Validate runs validation on every stage in order.
// It returns a slice of ValidationErrors for all failing stages.
func (c *Chain) Validate() []error {
	var errs []error
	for _, stage := range c.Stages {
		if err := stage.EnvSet.Validate(); err != nil {
			errs = append(errs, &ValidationError{Stage: stage.Name, Err: err})
		}
	}
	return errs
}

// ValidateUpTo validates stages in order, stopping at the first failure.
func (c *Chain) ValidateUpTo(stageName string) []error {
	var errs []error
	for _, stage := range c.Stages {
		if err := stage.EnvSet.Validate(); err != nil {
			errs = append(errs, &ValidationError{Stage: stage.Name, Err: err})
			break
		}
		if stage.Name == stageName {
			break
		}
	}
	return errs
}
