package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

type EnvVarHandler struct {
	envVarService *service.EnvVarService
}

func NewEnvVarHandler(envVarService *service.EnvVarService) *EnvVarHandler {
	return &EnvVarHandler{
		envVarService: envVarService,
	}
}

/*
Create an environment variable for a repository

API URL: POST /env

Args:

	key: the variable name (JSON body)
	value: the variable value (JSON body)
	repo_name: the full repository name, e.g. SudiMango/podlet-test (JSON body)

Returns:

	the created environment variable, without its value (HTTP 201).

Raises:

	400: invalid body, key, value or repository name
	401: unauthorized user
	500: unable to create the environment variable
*/
func (h *EnvVarHandler) CreateEnvVarByUserIdAndRepoName(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	var request contract.CreateEnvVarRequest
	if err := util.DecodeJSON(r, &request); err != nil {
		util.SendError(w, http.StatusBadRequest, err)
		return
	}

	envVar, err := h.envVarService.CreateEnvVarByUserIdAndRepoName(r.Context(), userId, request)
	if err != nil {
		sendEnvVarError(w, err)
		return
	}

	util.SendJSON(w, http.StatusCreated, envVar)
}

/*
List the environment variables of a repository

API URL: GET /env?repo_name={owner/repo}

Args:

	repo_name: the full repository name, e.g. SudiMango/podlet-test (query string)

Returns:

	the repository's environment variables, without values (HTTP 200).

Raises:

	400: invalid repository name
	401: unauthorized user
	500: unable to load the environment variables
*/
func (h *EnvVarHandler) GetEnvVarsByUserIdAndRepoName(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	envVars, err := h.envVarService.GetEnvVarsByUserIdAndRepoName(r.Context(), userId, r.URL.Query().Get("repo_name"))
	if err != nil {
		sendEnvVarError(w, err)
		return
	}

	util.SendJSON(w, http.StatusOK, envVars)
}

/*
Delete one environment variable of a repository

API URL: DELETE /env?repo_name={owner/repo}&key={key}

Args:

	repo_name: the full repository name, e.g. SudiMango/podlet-test (query string)
	key: the variable name (query string)

Returns:

	no content (HTTP 204).

Raises:

	400: invalid repository name
	401: unauthorized user
	404: the environment variable doesn't exist
	500: unable to delete the environment variable
*/
func (h *EnvVarHandler) DeleteEnvVarByUserIdAndRepoName(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	query := r.URL.Query()
	if err := h.envVarService.DeleteEnvVarByUserIdAndRepoName(r.Context(), userId, query.Get("repo_name"), query.Get("key")); err != nil {
		sendEnvVarError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Delete all environment variables of a repository

API URL: DELETE /env/all?repo_name={owner/repo}

Args:

	repo_name: the full repository name, e.g. SudiMango/podlet-test (query string)

Returns:

	no content (HTTP 204).

Raises:

	400: invalid repository name
	401: unauthorized user
	500: unable to delete the environment variables
*/
func (h *EnvVarHandler) DeleteAllEnvVarsByUserIdAndRepoName(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	if err := h.envVarService.DeleteAllEnvVarsByUserIdAndRepoName(r.Context(), userId, r.URL.Query().Get("repo_name")); err != nil {
		sendEnvVarError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

/*
Get the decrypted environment variables of a repository for the authenticated server agent

API URL: GET /agent/env?repository_id={id}

Args:

	repository_id: the GitHub repository ID (query string)

Returns:

	the repository's environment variables with decrypted values (HTTP 200).

Raises:

	400: missing or invalid repository_id
	401: invalid agent identity or authentication token
	500: unable to load or decrypt the environment variables
*/
func (h *EnvVarHandler) GetEnvVarsForAgent(w http.ResponseWriter, r *http.Request) {
	agentIdentity, ok := auth.AgentIdentityFromContext(r.Context())
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Invalid agent identity"))
		return
	}

	repositoryId, err := strconv.ParseInt(r.URL.Query().Get("repository_id"), 10, 64)
	if err != nil {
		util.SendError(w, http.StatusBadRequest, errors.New("repository_id is required and must be an integer"))
		return
	}

	values, err := h.envVarService.GetDecryptedEnvVarsByUserIdAndRepoId(r.Context(), agentIdentity.UserID, repositoryId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	response := make(contract.EnvVarsForAgentResponse, 0, len(values))
	for key, value := range values {
		response = append(response, contract.EnvVarForAgentResponse{Key: key, Value: value})
	}

	util.SendJSON(w, http.StatusOK, response)
}

/*
	Helpers
*/

func sendEnvVarError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidEnvKey),
		errors.Is(err, service.ErrInvalidEnvValue),
		errors.Is(err, service.ErrInvalidRepoName):
		util.SendError(w, http.StatusBadRequest, err)
	case errors.Is(err, service.ErrEnvVarNotFound):
		util.SendError(w, http.StatusNotFound, err)
	default:
		util.SendError(w, http.StatusInternalServerError, err)
	}
}
