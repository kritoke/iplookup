package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/njreid/gokdl2"
)

type MaxMindConfig struct {
	AccountID  string `kdl:"account-id"`
	LicenseKey string `kdl:"license-key"`
	DbPath     string `kdl:"db-path"`
}

type Config struct {
	MaxMind MaxMindConfig `kdl:"maxmind"`
}

func defaultConfig() *Config {
	return &Config{MaxMind: MaxMindConfig{DbPath: "."}}
}

func applyDefaults(cfg *Config) {
	def := defaultConfig()
	if cfg.MaxMind.DbPath == "" {
		cfg.MaxMind.DbPath = def.MaxMind.DbPath
	}
}

func createDefaultConfig(path string) (*Config, error) {
	cfg := defaultConfig()
	data, err := kdl.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("marshaling default config: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return nil, fmt.Errorf("writing default config to %q: %w", path, err)
	}
	return cfg, nil
}

func LoadOrCreateConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("config file %q not found; creating a default", path)
			return createDefaultConfig(path)
		}
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := kdl.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config file %q: %w", path, err)
	}
	applyDefaults(&cfg)
	return &cfg, nil
}
