package uninstall

import (
	"context"
	"fmt"

	"github.com/SudiMango/podlet/internal/contract"
)

type RemoveBuildkitJob struct{}

func (RemoveBuildkitJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if err := removeDockerResource(ctx, "container", "podlet-buildkit", true); err != nil {
		return fmt.Errorf("remove BuildKit container: %w", err)
	}
	if err := removeDockerResource(ctx, "image", "moby/buildkit:latest", false); err != nil {
		return fmt.Errorf("remove BuildKit image: %w", err)
	}

	return nil
}
