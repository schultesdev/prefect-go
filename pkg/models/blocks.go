package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
