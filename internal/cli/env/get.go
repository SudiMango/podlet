package env

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type GetCommand struct{}

func (GetCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "List the environment variables of a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, _ := cmd.Flags().GetString("repo")
			return util.Display(getEnvVars(cmd, repo))
		},
	}
	cmd.Flags().String("repo", "", "Full repository name, e.g. SudiMango/podlet-test")
	_ = cmd.MarkFlagRequired("repo")

	return cmd
}

func getEnvVars(cmd *cobra.Command, repo string) error {
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	var response contract.EnvVarsResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/env?repo_name="+url.QueryEscape(repo),
		nil,
		&response,
		token,
	); err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if len(response) == 0 {
		fmt.Fprintln(out, pterm.FgGray.Sprint("No env vars found for "+repo+"."))
		return nil
	}

	rows := make([][]string, len(response))
	for i, envVar := range response {
		rows[i] = []string{envVar.Key, util.FormatTime(envVar.CreatedAt)}
	}

	return util.PrintTable(out, []string{"KEY", "CREATED"}, rows)
}
