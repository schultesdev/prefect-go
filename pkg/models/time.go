package models

import (
	"encoding/json"
	"time"
)

// parseOptionalTime parses a time string, returning nil for empty strings.
// Supports RFC3339Nano and RFC3339 formats used by the Prefect API.
func parseOptionalTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// optTime handles JSON unmarshaling of nullable time fields that may be
// represented as empty strings or JSON null by the Prefect API.
type optTime struct {
	V *time.Time
}

func (o *optTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.V = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	v, err := parseOptionalTime(s)
	if err != nil {
		return err
	}
	o.V = v
	return nil
}
