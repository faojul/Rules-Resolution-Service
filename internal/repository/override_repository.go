package repository

import (
	"context"
	"rules-resolution-service/internal/domain"
)

type OverrideRepository interface {
    FindByStepAndTrait(ctx context.Context, step, trait string) ([]domain.Override, error)
    FindAllOverrides(ctx context.Context) ([]domain.Override, error)
    List(ctx context.Context, filter domain.OverrideFilter) ([]domain.Override, error)
    GetByID(ctx context.Context, id string) (*domain.Override, error)
    Create(ctx context.Context, o domain.Override) error
    Update(ctx context.Context, o domain.Override) error
    UpdateStatus(ctx context.Context, id string, status string) error
    InsertHistory(ctx context.Context, before, after domain.Override) error
    GetHistory(ctx context.Context, id string) ([]domain.OverrideHistory, error)
}