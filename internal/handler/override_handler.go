package handler

import (
	"net/http"
	"rules-resolution-service/internal/service"
)

type OverrideHandler struct {
	service *service.OverrideService
}

func NewOverrideHandler(service *service.OverrideService) *OverrideHandler {
    return &OverrideHandler{service: service}
}


func (h *OverrideHandler) List(w http.ResponseWriter, r *http.Request)
func (h *OverrideHandler) GetByID(w http.ResponseWriter, r *http.Request)
func (h *OverrideHandler) Create(w http.ResponseWriter, r *http.Request)
func (h *OverrideHandler) Update(w http.ResponseWriter, r *http.Request)
func (h *OverrideHandler) UpdateStatus(w http.ResponseWriter, r *http.Request)