package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envchain/internal/differ"
	"github.com/yourorg/envchain/internal/loader"
)

// runDiff loads the config file and prints differences between adjacent stages.
func runDiff(configPath string, out io.Writer) error {
	cfg, err := loader.LoadFile(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	envSets, err := loader.ToEnvSets(cfg)
	if err != nil {
		return fmt.Errorf("building env sets: %w", err)
	}

	if len(cfg.Chain) < 2 {
		fmt.Fprintln(out, "no adjacent stages to compare")
		return nil
	}

	for i := 0; i < len(cfg.Chain)-1; i++ {
		fromName := cfg.Chain[i]
		toName := cfg.Chain[i+1]

		fromSet, ok1 := envSets[fromName]
		toSet, ok2 := envSets[toName]
		if !ok1 || !ok2 {
			continue
		}

		result := differ.Diff(fromName, fromSet, toName, toSet)
		printDiffResult(out, result)
	}

	return nil
}

func printDiffResult(out io.Writer, result differ.Result) {
	fmt.Fprintf(out, "diff %s → %s\n", result.FromStage, result.ToStage)
	if !result.HasChanges() {
		fmt.Fprintln(out, "  (no changes)")
		return
	}
	for _, c := range result.Changes {
		switch c.Kind {
		case differ.Added:
			fmt.Fprintf(out, "  + %s\n", c.Key)
		case differ.Removed:
			fmt.Fprintf(out, "  - %s\n", c.Key)
		case differ.Changed:
			fmt.Fprintf(out, "  ~ %s (%s → %s)\n", c.Key, c.From, c.To)
		}
	}
}

// diffMain is the entry point when the diff sub-command is invoked.
func diffMain(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: envchain diff <config.yaml>")
		os.Exit(1)
	}
	if err := runDiff(args[0], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
