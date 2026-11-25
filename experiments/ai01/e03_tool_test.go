package ai01

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"

	"campusly.local/ai01/internal/agentrun"
	"campusly.local/ai01/internal/fakemodel"
	"campusly.local/ai01/internal/fixture"
	"campusly.local/ai01/internal/lookup"
)

func TestE03ToolArgumentSchema(t *testing.T) {
	info, err := lookup.New(fixture.New()).Info(context.Background())
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if info.Name != "lookup_activity" {
		t.Fatalf("tool name = %q", info.Name)
	}
	schemaJSON, err := info.ParamsOneOf.ToJSONSchema()
	if err != nil {
		t.Fatalf("ToJSONSchema: %v", err)
	}
	encoded, err := json.Marshal(schemaJSON)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	var parsed struct {
		Required []string                  `json:"required"`
		Props    map[string]map[string]any `json:"properties"`
	}
	if err := json.Unmarshal(encoded, &parsed); err != nil {
		t.Fatalf("parse schema: %v", err)
	}
	for _, field := range []string{"activity_id", "organization", "actor"} {
		found := false
		for _, required := range parsed.Required {
			if required == field {
				found = true
			}
		}
		if !found {
			t.Errorf("%q is not required in %s", field, encoded)
		}
		prop, ok := parsed.Props[field]
		if !ok || prop["type"] != "string" {
			t.Errorf("%q schema = %v, want a string property", field, prop)
		}
	}
}

func TestE03AgentUsesToolAndCorrelatesCall(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:venue-tool"
	const callID = "call-lookup-1"
	toolCall := fakemodel.ToolCallStep("lookup_activity", callID,
		`{"activity_id":"north-workshop","organization":"north","actor":"mira"}`)
	model := fakemodel.New(fakemodel.Script{
		prompt: {toolCall, {Content: "Room 12 — Intro to Python (north)."}},
	})
	repo := fixture.New()
	agent, err := agentrun.NewAgent(ctx, model, []tool.BaseTool{lookup.New(repo)})
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	events := agentrun.Collect(agentrun.NewRunner(agent).Query(ctx, prompt))
	t.Logf("agent events: %d", len(events))
	logEvents(t, events)
	text, err := agentrun.Outcome(events)
	if err != nil {
		logEvents(t, events)
		t.Fatalf("agent error: %v", err)
	}
	if !strings.Contains(text, "Room 12") {
		logEvents(t, events)
		t.Fatalf("final answer = %q, want the synthetic venue", text)
	}

	var sawCall, sawResult bool
	for _, event := range events {
		if event == nil || event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		variant := event.Output.MessageOutput
		message, err := variant.GetMessage()
		if err != nil {
			continue
		}
		if variant.Role == schema.Assistant {
			for _, call := range message.ToolCalls {
				if call.ID == callID && call.Function.Name == "lookup_activity" {
					sawCall = true
				}
			}
		}
		if variant.Role == schema.Tool && variant.ToolName == "lookup_activity" && message.ToolCallID == callID {
			sawResult = true
			if !strings.Contains(message.Content, "Room 12") {
				t.Fatalf("tool result = %q, want the synthetic venue", message.Content)
			}
		}
	}
	if !sawCall || !sawResult {
		logEvents(t, events)
		t.Fatalf("correlation: saw call=%v saw result=%v", sawCall, sawResult)
	}
	for _, name := range model.ToolNames() {
		if name == "lookup_activity" {
			return
		}
	}
	t.Fatalf("model was not given the tool schema; names=%v", model.ToolNames())
}

func TestE03RejectsBadArgumentsAndScopes(t *testing.T) {
	toolInstance := lookup.New(fixture.New())
	cases := []struct {
		name       string
		arguments  string
		wantErr    error
		wantInText string
	}{
		{"malformed JSON", `{"activity_id":`, fixture.ErrInvalidArguments, ""},
		{"missing actor", `{"activity_id":"north-workshop","organization":"north"}`, fixture.ErrInvalidArguments, ""},
		{"unknown field", `{"activity_id":"north-workshop","organization":"north","actor":"mira","action":"write"}`, fixture.ErrInvalidArguments, ""},
		{"wrong type", `{"activity_id":"north-workshop","organization":"north","actor":7}`, fixture.ErrInvalidArguments, ""},
		{"unknown id", `{"activity_id":"missing","organization":"north","actor":"mira"}`, fixture.ErrNotFound, ""},
		{"malformed id", `{"activity_id":"north workshop!","organization":"north","actor":"mira"}`, fixture.ErrInvalidArguments, ""},
		{"wrong tenant", `{"activity_id":"south-workshop","organization":"north","actor":"mira"}`, fixture.ErrNotFound, "Room 9"},
		{"unauthorized actor", `{"activity_id":"north-workshop","organization":"north","actor":"lee"}`, fixture.ErrDenied, ""},
		{"unauthorized org", `{"activity_id":"north-workshop","organization":"south","actor":"mira"}`, fixture.ErrDenied, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := toolInstance.InvokableRun(context.Background(), tc.arguments)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantInText != "" && strings.Contains(result, tc.wantInText) {
				t.Fatalf("denied lookup disclosed %q: %s", tc.wantInText, result)
			}
			if strings.Contains(result, "Room") && tc.wantErr != nil {
				t.Fatalf("denied lookup returned a venue: %s", result)
			}
		})
	}
}

func TestE03RepeatedToolCallsStopAtIterationCap(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:endless-tool"
	model := fakemodel.New(fakemodel.Script{
		prompt: {fakemodel.ToolCallStep("lookup_activity", "call-loop",
			`{"activity_id":"north-workshop","organization":"north","actor":"mira"}`)},
	})
	agent, err := agentrun.NewAgent(ctx, model, []tool.BaseTool{lookup.New(fixture.New())})
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	events := agentrun.Collect(agentrun.NewRunner(agent).Query(ctx, prompt))
	if calls := model.Attempts(prompt); calls > agentrun.MaxIterations {
		t.Fatalf("model attempts = %d, above the %d iteration cap", calls, agentrun.MaxIterations)
	} else if calls == 0 {
		t.Fatal("model was never called")
	}
	_, err = agentrun.Outcome(events)
	t.Logf("bounded loop stopped after %d model attempts: %v", model.Attempts(prompt), err)
	if err == nil {
		t.Fatalf("expected the bounded loop to fail; events: %+v", events)
	}
}

func logEvents(t *testing.T, events []*adk.AgentEvent) {
	t.Helper()
	for i, event := range events {
		if event == nil {
			t.Logf("event %d: nil", i)
			continue
		}
		role := ""
		content := ""
		if event.Output != nil && event.Output.MessageOutput != nil {
			role = string(event.Output.MessageOutput.Role)
			if message, err := event.Output.MessageOutput.GetMessage(); err == nil && message != nil {
				content = message.Content
			}
		}
		t.Logf("event %d: role=%s output=%+v content=%q err=%v", i, role, event.Output, content, event.Err)
	}
}
