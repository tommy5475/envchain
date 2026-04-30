// Package auditor provides audit logging for envchain validation runs.
// It records which stages were validated, their outcomes, and timestamps
// so operators can review historical validation activity.
package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log record for one validation run.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	ConfigFile string         `json:"config_file"`
	Results    []StageResult  `json:"results"`
	Passed     bool           `json:"passed"`
}

// StageResult holds the outcome of validating a single stage.
type StageResult struct {
	Stage   string   `json:"stage"`
	Passed  bool     `json:"passed"`
	Missing []string `json:"missing,omitempty"`
}

// Auditor writes audit entries to a log file in newline-delimited JSON.
type Auditor struct {
	path string
}

// New creates an Auditor that appends entries to the file at path.
func New(path string) *Auditor {
	return &Auditor{path: path}
}

// Record appends an audit Entry to the log file.
func (a *Auditor) Record(entry Entry) error {
	f, err := os.OpenFile(a.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("auditor: open log file: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("auditor: marshal entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("auditor: write entry: %w", err)
	}
	return nil
}

// ReadAll reads all audit entries from the log file.
func (a *Auditor) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(a.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("auditor: read log file: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(
		// wrap raw bytes in a reader via strings trick
		newBytesReader(data),
	)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("auditor: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// bytesReader wraps []byte to satisfy io.Reader for json.NewDecoder.
type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader { return &bytesReader{data: data} }

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}
