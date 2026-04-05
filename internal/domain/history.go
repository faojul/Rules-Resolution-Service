package domain

import "time"

type OverrideHistory struct {
	ID         int
	OverrideID string
	Before     any
	After      any
	ChangedAt  time.Time
	ChangedBy  string
}