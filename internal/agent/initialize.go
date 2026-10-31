package agent

import (
	"context"
	"errors"

	"github.com/SudiMango/podlet/internal/agent/jobs"
	"github.com/SudiMango/podlet/internal/contract"
)

func (a *Agent) Initialize(ctx context.Context) error {
	// Docker setup
	e := a.executors[contract.JobTypeEnsureDocker]
	if err := e.Execute(ctx, nil); err != nil {
		e = a.executors[contract.JobTypeInstallDocker]

		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	// Docker network setup
	e = a.executors[contract.JobTypeEnsurePodletNetwork]
	if err := e.Execute(ctx, nil); err != nil {
		e = a.executors[contract.JobTypeCreatePodletNetwork]

		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	// CF tunnel setup
	e = a.executors[contract.JobTypeEnsureCfTunnel]
	if err := e.Execute(ctx, nil); err != nil {
		if !errors.Is(err, jobs.ErrCfTunnelNotFound) {
			return err
		}

		e = a.executors[contract.JobTypeCreateCfTunnel]
		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	// CF container setup
	e = a.executors[contract.JobTypeEnsureCfContainer]
	if err := e.Execute(ctx, nil); err != nil {
		e = a.executors[contract.JobTypeCreateCfContainer]

		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	// Railpack setup
	e = a.executors[contract.JobTypeEnsureRailpack]
	if err := e.Execute(ctx, nil); err != nil {
		e = a.executors[contract.JobTypeInstallRailpack]

		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	// Buildkit setup
	e = a.executors[contract.JobTypeEnsureBuildkit]
	if err := e.Execute(ctx, nil); err != nil {
		e = a.executors[contract.JobTypeSetupBuildkit]

		if err := e.Execute(ctx, nil); err != nil {
			return err
		}
	}

	return nil
}
