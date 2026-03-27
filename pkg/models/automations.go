package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AutomationSort represents sort options for automations.
type AutomationSort string

const (
	AutomationSortCreatedDesc AutomationSort = "CREATED_DESC"
	AutomationSortUpdatedDesc AutomationSort = "UPDATED_DESC"
	AutomationSortNameAsc     AutomationSort = "NAME_ASC"
	AutomationSortNameDesc    AutomationSort = "NAME_DESC"
)

// TriggerPosture represents the posture of an event trigger.
type TriggerPosture string

const (
	TriggerPostureReactive  TriggerPosture = "Reactive"
	TriggerPostureProactive TriggerPosture = "Proactive"
)

// ActionSource represents the source type for automation actions.
type ActionSource string

const (
	ActionSourceSelected ActionSource = "selected"
	ActionSourceInferred ActionSource = "inferred"
)

// ResourceSpecification represents resource labels to match.
type ResourceSpecification map[string]string

// Automation represents a Prefect automation.
type Automation struct {
	ID               uuid.UUID         `json:"id"`
	Created          *time.Time        `json:"created"`
	Updated          *time.Time        `json:"updated"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Enabled          bool              `json:"enabled"`
	Tags             []string          `json:"tags,omitempty"`
	Trigger          json.RawMessage   `json:"trigger"`
	Actions          []json.RawMessage `json:"actions"`
	ActionsOnTrigger []json.RawMessage `json:"actions_on_trigger,omitempty"`
	ActionsOnResolve []json.RawMessage `json:"actions_on_resolve,omitempty"`
}

// AutomationCreate represents the request to create an automation.
type AutomationCreate struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Trigger          json.RawMessage   `json:"trigger"`
	Actions          []json.RawMessage `json:"actions"`
	ActionsOnTrigger []json.RawMessage `json:"actions_on_trigger,omitempty"`
	ActionsOnResolve []json.RawMessage `json:"actions_on_resolve,omitempty"`
}

// AutomationUpdate represents the request to fully update an automation.
type AutomationUpdate struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Enabled          *bool             `json:"enabled,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Trigger          json.RawMessage   `json:"trigger"`
	Actions          []json.RawMessage `json:"actions"`
	ActionsOnTrigger []json.RawMessage `json:"actions_on_trigger,omitempty"`
	ActionsOnResolve []json.RawMessage `json:"actions_on_resolve,omitempty"`
}

// AutomationPartialUpdate represents the request to partially update an automation.
type AutomationPartialUpdate struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// AutomationFilter represents filter criteria for querying automations.
type AutomationFilter struct {
	Name   *string  `json:"name,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}

// EventTrigger represents an event-based automation trigger.
type EventTrigger struct {
	Type         string                `json:"type"`
	ID           uuid.UUID             `json:"id"`
	Match        ResourceSpecification `json:"match,omitempty"`
	MatchRelated json.RawMessage       `json:"match_related,omitempty"`
	After        []string              `json:"after,omitempty"`
	Expect       []string              `json:"expect,omitempty"`
	ForEach      []string              `json:"for_each,omitempty"`
	Posture      TriggerPosture        `json:"posture"`
	Threshold    int                   `json:"threshold,omitempty"`
	Within       float64               `json:"within,omitempty"`
}

// CompoundTrigger represents a compound automation trigger.
type CompoundTrigger struct {
	Type     string          `json:"type"`
	ID       uuid.UUID       `json:"id"`
	Triggers json.RawMessage `json:"triggers"`
	Within   *float64        `json:"within"`
	Require  json.RawMessage `json:"require"`
}

// SequenceTrigger represents a sequence automation trigger.
type SequenceTrigger struct {
	Type     string          `json:"type"`
	ID       uuid.UUID       `json:"id"`
	Triggers json.RawMessage `json:"triggers"`
	Within   *float64        `json:"within"`
}

// DoNothing represents a no-op automation action.
type DoNothing struct {
	Type string `json:"type"`
}

// RunDeployment represents an action to run a deployment.
type RunDeployment struct {
	Type         string                 `json:"type"`
	Source       ActionSource           `json:"source,omitempty"`
	DeploymentID *uuid.UUID             `json:"deployment_id,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	JobVariables map[string]interface{} `json:"job_variables,omitempty"`
}

