package auth

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
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

			util.Tick.WithWriter(cmd.OutOrStdout()).Println("Revoked all sessions and logged out")
			return nil
		},
	}
}

func revoke_all(cmd *cobra.Command) error {
	// Delete token from server
	api, token, err := util.NewAuthedClient()
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
