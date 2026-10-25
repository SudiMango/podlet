package auth

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type ReposCommand struct{}

func (ReposCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "repos",
		Short: "Authorize what GitHub repositories you want Podlet to have access to",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startInstallation(cmd)
		},
	}
}

func startInstallation(cmd *cobra.Command) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	// Start installation
	api, err := client.New()
	if err != nil {
		return err
	}

	var response contract.GithubInstallationStartResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/github/installations/start",
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("github installations start: %w", err)
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Your Installation URL: %s. It will expire by %s.\n", response.InstallationUrl, response.ExpiresAt.Local().Format("2006-01-02 15:04"))

	return nil
}
