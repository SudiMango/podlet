package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type MeCommand struct{}

func (MeCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show the currently authenticated user",
		RunE: func(cmd *cobra.Command, args []string) error {
			profile, err := getCurrentProfile(cmd.Context())
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			return encoder.Encode(profile)
		},
	}
}

func getCurrentProfile(ctx context.Context) (*contract.PodletProfileResponse, error) {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	// Fetch current user
	api, err := client.New()
	if err != nil {
		return nil, err
	}

	var profile contract.PodletProfileResponse
	if _, err := api.DoJSON(
		ctx,
		http.MethodGet,
		"/profile",
		nil,
		&profile,
		token,
	); err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}

	return &profile, nil
}
