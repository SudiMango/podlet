package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/spf13/cobra"
)

type AddServerCommand struct{}

func (AddServerCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new server to your Podlet account",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := registerPendingServer(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func registerPendingServer(cmd *cobra.Command) error {
	// Create pending server
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.PendingServerResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/servers/add",
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("create pending server: %w", err)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Created pending server with id %s. Use the following token when registering your server: %s. It will expire in 10 minutes.\n", response.Id, *response.RegistrationToken)

	return nil
}
