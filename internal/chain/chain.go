// Package chain orchestrates validation across ordered deployment stages.
package chain

import (
	"fmt"

	"github.com/user/envchain/internal/envset"
)

// StageResult holds the outcome of validating a single stage.
type StageResult struct {
	Stage   string
	OK      bool
	Missing []string
	Empty   []string
}

// Chain represents an ordered sequence of named environment stages.
type Chain struct {
	stages []stage
}

type stage struct {
	name string
	set  *envset.EnvSet
}

// New creates a Chain from an ordered list of (name, EnvSet) pairs.
func New(pairs ...interface{}) (*Chain, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("chain.New: arguments must be name/EnvSet pairs")
	}
	c := &Chain{}
	for i := 0; i < len(pairs); i += 2 {
		name, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("chain.New: expected string name at index %d", i)
		}
		set, ok := pairs[i+1].(*envset.EnvSet)
		if !ok {
			return nil, fmt.Errorf("chain.New: expected *envset.EnvSet at index %d", i+1)
		}
		c.stages = append(c.stages, stage{name: name, set: set})
	}
	return c, nil
}

// ValidateAll runs validation on every stage and returns all results.
func (c *Chain) ValidateAll() []StageResult {
	results := make([]StageResult, 0, len(c.stages))
	for _, s := range c.stages {
		res := s.set.Validate()
		results = append(results, StageResult{
			Stage:   s.name,
			OK:      res.OK,
			Missing: res.Missing,
			Empty:   res.Empty,
		})
	}
	return results
}

// ValidateUpTo runs validation stopping after the named stage or on first failure.
func (c *Chain) ValidateUpTo(stageName string) []StageResult {
	results := make([]StageResult, 0, len(c.stages))
	for _, s := range c.stages {
		res := s.set.Validate()
		r := StageResult{
			Stage:   s.name,
			OK:      res.OK,
			Missing: res.Missing,
			Empty:   res.Empty,
		}
		results = append(results, r)
		if !r.OK {
			break
		}
		if s.name == stageName {
			break
		}
	}
	return results
}
