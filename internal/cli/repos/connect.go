package repos

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/spf13/cobra"
)

type ConnectCommand struct{}

func (ConnectCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "Authorize what GitHub repositories Podlet can read",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startInstallation(cmd)
		},
	}
}

func startInstallation(cmd *cobra.Command) error {
	// Start installation
	api, token, err := util.NewAuthedClient()
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
	fmt.Fprintf(out, "Your Installation URL: %s. It will expire by %s.\n", response.InstallationUrl, util.FormatTime(response.ExpiresAt))

	return nil
}
