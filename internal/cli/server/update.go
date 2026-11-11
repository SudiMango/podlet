package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/spf13/cobra"
)

type UpdateServerCommand struct{}

func (UpdateServerCommand) Command() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "update [server-id]",
		Short: "Update a registered server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := contract.UpdateServerRequest{
				Name: &name,
			}

			err := updateRegisteredServer(cmd, args[0], request)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(
		&name,
		"name",
		"",
		"Update the name of a registered server",
	)

	cmd.MarkFlagRequired("name")

	return cmd
}

func updateRegisteredServer(cmd *cobra.Command, serverId string, request contract.UpdateServerRequest) error {
	// Update server
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.ServerResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPatch,
		"/servers/"+serverId,
		request,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("update server: %w", err)
	}

	// Print server to console
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "Server updated successfully!")

	return util.PrintFields(out,
		"ID", response.Id,
		"Name", response.Name,
		"Lifecycle State", string(response.LifecycleState),
		"Last Seen At", util.FormatLastSeenAt(response.LastSeenAt),
		"Created At", util.FormatTime(response.CreatedAt),
		"Updated At", util.FormatTime(response.UpdatedAt),
	)
}
