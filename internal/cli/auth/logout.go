package auth

import (
	"fmt"

	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type LogoutCommand struct{}

func (LogoutCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Logout of Podlet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := keyring.DeleteToken(); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Failed to logout: %v\n", err)
				return err
			}

			fmt.Println("Logged out of Podlet successfully!")
			return nil
		},
	}
}
