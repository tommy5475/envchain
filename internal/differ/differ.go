// Package differ compares environment variable sets across two stages
// and reports additions, removals, and changes between them.
package differ

import (
	"fmt"

	"github.com/yourorg/envchain/internal/envset"
)

// ChangeKind describes the type of difference between two stages.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single variable difference between two stages.
type Change struct {
	Key  string
	Kind ChangeKind
	// From is the previous value (empty string if Added).
	From string
	// To is the new value (empty string if Removed).
	To string
}

// Result holds the diff output between two named stages.
type Result struct {
	FromStage string
	ToStage   string
	Changes   []Change
}

// HasChanges returns true when at least one difference was detected.
func (r Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Diff compares two EnvSets and returns a Result describing the differences.
// Only keys that are present in either set are compared; the resolved
// environment (os.Getenv) is used for actual values.
func Diff(fromName string, from envset.EnvSet, toName string, to envset.EnvSet) Result {
	result := Result{
		FromStage: fromName,
		ToStage:   toName,
	}

	fromKeys := keySet(from)
	toKeys := keySet(to)

	for key := range toKeys {
		if _, exists := fromKeys[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:  key,
				Kind: Added,
				From: "",
				To:   resolveValue(to, key),
			})
		}
	}

	for key := range fromKeys {
		if _, exists := toKeys[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:  key,
				Kind: Removed,
				From: resolveValue(from, key),
				To:   "",
			})
		} else {
			fromVal := resolveValue(from, key)
			toVal := resolveValue(to, key)
			if fromVal != toVal {
				result.Changes = append(result.Changes, Change{
					Key:  key,
					Kind: Changed,
					From: fromVal,
					To:   toVal,
				})
			}
		}
	}

	return result
}

func keySet(es envset.EnvSet) map[string]struct{} {
	keys := make(map[string]struct{})
	for _, v := range es.Required {
		keys[v] = struct{}{}
	}
	for _, v := range es.Optional {
		keys[v] = struct{}{}
	}
	return keys
}

func resolveValue(es envset.EnvSet, key string) string {
	if val, ok := es.Defaults[key]; ok {
		return fmt.Sprintf("<default:%s>", val)
	}
	return "<unset>"
}
