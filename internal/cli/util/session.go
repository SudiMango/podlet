package util

import (
	"fmt"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/keyring"
)

func NewAuthedClient() (*client.Client, string, error) {
	token, err := keyring.LoadToken()
	if err != nil {
		return nil, "", fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	api, err := client.New()
	if err != nil {
		return nil, "", err
	}

	return api, token, nil
}
