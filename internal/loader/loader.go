// Package loader provides functionality for loading envchain
// configuration from YAML files.
package loader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/yourorg/envchain/internal/envset"
)

// StageConfig represents a single deployment stage in the config file.
type StageConfig struct {
	Name     string   `yaml:"name"`
	Required []string `yaml:"required"`
	Optional []string `yaml:"optional"`
}

// Config represents the top-level envchain configuration file.
type Config struct {
	Chain  string        `yaml:"chain"`
	Stages []StageConfig `yaml:"stages"`
}

// LoadFile reads and parses an envchain YAML config from the given path.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: reading config file %q: %w", path, err)
	}
	return Parse(data)
}

// Parse unmarshals raw YAML bytes into a Config.
func Parse(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("loader: parsing config: %w", err)
	}
	if cfg.Chain == "" {
		return nil, fmt.Errorf("loader: config missing required field 'chain'")
	}
	if len(cfg.Stages) == 0 {
		return nil, fmt.Errorf("loader: config must define at least one stage")
	}
	return &cfg, nil
}

// ToEnvSets converts a Config into a slice of envset.EnvSet values,
// preserving stage order.
func ToEnvSets(cfg *Config) []envset.EnvSet {
	sets := make([]envset.EnvSet, 0, len(cfg.Stages))
	for _, s := range cfg.Stages {
		sets = append(sets, envset.EnvSet{
			Name:     s.Name,
			Required: s.Required,
			Optional: s.Optional,
		})
	}
	return sets
}
