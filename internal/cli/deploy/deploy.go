package deploy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"text/tabwriter"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type DeployCommand struct{}

func (DeployCommand) Command() *cobra.Command {
	var detach bool

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a repository to a server, or redeploy a repository",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(cmd, detach)
		},
	}

	cmd.Flags().BoolVarP(&detach, "detach", "d", false, "Queue the deployment and exit without watching it")

	return cmd
}

func deploy(cmd *cobra.Command, detach bool) error {
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w; run `podlet auth login`", err)
	}

	api, err := client.New()
	if err != nil {
		return err
	}

	repoName, err := util.PickRepo(cmd.Context(), api, token)
	if err != nil {
		return err
	}

	serverId, err := util.PickServer(cmd.Context(), api, token)
	if err != nil {
		return err
	}

	util.PrintDivider()

	job, err := postDeploymentRequest(cmd, api, token, serverId, repoName)
	if err != nil {
		return err
	}

	if detach {
		printQueued(cmd, job, repoName, serverId)
		return nil
	}

	return watchJob(cmd, api, token, job.JobId)
}

func postDeploymentRequest(
	cmd *cobra.Command,
	api *client.Client,
	token string,
	serverId string,
	repoName string,
) (*contract.CreateJobResponse, error) {
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
		return nil, fmt.Errorf("post deploy request: %w", err)
	}

	return &response, nil
}

func printQueued(cmd *cobra.Command, job *contract.CreateJobResponse, repoName string, serverId string) {
	out := cmd.OutOrStdout()

	util.Tick.WithWriter(out).Println("Deployment queued")
	fmt.Fprintln(out)

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "  Repository\t%s\n", repoName)
	fmt.Fprintf(w, "  Server\t%s\n", serverId)
	fmt.Fprintf(w, "  Job\t%s\n", job.JobId)
	fmt.Fprintf(w, "  Status\t%s\n", job.Status)
	w.Flush()

	fmt.Fprintln(out, pterm.FgGray.Sprintf("\nCheck progress with `podlet jobs get %s`", job.JobId))
}

/*
	Job watching
*/

type jobView struct {
	current  int
	spinner  *pterm.SpinnerPrinter
	finished bool
}

func watchJob(cmd *cobra.Command, api *client.Client, token string, jobId string) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer stop()

	view := &jobView{current: -1}
	view.begin("Waiting for the server to pick up the job...")

	err := api.Stream(ctx, "/jobs/"+jobId+"/stream", token, func(data []byte) (bool, error) {
		var job contract.JobResponse
		if err := json.Unmarshal(data, &job); err != nil {
			return false, fmt.Errorf("decode job update: %w", err)
		}

		return view.update(job)
	})
	view.stop()

	switch {
	case ctx.Err() != nil:
		fmt.Fprintf(cmd.OutOrStdout(), "\nStopped watching. The deployment keeps running; check it with `podlet jobs get %s`.\n", jobId)
	case err != nil:
		return err
	case !view.finished:
		return fmt.Errorf("connection closed before the job finished; check it with `podlet jobs get %s`", jobId)
	}

	return nil
}

func (v *jobView) label(i int) string {
	return fmt.Sprintf("[%d/%d] %s", i+1, len(contract.ChainDeploySteps), contract.ChainDeploySteps[i])
}

func (v *jobView) begin(text string) {
	spinner := pterm.DefaultSpinner
	spinner.SuccessPrinter = &util.Tick
	spinner.FailPrinter = &util.Cross
	v.spinner, _ = spinner.Start(text)
}

func (v *jobView) stop() {
	if v.spinner != nil {
		v.spinner.Stop()
		v.spinner = nil
	}
}

func (v *jobView) advanceTo(i int) {
	if v.spinner != nil {
		if v.current < 0 {
			v.spinner.Success("Picked up by the server")
		} else {
			v.spinner.Success()
		}
		v.spinner = nil
	}

	for k := v.current + 1; k < i; k++ {
		util.Tick.Println(v.label(k))
	}

	v.current = i
	if i < len(contract.ChainDeploySteps) {
		v.begin(v.label(i))
	}
}

func (v *jobView) update(job contract.JobResponse) (done bool, err error) {
	if i := slices.Index(contract.ChainDeploySteps, job.Step); i > v.current {
		v.advanceTo(i)
	}

	switch job.JobStatus {
	case contract.StatusSuccess:
		v.finished = true
		v.advanceTo(len(contract.ChainDeploySteps))
		return true, nil
	case contract.StatusFailed:
		v.finished = true
		if v.spinner != nil {
			text := "Job failed"
			if v.current >= 0 {
				text = v.label(v.current) + " failed"
			}
			v.spinner.Fail(text)
			v.spinner = nil
		}

		message := "unknown error"
		if job.ErrorMessage != nil {
			message = *job.ErrorMessage
		}
		return true, fmt.Errorf("deploy failed: %s", message)
	}

	return false, nil
}
