package service

import (
	"context"
	"errors"
	"time"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/repository"
)

type OverrideService struct {
	repo repository.OverrideRepository
}

func NewOverrideService(r repository.OverrideRepository) *OverrideService {
	return &OverrideService{repo: r}
}

func (s *OverrideService) List(filter domain.OverrideFilter) ([]domain.Override, error) {
	return s.repo.List(context.Background(), filter)
}

func (s *OverrideService) GetByID(id string) (*domain.Override, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetByID(context.Background(), id)
}

func (s *OverrideService) Create(o domain.Override) (*domain.Override, error) {

	// validation
	if err := validateOverride(o); err != nil {
		return nil, err
	}

	now := time.Now()

	o.Specificity = domain.ComputeSpecificity(o.Selector)
	o.CreatedAt = &now
	o.UpdatedAt = &now

	if o.Status == "" {
		o.Status = "draft"
	}

	err := s.repo.Create(context.Background(), o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *OverrideService) Update(o domain.Override) (*domain.Override, error) {

	if o.ID == "" {
		return nil, errors.New("id is required")
	}

	if err := validateOverride(o); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByID(context.Background(), o.ID)
	if err != nil {
		return nil, err
	}

	o.Specificity = domain.ComputeSpecificity(o.Selector)
	o.CreatedAt = existing.CreatedAt
		now := time.Now()
		o.UpdatedAt = &now

	err = s.repo.Update(context.Background(), o)
	if err != nil {
		return nil, err
	}

	// audit trail
	_ = s.repo.InsertHistory(context.Background(), *existing, o)

	return &o, nil
}

func (s *OverrideService) GetConflicts() ([]Conflict, error) {

	overrides, err := s.repo.FindAllOverrides(context.Background())
	if err != nil {
		return nil, err
	}

	conflicts := DetectConflicts(overrides)

	return conflicts, nil
}

func (s *OverrideService) GetHistory(id string) ([]domain.OverrideHistory, error) {

	if id == "" {
		return nil, errors.New("id is required")
	}

	return s.repo.GetHistory(context.Background(), id)
}

func validateOverride(o domain.Override) error {

	if o.StepKey == "" {
		return errors.New("stepKey is required")
	}

	if o.TraitKey == "" {
		return errors.New("traitKey is required")
	}

	if o.EffectiveDate.IsZero() {
		return errors.New("effectiveDate is required")
	}

	if o.ExpiresDate != nil && o.ExpiresDate.Before(o.EffectiveDate) {
		return errors.New("expiresDate must be after effectiveDate")
	}

	return nil
}

func (s *OverrideService) UpdateStatus(id string, status string) error {
	return s.repo.UpdateStatus(context.Background(), id, status)
}
