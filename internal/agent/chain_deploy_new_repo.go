package agent

import (
	"context"

	"github.com/SudiMango/podlet/internal/agent/dto"
	"github.com/SudiMango/podlet/internal/contract"
)

type DeployNewRepositoryChain struct{}

func (DeployNewRepositoryChain) Execute(ctx context.Context, a *Agent, job *contract.ClaimJobResponse) error {
	steps := []dto.ChainStep{
		{Name: contract.StepCloneRepository, JobType: dto.JobTypeCloneRepository},
		{Name: contract.StepBuildImage, JobType: dto.JobTypeBuildRepository},
		{Name: contract.StepCreateDnsRecord, JobType: dto.JobTypeCreateDnsRecord},
		{Name: contract.StepRemovePrevDeploy, JobType: dto.JobTypeRemovePrevDeploy},
	}

	for _, s := range steps {
		a.reportStep(ctx, job.Id, s.Name)
		if err := a.executors[s.JobType].Execute(ctx, job); err != nil {
			return err
		}
	}

	return nil
}
