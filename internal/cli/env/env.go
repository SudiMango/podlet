package env

import "github.com/spf13/cobra"

type EnvCommand struct{}

func (EnvCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage environment variables for a repository",
	}

	cmd.AddCommand(
		AddCommand{}.Command(),
		GetCommand{}.Command(),
		DeleteCommand{}.Command(),
	)

	return cmd
}
