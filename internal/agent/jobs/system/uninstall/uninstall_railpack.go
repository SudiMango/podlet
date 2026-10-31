package uninstall

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/SudiMango/podlet/internal/contract"
)

const railpackBinaryPath = "/usr/local/bin/railpack"

type UninstallRailpackJob struct{}

func (UninstallRailpackJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	if err := os.Remove(railpackBinaryPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("uninstall Railpack: %w", err)
	}

	return nil
}
