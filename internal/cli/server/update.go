package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpdateServerCommand struct{}

func (UpdateServerCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Rename a registered server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateRegisteredServer(cmd)
		},
	}
}

func updateRegisteredServer(cmd *cobra.Command) error {
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	serverId, err := util.PickServer(cmd.Context(), api, token)
	if err != nil {
		return err
	}

	name, err := util.Input("New name (leave blank to keep current)")
	if err != nil {
		return err
	}

	util.PrintDivider()

	request := contract.UpdateServerRequest{}
	if name = strings.TrimSpace(name); name != "" {
		request.Name = &name
	}

	if request.Name == nil {
		fmt.Fprintln(cmd.OutOrStdout(), pterm.FgGray.Sprint("No changes made: name left blank, so it stays the same."))
		return nil
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

	out := cmd.OutOrStdout()
	util.Tick.WithWriter(out).Println("Server updated")
	fmt.Fprintln(out)

	return printServer(out, response)
}
