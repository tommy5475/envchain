package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
)

type config struct {
	configFile string
	upTo       string
	format     string
}

func parseFlags(args []string) (*config, error) {
	fs := flag.NewFlagSet("envchain", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &config{}

	fs.StringVar(&cfg.configFile, "config", "envchain.yaml", "path to envchain config file")
	fs.StringVar(&cfg.upTo, "up-to", "", "validate only up to and including this stage")
	fs.StringVar(&cfg.format, "format", "text", "output format: text or json")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, err
		}
		return nil, fmt.Errorf("parsing flags: %w", err)
	}

	if cfg.format != "text" && cfg.format != "json" {
		return nil, fmt.Errorf("invalid format %q: must be text or json", cfg.format)
	}

	return cfg, nil
}
