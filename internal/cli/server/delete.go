package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type DeleteServerCommand struct{}

func (DeleteServerCommand) Command() *cobra.Command {
	var pending bool

	cmd := &cobra.Command{
		Use:   "delete [server-id]",
		Short: "Delete a registered or pending server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteServer(cmd, args[0], pending)
		},
	}

	cmd.Flags().BoolVarP(
		&pending,
		"pending",
		"p",
		false,
		"Delete a pending server instead of a registered server",
	)

	return cmd
}

func deleteServer(cmd *cobra.Command, serverId string, pending bool) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Delete server
	api, err := client.New()
	if err != nil {
		return err
	}

	path := "/servers/" + serverId
	serverType := "registered"
	if pending {
		path = "/servers/pending/" + serverId
		serverType = "pending"
	}

	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodDelete,
		path,
		nil,
		nil,
		token,
	); err != nil {
		return fmt.Errorf("delete %s server: %w", serverType, err)
	}

	if pending {
		fmt.Fprintf(cmd.OutOrStdout(), "pending server %s deleted successfully.\n", serverId)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "registered server %s marked for deletion. It will be removed when its agent checks in.\n", serverId)
	return nil
}
