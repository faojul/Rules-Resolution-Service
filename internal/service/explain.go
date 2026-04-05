package service

import (
	"context"
	"sort"
	"time"

	"rules-resolution-service/internal/domain"
)
type ExplainResult struct {
    Step          string                 `json:"step"`
    Trait         string                 `json:"trait"`
    ResolvedValue any                    `json:"resolvedValue"`
    ResolvedFrom  *domain.Override       `json:"resolvedFrom"`
    Candidates    []CandidateExplanation `json:"candidates"`
}

type CandidateExplanation struct {
    OverrideID   string            `json:"overrideId"`
    Specificity  int               `json:"specificity"`
    EffectiveDate time.Time        `json:"effectiveDate"`
    Value        any               `json:"value"`
    Outcome      string            `json:"outcome"`
}

func (r *Resolver) Explain(ctx domain.Context) ([]ExplainResult, error) {

    steps := []string{"file-complaint"} // ideally from DB
    traits := []string{"slaHours"}

    var results []ExplainResult

    for _, step := range steps {
        for _, trait := range traits {

            overrides, _ := r.repo.FindByStepAndTrait(context.Background(), step, trait)

            explain := r.explainTrait(ctx, step, trait, overrides)

            results = append(results, explain)
        }
    }

    return results, nil
}

func (r *Resolver) explainTrait(
    ctx domain.Context,
    step string,
    trait string,
    overrides []domain.Override,
) ExplainResult {

    var candidates []domain.Override

    // STEP 1: Filter matching overrides
    for _, o := range overrides {
        if o.Matches(ctx) && o.IsActive(ctx.AsOfDate) {
            candidates = append(candidates, o)
        }
    }

    // STEP 2: Sort by priority
    sort.Slice(candidates, func(i, j int) bool {
        if candidates[i].Specificity != candidates[j].Specificity {
            return candidates[i].Specificity > candidates[j].Specificity
        }
        return candidates[i].EffectiveDate.After(candidates[j].EffectiveDate)
    })

    var explanations []CandidateExplanation

    if len(candidates) == 0 {
        return ExplainResult{
            Step: step,
            Trait: trait,
            ResolvedValue: nil,
            Candidates: []CandidateExplanation{},
        }
    }

    winner := candidates[0]

    // STEP 3: Build explanation list
    for i, c := range candidates {

        outcome := ""

        if i == 0 {
            outcome = "SELECTED — highest priority"
        } else {
            outcome = "SHADOWED — lower priority"
        }

        explanations = append(explanations, CandidateExplanation{
            OverrideID:  c.ID,
            Specificity: c.Specificity,
            EffectiveDate: c.EffectiveDate,
            Value: c.Value,
            Outcome: outcome,
        })
    }

    return ExplainResult{
        Step: step,
        Trait: trait,
        ResolvedValue: winner.Value,
        ResolvedFrom: &winner,
        Candidates: explanations,
    }
}