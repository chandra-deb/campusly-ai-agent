// Package draft exercises framework-derived schema against local validation.
// Eino can describe and generate a JSON schema; it does not enforce that a
// model's free-form output satisfies it. The application still validates.
package draft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/components/tool/utils"
)

type Draft struct {
	Title    string `json:"title" jsonschema:"required,description=Activity title"`
	Capacity int    `json:"capacity" jsonschema:"required,description=Bounded positive capacity"`
}

var (
	ErrInvalidDraft = errors.New("invalid activity draft")
	MaxCapacity     = 200
)

// SchemaJSON returns the Eino-generated JSON schema as plain data so a test
// can inspect what the framework contributes (types and required fields).
func SchemaJSON() (map[string]any, error) {
	params, err := utils.GoStruct2ParamsOneOf[Draft]()
	if err != nil {
		return nil, err
	}
	schema, err := params.ToJSONSchema()
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(encoded, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Validate is application-owned: strict decoding plus explicit content rules.
func Validate(raw string) (Draft, error) {
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	var d Draft
	if err := decoder.Decode(&d); err != nil {
		return Draft{}, fmt.Errorf("%w: %v", ErrInvalidDraft, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Draft{}, fmt.Errorf("%w: trailing JSON content", ErrInvalidDraft)
	}
	if strings.TrimSpace(d.Title) == "" {
		return Draft{}, fmt.Errorf("%w: empty title", ErrInvalidDraft)
	}
	if d.Capacity < 1 || d.Capacity > MaxCapacity {
		return Draft{}, fmt.Errorf("%w: capacity outside [1, %d]", ErrInvalidDraft, MaxCapacity)
	}
	return d, nil
}
