package agent

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SudiMango/podlet/internal/contract"
)

func (a *Agent) claimJob(ctx context.Context) (string, error) {
	var job *contract.ClaimJobResponse
	if err := a.client.DoJSON(
		ctx,
		http.MethodGet,
		"/agent/jobs/claim",
		nil,
		&job,
	); err != nil {
		return "", err
	}

	if job == nil {
		return "", nil
	}

	if err := a.reportJobClaimed(ctx, job.JobId); err != nil {
		return "", err
	}

	a.busy.Store(true)
	go func() {
		a.executeJob(ctx, *job)
	}()

	return job.JobId, nil
}

func (a *Agent) executeJob(ctx context.Context, job contract.ClaimJobResponse) {
	defer a.busy.Store(false)

	// Validate job type and perform job
	var err error
	switch job.JobType {
	case contract.JobTypeCloneRepository:
		err = a.cloneRepository(ctx, job)
	default:
		err = fmt.Errorf("unsupported job type %q", job.JobType)
	}

	// Report errors
	if err != nil {
		if reportErr := a.reportJobFailed(ctx, job.JobId, truncateError(err.Error())); reportErr != nil {
			fmt.Fprintf(os.Stderr, "report failed job %s: %v\n", job.JobId, reportErr)
		}
		return
	}

	// Report success
	if err := a.reportJobSucceeded(ctx, job.JobId); err != nil {
		fmt.Fprintf(os.Stderr, "report successful job %s: %v\n", job.JobId, err)
	}
}

func (a *Agent) cloneRepository(ctx context.Context, job contract.ClaimJobResponse) error {
	if job.JobId == "" || filepath.Base(job.JobId) != job.JobId {
		return errors.New("invalid job id")
	}
	if job.GithubToken == "" {
		return errors.New("missing GitHub token")
	}

	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}

	jobDir := filepath.Join(dataDir, "podlet", "deployments", job.JobId)
	sourceDir := filepath.Join(jobDir, "source")
	if err := os.MkdirAll(jobDir, 0700); err != nil {
		return err
	}

	askpassPath := filepath.Join(jobDir, "git-askpass.sh")
	askpass := "#!/bin/sh\ncase \"$1\" in\n  *Username*) printf '%s\\n' 'x-access-token' ;;\n  *Password*) printf '%s\\n' \"$PODLET_GITHUB_TOKEN\" ;;\nesac\n"
	if err := os.WriteFile(askpassPath, []byte(askpass), 0700); err != nil {
		return err
	}
	defer os.Remove(askpassPath)

	command := exec.CommandContext(ctx, "git", "clone", "--depth", "1", "--", job.CloneUrl, sourceDir)
	command.Env = append(os.Environ(),
		"GIT_ASKPASS="+askpassPath,
		"GIT_TERMINAL_PROMPT=0",
		"PODLET_GITHUB_TOKEN="+job.GithubToken,
	)

	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Run(); err != nil {
		_ = os.RemoveAll(sourceDir)
		message := strings.TrimSpace(output.String())
		if message == "" {
			return err
		}
		return fmt.Errorf("git clone failed: %s", truncateError(message))
	}

	return nil
}

func (a *Agent) reportJob(ctx context.Context, report contract.ReportJobRequest) error {
	if err := a.client.DoJSON(
		ctx,
		http.MethodPost,
		"/agent/jobs/report",
		report,
		nil,
	); err != nil {
		return err
	}

	return nil
}

/*
	Helpers
*/

func (a *Agent) reportJobClaimed(ctx context.Context, jobId string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusOngoing,
		Error:     "",
	})
}

func (a *Agent) reportJobFailed(ctx context.Context, jobId string, err string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusFailed,
		Error:     err,
	})
}

func (a *Agent) reportJobSucceeded(ctx context.Context, jobId string) error {
	return a.reportJob(ctx, contract.ReportJobRequest{
		JobId:     jobId,
		JobStatus: contract.StatusSuccess,
		Error:     "",
	})
}

func truncateError(message string) string {
	const maxErrorLength = 4000
	if len(message) <= maxErrorLength {
		return message
	}

	return message[:maxErrorLength]
}
