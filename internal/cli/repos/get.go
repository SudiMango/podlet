package repos

import (
	"fmt"
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type ListCommand struct{}

func (ListCommand) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "List which Github repositories Podlet can read",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRepos(cmd)
		},
	}
}

func listRepos(cmd *cobra.Command) error {
	// Get repositories
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.GithubRepositoriesResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/github/repositories",
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("list GitHub repositories: %w", err)
	}

	// Print repositories to the console
	out := cmd.OutOrStdout()
	if len(response) == 0 {
		fmt.Fprintln(out, pterm.FgGray.Sprint("No connected GitHub repositories found. Run `podlet repos connect` to grant access."))
		return nil
	}

	rows := make([][]string, len(response))
	for i, repository := range response {
		visibility := "public"
		if repository.Private {
			visibility = "private"
		}

		rows[i] = []string{repository.FullName, visibility, repository.DefaultBranch, repository.CloneUrl}
	}

	return util.PrintTable(out, []string{"REPOSITORY", "VISIBILITY", "DEFAULT BRANCH", "CLONE URL"}, rows)
}
