package handler

import (
	"encoding/json"
	"net/http"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type OverrideHandler struct {
	service *service.OverrideService
}

func NewOverrideHandler(service *service.OverrideService) *OverrideHandler {
	return &OverrideHandler{service: service}
}

// Used for PATCH status
type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// List godoc
// @Summary List overrides
// @Description Get list of overrides with optional filters
// @Tags Overrides
// @Accept json
// @Produce json
// @Param state query string false "State"
// @Param client query string false "Client"
// @Param investor query string false "Investor"
// @Param caseType query string false "Case Type"
// @Param stepKey query string false "Step Key"
// @Param traitKey query string false "Trait Key"
// @Param status query string false "Status"
// @Success 200 {array} domain.Override
// @Failure 500 {string} string "internal server error"
// @Router /api/overrides [get]
func (h *OverrideHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := domain.OverrideFilter{
		State:     strPtr(q.Get("state")),
		Client:    strPtr(q.Get("client")),
		Investor:  strPtr(q.Get("investor")),
		CaseType:  strPtr(q.Get("caseType")),
		StepKey:   strPtr(q.Get("stepKey")),
		TraitKey:  strPtr(q.Get("traitKey")),
		Status:    strPtr(q.Get("status")),
	}

	data, err := h.service.List(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// GetByID godoc
// @Summary Get override by ID
// @Description Get a single override by its ID
// @Tags Overrides
// @Produce json
// @Param id path string true "Override ID"
// @Success 200 {object} domain.Override
// @Failure 404 {string} string "not found"
// @Router /api/overrides/{id} [get]
func (h *OverrideHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Create godoc
// @Summary Create override
// @Description Create a new override
// @Tags Overrides
// @Accept json
// @Produce json
// @Param request body domain.Override true "Override payload"
// @Success 201 {string} string "created"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/overrides [post]
func (h *OverrideHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.Override

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Specificity = domain.ComputeSpecificity(req.Selector)
	_, err := h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update godoc
// @Summary Update override
// @Description Update an existing override
// @Tags Overrides
// @Accept json
// @Produce json
// @Param id path string true "Override ID"
// @Param request body domain.Override true "Override payload"
// @Success 200 {string} string "updated"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/overrides/{id} [put]
func (h *OverrideHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req domain.Override
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.ID = id

	_,err := h.service.Update(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateStatus godoc
// @Summary Update override status
// @Description Update only the status of an override
// @Tags Overrides
// @Accept json
// @Produce json
// @Param id path string true "Override ID"
// @Param request body UpdateStatusRequest true "Status payload"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid request"
// @Failure 500 {string} string "internal server error"
// @Router /api/overrides/{id}/status [patch]
func (h *OverrideHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var body UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.UpdateStatus(id, body.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}