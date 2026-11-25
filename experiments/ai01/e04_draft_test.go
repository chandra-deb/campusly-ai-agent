package ai01

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"campusly.local/ai01/internal/draft"
)

func TestE04FrameworkSchemaDescribesDraftOnly(t *testing.T) {
	schemaJSON, err := draft.SchemaJSON()
	if err != nil {
		t.Fatalf("SchemaJSON: %v", err)
	}
	encoded, err := json.Marshal(schemaJSON)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	var parsed struct {
		Required []string       `json:"required"`
		Props    map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	for _, field := range []string{"title", "capacity"} {
		found := false
		for _, required := range parsed.Required {
			if required == field {
				found = true
			}
		}
		if !found {
			t.Errorf("%q is not required in %s", field, encoded)
		}
		if _, ok := parsed.Props[field]; !ok {
			t.Errorf("%q is absent from the properties in %s", field, encoded)
		}
	}
	if _, ok := parsed.Props["action"]; ok {
		t.Fatalf("framework schema unexpectedly exposes an action field: %s", encoded)
	}
}

func TestE04ValidDraft(t *testing.T) {
	got, err := draft.Validate(`{"title":"Intro to Python","capacity":60}`)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if got.Title != "Intro to Python" || got.Capacity != 60 {
		t.Fatalf("draft = %+v", got)
	}
}

func TestE04RejectsMalformedDrafts(t *testing.T) {
	cases := []struct {
		name string
		raw  string
	}{
		{"malformed JSON", `{"title":`},
		{"missing title", `{"capacity":20}`},
		{"missing capacity", `{"title":"Workshop"}`},
		{"wrong capacity type", `{"title":"Workshop","capacity":"60"}`},
		{"zero capacity", `{"title":"Workshop","capacity":0}`},
		{"negative capacity", `{"title":"Workshop","capacity":-4}`},
		{"capacity above bound", `{"title":"Workshop","capacity":201}`},
		{"blank title", `{"title":"   ","capacity":20}`},
		{"unexpected action field", `{"title":"Workshop","capacity":20,"action":"publish"}`},
		{"trailing JSON", `{"title":"Workshop","capacity":20}{"x":1}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := draft.Validate(tc.raw)
			if !errors.Is(err, draft.ErrInvalidDraft) {
				t.Fatalf("error = %v, want ErrInvalidDraft; draft = %+v", err, got)
			}
		})
	}
}

func TestE04DraftNeverWritesAnything(t *testing.T) {
	// The validation result is data only: no create/publish call exists in this
	// module. The strict decoder also rejects any write instruction field.
	_, err := draft.Validate(`{"title":"Workshop","capacity":20,"action":"create_activity"}`)
	if err == nil || !strings.Contains(err.Error(), "action") {
		t.Fatalf("expected the action field to be rejected, got %v", err)
	}
}
