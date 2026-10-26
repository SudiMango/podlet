package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
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
