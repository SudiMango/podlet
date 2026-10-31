package uninstall

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func removeDockerResource(ctx context.Context, resourceType, name string, force bool) error {
	if _, err := exec.LookPath("docker"); err != nil {
		return nil
	}

	inspectOutput, err := exec.CommandContext(ctx, "docker", resourceType, "inspect", name).CombinedOutput()
	if err != nil {
		if dockerResourceNotFound(inspectOutput) {
			return nil
		}
		return fmt.Errorf("inspect Docker %s %s: %w: %s", resourceType, name, err, strings.TrimSpace(string(inspectOutput)))
	}

	args := []string{resourceType, "rm"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, name)

	output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("remove Docker %s %s: %w: %s", resourceType, name, err, strings.TrimSpace(string(output)))
	}

	return nil
}

func dockerResourceNotFound(output []byte) bool {
	message := strings.ToLower(string(output))
	return strings.Contains(message, "no such") || strings.Contains(message, "not found")
}
