package ai01

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"

	"campusly.local/ai01/internal/fakemodel"
)

func TestE01CompleteResponse(t *testing.T) {
	ctx := context.Background()
	const prompt = "synthetic:venue-question"
	const expected = "Room 12 — Intro to Python (north)."

	model := fakemodel.New(fakemodel.Script{
		prompt: {{Content: expected}},
	})

	response, err := model.Generate(ctx, []*schema.Message{schema.UserMessage(prompt)})
	if err != nil {
		t.Fatalf("Generate returned an error: %v", err)
	}
	if response == nil {
		t.Fatal("Generate returned a nil message")
	}
	if response.Content != expected {
		t.Fatalf("content mismatch: got %q, want %q", response.Content, expected)
	}
	if len(response.ToolCalls) != 0 {
		t.Fatalf("unexpected tool calls in a complete answer: %+v", response.ToolCalls)
	}
	if model.StreamCalls() != 0 {
		t.Fatalf("stream interface was called %d times; complete responses must not stream", model.StreamCalls())
	}
	if model.Calls() != 1 {
		t.Fatalf("model calls = %d, want 1", model.Calls())
	}
}
