package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WorkQueue represents a Prefect work queue.
type WorkQueue struct {
	ID          uuid.UUID  `json:"id"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	IsPaused    bool       `json:"is_paused"`
	Priority    int        `json:"priority"`
}

// WorkQueueCreate represents the request to create a work queue.
type WorkQueueCreate struct {
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	IsPaused         bool   `json:"is_paused,omitempty"`
	ConcurrencyLimit *int   `json:"concurrency_limit,omitempty"`
	Priority         *int   `json:"priority,omitempty"`
}

// WorkQueueUpdate represents the request to update a work queue.
type WorkQueueUpdate struct {
	Name             *string `json:"name,omitempty"`
	Description      *string `json:"description,omitempty"`
	IsPaused         *bool   `json:"is_paused,omitempty"`
	ConcurrencyLimit *int    `json:"concurrency_limit,omitempty"`
	Priority         *int    `json:"priority,omitempty"`
}

// WorkQueueFilter represents filter criteria for querying work queues.
type WorkQueueFilter struct {
	Name   *string `json:"name,omitempty"`
	Offset int     `json:"offset,omitempty"`
	Limit  int     `json:"limit,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (wq *WorkQueue) UnmarshalJSON(data []byte) error {
	type Alias WorkQueue
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(wq),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	wq.Created = aux.Created.V
	wq.Updated = aux.Updated.V
	return nil
}
