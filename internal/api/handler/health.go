package handler

import (
	"net/http"

	"github.com/SudiMango/podlet/internal/api/service"
	"github.com/SudiMango/podlet/internal/api/util"
)

type HealthHandler struct {
	healthService *service.HealthService
}

func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

func (h *HealthHandler) IsDatabaseHealthy(w http.ResponseWriter, r *http.Request) {
	healthy, err := h.healthService.IsDatabaseHealthy(r.Context())

	if err != nil || !healthy {
		util.SendJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unhealthy",
		})
		return
	}

	util.SendJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}
