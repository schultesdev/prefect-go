package models

import (
	"time"

	"github.com/google/uuid"
)

// EventOrder represents sort order for events.
type EventOrder string

const (
	EventOrderAsc  EventOrder = "ASC"
	EventOrderDesc EventOrder = "DESC"
)

// Resource represents an observable business object.
type Resource map[string]string

// RelatedResource represents a resource with a specific role in an event.
type RelatedResource map[string]string

// Event represents a Prefect event.
type Event struct {
	Occurred time.Time              `json:"occurred"`
	Event    string                 `json:"event"`
	Resource Resource               `json:"resource"`
	Related  []RelatedResource      `json:"related,omitempty"`
	ID       uuid.UUID              `json:"id"`
	Payload  map[string]interface{} `json:"payload,omitempty"`
	Received *time.Time             `json:"received,omitempty"`
}

// EventFilter represents filter criteria for querying events.
type EventFilter struct {
	Occurred interface{} `json:"occurred,omitempty"`
	Event    interface{} `json:"event,omitempty"`
	Resource interface{} `json:"resource,omitempty"`
	Related  interface{} `json:"related,omitempty"`
	ID       interface{} `json:"id,omitempty"`
	Order    EventOrder  `json:"order,omitempty"`
}

// EventPage represents a page of events.
type EventPage struct {
	Events   []Event `json:"events"`
	Total    int     `json:"total"`
	NextPage *string `json:"next_page"`
}

// EventCount represents a count of events for a given filter value.
type EventCount struct {
	Value     string    `json:"value"`
	Label     string    `json:"label"`
	Count     int       `json:"count"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
