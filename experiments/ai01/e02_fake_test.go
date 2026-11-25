package ai01

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/eino/schema"

	"campusly.local/ai01/internal/agentrun"
	"campusly.local/ai01/internal/fakemodel"
)

func TestE02RepeatedInputsAreScripted(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:scripted"
	model := fakemodel.New(fakemodel.Script{
		prompt: {
			{Content: "reply-1"},
			{Content: "reply-2"},
			{Content: "reply-3"},
		},
	})
	want := []string{"reply-1", "reply-2", "reply-3"}
	for i, expected := range want {
		got, err := model.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
		if err != nil {
			t.Fatalf("attempt %d: %v", i+1, err)
		}
		if got.Content != expected {
			t.Fatalf("attempt %d: got %q, want %q", i+1, got.Content, expected)
		}
	}
	if model.Attempts(prompt) != len(want) {
		t.Fatalf("attempts = %d, want %d", model.Attempts(prompt), len(want))
	}
}

func TestE02TransientRetryStaysWithinCap(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:retry-transient"
	model := fakemodel.New(fakemodel.Script{
		prompt: {
			{Err: fakemodel.ErrTransient},
			{Err: fakemodel.ErrTransient},
			{Content: "recovered"},
		},
	})
	agent, err := agentrun.NewAgent(ctx, model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	events := agentrun.Collect(agentrun.NewRunner(agent).Query(ctx, prompt))
	text, err := agentrun.Outcome(events)
	if err != nil {
		t.Fatalf("agent returned %v; events: %+v", err, events)
	}
	if text != "recovered" {
		t.Fatalf("final answer = %q, want %q", text, "recovered")
	}
	if model.Attempts(prompt) != agentrun.MaxRetries+1 {
		t.Fatalf("attempts = %d, want %d (cap %d plus first call)", model.Attempts(prompt), agentrun.MaxRetries+1, agentrun.MaxRetries)
	}
}

func TestE02PermanentErrorDoesNotRetry(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:permanent-failure"
	model := fakemodel.New(fakemodel.Script{
		prompt: {{Err: fakemodel.ErrPermanent}},
	})
	agent, err := agentrun.NewAgent(ctx, model, nil)
	if err != nil {
		t.Fatalf("create agent: %v", err)
	}
	events := agentrun.Collect(agentrun.NewRunner(agent).Query(ctx, prompt))
	_, err = agentrun.Outcome(events)
	if !errors.Is(err, fakemodel.ErrPermanent) {
		t.Fatalf("error = %v, want ErrPermanent; events: %+v", err, events)
	}
	if model.Attempts(prompt) != 1 {
		t.Fatalf("attempts = %d, want 1 for a permanent error", model.Attempts(prompt))
	}
}

func TestE02NoCredentialsNeeded(t *testing.T) {
	// The fake model is the only inference implementation in this module.
	// A missing key or proxy cannot affect these tests: no provider client is
	// constructed and no HTTP request is made.
	if _, err := (&fakemodel.Model{}).Generate(context.Background(), nil); err == nil {
		t.Fatal("expected a fake-model error for an empty script")
	}
}
