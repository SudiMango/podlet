package agent

import (
	"context"
	"errors"

	"github.com/SudiMango/podlet/internal/agent/dto"
	initializejobs "github.com/SudiMango/podlet/internal/agent/jobs/system/initialize"
)

func (a *Agent) Initialize(ctx context.Context, showProgress bool) error {
	const total = 6

	if err := runStep(showProgress, 1, total, "Ensuring Docker", func() error {
		e := a.executors[dto.JobTypeEnsureDocker]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[dto.JobTypeInstallDocker].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 2, total, "Ensuring Podlet network", func() error {
		e := a.executors[dto.JobTypeEnsurePodletNetwork]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[dto.JobTypeCreatePodletNetwork].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 3, total, "Ensuring Cloudflare tunnel", func() error {
		e := a.executors[dto.JobTypeEnsureCfTunnel]
		if err := e.Execute(ctx, nil); err != nil {
			if !errors.Is(err, initializejobs.ErrCfTunnelNotFound) {
				return err
			}
			return a.executors[dto.JobTypeCreateCfTunnel].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 4, total, "Ensuring Cloudflare container", func() error {
		e := a.executors[dto.JobTypeEnsureCfContainer]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[dto.JobTypeCreateCfContainer].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 5, total, "Ensuring Railpack", func() error {
		e := a.executors[dto.JobTypeEnsureRailpack]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[dto.JobTypeInstallRailpack].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	if err := runStep(showProgress, 6, total, "Ensuring BuildKit", func() error {
		e := a.executors[dto.JobTypeEnsureBuildkit]
		if err := e.Execute(ctx, nil); err != nil {
			return a.executors[dto.JobTypeSetupBuildkit].Execute(ctx, nil)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
