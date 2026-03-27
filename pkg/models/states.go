package models

import (
	"time"

	"github.com/google/uuid"
)

// StateType represents the type of a flow or task run state.
type StateType string

const (
	StateTypePending    StateType = "PENDING"
	StateTypeRunning    StateType = "RUNNING"
	StateTypeCompleted  StateType = "COMPLETED"
	StateTypeFailed     StateType = "FAILED"
	StateTypeCancelled  StateType = "CANCELLED"
	StateTypeCrashed    StateType = "CRASHED"
	StateTypePaused     StateType = "PAUSED"
	StateTypeScheduled  StateType = "SCHEDULED"
	StateTypeCancelling StateType = "CANCELLING"
)

// State represents the state of a flow or task run.
type State struct {
	ID           uuid.UUID              `json:"id"`
	Type         StateType              `json:"type"`
	Name         *string                `json:"name"`
	Timestamp    time.Time              `json:"timestamp"`
	Message      *string                `json:"message"`
	Data         interface{}            `json:"data,omitempty"`
	StateDetails map[string]interface{} `json:"state_details,omitempty"`
}

// StateCreate represents the request to create a state.
type StateCreate struct {
	Type         StateType              `json:"type"`
	Name         *string                `json:"name,omitempty"`
	Message      *string                `json:"message,omitempty"`
	Data         interface{}            `json:"data,omitempty"`
	StateDetails map[string]interface{} `json:"state_details,omitempty"`
}

// StateDetails represents detailed information about a state.
type StateDetails struct {
	FlowRunID         *uuid.UUID `json:"flow_run_id,omitempty"`
	TaskRunID         *uuid.UUID `json:"task_run_id,omitempty"`
	ChildFlowRunID    *uuid.UUID `json:"child_flow_run_id,omitempty"`
	ScheduledTime     *time.Time `json:"scheduled_time,omitempty"`
	CacheKey          *string    `json:"cache_key,omitempty"`
	CacheExpiration   *time.Time `json:"cache_expiration,omitempty"`
	Deferred          *bool      `json:"deferred,omitempty"`
	UntrackableResult bool       `json:"untrackable_result,omitempty"`
	PauseTimeout      *time.Time `json:"pause_timeout,omitempty"`
}
