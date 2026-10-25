package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/dto"
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

func (h *GithubHandler) VerifyInstallation(w http.ResponseWriter, r *http.Request) {
	// GitHub includes the installation ID and our original state in its redirect.
	installationID := r.URL.Query().Get("installation_id")
	stateID := r.URL.Query().Get("state")
	if installationID == "" || stateID == "" {
		util.SendError(w, http.StatusBadRequest, errors.New("missing installation_id or state"))
		return
	}

	parsedInstallationID, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil || parsedInstallationID <= 0 {
		util.SendError(w, http.StatusBadRequest, errors.New("invalid installation_id"))
		return
	}

	// Resolve the short-lived state back to the Podlet user who started setup.
	state, err := h.githubService.GetStateByStateId(r.Context(), stateID)
	if err != nil {
		if errors.Is(err, service.ErrStateExpired) {
			util.SendError(w, http.StatusBadRequest, err)
			return
		}

		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	// Authenticate this API as the GitHub App.
	bearerToken, err := auth.GenerateGithubAppJWT()
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	// Ask GitHub for the canonical installation rather than trusting query data.
	endpoint := "https://api.github.com/app/installations/" + strconv.FormatInt(parsedInstallationID, 10)

	request, err := http.NewRequestWithContext(r.Context(), http.MethodGet, endpoint, nil)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	githubHTTPClient := &http.Client{Timeout: 15 * time.Second}
	response, err := githubHTTPClient.Do(request)
	if err != nil {
		util.SendError(w, http.StatusBadGateway, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			util.SendError(
				w,
				http.StatusBadGateway,
				errors.New("GitHub rejected the App credentials"),
			)
		case http.StatusNotFound:
			util.SendError(
				w,
				http.StatusBadRequest,
				errors.New("GitHub installation was not found"),
			)
		default:
			util.SendError(
				w,
				http.StatusBadGateway,
				fmt.Errorf("GitHub returned %s", response.Status),
			)
		}
		return
	}

	// Decode only the installation fields Podlet needs to validate and store.
	var installation dto.GithubInstallationResponse

	if err := json.NewDecoder(
		io.LimitReader(response.Body, 1<<20),
	).Decode(&installation); err != nil {
		util.SendError(w, http.StatusBadGateway, err)
		return
	}

	// Ensure GitHub returned an active installation for this exact App.
	expectedAppID, err := strconv.ParseInt(
		os.Getenv("GITHUB_APP_ID"),
		10,
		64,
	)
	if err != nil {
		util.SendError(
			w,
			http.StatusInternalServerError,
			errors.New("invalid GITHUB_APP_ID"),
		)
		return
	}

	if installation.ID != parsedInstallationID ||
		installation.AppID != expectedAppID ||
		installation.AppSlug != os.Getenv("GITHUB_APP_NAME") {
		util.SendError(
			w,
			http.StatusUnauthorized,
			errors.New("installation does not belong to this GitHub App"),
		)
		return
	}

	if installation.SuspendedAt != nil {
		util.SendError(
			w,
			http.StatusForbidden,
			errors.New("GitHub installation is suspended"),
		)
		return
	}

	if installation.Permissions.Contents != "read" &&
		installation.Permissions.Contents != "write" {
		util.SendError(
			w,
			http.StatusForbidden,
			errors.New("installation does not grant repository contents access"),
		)
		return
	}

	// Persist the verified installation against the user recovered from state.
	if _, err := h.githubService.UpsertGithubInstallationByUserId(
		r.Context(),
		state.UserId,
		installation,
	); err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	// Consume the state so the callback cannot be replayed after success.
	if err := h.githubService.DeleteGithubInstallationStateByStateId(
		r.Context(),
		stateID,
	); err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	// The callback opens in a browser, so finish with a small confirmation page.
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
