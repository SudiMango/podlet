package auth

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type LogoutCommand struct{}

func (LogoutCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout of Podlet",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := logout(cmd)
			if err != nil {
				return err
			}

			util.Tick.WithWriter(cmd.OutOrStdout()).Println("Logged out")
			return nil
		},
	}
}

func logout(cmd *cobra.Command) error {
	// Delete token from server
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/auth/signout",
		nil,
		nil,
		token,
	); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	// Delete token locally
	if err := keyring.DeleteToken(); err != nil {
		return fmt.Errorf("delete local session token: %w", err)
	}

	return nil
}
