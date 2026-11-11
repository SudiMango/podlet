package deploy

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeployCommand struct{}

func (DeployCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a repository to a server, or redeploy a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(cmd)
		},
	}
}

func deploy(cmd *cobra.Command) error {
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	api, err := client.New()
	if err != nil {
		return err
	}

	repoName, err := pickRepo(cmd, api, token)
	if err != nil {
		return err
	}

	serverId, err := pickServer(cmd, api, token)
	if err != nil {
		return err
	}

	return postDeploymentRequest(cmd, serverId, repoName)
}

func pickRepo(cmd *cobra.Command, api *client.Client, token string) (string, error) {
	var repos contract.GithubRepositoriesResponse
	if _, err := api.DoJSON(cmd.Context(), http.MethodGet, "/github/repositories", nil, &repos, token); err != nil {
		return "", fmt.Errorf("list GitHub repositories: %w", err)
	}

	options := make([]string, len(repos))
	for i, r := range repos {
		options[i] = r.FullName
	}

	return pick("Select a repository", options)
}

func pickServer(cmd *cobra.Command, api *client.Client, token string) (string, error) {
	var servers contract.ServersResponse
	if _, err := api.DoJSON(cmd.Context(), http.MethodGet, "/servers", nil, &servers, token); err != nil {
		return "", fmt.Errorf("get registered servers: %w", err)
	}

	options := make([]string, len(servers))
	for i, sv := range servers {
		options[i] = fmt.Sprintf("%s (%s)", sv.Name, sv.Id)
	}

	picked, err := pick("Select a server", options)
	if err != nil {
		return "", err
	}

	return servers[slices.Index(options, picked)].Id, nil
}

// pick shows an arrow-key menu (type to filter) and returns the chosen option.
func pick(prompt string, options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("%s: nothing to choose from", prompt)
	}
	return pterm.DefaultInteractiveSelect.WithDefaultText(prompt).WithOptions(options).Show()
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
