package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	BaseURL string `json:"base_url"`
}

func LoadBaseURL() (string, error) {
	if baseURL := strings.TrimRight(os.Getenv("BASE_URL"), "/"); baseURL != "" {
		return baseURL, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("find user configuration directory: %w", err)
	}

	configPath := filepath.Join(configDir, "podlet", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("Podlet CLI configuration not found at %s", configPath)
		}
		return "", fmt.Errorf("read Podlet CLI configuration: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("decode Podlet CLI configuration: %w", err)
	}

	baseURL := strings.TrimRight(config.BaseURL, "/")
	if baseURL == "" {
		return "", errors.New("Podlet CLI configuration is missing base_url")
	}

	return baseURL, nil
}
