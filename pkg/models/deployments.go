package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

// ConcurrencyOptions represents concurrency configuration for a deployment.
type ConcurrencyOptions struct {
	CollisionStrategy  CollisionStrategy `json:"collision_strategy"`
	GracePeriodSeconds *int              `json:"grace_period_seconds,omitempty"`
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

// DeploymentFilter represents filter criteria for querying deployments.
type DeploymentFilter struct {
	FlowID *uuid.UUID `json:"flow_id,omitempty"`
	Name   *string    `json:"name,omitempty"`
	Tags   []string   `json:"tags,omitempty"`
	Paused *bool      `json:"paused,omitempty"`
	Offset int        `json:"offset,omitempty"`
	Limit  int        `json:"limit,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (d *Deployment) UnmarshalJSON(data []byte) error {
	type Alias Deployment
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	d.Created = aux.Created.V
	d.Updated = aux.Updated.V
	return nil
}
