package main

import (
	"os"

	"github.com/njreid/gokdl2"
)

type Config struct {
	MaxMind struct {
		AccountID  string `kdl:"account-id"`
		LicenseKey string `kdl:"license-key"`
		DbPath     string `kdl:"db-path"`
	} `kdl:"maxmind"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := kdl.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
