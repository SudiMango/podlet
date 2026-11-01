package agent

import (
	"context"

	"github.com/SudiMango/podlet/internal/agent/dto"
	"github.com/SudiMango/podlet/internal/contract"
)

type DeployNewRepositoryChain struct{}

func (DeployNewRepositoryChain) Execute(ctx context.Context, a *Agent, job *contract.ClaimJobResponse) error {
	e := a.executors[dto.JobTypeCloneRepository]
	if err := e.Execute(ctx, job); err != nil {
		return err
	}

	e = a.executors[dto.JobTypeBuildRepository]
	if err := e.Execute(ctx, job); err != nil {
		return err
	}

	e = a.executors[dto.JobTypeCreateDnsRecord]
	if err := e.Execute(ctx, job); err != nil {
		return err
	}

	return nil
}
