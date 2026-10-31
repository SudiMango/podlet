package uninstall

import (
	"context"
	"fmt"

	"github.com/SudiMango/podlet/internal/contract"
)

type RemovePodletNetworkJob struct{}

func (RemovePodletNetworkJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if err := removeDockerResource(ctx, "network", "podlet", false); err != nil {
		return fmt.Errorf("remove Podlet network: %w", err)
	}

	return nil
}
