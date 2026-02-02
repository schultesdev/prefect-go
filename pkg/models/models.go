// Package models contains the data structures for the Prefect API.
package models

import (
	"encoding/json"
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

// DeploymentStatus represents the status of a deployment.
type DeploymentStatus string

const (
	DeploymentStatusReady    DeploymentStatus = "READY"
	DeploymentStatusNotReady DeploymentStatus = "NOT_READY"
)

// CollisionStrategy represents the strategy for handling concurrent flow runs.
type CollisionStrategy string

const (
	CollisionStrategyEnqueue   CollisionStrategy = "ENQUEUE"
	CollisionStrategyCancelNew CollisionStrategy = "CANCEL_NEW"
)

// Flow represents a Prefect flow.
type Flow struct {
	ID      uuid.UUID              `json:"id"`
	Created *time.Time             `json:"created"`
	Updated *time.Time             `json:"updated"`
	Name    string                 `json:"name"`
	Tags    []string               `json:"tags,omitempty"`
	Labels  map[string]interface{} `json:"labels,omitempty"`
}

// FlowCreate represents the request to create a flow.
type FlowCreate struct {
	Name   string                 `json:"name"`
	Tags   []string               `json:"tags,omitempty"`
	Labels map[string]interface{} `json:"labels,omitempty"`
}

// FlowUpdate represents the request to update a flow.
type FlowUpdate struct {
	Tags   *[]string               `json:"tags,omitempty"`
	Labels *map[string]interface{} `json:"labels,omitempty"`
}

// FlowFilter represents filter criteria for querying flows.
type FlowFilter struct {
	Name   *string  `json:"name,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}

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

// Deployment represents a Prefect deployment.
type Deployment struct {
	ID                     uuid.UUID              `json:"id"`
	Created                *time.Time             `json:"created"`
	Updated                *time.Time             `json:"updated"`
	Name                   string                 `json:"name"`
	FlowID                 uuid.UUID              `json:"flow_id"`
	Version                *string                `json:"version"`
	Description            *string                `json:"description"`
	Paused                 bool                   `json:"paused"`
	Parameters             map[string]interface{} `json:"parameters,omitempty"`
	Tags                   []string               `json:"tags,omitempty"`
	Labels                 map[string]interface{} `json:"labels,omitempty"`
	WorkQueueName          *string                `json:"work_queue_name"`
	WorkPoolName           *string                `json:"work_pool_name"`
	Path                   *string                `json:"path"`
	Entrypoint             *string                `json:"entrypoint"`
	Status                 *DeploymentStatus      `json:"status"`
	EnforceParameterSchema bool                   `json:"enforce_parameter_schema"`
}

// DeploymentCreate represents the request to create a deployment.
type DeploymentCreate struct {
	Name                   string                 `json:"name"`
	FlowID                 uuid.UUID              `json:"flow_id"`
	Paused                 bool                   `json:"paused,omitempty"`
	Description            *string                `json:"description,omitempty"`
	Version                *string                `json:"version,omitempty"`
	Parameters             map[string]interface{} `json:"parameters,omitempty"`
	Tags                   []string               `json:"tags,omitempty"`
	Labels                 map[string]interface{} `json:"labels,omitempty"`
	WorkPoolName           *string                `json:"work_pool_name,omitempty"`
	WorkQueueName          *string                `json:"work_queue_name,omitempty"`
	Path                   *string                `json:"path,omitempty"`
	Entrypoint             *string                `json:"entrypoint,omitempty"`
	EnforceParameterSchema bool                   `json:"enforce_parameter_schema,omitempty"`
}

// DeploymentUpdate represents the request to update a deployment.
type DeploymentUpdate struct {
	Paused      *bool   `json:"paused,omitempty"`
	Description *string `json:"description,omitempty"`
}

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

// Variable represents a Prefect variable.
type Variable struct {
	ID      uuid.UUID  `json:"id"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
	Name    string     `json:"name"`
	Value   string     `json:"value"`
	Tags    []string   `json:"tags,omitempty"`
}

// VariableCreate represents the request to create a variable.
type VariableCreate struct {
	Name  string   `json:"name"`
	Value string   `json:"value"`
	Tags  []string `json:"tags,omitempty"`
}

// VariableUpdate represents the request to update a variable.
type VariableUpdate struct {
	Value *string  `json:"value,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

// Log represents a Prefect log entry.
type Log struct {
	ID        uuid.UUID  `json:"id"`
	Created   *time.Time `json:"created"`
	Updated   *time.Time `json:"updated"`
	Name      string     `json:"name"`
	Level     int        `json:"level"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
	FlowRunID *uuid.UUID `json:"flow_run_id"`
	TaskRunID *uuid.UUID `json:"task_run_id"`
}

// LogCreate represents the request to create a log.
type LogCreate struct {
	Name      string     `json:"name"`
	Level     int        `json:"level"`
	Message   string     `json:"message"`
	Timestamp time.Time  `json:"timestamp"`
	FlowRunID *uuid.UUID `json:"flow_run_id,omitempty"`
	TaskRunID *uuid.UUID `json:"task_run_id,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (f *Flow) UnmarshalJSON(data []byte) error {
	type Alias Flow
	aux := &struct {
		Created string `json:"created"`
		Updated string `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Created != "" {
		t, err := time.Parse(time.RFC3339, aux.Created)
		if err != nil {
			return err
		}
		f.Created = &t
	}

	if aux.Updated != "" {
		t, err := time.Parse(time.RFC3339, aux.Updated)
		if err != nil {
			return err
		}
		f.Updated = &t
	}

	return nil
}
