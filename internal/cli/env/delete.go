package env

import (
	"net/http"
	"net/url"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/spf13/cobra"
)

type DeleteCommand struct{}

func (DeleteCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete one or all environment variables of a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, _ := cmd.Flags().GetString("repo")
			key, _ := cmd.Flags().GetString("key")
			return util.Display(deleteEnvVars(cmd, repo, key))
		},
	}
	cmd.Flags().String("repo", "", "Full repository name, e.g. SudiMango/podlet-test")
	cmd.Flags().String("key", "", "Variable to delete (omit to delete all)")
	_ = cmd.MarkFlagRequired("repo")

	return cmd
}

func deleteEnvVars(cmd *cobra.Command, repo string, key string) error {
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	query := "?repo_name=" + url.QueryEscape(repo)
	path := "/env/all" + query
	prompt := "Delete ALL env vars for " + repo + "?"
	done := "All env vars deleted from " + repo
	if key != "" {
		path = "/env" + query + "&key=" + url.QueryEscape(key)
		prompt = "Delete " + key + " from " + repo + "?"
		done = "Env var " + key + " deleted from " + repo
	}

	ok, err := util.Confirm(prompt)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	util.PrintDivider()

	if _, err := api.DoJSON(cmd.Context(), http.MethodDelete, path, nil, nil, token); err != nil {
		return err
	}

	util.Tick.WithWriter(cmd.OutOrStdout()).Printfln("%s", done)
	return nil
}
