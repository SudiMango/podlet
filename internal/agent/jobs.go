package agent

import (
	"context"
	"net/http"
	"time"

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

	go func() {
		a.busy = true
		a.executeJob(ctx, *job)
	}()

	return job.JobId, nil
}

func (a *Agent) executeJob(ctx context.Context, job contract.ClaimJobResponse) {
	time.Sleep(10 * time.Second) // For testing only
	a.reportJobFailed(ctx, job.JobId, "error")
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

	a.busy = false
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
