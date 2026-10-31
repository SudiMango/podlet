package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ID               string `json:"id"`
	Token            string `json:"token"`
	APIURL           string `json:"api_url"`
	CfConnectorToken string `json:"connector_token,omitempty"`
}

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

func SaveAgentConfig(config *Config) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("get config directory: %w", err)
	}

	podletDir := filepath.Join(configDir, "podlet")
	if err := os.MkdirAll(podletDir, 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encode agent configuration: %w", err)
	}

	configPath := filepath.Join(podletDir, "server.json")
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("write agent configuration: %w", err)
	}

	return nil
}
