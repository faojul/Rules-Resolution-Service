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
// Resolve godoc
// @Summary Resolve rules
// @Description Resolve rules based on input context
// @Tags Resolve
// @Accept json
// @Produce json
// @Param request body domain.Context true "Resolution context"
// @Success 200 {object} map[string]interface{} "resolved result"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/resolve [post]
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

// Bulk Resolve godoc
// @Summary Resolve rules in Bulk
// @Description Resolve rules bulk based on input contexts
// @Tags Bulk Resolve
// @Accept json
// @Produce json
// @Param request body []domain.Context true "Resolution contexts"
// @Success 200 {object} map[string]interface{} "resolved result"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/resolve [post]
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