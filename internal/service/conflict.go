package service

import (
	"rules-resolution-service/internal/domain"
)

type Conflict struct {
    OverrideA string `json:"overrideA"`
    OverrideB string `json:"overrideB"`
    StepKey   string `json:"stepKey"`
    TraitKey  string `json:"traitKey"`
    Reason    string `json:"reason"`
}

func DetectConflicts(overrides []domain.Override) []Conflict {

    var conflicts []Conflict

    for i := 0; i < len(overrides); i++ {
        for j := i + 1; j < len(overrides); j++ {

            a := overrides[i]
            b := overrides[j]

            // 1. same target
            if a.StepKey != b.StepKey || a.TraitKey != b.TraitKey {
                continue
            }

            // 2. same specificity
            if a.Specificity != b.Specificity {
                continue
            }

            // 3. overlapping date range
            if !dateOverlap(a, b) {
                continue
            }

            // 4. selector overlap
            if !selectorOverlap(a, b) {
                continue
            }

            conflicts = append(conflicts, Conflict{
                OverrideA: a.ID,
                OverrideB: b.ID,
                StepKey:   a.StepKey,
                TraitKey:  a.TraitKey,
                Reason:    "Same specificity + overlapping dates + overlapping selector",
            })
        }
    }

    return conflicts
}

func dateOverlap(a, b domain.Override) bool {

    aEnd := a.ExpiresDate
    bEnd := b.ExpiresDate

    if aEnd == nil && bEnd == nil {
        return true
    }

    if aEnd != nil && aEnd.Before(b.EffectiveDate) {
        return false
    }

    if bEnd != nil && bEnd.Before(a.EffectiveDate) {
        return false
    }

    return true
}

func selectorOverlap(a, b domain.Override) bool {

    if !fieldOverlap(a.State, b.State) {
        return false
    }
    if !fieldOverlap(a.Client, b.Client) {
        return false
    }
    if !fieldOverlap(a.Investor, b.Investor) {
        return false
    }
    if !fieldOverlap(a.CaseType, b.CaseType) {
        return false
    }

    return true
}

func fieldOverlap(a, b *string) bool {

    // both wildcard
    if a == nil && b == nil {
        return true
    }

    // one wildcard → overlap possible
    if a == nil || b == nil {
        return true
    }

    // both set → must match
    return *a == *b
}