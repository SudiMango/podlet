package cli

import (
	"fmt"
	"os"

	authcmd "github.com/SudiMango/podlet/internal/cli/auth"
	reposcmd "github.com/SudiMango/podlet/internal/cli/repos"
	servercmd "github.com/SudiMango/podlet/internal/cli/server"
	"github.com/spf13/cobra"
)

type Command interface {
	Command() *cobra.Command
}

type RootCommand struct{}

func (RootCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "podlet",
		Short:         "Podlet",
		Long:          "Podlet",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(
		HealthCommand{}.Command(),
		UninstallCommand{}.Command(),
		authcmd.AuthCommand{}.Command(),
		servercmd.ServerCommand{}.Command(),
		reposcmd.ReposCommand{}.Command(),
	)

	return cmd
}

func (r RootCommand) Execute() {
	if err := r.Command().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
