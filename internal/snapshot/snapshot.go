// Package snapshot provides functionality for saving and loading
// environment variable snapshots to disk for later comparison or auditing.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of environment variables at a point in time.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Stage     string            `json:"stage"`
	Vars      map[string]string `json:"vars"`
}

// New creates a new Snapshot for the given stage and variable map.
func New(stage string, vars map[string]string) *Snapshot {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Stage:     stage,
		Vars:      copy,
	}
}

// Save writes the snapshot as JSON to the given file path.
func Save(path string, s *Snapshot) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}
	if s.Stage == "" {
		return nil, fmt.Errorf("snapshot: missing stage field")
	}
	if s.Vars == nil {
		s.Vars = make(map[string]string)
	}
	return &s, nil
}
