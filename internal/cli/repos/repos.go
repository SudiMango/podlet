package repos

import "github.com/spf13/cobra"

type ReposCommand struct{}

func (ReposCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "Manage Github repositories",
	}

	cmd.AddCommand(
		ConnectCommand{}.Command(),
	)

	return cmd
}
