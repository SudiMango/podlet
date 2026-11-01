package agent

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/SudiMango/podlet/internal/agent/client"
	"github.com/SudiMango/podlet/internal/agent/config"
	"github.com/SudiMango/podlet/internal/agent/dto"
	runtimejobs "github.com/SudiMango/podlet/internal/agent/jobs/runtime"
	initializejobs "github.com/SudiMango/podlet/internal/agent/jobs/system/initialize"
	uninstalljobs "github.com/SudiMango/podlet/internal/agent/jobs/system/uninstall"
	"github.com/SudiMango/podlet/internal/contract"
)

type JobChain interface {
	Execute(ctx context.Context, agent *Agent, job *contract.ClaimJobResponse) error
}

type JobExecutor interface {
	Execute(ctx context.Context, job *contract.ClaimJobResponse) error
}

func JobRegistry(c *client.Client, cfg *config.Config) map[dto.JobType]JobExecutor {
	return map[dto.JobType]JobExecutor{
		dto.JobTypeEnsureDocker:        initializejobs.EnsureDockerJob{},
		dto.JobTypeInstallDocker:       initializejobs.InstallDockerJob{},
		dto.JobTypeEnsurePodletNetwork: initializejobs.EnsurePodletNetworkJob{},
		dto.JobTypeCreatePodletNetwork: initializejobs.CreatePodletNetworkJob{},
		dto.JobTypeEnsureCfTunnel:      initializejobs.NewEnsureCfTunnelJob(c, cfg),
		dto.JobTypeCreateCfTunnel:      initializejobs.NewCreateCfTunnelJob(c, cfg),
		dto.JobTypeEnsureCfContainer:   initializejobs.EnsureCfContainerJob{},
		dto.JobTypeCreateCfContainer:   initializejobs.NewCreateCfContainerJob(cfg),
		dto.JobTypeEnsureRailpack:      initializejobs.EnsureRailpackJob{},
		dto.JobTypeInstallRailpack:     initializejobs.InstallRailpackJob{},
		dto.JobTypeEnsureBuildkit:      initializejobs.EnsureBuildkitJob{},
		dto.JobTypeSetupBuildkit:       initializejobs.SetupBuildkitJob{},
		dto.JobTypeRemoveCfContainer:   uninstalljobs.RemoveCfContainerJob{},
		dto.JobTypeRemoveBuildkit:      uninstalljobs.RemoveBuildkitJob{},
		dto.JobTypeRemovePodletNetwork: uninstalljobs.RemovePodletNetworkJob{},
		dto.JobTypeUninstallRailpack:   uninstalljobs.UninstallRailpackJob{},
		dto.JobTypeCloneRepository:     runtimejobs.CloneRepositoryJob{},
		dto.JobTypeBuildRepository:     runtimejobs.BuildRepositoryJob{},
		dto.JobTypeCreateDnsRecord:     runtimejobs.NewCreateDNSRecordJob(c),
	}
}

func ChainRegistry() map[contract.JobType]JobChain {
	return map[contract.JobType]JobChain{
		contract.JobTypeDeployNewRepository: DeployNewRepositoryChain{},
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
	defer func() {
		a.busy.Store(false)
		select {
		case a.jobFinished <- struct{}{}:
		default:
		}
	}()

	executor := a.chains[job.Type]
	err := executor.Execute(ctx, a, &job)

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
