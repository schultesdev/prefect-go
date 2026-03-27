package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeploymentSchedule represents a schedule for a deployment.
type DeploymentSchedule struct {
	ID               uuid.UUID              `json:"id"`
	Created          *time.Time             `json:"created"`
	Updated          *time.Time             `json:"updated"`
	DeploymentID     *uuid.UUID             `json:"deployment_id"`
	Schedule         json.RawMessage        `json:"schedule"`
	Active           bool                   `json:"active"`
	MaxScheduledRuns *int                   `json:"max_scheduled_runs,omitempty"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
	Slug             *string                `json:"slug,omitempty"`
}

// DeploymentScheduleCreate represents the request to create a deployment schedule.
type DeploymentScheduleCreate struct {
	Active           *bool                  `json:"active,omitempty"`
	Schedule         json.RawMessage        `json:"schedule"`
	MaxScheduledRuns *int                   `json:"max_scheduled_runs,omitempty"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
	Slug             *string                `json:"slug,omitempty"`
}

// DeploymentScheduleUpdate represents the request to update a deployment schedule.
type DeploymentScheduleUpdate struct {
	Active           *bool                  `json:"active,omitempty"`
	Schedule         json.RawMessage        `json:"schedule,omitempty"`
	MaxScheduledRuns *int                   `json:"max_scheduled_runs,omitempty"`
	Parameters       map[string]interface{} `json:"parameters,omitempty"`
	Slug             *string                `json:"slug,omitempty"`
}

// IntervalSchedule represents an interval-based schedule.
type IntervalSchedule struct {
	Interval   float64 `json:"interval"`
	AnchorDate *string `json:"anchor_date,omitempty"`
	Timezone   *string `json:"timezone,omitempty"`
}

// CronSchedule represents a cron-based schedule.
type CronSchedule struct {
	Cron     string  `json:"cron"`
	Timezone *string `json:"timezone,omitempty"`
	DayOr    *bool   `json:"day_or,omitempty"`
}

// RRuleSchedule represents an RRule-based schedule.
type RRuleSchedule struct {
	RRule    string  `json:"rrule"`
	Timezone *string `json:"timezone,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (ds *DeploymentSchedule) UnmarshalJSON(data []byte) error {
	type Alias DeploymentSchedule
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(ds),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	ds.Created = aux.Created.V
	ds.Updated = aux.Updated.V
	return nil
}
