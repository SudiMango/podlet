package agent

import (
	"context"
	"errors"

	"github.com/SudiMango/podlet/internal/agent/jobs"
	"github.com/SudiMango/podlet/internal/contract"
)

func (a *Agent) Initialize(ctx context.Context, showProgress bool) error {
	const total = 6

	if err := runStep(showProgress, 1, total, "Ensuring Docker", func() error {
		e := a.executors[contract.JobTypeEnsureDocker]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[contract.JobTypeInstallDocker].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 2, total, "Ensuring Podlet network", func() error {
		e := a.executors[contract.JobTypeEnsurePodletNetwork]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[contract.JobTypeCreatePodletNetwork].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 3, total, "Ensuring Cloudflare tunnel", func() error {
		e := a.executors[contract.JobTypeEnsureCfTunnel]
		if err := e.Execute(ctx, nil); err != nil {
			if !errors.Is(err, jobs.ErrCfTunnelNotFound) {
				return err
			}
			return a.executors[contract.JobTypeCreateCfTunnel].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 4, total, "Ensuring Cloudflare container", func() error {
		e := a.executors[contract.JobTypeEnsureCfContainer]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[contract.JobTypeCreateCfContainer].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 5, total, "Ensuring Railpack", func() error {
		e := a.executors[contract.JobTypeEnsureRailpack]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[contract.JobTypeInstallRailpack].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 6, total, "Ensuring BuildKit", func() error {
		e := a.executors[contract.JobTypeEnsureBuildkit]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[contract.JobTypeSetupBuildkit].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
