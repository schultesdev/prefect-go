package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WorkPoolStatus represents the status of a work pool.
type WorkPoolStatus string

const (
	WorkPoolStatusReady    WorkPoolStatus = "READY"
	WorkPoolStatusNotReady WorkPoolStatus = "NOT_READY"
	WorkPoolStatusPaused   WorkPoolStatus = "PAUSED"
)

// WorkPool represents a Prefect work pool.
type WorkPool struct {
	ID          uuid.UUID  `json:"id"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Type        string     `json:"type"`
	IsPaused    bool       `json:"is_paused"`
}

// WorkPoolCreate represents the request to create a work pool.
type WorkPoolCreate struct {
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	Type             string                 `json:"type"`
	BaseJobTemplate  map[string]interface{} `json:"base_job_template,omitempty"`
	IsPaused         bool                   `json:"is_paused,omitempty"`
	ConcurrencyLimit *int                   `json:"concurrency_limit,omitempty"`
}

// WorkPoolUpdate represents the request to update a work pool.
type WorkPoolUpdate struct {
	Description      *string                `json:"description,omitempty"`
	IsPaused         *bool                  `json:"is_paused,omitempty"`
	BaseJobTemplate  map[string]interface{} `json:"base_job_template,omitempty"`
	ConcurrencyLimit *int                   `json:"concurrency_limit,omitempty"`
}

// WorkPoolFilter represents filter criteria for querying work pools.
type WorkPoolFilter struct {
	Name   *string `json:"name,omitempty"`
	Type   *string `json:"type,omitempty"`
	Offset int     `json:"offset,omitempty"`
	Limit  int     `json:"limit,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (wp *WorkPool) UnmarshalJSON(data []byte) error {
	type Alias WorkPool
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(wp),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	wp.Created = aux.Created.V
	wp.Updated = aux.Updated.V
	return nil
}
