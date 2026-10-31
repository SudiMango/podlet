package uninstall

import (
	"context"
	"fmt"

	"github.com/SudiMango/podlet/internal/contract"
)

type RemoveCfContainerJob struct{}

func (RemoveCfContainerJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if err := removeDockerResource(ctx, "container", "podlet-cloudflared", true); err != nil {
		return fmt.Errorf("remove Cloudflare container: %w", err)
	}
	if err := removeDockerResource(ctx, "image", "cloudflare/cloudflared:latest", false); err != nil {
		return fmt.Errorf("remove Cloudflare image: %w", err)
	}

	return nil
}
