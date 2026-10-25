package handler

import (
	"errors"
	"net/http"
	"os"

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
