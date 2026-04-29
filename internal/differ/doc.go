// Package differ provides utilities for comparing two EnvSets and
// surfacing the differences between them as structured Change values.
//
// It is used by the envchain CLI to implement the `diff` sub-command,
// which lets operators understand what environment variables are
// added, removed, or modified when moving from one deployment stage
// to the next.
//
// Example usage:
//
//	result := differ.Diff("staging", stagingSet, "production", productionSet)
//	for _, c := range result.Changes {
//		fmt.Printf("%s %s\n", c.Kind, c.Key)
//	}
package differ
