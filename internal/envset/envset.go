package envset

import (
	"errors"
	"fmt"
	"os"
)

// EnvSet represents a named set of environment variables required for a deployment stage.
type EnvSet struct {
	Name     string
	Required []string
	Optional []string
}

// ValidationError holds details about a missing or invalid environment variable.
type ValidationError struct {
	Key     string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("env var %q: %s", e.Key, e.Message)
}

// Validate checks that all required environment variables in the set are present and non-empty.
// Returns a slice of ValidationErrors, one per missing variable.
func (s *EnvSet) Validate() []error {
	var errs []error
	for _, key := range s.Required {
		val, ok := os.LookupEnv(key)
		if !ok {
			errs = append(errs, &ValidationError{Key: key, Message: "is not set"})
		} else if val == "" {
			errs = append(errs, &ValidationError{Key: key, Message: "is set but empty"})
		}
	}
	return errs
}

// Resolve returns a map of key→value for all variables (required + optional) that are present.
func (s *EnvSet) Resolve() map[string]string {
	result := make(map[string]string)
	for _, key := range append(s.Required, s.Optional...) {
		if val, ok := os.LookupEnv(key); ok {
			result[key] = val
		}
	}
	return result
}

// ErrValidation is a sentinel used to identify validation failures.
var ErrValidation = errors.New("envset validation failed")
