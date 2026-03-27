package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FlowRun represents a Prefect flow run.
type FlowRun struct {
	ID                     uuid.UUID              `json:"id"`
	Created                *time.Time             `json:"created"`
	Updated                *time.Time             `json:"updated"`
	FlowID                 uuid.UUID              `json:"flow_id"`
	Name                   string                 `json:"name"`
	StateID                *uuid.UUID             `json:"state_id"`
	DeploymentID           *uuid.UUID             `json:"deployment_id"`
	WorkQueueID            *uuid.UUID             `json:"work_queue_id"`
	WorkQueueName          *string                `json:"work_queue_name"`
	FlowVersion            *string                `json:"flow_version"`
	Parameters             map[string]interface{} `json:"parameters,omitempty"`
	IdempotencyKey         *string                `json:"idempotency_key"`
	Context                map[string]interface{} `json:"context,omitempty"`
	Tags                   []string               `json:"tags,omitempty"`
	Labels                 map[string]interface{} `json:"labels,omitempty"`
	ParentTaskRunID        *uuid.UUID             `json:"parent_task_run_id"`
	StateType              *StateType             `json:"state_type"`
	StateName              *string                `json:"state_name"`
	RunCount               int                    `json:"run_count"`
	ExpectedStartTime      *time.Time             `json:"expected_start_time"`
	NextScheduledStartTime *time.Time             `json:"next_scheduled_start_time"`
	StartTime              *time.Time             `json:"start_time"`
	EndTime                *time.Time             `json:"end_time"`
	TotalRunTime           float64                `json:"total_run_time"`
	State                  *State                 `json:"state"`
}

// FlowRunCreate represents the request to create a flow run.
type FlowRunCreate struct {
	FlowID         uuid.UUID              `json:"flow_id"`
	Name           string                 `json:"name,omitempty"`
	State          *StateCreate           `json:"state,omitempty"`
	Parameters     map[string]interface{} `json:"parameters,omitempty"`
	Context        map[string]interface{} `json:"context,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Labels         map[string]interface{} `json:"labels,omitempty"`
	IdempotencyKey *string                `json:"idempotency_key,omitempty"`
	WorkPoolName   *string                `json:"work_pool_name,omitempty"`
	WorkQueueName  *string                `json:"work_queue_name,omitempty"`
	DeploymentID   *uuid.UUID             `json:"deployment_id,omitempty"`
}

// FlowRunUpdate represents the request to update a flow run.
type FlowRunUpdate struct {
	Name *string `json:"name,omitempty"`
}

// FlowRunFilter represents filter criteria for querying flow runs.
type FlowRunFilter struct {
	FlowID       *uuid.UUID `json:"flow_id,omitempty"`
	DeploymentID *uuid.UUID `json:"deployment_id,omitempty"`
	StateType    *StateType `json:"state_type,omitempty"`
	Tags         []string   `json:"tags,omitempty"`
	Offset       int        `json:"offset,omitempty"`
	Limit        int        `json:"limit,omitempty"`
}

// FlowRunPolicy represents a flow run's retry policy.
type FlowRunPolicy struct {
	MaxRetries        int      `json:"max_retries,omitempty"`
	RetryDelaySeconds float64  `json:"retry_delay_seconds,omitempty"`
	Retries           *int     `json:"retries,omitempty"`
	RetryDelay        *int     `json:"retry_delay,omitempty"`
	PauseKeys         []string `json:"pause_keys,omitempty"`
	Resuming          *bool    `json:"resuming,omitempty"`
}

// CreatedBy represents information about who created an object.
type CreatedBy struct {
	ID           *uuid.UUID `json:"id,omitempty"`
	Type         *string    `json:"type,omitempty"`
	DisplayValue *string    `json:"display_value,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (fr *FlowRun) UnmarshalJSON(data []byte) error {
	type Alias FlowRun
	aux := &struct {
		Created                optTime `json:"created"`
		Updated                optTime `json:"updated"`
		ExpectedStartTime      optTime `json:"expected_start_time"`
		NextScheduledStartTime optTime `json:"next_scheduled_start_time"`
		StartTime              optTime `json:"start_time"`
		EndTime                optTime `json:"end_time"`
		*Alias
	}{
		Alias: (*Alias)(fr),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	fr.Created = aux.Created.V
	fr.Updated = aux.Updated.V
	fr.ExpectedStartTime = aux.ExpectedStartTime.V
	fr.NextScheduledStartTime = aux.NextScheduledStartTime.V
	fr.StartTime = aux.StartTime.V
	fr.EndTime = aux.EndTime.V
	return nil
}
