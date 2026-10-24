package server

import "github.com/spf13/cobra"

type ServerCommand struct{}

func (ServerCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "Manage servers",
	}

	cmd.AddCommand(
		AddServerCommand{}.Command(),
		GetServerCommand{}.Command(),
		UpdateServerCommand{}.Command(),
		DeleteServerCommand{}.Command(),
	)

	return cmd
}
