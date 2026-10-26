package repos

import (
	"fmt"
	"net/http"
	"text/tabwriter"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
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
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	// Get repositories
	api, err := client.New()
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
		fmt.Fprintln(out, "No connected GitHub repositories found. Run `podlet repos connect` to grant access.")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "REPOSITORY\tVISIBILITY\tDEFAULT BRANCH\tCLONE URL")

	for _, repository := range response {
		visibility := "public"
		if repository.Private {
			visibility = "private"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", repository.FullName, visibility, repository.DefaultBranch, repository.CloneUrl)
	}

	return w.Flush()
}
