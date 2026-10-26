package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
	"github.com/SudiMango/podlet/internal/contract"
)

type GithubHandler struct {
	githubService *service.GithubService
}

func NewGithubHandler(githubService *service.GithubService) *GithubHandler {
	return &GithubHandler{
		githubService: githubService,
	}
}

/*
Start a GitHub App installation for the current user

API URL: POST /github/installations/start

Args:

	none

Returns:

	a GitHub installation URL and its expiration time (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to create the installation state
*/
func (h *GithubHandler) StartInstallation(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	created, err := h.githubService.CreateInstallationStateByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, contract.GithubInstallationStartResponse{
		InstallationUrl: "https://github.com/apps/" + os.Getenv("GITHUB_APP_NAME") + "/installations/new?state=" + created.Id,
		ExpiresAt:       created.ExpiresAt,
	})
}

/*
List the GitHub repositories available to the current user through installed GitHub Apps

API URL: GET /github/repositories

Args:

	none

Returns:

	the repositories accessible to the user's GitHub App installations (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to load installations or repositories
*/
func (h *GithubHandler) GetRepositoriesByUserId(w http.ResponseWriter, r *http.Request) {
	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	repositories, err := h.githubService.GetGithubRepositoriesByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, repositories)
}

/*
Verify and save a GitHub App installation after GitHub redirects the user

API URL: GET /github/installations/setup

Args:

	installation_id: the GitHub installation ID from the query string
	state: the temporary installation state from the query string

Returns:

	an HTML installation confirmation page (HTTP 200).

Raises:

	400: missing or invalid query parameters, expired state, or invalid installation
	403: installation is suspended or lacks repository contents access
	500: unable to verify or save the installation
*/
func (h *GithubHandler) VerifyInstallation(w http.ResponseWriter, r *http.Request) {
	stateID := r.URL.Query().Get("state")
	rawInstallationId := r.URL.Query().Get("installation_id")

	// Verify query params
	if stateID == "" || rawInstallationId == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("missing installation_id or state"))
		return
	}

	installationId, err := strconv.ParseInt(rawInstallationId, 10, 64)
	if err != nil || installationId <= 0 {
		util.SendError(w, http.StatusBadRequest, errors.New("invalid installation_id"))
		return
	}

	// Verify installation
	_, err = h.githubService.VerifyInstallation(r.Context(), stateID, installationId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStateExpired), errors.Is(err, service.ErrInvalidInstallation):
			util.SendError(w, http.StatusBadRequest, err)

		case errors.Is(err, service.ErrInstallationSuspended), errors.Is(err, service.ErrContentsPermissionRequired):
			util.SendError(w, http.StatusForbidden, err)

		default:
			util.SendError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Show success confirmation to user
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `<!doctype html>
		<html>
		<body>
			<h1>GitHub connected</h1>
			<p>Podlet can now access your selected repositories.</p>
			<p>You can close this window and return to the terminal.</p>
		</body>
		</html>`)
}
