package job

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/SudiMango/podlet/internal/cli/util"
	"github.com/SudiMango/podlet/internal/contract"
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
	// Get jobs
	api, token, err := util.NewAuthedClient()
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

	rows := make([][]string, len(response))
	for i, job := range response {
		rows[i] = []string{
			job.Id,
			job.ServerId,
			job.RepoName,
			string(job.JobType),
			string(job.JobStatus),
			util.FormatTime(job.CreatedAt),
			util.FormatTime(job.UpdatedAt),
		}
	}

	return util.PrintTable(out, []string{"ID", "SERVER ID", "REPOSITORY", "TYPE", "STATUS", "CREATED AT", "UPDATED AT"}, rows)
}

func getJob(cmd *cobra.Command, jobId string) error {
	// Get job
	api, token, err := util.NewAuthedClient()
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
	fields := []string{
		"ID", response.Id,
		"Server ID", response.ServerId,
		"Repository", response.RepoName,
		"Type", string(response.JobType),
		"Status", string(response.JobStatus),
		"Created At", util.FormatTime(response.CreatedAt),
		"Updated At", util.FormatTime(response.UpdatedAt),
	}
	if response.ErrorMessage != nil {
		fields = append(fields, "Error", *response.ErrorMessage)
	}

	return util.PrintFields(out, fields...)
}
