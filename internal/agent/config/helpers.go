package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadAgentConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(
		filepath.Join(configDir, "podlet", "server.json"),
	)
	if err != nil {
		return nil, fmt.Errorf("read agent configuration: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("decode agent configuration: %w", err)
	}

	if config.ID == "" || config.Token == "" || config.APIURL == "" {
		return nil, fmt.Errorf("agent configuration is incomplete")
	}

	return &config, nil
}
