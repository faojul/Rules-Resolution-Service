package service

import (
	"context"
	"sort"
	"time"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/repository"
)

type Resolver struct {
    repo repository.OverrideRepository
}

func NewResolver(repo repository.OverrideRepository) *Resolver {
    return &Resolver{repo: repo}
}

func (r *Resolver) Resolve(ctx domain.Context) (map[string]any, error) {

    steps := []string{"title-search", "file-complaint"}
    traits := []string{"slaHours", "feeAmount"}

    result := make(map[string]any)

    for _, step := range steps {
        result[step] = make(map[string]any)

        for _, trait := range traits {

            overrides, err := r.repo.FindByStepAndTrait(context.Background(), step, trait)
            if err != nil {
                return nil, err
            }

            val, _, _ := resolveTrait(ctx, 0, overrides)

            result[step].(map[string]any)[trait] = val
        }
    }

    return result, nil
}

func resolveTrait(ctx domain.Context, defaultVal any, overrides []domain.Override) (any, *domain.Override, error) {

    var candidates []domain.Override

    for _, o := range overrides {
        if o.Matches(ctx) && o.IsActive(ctx.AsOfDate) {
            candidates = append(candidates, o)
        }
    }

    if len(candidates) == 0 {
        return defaultVal, nil, nil
    }

    sort.Slice(candidates, func(i, j int) bool {
        if candidates[i].Specificity != candidates[j].Specificity {
            return candidates[i].Specificity > candidates[j].Specificity
        }
        return candidates[i].EffectiveDate.After(candidates[j].EffectiveDate)
    })

    return candidates[0].Value, &candidates[0], nil
}

func (r *Resolver) BulkResolve(contexts []domain.Context) ([]map[string]any, error) {

    overrides, err := r.repo.FindAllOverrides(context.Background())
    if err != nil {
        return nil, err
    }

    // group in memory
    grouped := make(map[string][]domain.Override)

    for _, o := range overrides {
        key := o.StepKey + ":" + o.TraitKey
        grouped[key] = append(grouped[key], o)
    }

    var results []map[string]any

    for _, ctx := range contexts {

        if ctx.AsOfDate.IsZero() {
            ctx.AsOfDate = time.Now()
        }

        res := make(map[string]any)

        for key, overrides := range grouped {

            val, _, _ := resolveTrait(ctx, 0, overrides)

            res[key] = val
        }

        results = append(results, res)
    }

    return results, nil
}