package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultDirname = ".hiden"

type Config struct {
	Dirname string `json:"dirname"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Dirname: defaultDirname,
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil
	}

	configPath := filepath.Join(homeDir, ".config", "hiden", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.Dirname == "" {
		cfg.Dirname = defaultDirname
	}

	return cfg, nil
}
