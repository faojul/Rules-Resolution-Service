package domain

import "time"

type Override struct {
    ID string

    StepKey  string
    TraitKey string

    State    *string
    Client   *string
    Investor *string
    CaseType *string

    Value any

    Specificity int

    EffectiveDate time.Time
    ExpiresDate   *time.Time

    Status string

    CreatedAt   *time.Time
    UpdatedAt   *time.Time
}

type OverrideFilter struct {
	StepKey  *string
	TraitKey *string
	State    *string
	Client   *string
	Investor *string
	CaseType *string
	Status   *string
}

func (o Override) Matches(ctx Context) bool {
    if o.State != nil && *o.State != ctx.State {
        return false
    }
    if o.Client != nil && *o.Client != ctx.Client {
        return false
    }
    if o.Investor != nil && *o.Investor != ctx.Investor {
        return false
    }
    if o.CaseType != nil && *o.CaseType != ctx.CaseType {
        return false
    }
    return true
}

func (o Override) IsActive(asOf time.Time) bool {
    if o.Status != "active" {
        return false
    }
    if o.EffectiveDate.After(asOf) {
        return false
    }
    if o.ExpiresDate != nil && !o.ExpiresDate.After(asOf) {
        return false
    }
    return true
}