package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/agent/jobs"
	"github.com/SudiMango/podlet/internal/contract"
)

type JobExecutor interface {
	Execute(ctx context.Context, job *contract.ClaimJobResponse) error
}

func JobRegistry(c *client.Client, cfg *config.Config) map[contract.JobType]JobExecutor {
	return map[contract.JobType]JobExecutor{
		contract.JobTypeEnsureDocker:        jobs.EnsureDockerJob{},
		contract.JobTypeInstallDocker:       jobs.InstallDockerJob{},
		contract.JobTypeEnsurePodletNetwork: jobs.EnsurePodletNetworkJob{},
		contract.JobTypeCreatePodletNetwork: jobs.CreatePodletNetworkJob{},
		contract.JobTypeEnsureCfTunnel:      jobs.NewEnsureCfTunnelJob(c, cfg),
		contract.JobTypeCreateCfTunnel:      jobs.NewCreateCfTunnelJob(c, cfg),
		contract.JobTypeEnsureCfContainer:   jobs.EnsureCfContainerJob{},
		contract.JobTypeCreateCfContainer:   jobs.NewCreateCfContainerJob(cfg),
		contract.JobTypeEnsureRailpack:      jobs.EnsureRailpackJob{},
		contract.JobTypeInstallRailpack:     jobs.InstallRailpackJob{},
		contract.JobTypeEnsureBuildkit:      jobs.EnsureBuildkitJob{},
		contract.JobTypeSetupBuildkit:       jobs.SetupBuildkitJob{},
	}
}

func (a *Agent) claimJob(ctx context.Context) (string, error) {
	var job *contract.ClaimJobResponse
	if err := a.client.DoJSON(
		ctx,
		http.MethodGet,
		"/agent/jobs/claim",
		nil,
		&job,
	); err != nil {
		return "", err
	}

	if job == nil {
		return "", nil
	}

	if err := a.reportJobClaimed(ctx, job.Id); err != nil {
		return "", err
	}

	a.busy.Store(true)
	go func() {
		a.executeJob(ctx, *job)
	}()

	return job.Id, nil
}

func (a *Agent) executeJob(ctx context.Context, job contract.ClaimJobResponse) {
	defer a.busy.Store(false)

	// Validate job type and perform job
	var err error
	switch job.Type {
	case contract.JobTypeCloneRepository:
		fmt.Print("asd")
	default:
		err = fmt.Errorf("unsupported job type %q", job.Type)
	}

	// Report errors
	if err != nil {
		if reportErr := a.reportJobFailed(ctx, job.Id, truncateError(err.Error())); reportErr != nil {
			fmt.Fprintf(os.Stderr, "report failed job %s: %v\n", job.Id, reportErr)
		}
		return
	}

	// Report success
	if err := a.reportJobSucceeded(ctx, job.Id); err != nil {
		fmt.Fprintf(os.Stderr, "report successful job %s: %v\n", job.Id, err)
	}
}

func (a *Agent) reportJob(ctx context.Context, report contract.ReportJobRequest) error {
	if err := a.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/jobs/report",
		report,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
	Helpers
*/

func (a *Agent) reportJobClaimed(ctx context.Context, jobId string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusOngoing,
		Error:     "",
	})
}

func (a *Agent) reportJobFailed(ctx context.Context, jobId string, err string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusFailed,
		Error:     err,
	})
}

func (a *Agent) reportJobSucceeded(ctx context.Context, jobId string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusSuccess,
		Error:     "",
	})
}

func truncateError(message string) string {
	const maxErrorLength = 4000
	if len(message) <= maxErrorLength {
		return message
	}

	return message[:maxErrorLength]
}
