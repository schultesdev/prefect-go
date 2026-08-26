package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// Prefect accepts any JSON type for a variable's value. An object must reach the
// API as an object, not as a JSON-encoded string.
func TestVariablesService_CreateSendsObjectValue(t *testing.T) {
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, want POST", r.Method)
		}
		if r.URL.Path != "/api/variables/" {
			t.Errorf("path = %v, want /api/variables/", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"id":"` + uuid.New().String() + `","name":"cfg","value":{"path":"/data"},"created":null,"updated":null}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))

	value, err := models.VariableValue(map[string]any{"path": "/data"})
	if err != nil {
		t.Fatalf("VariableValue: %v", err)
	}
	variable, err := client.Variables.Create(context.Background(), &models.VariableCreate{
		Name:  "cfg",
		Value: value,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var sent struct {
		Name  string          `json:"name"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(gotBody, &sent); err != nil {
		t.Fatalf("decode sent body %s: %v", gotBody, err)
	}
	if got := string(sent.Value); got != `{"path":"/data"}` {
		t.Errorf("sent value = %s, want %s", got, `{"path":"/data"}`)
	}
	if got := string(variable.Value); got != `{"path":"/data"}` {
		t.Errorf("decoded value = %s, want %s", got, `{"path":"/data"}`)
	}
}

// A plain string value must be sent JSON-encoded, i.e. quoted.
func TestVariablesService_CreateSendsStringValue(t *testing.T) {
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		gotBody = body
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"id":"` + uuid.New().String() + `","name":"p","value":"/data/imports","created":null,"updated":null}`)); err != nil {
			t.Errorf("write response: %v", err)
		}
	}))
	defer server.Close()

	client, _ := NewClient(WithBaseURL(server.URL + "/api"))

	value, err := models.VariableValue("/data/imports")
	if err != nil {
		t.Fatalf("VariableValue: %v", err)
	}
	variable, err := client.Variables.Create(context.Background(), &models.VariableCreate{Name: "p", Value: value})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var sent struct {
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(gotBody, &sent); err != nil {
		t.Fatalf("decode sent body %s: %v", gotBody, err)
	}
	if got := string(sent.Value); got != `"/data/imports"` {
		t.Errorf("sent value = %s, want %s", got, `"/data/imports"`)
	}
	if s, ok := variable.StringValue(); !ok || s != "/data/imports" {
		t.Errorf("StringValue() = %q, %v; want %q, true", s, ok, "/data/imports")
	}
}

// Every JSON type Prefect allows must decode without error. Before Value became
// json.RawMessage, the object and array cases failed with "cannot unmarshal
// object into Go struct field .Alias.value of type string".
func TestVariablesService_GetDecodesAllJSONTypes(t *testing.T) {
	cases := []struct {
		name      string
		value     string
		wantValue string
		wantStr   string
		wantIsStr bool
	}{
		{"object", `{"key":"value"}`, `{"key":"value"}`, "", false},
		{"array", `[1,2,3]`, `[1,2,3]`, "", false},
		{"nested", `{"a":{"b":[1,null]}}`, `{"a":{"b":[1,null]}}`, "", false},
		{"string", `"/data/imports"`, `"/data/imports"`, "/data/imports", true},
		{"number", `42`, `42`, "", false},
		{"large integer", `9007199254740993`, `9007199254740993`, "", false},
		{"bool", `true`, `true`, "", false},
		{"null", `null`, `null`, "", false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("method = %v, want GET", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				body := `{"id":"` + id.String() + `","name":"a","value":` + tc.value + `,"created":null,"updated":null}`
				if _, err := w.Write([]byte(body)); err != nil {
					t.Errorf("write response: %v", err)
				}
			}))
			defer server.Close()

			client, _ := NewClient(WithBaseURL(server.URL + "/api"))
			variable, err := client.Variables.Get(context.Background(), id)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := string(variable.Value); got != tc.wantValue {
				t.Errorf("value = %s, want %s", got, tc.wantValue)
			}
			s, ok := variable.StringValue()
			if ok != tc.wantIsStr || s != tc.wantStr {
				t.Errorf("StringValue() = %q, %v; want %q, %v", s, ok, tc.wantStr, tc.wantIsStr)
			}
		})
	}
}

// A nil Value must be omitted so the stored value stays untouched; an explicit
// JSON null must be sent.
func TestVariablesService_UpdateOmitsNilValue(t *testing.T) {
	cases := []struct {
		name         string
		req          models.VariableUpdate
		wantValueKey bool
		wantValue    string
	}{
		{"object value", models.VariableUpdate{Value: json.RawMessage(`{"a":1}`)}, true, `{"a":1}`},
		{"explicit null", models.VariableUpdate{Value: json.RawMessage(`null`)}, true, `null`},
		{"tags only", models.VariableUpdate{Tags: []string{"dts"}}, false, ""},
		{"empty request", models.VariableUpdate{}, false, ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()
			var gotBody []byte
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPatch {
					t.Errorf("method = %v, want PATCH", r.Method)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("read body: %v", err)
				}
				gotBody = body
				// Prefect answers PATCH /variables/{id} with 204 and no body.
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			client, _ := NewClient(WithBaseURL(server.URL + "/api"))
			if _, err := client.Variables.Update(context.Background(), id, &tc.req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var sent map[string]json.RawMessage
			if err := json.Unmarshal(gotBody, &sent); err != nil {
				t.Fatalf("decode sent body %s: %v", gotBody, err)
			}
			raw, ok := sent["value"]
			if ok != tc.wantValueKey {
				t.Fatalf("value key present = %v, want %v (body %s)", ok, tc.wantValueKey, gotBody)
			}
			if tc.wantValueKey && string(raw) != tc.wantValue {
				t.Errorf("value = %s, want %s", raw, tc.wantValue)
			}
		})
	}
}
