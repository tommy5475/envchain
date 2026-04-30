package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/yourorg/envchain/internal/auditor"
)

// runAudit prints the audit log stored at auditPath to w.
// If format is "json" the raw entries are printed as a JSON array;
// otherwise a human-readable table is shown.
func runAudit(auditPath, format string, w io.Writer) error {
	a := auditor.New(auditPath)

	entries, err := a.ReadAll()
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	if len(entries) == 0 {
		fmt.Fprintln(w, "No audit entries found.")
		return nil
	}

	if format == "json" {
		return printAuditJSON(entries, w)
	}
	return printAuditText(entries, w)
}

func printAuditText(entries []auditor.Entry, w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tCONFIG\tPASSED\tSTAGES")
	for _, e := range entries {
		status := "✓"
		if !e.Passed {
			status = "✗"
		}
		stageCount := len(e.Results)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n",
			e.Timestamp.Format(time.RFC3339),
			e.ConfigFile,
			status,
			stageCount,
		)
	}
	return tw.Flush()
}

func printAuditJSON(entries []auditor.Entry, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(entries)
}

// auditMain is the entry point called from main when the "audit" subcommand
// is detected. It reads -audit-log and -format flags from args.
func auditMain(args []string, w io.Writer) int {
	auditLog := os.Getenv("ENVCHAIN_AUDIT_LOG")
	format := "text"

	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "-audit-log":
			auditLog = args[i+1]
		case "-format":
			format = args[i+1]
		}
	}

	if auditLog == "" {
		fmt.Fprintln(w, "error: -audit-log or ENVCHAIN_AUDIT_LOG required")
		return 1
	}

	if err := runAudit(auditLog, format, w); err != nil {
		fmt.Fprintf(w, "error: %v\n", err)
		return 1
	}
	return 0
}
