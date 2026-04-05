package domain

import "time"

type Context struct {
    State    string    `json:"state"`
    Client   string    `json:"client"`
    Investor string    `json:"investor"`
    CaseType string    `json:"caseType"`
    AsOfDate time.Time `json:"asOfDate"`
}