package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envchain/internal/chain"
	"github.com/yourorg/envchain/internal/loader"
	"github.com/yourorg/envchain/internal/reporter"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	def, err := loader.LoadFile(cfg.configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	envSets, err := loader.ToEnvSets(def)
	if err != nil {
		return fmt.Errorf("building env sets: %w", err)
	}

	c, err := chain.New(def.Chain, envSets)
	if err != nil {
		return fmt.Errorf("building chain: %w", err)
	}

	var results []chain.StageResult
	if cfg.upTo != "" {
		results = c.ValidateUpTo(cfg.upTo)
	} else {
		results = c.ValidateAll()
	}

	r := reporter.New(os.Stdout, cfg.format)
	r.Write(results)

	for _, res := range results {
		if !res.Passed {
			return fmt.Errorf("validation failed at stage: %s", res.Stage)
		}
	}
	return nil
}