// PauseDeploymentAction represents an action to pause a deployment.
type PauseDeploymentAction struct {
	Type         string       `json:"type"`
	Source       ActionSource `json:"source,omitempty"`
	DeploymentID *uuid.UUID   `json:"deployment_id,omitempty"`
}

// ResumeDeploymentAction represents an action to resume a deployment.
type ResumeDeploymentAction struct {
	Type         string       `json:"type"`
	Source       ActionSource `json:"source,omitempty"`
	DeploymentID *uuid.UUID   `json:"deployment_id,omitempty"`
}

// CancelFlowRun represents an action to cancel a flow run.
type CancelFlowRun struct {
	Type string `json:"type"`
}

// SuspendFlowRun represents an action to suspend a flow run.
type SuspendFlowRun struct {
	Type string `json:"type"`
}

// ResumeFlowRunAction represents an action to resume a flow run.
type ResumeFlowRunAction struct {
	Type string `json:"type"`
}

// ChangeFlowRunState represents an action to change a flow run's state.
type ChangeFlowRunState struct {
	Type    string    `json:"type"`
	Name    *string   `json:"name,omitempty"`
	State   StateType `json:"state"`
	Message *string   `json:"message,omitempty"`
	Force   bool      `json:"force,omitempty"`
}

// PauseWorkQueueAction represents an action to pause a work queue.
type PauseWorkQueueAction struct {
	Type        string       `json:"type"`
	Source      ActionSource `json:"source,omitempty"`
	WorkQueueID *uuid.UUID   `json:"work_queue_id,omitempty"`
}

// ResumeWorkQueueAction represents an action to resume a work queue.
type ResumeWorkQueueAction struct {
	Type        string       `json:"type"`
	Source      ActionSource `json:"source,omitempty"`
	WorkQueueID *uuid.UUID   `json:"work_queue_id,omitempty"`
}

// PauseWorkPoolAction represents an action to pause a work pool.
type PauseWorkPoolAction struct {
	Type       string       `json:"type"`
	Source     ActionSource `json:"source,omitempty"`
	WorkPoolID *uuid.UUID   `json:"work_pool_id,omitempty"`
}

// ResumeWorkPoolAction represents an action to resume a work pool.
type ResumeWorkPoolAction struct {
	Type       string       `json:"type"`
	Source     ActionSource `json:"source,omitempty"`
	WorkPoolID *uuid.UUID   `json:"work_pool_id,omitempty"`
}

// PauseAutomationAction represents an action to pause an automation.
type PauseAutomationAction struct {
	Type         string       `json:"type"`
	Source       ActionSource `json:"source,omitempty"`
	AutomationID *uuid.UUID   `json:"automation_id,omitempty"`
}

// ResumeAutomationAction represents an action to resume an automation.
type ResumeAutomationAction struct {
	Type         string       `json:"type"`
	Source       ActionSource `json:"source,omitempty"`
	AutomationID *uuid.UUID   `json:"automation_id,omitempty"`
}

// SendNotification represents an action to send a notification.
type SendNotification struct {
	Type            string    `json:"type"`
	BlockDocumentID uuid.UUID `json:"block_document_id"`
	Subject         string    `json:"subject,omitempty"`
	Body            string    `json:"body"`
}

// CallWebhook represents an action to call a webhook.
type CallWebhook struct {
	Type            string    `json:"type"`
	BlockDocumentID uuid.UUID `json:"block_document_id"`
	Payload         string    `json:"payload,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (a *Automation) UnmarshalJSON(data []byte) error {
	type Alias Automation
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Created = aux.Created.V
	a.Updated = aux.Updated.V
	return nil
}
