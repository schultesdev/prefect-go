package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ConcurrencyLimit represents a Prefect v1 concurrency limit.
type ConcurrencyLimit struct {
	ID               uuid.UUID  `json:"id"`
	Created          *time.Time `json:"created"`
	Updated          *time.Time `json:"updated"`
	Tag              string     `json:"tag"`
	ConcurrencyLimit int        `json:"concurrency_limit"`
	ActiveSlots      []string   `json:"active_slots,omitempty"`
}

// ConcurrencyLimitCreate represents the request to create a v1 concurrency limit.
type ConcurrencyLimitCreate struct {
	Tag              string `json:"tag"`
	ConcurrencyLimit int    `json:"concurrency_limit"`
}

// ConcurrencyLimitV2 represents a Prefect v2 concurrency limit.
type ConcurrencyLimitV2 struct {
	ID                      uuid.UUID  `json:"id"`
	Created                 *time.Time `json:"created"`
	Updated                 *time.Time `json:"updated"`
	Active                  bool       `json:"active"`
	Name                    string     `json:"name"`
	Limit                   int        `json:"limit"`
	ActiveSlots             int        `json:"active_slots"`
	DeniedSlots             int        `json:"denied_slots"`
	SlotDecayPerSecond      float64    `json:"slot_decay_per_second"`
	AvgSlotOccupancySeconds float64    `json:"avg_slot_occupancy_seconds"`
}

// ConcurrencyLimitV2Create represents the request to create a v2 concurrency limit.
type ConcurrencyLimitV2Create struct {
	Active             *bool    `json:"active,omitempty"`
	Name               string   `json:"name"`
	Limit              int      `json:"limit"`
	ActiveSlots        *int     `json:"active_slots,omitempty"`
	DeniedSlots        *int     `json:"denied_slots,omitempty"`
	SlotDecayPerSecond *float64 `json:"slot_decay_per_second,omitempty"`
}

// ConcurrencyLimitV2Update represents the request to update a v2 concurrency limit.
type ConcurrencyLimitV2Update struct {
	Active             *bool    `json:"active,omitempty"`
	Name               *string  `json:"name,omitempty"`
	Limit              *int     `json:"limit,omitempty"`
	ActiveSlots        *int     `json:"active_slots,omitempty"`
	DeniedSlots        *int     `json:"denied_slots,omitempty"`
	SlotDecayPerSecond *float64 `json:"slot_decay_per_second,omitempty"`
}

// GlobalConcurrencyLimitResponse represents the response for global concurrency limits.
type GlobalConcurrencyLimitResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Created            *time.Time `json:"created"`
	Updated            *time.Time `json:"updated"`
	Active             bool       `json:"active"`
	Name               string     `json:"name"`
	Limit              int        `json:"limit"`
	ActiveSlots        int        `json:"active_slots"`
	SlotDecayPerSecond float64    `json:"slot_decay_per_second"`
}

// MinimalConcurrencyLimitResponse represents a minimal concurrency limit response.
type MinimalConcurrencyLimitResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Limit int       `json:"limit"`
}

// ConcurrencyLimitWithLeaseResponse represents a concurrency limit response with lease info.
type ConcurrencyLimitWithLeaseResponse struct {
	LeaseID uuid.UUID                         `json:"lease_id"`
	Limits  []MinimalConcurrencyLimitResponse `json:"limits"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (cl *ConcurrencyLimit) UnmarshalJSON(data []byte) error {
	type Alias ConcurrencyLimit
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(cl),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	cl.Created = aux.Created.V
	cl.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (cl *ConcurrencyLimitV2) UnmarshalJSON(data []byte) error {
	type Alias ConcurrencyLimitV2
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(cl),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	cl.Created = aux.Created.V
	cl.Updated = aux.Updated.V
	return nil
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (g *GlobalConcurrencyLimitResponse) UnmarshalJSON(data []byte) error {
	type Alias GlobalConcurrencyLimitResponse
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	g.Created = aux.Created.V
	g.Updated = aux.Updated.V
	return nil
}
