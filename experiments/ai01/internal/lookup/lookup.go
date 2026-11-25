// Package lookup implements the one read-only tool used in AI01.
package lookup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"campusly.local/ai01/internal/fixture"
)

type Input struct {
	ActivityID   string `json:"activity_id"`
	Organization string `json:"organization"`
	Actor        string `json:"actor"`
}

type Tool struct {
	repo *fixture.Repository
}

var _ tool.InvokableTool = (*Tool)(nil)

func New(repo *fixture.Repository) *Tool {
	return &Tool{repo: repo}
}

func (t *Tool) Info(context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "lookup_activity",
		Desc: "Look up a synthetic activity's title and venue for an authorized organizer. Read-only.",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"activity_id":  {Type: schema.String, Desc: "Synthetic activity id", Required: true},
			"organization": {Type: schema.String, Desc: "Organization scope of the request", Required: true},
			"actor":        {Type: schema.String, Desc: "Actor requesting the read", Required: true},
		}),
	}, nil
}

func (t *Tool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...tool.Option) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	decoder := json.NewDecoder(strings.NewReader(argumentsInJSON))
	decoder.DisallowUnknownFields()
	var in Input
	if err := decoder.Decode(&in); err != nil {
		return "", fmt.Errorf("%w: %v", fixture.ErrInvalidArguments, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("%w: trailing JSON content", fixture.ErrInvalidArguments)
	}
	activity, err := t.repo.Lookup(in.ActivityID, in.Organization, in.Actor)
	if err != nil {
		return "", err
	}
	out, err := json.Marshal(activity)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
