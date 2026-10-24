package auth

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type RevokeSessionsCommand struct{}

func (RevokeSessionsCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "revokeall",
		Short: "Revoke all your Podlet active sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := revoke_all(cmd)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Revoked all sessions and logged out successfully!")
			return nil
		},
	}
}

func revoke_all(cmd *cobra.Command) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Delete token from server
	api, err := client.New()
	if err != nil {
		return err
	}

	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/auth/revoke-sessions",
		nil,
		nil,
		token,
	); err != nil {
		return fmt.Errorf("revoke sessions: %w", err)
	}

	// Delete token locally
	if err := keyring.DeleteToken(); err != nil {
		return fmt.Errorf("delete local session token: %w", err)
	}

	return nil
}
