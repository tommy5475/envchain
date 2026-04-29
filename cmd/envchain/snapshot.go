package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourorg/envchain/internal/loader"
	"github.com/yourorg/envchain/internal/snapshot"
)

// runSnapshot handles the 'snapshot' subcommand: saves the resolved env vars
// for a given stage to a JSON file for later diffing or auditing.
func runSnapshot(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)
	configFile := fs.String("config", "envchain.yaml", "path to config file")
	stage := fs.String("stage", "", "stage name to snapshot (required)")
	output := fs.String("output", "", "output file path (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *stage == "" {
		return fmt.Errorf("snapshot: --stage is required")
	}
	if *output == "" {
		return fmt.Errorf("snapshot: --output is required")
	}

	cfg, err := loader.LoadFile(*configFile)
	if err != nil {
		return fmt.Errorf("snapshot: failed to load config: %w", err)
	}

	envSets := loader.ToEnvSets(cfg)
	set, ok := envSets[*stage]
	if !ok {
		return fmt.Errorf("snapshot: stage %q not found in config", *stage)
	}

	resolved := make(map[string]string)
	for _, v := range set.Required {
		resolved[v] = os.Getenv(v)
	}
	for _, v := range set.Optional {
		if val := os.Getenv(v); val != "" {
			resolved[v] = val
		}
	}

	snap := snapshot.New(*stage, resolved)
	if err := snapshot.Save(*output, snap); err != nil {
		return fmt.Errorf("snapshot: %w", err)
	}

	fmt.Fprintf(stdout, "snapshot saved: stage=%s vars=%d file=%s\n",
		*stage, len(resolved), *output)
	return nil
}
