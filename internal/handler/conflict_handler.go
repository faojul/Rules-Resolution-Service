package handler

import (
	"encoding/json"
	"net/http"

	"rules-resolution-service/internal/service"
)

type ConflictHandler struct {
    service *service.OverrideService
}

func NewConflictHandler(s *service.OverrideService) *ConflictHandler {
    return &ConflictHandler{service: s}
}
// GetConflicts godoc
// @Summary Get conflicts
// @Description Retrieve all rule conflicts detected in the system
// @Tags Conflicts
// @Produce json
// @Success 200 {object} map[string]interface{} "conflicts response"
// @Failure 500 {string} string "internal server error"
// @Router /api/conflicts [get]
func (h *ConflictHandler) GetConflicts(w http.ResponseWriter, r *http.Request) {

    conflicts, err := h.service.GetConflicts()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]any{
        "conflicts": conflicts,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}