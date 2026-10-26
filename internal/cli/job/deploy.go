package job

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type DeployCommand struct{}

func (DeployCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy [server-id] [repo-name]",
		Short: "Deploy a specific repo to a specific server",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return postDeploymentRequest(cmd, args[0], args[1])
		},
	}
}

func postDeploymentRequest(cmd *cobra.Command, serverId string, repoName string) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Post a deploy request
	api, err := client.New()
	if err != nil {
		return err
	}

	var response contract.CreateJobResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/servers/deploy",
		contract.CreateJobRequest{
			ServerId: serverId,
			RepoName: repoName,
		},
		&response,
		token,
	); err != nil {
		return fmt.Errorf("post deploy request: %w", err)
	}

	// Print servers to console
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Deployment request for repository %s on server %s posted successfully!\n", repoName, serverId)
	fmt.Fprintln(out, "--------------------")
	fmt.Fprintf(out, "Job ID: %s\n", response.JobId)
	fmt.Fprintf(out, "Status: %s\n", response.Status)
	fmt.Fprintln(out, "--------------------")

	return nil
}
