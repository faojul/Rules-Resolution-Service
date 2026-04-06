package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/service"
)

type ExplainHandler struct {
    resolver *service.Resolver
}

func NewExplainHandler(resolver *service.Resolver) *ExplainHandler {
    return &ExplainHandler{resolver: resolver}
}
// Explain godoc
// @Summary Explain resolution
// @Description Explain how a rule was resolved (step-by-step reasoning)
// @Tags Explain
// @Accept json
// @Produce json
// @Param request body domain.Context true "Resolution context"
// @Success 200 {object} map[string]interface{} "explanation result"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/explain [post]
func (h *ExplainHandler) Explain(w http.ResponseWriter, r *http.Request) {

    var ctx domain.Context

    if err := json.NewDecoder(r.Body).Decode(&ctx); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if ctx.AsOfDate.IsZero() {
        ctx.AsOfDate = time.Now()
    }

    result, err := h.resolver.Explain(ctx)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}