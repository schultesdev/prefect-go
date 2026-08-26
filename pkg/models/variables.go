package models

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Variable represents a Prefect variable.
//
// Value holds the raw JSON of the variable's value. Prefect 3.x accepts any JSON
// type here - string, integer, number, boolean, object, array or null - so the
// value is kept as json.RawMessage rather than being forced into a Go string.
// Decode it with json.Unmarshal, or use StringValue for the common string case.
type Variable struct {
	ID      uuid.UUID       `json:"id"`
	Created *time.Time      `json:"created"`
	Updated *time.Time      `json:"updated"`
	Name    string          `json:"name"`
	Value   json.RawMessage `json:"value"`
	Tags    []string        `json:"tags,omitempty"`
}

// VariableCreate represents the request to create a variable.
//
// Value must be valid JSON. Use VariableValue to build it from a Go value; note
// that a bare string has to be JSON-encoded (`"text"`, not `text`).
type VariableCreate struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
	Tags  []string        `json:"tags,omitempty"`
}

// VariableUpdate represents the request to update a variable.
//
// A nil Value is omitted from the request, leaving the stored value untouched.
// To set the value to JSON null explicitly, pass json.RawMessage("null").
type VariableUpdate struct {
	Value json.RawMessage `json:"value,omitempty"`
	Tags  []string        `json:"tags,omitempty"`
}

// VariableFilter represents filter criteria for querying variables.
type VariableFilter struct {
	Name   *string  `json:"name,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Offset int      `json:"offset,omitempty"`
	Limit  int      `json:"limit,omitempty"`
}

// VariableValue encodes an arbitrary Go value into the raw JSON expected by
// Variable.Value, VariableCreate.Value and VariableUpdate.Value.
//
//	v, err := models.VariableValue(map[string]any{"path": "/data", "limit": 500})
//	v, err := models.VariableValue("plain text")
func VariableValue(v any) (json.RawMessage, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

// StringValue returns the variable's value as a string when it is a JSON string,
// together with true. For any other JSON type - object, array, number, boolean,
// null or an absent value - it returns false, so callers can tell an actual
// string apart from a value that merely renders like one.
func (v Variable) StringValue() (string, bool) {
	raw := bytes.TrimSpace(v.Value)
	// A JSON null unmarshals into a string without error, so the type has to be
	// checked explicitly rather than inferred from Unmarshal succeeding.
	if len(raw) == 0 || raw[0] != '"' {
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
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
