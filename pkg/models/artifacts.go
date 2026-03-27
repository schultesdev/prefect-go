package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ArtifactSort represents sort options for artifacts.
type ArtifactSort string

const (
	ArtifactSortCreatedDesc ArtifactSort = "CREATED_DESC"
	ArtifactSortUpdatedDesc ArtifactSort = "UPDATED_DESC"
)

// ArtifactCollectionSort represents sort options for artifact collections.
type ArtifactCollectionSort string

const (
	ArtifactCollectionSortCreatedDesc ArtifactCollectionSort = "CREATED_DESC"
	ArtifactCollectionSortUpdatedDesc ArtifactCollectionSort = "UPDATED_DESC"
)

// Artifact represents a Prefect artifact.
type Artifact struct {
	ID          uuid.UUID         `json:"id"`
	Created     *time.Time        `json:"created"`
	Updated     *time.Time        `json:"updated"`
	Key         *string           `json:"key"`
	Type        *string           `json:"type"`
	Description *string           `json:"description"`
	Data        interface{}       `json:"data,omitempty"`
	Metadata    map[string]string `json:"metadata_,omitempty"`
	FlowRunID   *uuid.UUID        `json:"flow_run_id"`
	TaskRunID   *uuid.UUID        `json:"task_run_id"`
}

// ArtifactCreate represents the request to create an artifact.
type ArtifactCreate struct {
	Key         *string           `json:"key,omitempty"`
	Type        *string           `json:"type,omitempty"`
	Description *string           `json:"description,omitempty"`
	Data        interface{}       `json:"data,omitempty"`
	Metadata    map[string]string `json:"metadata_,omitempty"`
	FlowRunID   *uuid.UUID        `json:"flow_run_id,omitempty"`
	TaskRunID   *uuid.UUID        `json:"task_run_id,omitempty"`
}

// ArtifactUpdate represents the request to update an artifact.
type ArtifactUpdate struct {
	Data        interface{}       `json:"data,omitempty"`
	Description *string           `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata_,omitempty"`
}

// ArtifactFilter represents filter criteria for querying artifacts.
type ArtifactFilter struct {
	Key       *string    `json:"key,omitempty"`
	FlowRunID *uuid.UUID `json:"flow_run_id,omitempty"`
	TaskRunID *uuid.UUID `json:"task_run_id,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Offset    int        `json:"offset,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// ArtifactCollection represents a collection of artifacts with the same key.
type ArtifactCollection struct {
	ID          uuid.UUID         `json:"id"`
	Created     *time.Time        `json:"created"`
	Updated     *time.Time        `json:"updated"`
	Key         string            `json:"key"`
	LatestID    uuid.UUID         `json:"latest_id"`
	Type        *string           `json:"type"`
	Description *string           `json:"description"`
	Data        interface{}       `json:"data,omitempty"`
	Metadata    map[string]string `json:"metadata_,omitempty"`
	FlowRunID   *uuid.UUID        `json:"flow_run_id"`
	TaskRunID   *uuid.UUID        `json:"task_run_id"`
}

// ArtifactCollectionFilter represents filter criteria for querying artifact collections.
type ArtifactCollectionFilter struct {
	Key       *string    `json:"key,omitempty"`
	FlowRunID *uuid.UUID `json:"flow_run_id,omitempty"`
	TaskRunID *uuid.UUID `json:"task_run_id,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Offset    int        `json:"offset,omitempty"`
	Limit     int        `json:"limit,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (a *Artifact) UnmarshalJSON(data []byte) error {
	type Alias Artifact
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

// UnmarshalJSON implements custom JSON unmarshaling for time fields.
func (ac *ArtifactCollection) UnmarshalJSON(data []byte) error {
	type Alias ArtifactCollection
	aux := &struct {
		Created optTime `json:"created"`
		Updated optTime `json:"updated"`
		*Alias
	}{
		Alias: (*Alias)(ac),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	ac.Created = aux.Created.V
	ac.Updated = aux.Updated.V
	return nil
}
