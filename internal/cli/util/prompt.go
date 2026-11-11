package util

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
)

func Pick(prompt string, options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("%s: nothing to choose from", prompt)
	}
	return pterm.DefaultInteractiveSelect.WithDefaultText(prompt).WithOptions(options).Show()
}

func PickRepo(ctx context.Context, api *client.Client, token string) (string, error) {
	var repos contract.GithubRepositoriesResponse
	if _, err := api.DoJSON(ctx, http.MethodGet, "/github/repositories", nil, &repos, token); err != nil {
		return "", fmt.Errorf("list GitHub repositories: %w", err)
	}

	options := make([]string, len(repos))
	for i, r := range repos {
		options[i] = r.FullName
	}

	return Pick("Select a repository", options)
}

func PickServer(ctx context.Context, api *client.Client, token string) (string, error) {
	var servers contract.ServersResponse
	if _, err := api.DoJSON(ctx, http.MethodGet, "/servers", nil, &servers, token); err != nil {
		return "", fmt.Errorf("get registered servers: %w", err)
	}

	options := make([]string, len(servers))
	for i, sv := range servers {
		options[i] = fmt.Sprintf("%s (%s)", sv.Name, sv.Id)
	}

	picked, err := Pick("Select a server", options)
	if err != nil {
		return "", err
	}

	return servers[slices.Index(options, picked)].Id, nil
}
