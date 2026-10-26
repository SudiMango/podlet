package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

/*
Create a pending repository deployment job for one of the user's servers

API URL: POST /servers/deploy

Args:

	server_id: the ID of the server that should run the job
	repo_name: the GitHub repository in Owner/Repository format

Returns:

	the created job ID and its pending status (HTTP 200).

Raises:

	400: invalid request body, server ID, or repository name
	401: unauthorized user
	500: unable to create the job
*/
func (h *JobHandler) CreatePendingJobForServer(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	var request contract.CreateJobRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	result, err := h.jobService.CreatePendingJobForServer(r.Context(), userId, request)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidServerId), errors.Is(err, service.ErrInvalidRepoName):
			util.SendError(w, http.StatusBadRequest, err)
		default:
			util.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	util.SendJSON(w, http.StatusOK, result)
}

/*
Get all jobs owned by the current user

API URL: GET /jobs

Args:

	status: optional pending, ongoing, success, or failed query filter

Returns:

	the user's jobs ordered from newest to oldest (HTTP 200).

Raises:

	400: invalid status filter
	401: unauthorized user
	500: unable to get jobs
*/
func (h *JobHandler) GetAllJobsByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	jobs, err := h.jobService.GetAllJobsByUserId(r.Context(), userId, r.URL.Query().Get("status"))
	if err != nil {
		if errors.Is(err, service.ErrInvalidJobFilter) {
			util.SendError(w, http.StatusBadRequest, err)
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, jobs)
}

/*
Get all jobs for one of the current user's servers

API URL: GET /jobs/{serverId}/jobs

Args:

	serverId: the server ID from the URL path
	status: optional pending, ongoing, success, or failed query filter

Returns:

	the server's jobs ordered from newest to oldest (HTTP 200).

Raises:

	400: serverId is missing or the status filter is invalid
	401: unauthorized user
	500: unable to get jobs
*/
func (h *JobHandler) GetAllJobsByServerIdAndUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	serverId := chi.URLParam(r, "serverId")
	if serverId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("serverId is required"))
		return
	}

	jobs, err := h.jobService.GetAllJobsByServerIdAndUserId(r.Context(), userId, serverId, r.URL.Query().Get("status"))
	if err != nil {
		if errors.Is(err, service.ErrInvalidJobFilter) {
			util.SendError(w, http.StatusBadRequest, err)
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, jobs)
}

/*
Get a job owned by the current user

API URL: GET /jobs/{jobId}

Args:

	jobId: the job ID from the URL path

Returns:

	the requested job (HTTP 200).

Raises:

	400: jobId is missing
	401: unauthorized user
	404: job not found
	500: unable to get the job
*/
func (h *JobHandler) GetJobByJobIdAndUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	jobId := chi.URLParam(r, "jobId")
	if jobId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("jobId is required"))
		return
	}

	job, err := h.jobService.GetJobByJobIdAndUserId(r.Context(), userId, jobId)
	if err != nil {
		if errors.Is(err, service.ErrJobNotFound) {
			util.SendError(w, http.StatusNotFound, err)
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, job)
}

/*
Claim the oldest pending job assigned to the authenticated server agent

API URL: GET /agent/jobs/claim

Args:

	none

Returns:

	the claimed job and temporary GitHub credentials (HTTP 200), or
	an empty response when no jobs are available (HTTP 204).

Raises:

	401: invalid agent identity
	500: unable to claim the job or create GitHub credentials
*/
func (h *JobHandler) ClaimJob(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	result, err := h.jobService.ClaimJob(r.Context(), agentIdentity.ServerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoJobsAvailable):
			util.SendJSON(w, http.StatusNoContent, nil)
		default:
			util.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	util.SendJSON(w, http.StatusOK, result)
}

/*
Report the current status of a job assigned to the authenticated server agent

API URL: POST /agent/jobs/report

Args:

	job_id: the ID of the assigned job
	job_status: ongoing, success, or failed
	error: the failure message when job_status is failed

Returns:

	an empty response (HTTP 204).

Raises:

	400: invalid request body or job status
	401: invalid agent identity
	404: job not found, owned by another server, or cannot transition status
	500: unable to update the job
*/
func (h *JobHandler) ReportJobByServerIdAndJobId(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	var request contract.ReportJobRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("Invalid request body"))
		return
	}

	if err := h.jobService.ReportJobByServerIdAndJobId(r.Context(), agentIdentity.ServerID, request); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidJobStatus):
			util.SendError(w, http.StatusBadRequest, err)
		case errors.Is(err, service.ErrJobNotFound):
			util.SendError(w, http.StatusNotFound, err)
		default:
			util.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
