package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
