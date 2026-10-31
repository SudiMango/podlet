package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	initializejobs "github.com/SudiMango/podlet/internal/agent/jobs/system/initialize"
	uninstalljobs "github.com/SudiMango/podlet/internal/agent/jobs/system/uninstall"
	"github.com/SudiMango/podlet/internal/contract"
)

type JobExecutor interface {
	Execute(ctx context.Context, job *contract.ClaimJobResponse) error
}

func JobRegistry(c *client.Client, cfg *config.Config) map[contract.JobType]JobExecutor {
	return map[contract.JobType]JobExecutor{
		contract.JobTypeEnsureDocker:        initializejobs.EnsureDockerJob{},
		contract.JobTypeInstallDocker:       initializejobs.InstallDockerJob{},
		contract.JobTypeEnsurePodletNetwork: initializejobs.EnsurePodletNetworkJob{},
		contract.JobTypeCreatePodletNetwork: initializejobs.CreatePodletNetworkJob{},
		contract.JobTypeEnsureCfTunnel:      initializejobs.NewEnsureCfTunnelJob(c, cfg),
		contract.JobTypeCreateCfTunnel:      initializejobs.NewCreateCfTunnelJob(c, cfg),
		contract.JobTypeEnsureCfContainer:   initializejobs.EnsureCfContainerJob{},
		contract.JobTypeCreateCfContainer:   initializejobs.NewCreateCfContainerJob(cfg),
		contract.JobTypeEnsureRailpack:      initializejobs.EnsureRailpackJob{},
		contract.JobTypeInstallRailpack:     initializejobs.InstallRailpackJob{},
		contract.JobTypeEnsureBuildkit:      initializejobs.EnsureBuildkitJob{},
		contract.JobTypeSetupBuildkit:       initializejobs.SetupBuildkitJob{},
		contract.JobTypeRemoveCfContainer:   uninstalljobs.RemoveCfContainerJob{},
		contract.JobTypeRemoveBuildkit:      uninstalljobs.RemoveBuildkitJob{},
		contract.JobTypeRemovePodletNetwork: uninstalljobs.RemovePodletNetworkJob{},
		contract.JobTypeUninstallRailpack:   uninstalljobs.UninstallRailpackJob{},
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

	a.busy.Store(true)
	go a.executeJob(ctx, *job)

	return job.Id, nil
}

func (a *Agent) executeJob(ctx context.Context, job contract.ClaimJobResponse) {
	defer a.busy.Store(false)

	executor := a.executors[job.Type]
	err := executor.Execute(ctx, &job)

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
