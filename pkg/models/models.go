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

// BlockDocumentSort represents sort options for block documents.
type BlockDocumentSort string

const (
	BlockDocumentSortNameAsc  BlockDocumentSort = "NAME_ASC"
	BlockDocumentSortNameDesc BlockDocumentSort = "NAME_DESC"
)

// BlockType represents a Prefect block type.
type BlockType struct {
	ID               uuid.UUID  `json:"id"`
	Created          *time.Time `json:"created"`
	Updated          *time.Time `json:"updated"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	LogoURL          *string    `json:"logo_url"`
	DocumentationURL *string    `json:"documentation_url"`
	Description      *string    `json:"description"`
	CodeExample      *string    `json:"code_example"`
	IsProtected      bool       `json:"is_protected"`
}

// BlockTypeCreate represents the request to create a block type.
type BlockTypeCreate struct {
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	LogoURL          *string `json:"logo_url,omitempty"`
	DocumentationURL *string `json:"documentation_url,omitempty"`
	Description      *string `json:"description,omitempty"`
	CodeExample      *string `json:"code_example,omitempty"`
}

// BlockTypeUpdate represents the request to update a block type.
type BlockTypeUpdate struct {
	LogoURL          *string `json:"logo_url,omitempty"`
	DocumentationURL *string `json:"documentation_url,omitempty"`
	Description      *string `json:"description,omitempty"`
	CodeExample      *string `json:"code_example,omitempty"`
}

// BlockTypeFilter represents filter criteria for querying block types.
type BlockTypeFilter struct {
	Name         *string  // Filter by name (partial match)
	Slugs        []string // Filter by slugs
	Capabilities []string // Filter by block capabilities
}

// BlockSchema represents a Prefect block schema.
type BlockSchema struct {
	ID           uuid.UUID              `json:"id"`
	Created      *time.Time             `json:"created"`
	Updated      *time.Time             `json:"updated"`
	Checksum     string                 `json:"checksum"`
	Fields       map[string]interface{} `json:"fields"`
	BlockTypeID  *uuid.UUID             `json:"block_type_id"`
	BlockType    *BlockType             `json:"block_type"`
	Capabilities []string               `json:"capabilities"`
	Version      string                 `json:"version"`
}

// BlockSchemaCreate represents the request to create a block schema.
type BlockSchemaCreate struct {
	Fields       map[string]interface{} `json:"fields,omitempty"`
	BlockTypeID  uuid.UUID              `json:"block_type_id"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Version      *string                `json:"version,omitempty"`
}

// BlockSchemaFilter represents filter criteria for querying block schemas.
type BlockSchemaFilter struct {
	BlockTypeID  *uuid.UUID // Filter by block type ID
	Capabilities []string   // Filter by required capabilities
	Version      *string    // Filter by schema version
}

// BlockDocument represents a Prefect block document.
type BlockDocument struct {
	ID                      uuid.UUID              `json:"id"`
	Created                 *time.Time             `json:"created"`
	Updated                 *time.Time             `json:"updated"`
	Name                    *string                `json:"name"`
	Data                    map[string]interface{} `json:"data"`
	BlockSchemaID           uuid.UUID              `json:"block_schema_id"`
	BlockSchema             *BlockSchema           `json:"block_schema"`
	BlockTypeID             uuid.UUID              `json:"block_type_id"`
	BlockTypeName           *string                `json:"block_type_name"`
	BlockType               *BlockType             `json:"block_type"`
	BlockDocumentReferences map[string]interface{} `json:"block_document_references"`
	IsAnonymous             bool                   `json:"is_anonymous"`
}

// BlockDocumentCreate represents the request to create a block document.
type BlockDocumentCreate struct {
	Name          *string                `json:"name,omitempty"`
	Data          map[string]interface{} `json:"data,omitempty"`
	BlockSchemaID uuid.UUID              `json:"block_schema_id"`
	BlockTypeID   uuid.UUID              `json:"block_type_id"`
	IsAnonymous   bool                   `json:"is_anonymous,omitempty"`
}

// BlockDocumentUpdate represents the request to update a block document.
type BlockDocumentUpdate struct {
	BlockSchemaID *uuid.UUID             `json:"block_schema_id,omitempty"`
	Data          map[string]interface{} `json:"data"`
	MergeExisting *bool                  `json:"merge_existing_data,omitempty"`
}

// BlockDocumentFilter represents filter criteria for querying block documents.
type BlockDocumentFilter struct {
	BlockTypeID *uuid.UUID // Filter by block type ID
	Name        *string    // Filter by document name
	IsAnonymous *bool      // Filter by anonymity
}

// parseOptionalTime parses a time string, returning nil for empty strings.
// Supports RFC3339Nano and RFC3339 formats used by the Prefect API.
func parseOptionalTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// optTime handles JSON unmarshaling of nullable time fields that may be
// represented as empty strings or JSON null by the Prefect API.
type optTime struct {
	V *time.Time
}

func (o *optTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.V = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := parseOptionalTime(s)
	if err != nil {
		return err
	}
	o.V = v
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (f *Flow) UnmarshalJSON(data []byte) error {
	type Alias Flow
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(f),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	f.Created = aux.Created.V
	f.Updated = aux.Updated.V
	return nil
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

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (v *Variable) UnmarshalJSON(data []byte) error {
	type Alias Variable
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v.Created = aux.Created.V
	v.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (l *Log) UnmarshalJSON(data []byte) error {
	type Alias Log
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	l.Created = aux.Created.V
	l.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (bt *BlockType) UnmarshalJSON(data []byte) error {
	type Alias BlockType
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(bt),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	bt.Created = aux.Created.V
	bt.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (bs *BlockSchema) UnmarshalJSON(data []byte) error {
	type Alias BlockSchema
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(bs),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	bs.Created = aux.Created.V
	bs.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (bd *BlockDocument) UnmarshalJSON(data []byte) error {
	type Alias BlockDocument
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(bd),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	bd.Created = aux.Created.V
	bd.Updated = aux.Updated.V
	return nil
}
