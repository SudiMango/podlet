package server

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeleteServerCommand struct{}

func (DeleteServerCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete a registered or pending server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteServer(cmd)
		},
	}
}

func deleteServer(cmd *cobra.Command) error {
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	kind, err := util.Pick("Which kind of server?", []string{"Registered", "Pending"})
	if err != nil {
		return err
	}
	pending := kind == "Pending"

	pick := util.PickServer
	path := "/servers/"
	if pending {
		pick = util.PickPendingServer
		path = "/servers/pending/"
	}

	serverId, err := pick(cmd.Context(), api, token)
	if err != nil {
		return err
	}

	ok, err := util.Confirm("Delete " + serverId + "?")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	util.PrintDivider()

	if _, err := api.DoJSON(cmd.Context(), http.MethodDelete, path+serverId, nil, nil, token); err != nil {
		return fmt.Errorf("delete %s server: %w", kind, err)
	}

	out := cmd.OutOrStdout()
	if pending {
		util.Tick.WithWriter(out).Printfln("Pending server %s deleted", serverId)
		return nil
	}

	util.Tick.WithWriter(out).Printfln("Server %s marked for deletion", serverId)
	fmt.Fprintln(out, pterm.FgGray.Sprint("It will be removed when its agent checks in."))
	return nil
}
