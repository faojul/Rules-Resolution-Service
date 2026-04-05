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