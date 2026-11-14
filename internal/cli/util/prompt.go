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

type EmptyError struct{ Msg string }

func (e *EmptyError) Error() string { return e.Msg }

type DisplayError struct{ Err error }

func (e *DisplayError) Error() string { return e.Err.Error() }
func (e *DisplayError) Unwrap() error { return e.Err }

func Display(err error) error {
	if err == nil {
		return nil
	}
	return &DisplayError{Err: err}
}

func Pick(prompt string, options []string) (string, error) {
	if len(options) == 0 {
		return "", &EmptyError{Msg: "Nothing to choose from."}
	}
	return pterm.DefaultInteractiveSelect.WithDefaultText(prompt).WithOptions(options).Show()
}

func PickRepo(ctx context.Context, api *client.Client, token string) (string, error) {
	var repos contract.GithubRepositoriesResponse
	if _, err := api.DoJSON(ctx, http.MethodGet, "/github/repositories", nil, &repos, token); err != nil {
		return "", fmt.Errorf("list GitHub repositories: %w", err)
	}

	if len(repos) == 0 {
		return "", &EmptyError{Msg: "No connected GitHub repositories found. Run `podlet repos connect` to grant access."}
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

	if len(servers) == 0 {
		return "", &EmptyError{Msg: "No registered servers found. Add one with `podlet servers add`."}
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

func PickPendingServer(ctx context.Context, api *client.Client, token string) (string, error) {
	var servers contract.PendingServersResponse
	if _, err := api.DoJSON(ctx, http.MethodGet, "/servers/pending", nil, &servers, token); err != nil {
		return "", fmt.Errorf("get pending servers: %w", err)
	}

	if len(servers) == 0 {
		return "", &EmptyError{Msg: "No pending servers found."}
	}

	options := make([]string, len(servers))
	for i, sv := range servers {
		options[i] = fmt.Sprintf("%s (expires %s)", sv.Id, FormatTime(sv.ExpiresAt))
	}

	picked, err := Pick("Select a pending server", options)
	if err != nil {
		return "", err
	}

	return servers[slices.Index(options, picked)].Id, nil
}

func Input(prompt string) (string, error) {
	return pterm.DefaultInteractiveTextInput.WithDefaultText(prompt).Show()
}

func Confirm(prompt string) (bool, error) {
	return pterm.DefaultInteractiveConfirm.WithDefaultText(prompt).Show()
}
