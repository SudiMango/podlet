package job

import (
	"fmt"
	"net/http"
	"net/url"
	"text/tabwriter"

	"github.com/SudiMango/podlet/internal/cli/client"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/SudiMango/podlet/internal/keyring"
	"github.com/spf13/cobra"
)

type GetJobsCommand struct{}

func (GetJobsCommand) Command() *cobra.Command {
	var status string
	var serverId string

	cmd := &cobra.Command{
		Use:   "get [job-id]",
		Short: "Check jobs",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return getJobs(cmd, serverId, status)
			}
			if status != "" || serverId != "" {
				return fmt.Errorf("--server and --status can only be used when listing jobs")
			}

			return getJob(cmd, args[0])
		},
	}

	cmd.Flags().StringVar(&serverId, "server", "", "Filter by server ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by pending, ongoing, success, or failed")
	return cmd
}

func getJobs(cmd *cobra.Command, serverId string, status string) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get jobs
	api, err := client.New()
	if err != nil {
		return err
	}

	path := "/jobs"
	if serverId != "" {
		path += "/" + url.PathEscape(serverId) + "/jobs"
	}
	if status != "" {
		path += "?" + url.Values{"status": []string{status}}.Encode()
	}

	var response contract.JobsResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		path,
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get jobs: %w", err)
	}

	// Print jobs to console
	out := cmd.OutOrStdout()
	if len(response) == 0 {
		fmt.Fprintln(out, "No jobs found.")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tSERVER ID\tREPOSITORY\tTYPE\tSTATUS\tCREATED AT\tUPDATED AT")

	for _, job := range response {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			job.Id,
			job.ServerId,
			job.RepoName,
			job.JobType,
			job.JobStatus,
			job.CreatedAt.Local().Format("2006-01-02 15:04"),
			job.UpdatedAt.Local().Format("2006-01-02 15:04"),
		)
	}

	return w.Flush()
}

func getJob(cmd *cobra.Command, jobId string) error {
	// Load existing token
	token, err := keyring.LoadToken()
	if err != nil {
		return fmt.Errorf("load session token: %w", err)
	}

	// Get job
	api, err := client.New()
	if err != nil {
		return err
	}

	var response contract.JobResponse
	if _, err := api.DoJSON(
		cmd.Context(),
		http.MethodGet,
		"/jobs/"+jobId,
		nil,
		&response,
		token,
	); err != nil {
		return fmt.Errorf("get job: %w", err)
	}

	// Print job to console
	out := cmd.OutOrStdout()
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", response.Id)
	fmt.Fprintf(w, "Server ID:\t%s\n", response.ServerId)
	fmt.Fprintf(w, "Repository:\t%s\n", response.RepoName)
	fmt.Fprintf(w, "Type:\t%s\n", response.JobType)
	fmt.Fprintf(w, "Status:\t%s\n", response.JobStatus)
	fmt.Fprintf(w, "Created At:\t%s\n", response.CreatedAt.Local().Format("2006-01-02 15:04"))
	fmt.Fprintf(w, "Updated At:\t%s\n", response.UpdatedAt.Local().Format("2006-01-02 15:04"))
	if response.ErrorMessage != nil {
		fmt.Fprintf(w, "Error:\t%s\n", *response.ErrorMessage)
	}

	return w.Flush()
}
