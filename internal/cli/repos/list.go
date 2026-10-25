package repos

import (
	"github.com/spf13/cobra"
)

type ListCommand struct{}

func (ListCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List which Github repositories Podlet can read",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRepos(cmd)
		},
	}
}

func listRepos(cmd *cobra.Command) error {
	// // Load existing token
	// token, err := keyring.LoadToken()
	// if err != nil {
	// 	return fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	// }

	// // Start installation
	// api, err := client.New()
	// if err != nil {
	// 	return err
	// }

	// var response contract.GithubInstallationStartResponse
	// if _, err := api.DoJSON(
	// 	cmd.Context(),
	// 	http.MethodPost,
	// 	"/github/installations/start",
	// 	nil,
	// 	&response,
	// 	token,
	// ); err != nil {
	// 	return fmt.Errorf("github installations start: %w", err)
	// }

	// out := cmd.OutOrStdout()
	// fmt.Fprintf(out, "Your Installation URL: %s. It will expire by %s.\n", response.InstallationUrl, response.ExpiresAt.Local().Format("2006-01-02 15:04"))

	return nil
}
