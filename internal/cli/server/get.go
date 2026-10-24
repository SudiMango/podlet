package server

import (
	"fmt"
	"net/http"
	"text/tabwriter"
	"time"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
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
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get servers
	api, err := client.New()
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

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tNAME\tLIFECYCLE STATE\tLAST SEEN AT\tCREATED AT\tUPDATED AT")

	for _, server := range response {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			server.Id,
			server.Name,
			server.LifecycleState,
			formatLastSeenAt(server.LastSeenAt),
			server.CreatedAt.Local().Format("2006-01-02 15:04"),
			server.UpdatedAt.Local().Format("2006-01-02 15:04"),
		)
	}

	return w.Flush()
}

func getRegisteredServer(cmd *cobra.Command, serverId string) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get server
	api, err := client.New()
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
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", response.Id)
	fmt.Fprintf(w, "Name:\t%s\n", response.Name)
	fmt.Fprintf(w, "Lifecycle State:\t%s\n", response.LifecycleState)
	fmt.Fprintf(w, "Last Seen At:\t%s\n", formatLastSeenAt(response.LastSeenAt))
	fmt.Fprintf(w, "Created At:\t%s\n", response.CreatedAt.Local().Format("2006-01-02 15:04"))
	fmt.Fprintf(w, "Updated At:\t%s\n", response.UpdatedAt.Local().Format("2006-01-02 15:04"))

	return w.Flush()
}

func formatLastSeenAt(lastSeenAt *time.Time) string {
	if lastSeenAt == nil {
		return "Never"
	}

	return lastSeenAt.Local().Format("2006-01-02 15:04")
}

func getPendingServers(cmd *cobra.Command) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get servers
	api, err := client.New()
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

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tNAME\tCREATED AT\tEXPIRES AT")

	for _, server := range response {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			server.Id,
			server.CreatedAt.Local().Format("2006-01-02 15:04"),
			server.ExpiresAt.Local().Format("2006-01-02 15:04"),
		)
	}

	return w.Flush()
}

func getPendingServer(cmd *cobra.Command, serverId string) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get server
	api, err := client.New()
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
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", response.Id)
	fmt.Fprintf(w, "Created At:\t%s\n", response.CreatedAt.Local().Format("2006-01-02 15:04"))
	fmt.Fprintf(w, "Expires At:\t%s\n", response.ExpiresAt.Local().Format("2006-01-02 15:04"))

	return w.Flush()
}
