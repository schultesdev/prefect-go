package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// TaskRun represents a Prefect task run.
type TaskRun struct {
	ID           uuid.UUID  `json:"id"`
	Created      *time.Time `json:"created"`
	Updated      *time.Time `json:"updated"`
	Name         string     `json:"name"`
	FlowRunID    uuid.UUID  `json:"flow_run_id"`
	TaskKey      string     `json:"task_key"`
	DynamicKey   string     `json:"dynamic_key"`
	CacheKey     *string    `json:"cache_key"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	TotalRunTime float64    `json:"total_run_time"`
	Status       *StateType `json:"status"`
	StateID      *uuid.UUID `json:"state_id"`
	Tags         []string   `json:"tags,omitempty"`
	State        *State     `json:"state"`
}

// TaskRunCreate represents the request to create a task run.
type TaskRunCreate struct {
	Name       string       `json:"name"`
	FlowRunID  uuid.UUID    `json:"flow_run_id"`
	TaskKey    string       `json:"task_key"`
	DynamicKey string       `json:"dynamic_key"`
	CacheKey   *string      `json:"cache_key,omitempty"`
	State      *StateCreate `json:"state,omitempty"`
	Tags       []string     `json:"tags,omitempty"`
}

// TaskRunUpdate represents the request to update a task run.
type TaskRunUpdate struct {
	Name *string `json:"name,omitempty"`
}

// TaskRunFilter represents filter criteria for querying task runs.
type TaskRunFilter struct {
	FlowRunID *uuid.UUID `json:"flow_run_id,omitempty"`
	StateType *StateType `json:"state_type,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
	Offset    int        `json:"offset,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// TaskRunPolicy represents a task run's retry policy.
type TaskRunPolicy struct {
	MaxRetries        int         `json:"max_retries,omitempty"`
	RetryDelaySeconds float64     `json:"retry_delay_seconds,omitempty"`
	Retries           *int        `json:"retries,omitempty"`
	RetryDelay        interface{} `json:"retry_delay,omitempty"`
	RetryJitterFactor *float64    `json:"retry_jitter_factor,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (tr *TaskRun) UnmarshalJSON(data []byte) error {
	type Alias TaskRun
	aux := &struct {
		Created   optTime `json:"created"`
		Updated   optTime `json:"updated"`
		StartTime optTime `json:"start_time"`
		EndTime   optTime `json:"end_time"`
		*Alias
	}{
		Alias: (*Alias)(tr),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	tr.Created = aux.Created.V
	tr.Updated = aux.Updated.V
	tr.StartTime = aux.StartTime.V
	tr.EndTime = aux.EndTime.V
	return nil
}
