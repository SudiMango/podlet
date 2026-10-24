package auth

import "github.com/spf13/cobra"

type AuthCommand struct{}

func (AuthCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}

	cmd.AddCommand(
		LoginCommand{}.Command(),
		LogoutCommand{}.Command(),
		MeCommand{}.Command(),
		RevokeSessionsCommand{}.Command(),
	)

	return cmd
}
