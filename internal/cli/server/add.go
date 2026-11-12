package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type AddServerCommand struct{}

func (AddServerCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a new server to your Podlet account",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return registerPendingServer(cmd)
		},
	}
}

func registerPendingServer(cmd *cobra.Command) error {
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
	util.Tick.WithWriter(out).Println("Pending server created")
	fmt.Fprintln(out)

	if err := util.PrintFields(out,
		"  ID", response.Id,
		"  Registration Token", *response.RegistrationToken,
		"  Expires At", util.FormatTime(response.ExpiresAt),
	); err != nil {
		return err
	}

	fmt.Fprintln(out, pterm.FgGray.Sprint("\nUse the registration token when installing the agent on your server."))
	return nil
}
