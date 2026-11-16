package env

import (
	"net/http"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type AddCommand struct{}

func (AddCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an environment variable to a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, _ := cmd.Flags().GetString("repo")
			return util.Display(addEnvVar(cmd, repo))
		},
	}
	cmd.Flags().String("repo", "", "Full repository name, e.g. SudiMango/podlet-test")
	_ = cmd.MarkFlagRequired("repo")

	return cmd
}

func addEnvVar(cmd *cobra.Command, repo string) error {
	api, token, err := util.NewAuthedClient()
	if err != nil {
		return err
	}

	key, err := util.Input("Key")
	if err != nil {
		return err
	}

	// Masked so the value never shows on screen or lands in shell history
	value, err := pterm.DefaultInteractiveTextInput.WithDefaultText("Value").WithMask("*").Show()
	if err != nil {
		return err
	}

	util.PrintDivider()

	var response contract.EnvVarResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodPost,
		"/env",
		contract.CreateEnvVarRequest{Key: key, Value: value, RepoName: repo},
		&response,
		token,
	); err != nil {
		return err
	}

	util.Tick.WithWriter(cmd.OutOrStdout()).Printfln("Env var %s added to %s", response.Key, response.RepoName)
	return nil
}
