package server

import (
	"fmt"
	"net/http"
	"text/tabwriter"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
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
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Update server
	api, err := client.New()
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
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "Server updated successfully!")
	fmt.Fprintf(w, "ID:\t%s\n", response.Id)
	fmt.Fprintf(w, "Name:\t%s\n", response.Name)
	fmt.Fprintf(w, "Created At:\t%s\n", response.CreatedAt.Local().Format("2006-01-02 15:04"))
	fmt.Fprintf(w, "Updated At:\t%s\n", response.UpdatedAt.Local().Format("2006-01-02 15:04"))

	return w.Flush()
}
