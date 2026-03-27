package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SavedSearch represents a Prefect saved search.
type SavedSearch struct {
	ID      uuid.UUID           `json:"id"`
	Created *time.Time          `json:"created"`
	Updated *time.Time          `json:"updated"`
	Name    string              `json:"name"`
	Filters []SavedSearchFilter `json:"filters"`
}

// SavedSearchCreate represents the request to create a saved search.
type SavedSearchCreate struct {
	Name    string              `json:"name"`
	Filters []SavedSearchFilter `json:"filters,omitempty"`
}

// SavedSearchFilter represents a filter definition within a saved search.
type SavedSearchFilter struct {
	Object    string      `json:"object"`
	Property  string      `json:"property"`
	Type      string      `json:"type"`
	Operation string      `json:"operation"`
	Value     interface{} `json:"value"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (ss *SavedSearch) UnmarshalJSON(data []byte) error {
	type Alias SavedSearch
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(ss),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	ss.Created = aux.Created.V
	ss.Updated = aux.Updated.V
	return nil
}
