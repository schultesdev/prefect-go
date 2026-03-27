package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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

// VariableFilter represents filter criteria for querying variables.
type VariableFilter struct {
	Name   *string  `json:"name,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Limit  int      `json:"limit,omitempty"`
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
