package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
