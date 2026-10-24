package cli

import (
	"fmt"
	"os"

	authcmd "github.com/SudiMango/podlet/internal/cli/auth"
	"github.com/spf13/cobra"
)

type Command interface {
	Command() *cobra.Command
}

type RootCommand struct{}

func (RootCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "podlet",
		Short: "Podlet",
		Long:  "Podlet",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(
		HealthCommand{}.Command(),
		authcmd.AuthCommand{}.Command(),
	)

	return cmd
}

func (r RootCommand) Execute() {
	if err := r.Command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
