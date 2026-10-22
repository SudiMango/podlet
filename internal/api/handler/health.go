package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SudiMango/podlet/internal/api/service"
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
	w.Header().Set("Content-Type", "application/json")

	healthy, err := h.healthService.IsDatabaseHealthy(r.Context())

	response := map[string]interface{}{
		"status": "healthy",
	}

	if err != nil || !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		response["status"] = "unhealthy"
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
