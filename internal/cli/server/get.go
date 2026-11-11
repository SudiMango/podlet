package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/spf13/cobra"
)

type GetServerCommand struct{}

func (GetServerCommand) Command() *cobra.Command {
	var pending bool

	cmd := &cobra.Command{
		Use:   "get [server-id]",
		Short: "Check your registed servers",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if pending {
				if len(args) == 0 {
					return getPendingServers(cmd)
				}

				return getPendingServer(cmd, args[0])
			}

			if len(args) == 0 {
				return getRegisteredServers(cmd)
			}

			return getRegisteredServer(cmd, args[0])
		},
	}

	cmd.Flags().BoolVarP(
		&pending,
		"pending",
		"p",
		false,
		"Get pending servers instead of registered servers",
	)

	return cmd
}

func getRegisteredServers(cmd *cobra.Command) error {
	// Get servers
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.ServersResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/servers",
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get registered servers: %w", err)
	}

	// Print servers to console
	out := cmd.OutOrStdout()
	if len(response) == 0 {
		fmt.Fprintln(out, "No registered servers found.")
		return nil
	}

	rows := make([][]string, len(response))
	for i, server := range response {
		rows[i] = []string{
			server.Id,
			server.Name,
			string(server.LifecycleState),
			util.FormatLastSeenAt(server.LastSeenAt),
			util.FormatTime(server.CreatedAt),
			util.FormatTime(server.UpdatedAt),
		}
	}

	return util.PrintTable(out, []string{"ID", "NAME", "LIFECYCLE STATE", "LAST SEEN AT", "CREATED AT", "UPDATED AT"}, rows)
}

func getRegisteredServer(cmd *cobra.Command, serverId string) error {
	// Get server
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.ServerResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/servers/"+serverId,
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get registered server: %w", err)
	}

	// Print server to console
	out := cmd.OutOrStdout()
	return util.PrintFields(out,
		"ID", response.Id,
		"Name", response.Name,
		"Lifecycle State", string(response.LifecycleState),
		"Last Seen At", util.FormatLastSeenAt(response.LastSeenAt),
		"Created At", util.FormatTime(response.CreatedAt),
		"Updated At", util.FormatTime(response.UpdatedAt),
	)
}

func getPendingServers(cmd *cobra.Command) error {
	// Get servers
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.PendingServersResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/servers/pending",
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get pending servers: %w", err)
	}

	// Print servers to console
	out := cmd.OutOrStdout()
	if len(response) == 0 {
		fmt.Fprintln(out, "No pending servers found.")
		return nil
	}

	rows := make([][]string, len(response))
	for i, server := range response {
		rows[i] = []string{server.Id, util.FormatTime(server.CreatedAt), util.FormatTime(server.ExpiresAt)}
	}

	return util.PrintTable(out, []string{"ID", "CREATED AT", "EXPIRES AT"}, rows)
}

func getPendingServer(cmd *cobra.Command, serverId string) error {
	// Get server
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.PendingServerResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/servers/pending/"+serverId,
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get pending server: %w", err)
	}

	// Print server to console
	out := cmd.OutOrStdout()
	return util.PrintFields(out,
		"ID", response.Id,
		"Created At", util.FormatTime(response.CreatedAt),
		"Expires At", util.FormatTime(response.ExpiresAt),
	)
}
