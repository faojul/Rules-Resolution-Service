package domain

import (
	"time"
)

type Override struct {
	ID string `json:"id"`

	StepKey  string `json:"stepKey"`
	TraitKey string `json:"traitKey"`

	Selector map[string]string `json:"selector"`

	Value any `json:"value"`

	Specificity int `json:"specificity"`

	EffectiveDate time.Time  `json:"effectiveDate"`
	ExpiresDate   *time.Time `json:"expiresDate,omitempty"`

	Status string `json:"status"`

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
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

func Matches(selector map[string]string, ctx Context) bool {
	for k, v := range selector {
		switch k {
		case "state":
			if ctx.State != v {
				return false
			}
		case "client":
			if ctx.Client != v {
				return false
			}
		case "investor":
			if ctx.Investor != v {
				return false
			}
		case "caseType":
			if ctx.CaseType != v {
				return false
			}
		}
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

func ComputeSpecificity(selector map[string]string) int {
    return len(selector)
}