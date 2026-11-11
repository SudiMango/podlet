package runtime

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

const repositoryRoot = "/var/lib/podlet/repositories"

type CloneRepositoryJob struct{}

func (CloneRepositoryJob) Execute(ctx context.Context, job *contract.ClaimJobResponse) error {
	var payload contract.CloneRepositoryJob
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return fmt.Errorf("decode clone payload: %w", err)
	}

	repositoryPath := filepath.Join(repositoryRoot, job.Id)

	if err := os.MkdirAll(repositoryRoot, 0700); err != nil {
		return fmt.Errorf("create repository directory: %w", err)
	}

	if err := os.RemoveAll(repositoryPath); err != nil {
		return fmt.Errorf("remove existing repository: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		"git",
		"clone",
		"--depth", "1",
		"--",
		payload.CloneUrl,
		repositoryPath,
	)
	credentials := base64.StdEncoding.EncodeToString(
		[]byte("x-access-token:" + payload.GithubToken),
	)

	cmd.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=http.extraHeader",
		"GIT_CONFIG_VALUE_0=Authorization: Basic "+credentials,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone: %w: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}
