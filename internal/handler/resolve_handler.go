package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/service"
)

type Handler struct {
    resolver *service.Resolver
}

func NewResolveHandler(r *service.Resolver) *Handler {
    return &Handler{resolver: r}
}

func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) {

    var req domain.Context
    json.NewDecoder(r.Body).Decode(&req)

    if req.AsOfDate.IsZero() {
        req.AsOfDate = time.Now()
    }

    result, err := h.resolver.Resolve(req)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(result)
}

func (h *Handler) BulkResolve(w http.ResponseWriter, r *http.Request) {

    var req struct {
        Contexts []domain.Context `json:"contexts"`
    }

    json.NewDecoder(r.Body).Decode(&req)

    result, err := h.resolver.BulkResolve(req.Contexts)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    json.NewEncoder(w).Encode(result)
}