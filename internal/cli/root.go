package cli

import (
	"errors"
	"fmt"
	"os"

	authcmd "github.com/SudiMango/podlet/internal/cli/auth"
	deploycmd "github.com/SudiMango/podlet/internal/cli/deploy"
	envcmd "github.com/SudiMango/podlet/internal/cli/env"
	jobcmd "github.com/SudiMango/podlet/internal/cli/job"
	reposcmd "github.com/SudiMango/podlet/internal/cli/repos"
	servercmd "github.com/SudiMango/podlet/internal/cli/server"
	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/pterm/pterm"
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
		envcmd.EnvCommand{}.Command(),
		jobcmd.JobCommand{}.Command(),
		deploycmd.DeployCommand{}.Command(),
	)

	return cmd
}

func (r RootCommand) Execute() {
	if err := r.Command().Execute(); err != nil {
		var empty *util.EmptyError
		if errors.As(err, &empty) {
			fmt.Fprintln(os.Stderr, pterm.FgGray.Sprint(empty.Msg))
			return
		}

		var display *util.DisplayError
		if errors.As(err, &display) {
			util.Cross.WithWriter(os.Stderr).Println(display.Error())
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
