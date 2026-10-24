package handler

import (
	"errors"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/auth"
	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
)

type ProfileHandler struct {
	profileService *service.ProfileService
}

func NewProfileHandler(profileService *service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

/*
Get the currently authenticated user's profile.

API URL: GET /profile

Args:

	None

Returns:

	response body containing the current user's profile (HTTP 200).

Raises:

	401: unauthorized user
	500: unable to load profile
*/
func (h *ProfileHandler) GetPodletProfileByUserId(w http.ResponseWriter, r *http.Request) {

	userId, ok := auth.CurrentUserID(r)
	if !ok {
		util.SendError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	user, err := h.profileService.GetPodletProfileByUserId(r.Context(), userId)
	if err != nil {
		util.SendError(w, http.StatusInternalServerError, err)
		return
	}

	util.SendJSON(w, http.StatusOK, *user)
}
