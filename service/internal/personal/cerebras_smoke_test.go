package personal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cloudwego/eino/schema"
)

type mockTransport func(*http.Request) (*http.Response, error)

func (mock mockTransport) RoundTrip(r *http.Request) (*http.Response, error) { return mock(r) }

func TestCerebrasSmokeOnlySendsTheApprovedPromptOnce(t *testing.T) {
	smoke := NewCerebrasSmoke("mock-key")
	calls := 0
	smoke.client.Transport = mockTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.String() != cerebrasCompletions ||
			r.Header.Get("Authorization") != "Bearer mock-key" || r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("request did not use the documented chat endpoint")
		}
		var input completionRequest
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Fatal(err)
		}
		if input.Stream || input.Model != "gpt-oss-120b" || input.ReasoningEffort != "low" ||
			input.MaxTokens != 256 || len(input.Messages) != 1 ||
			input.Messages[0].Role != "user" || input.Messages[0].Content != CerebrasSmokePrompt {
			t.Fatal("request was outside the approved synthetic chat scope")
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(
			`{"choices":[{"finish_reason":"stop","message":{"role":"assistant","content":"1. Book a room."}}]}`,
		))}, nil
	})
	ctx := context.Background()
	if _, err := smoke.Generate(ctx, []*schema.Message{schema.UserMessage("Unapproved input")}); err == nil || calls != 0 {
		t.Fatal("unapproved input reached the provider")
	}
	input := []*schema.Message{schema.UserMessage(CerebrasSmokePrompt)}
	withHistory := append([]*schema.Message{schema.AssistantMessage("Private saved history", nil)}, input...)
	if _, err := smoke.Generate(ctx, withHistory); err == nil || calls != 0 {
		t.Fatal("conversation history reached the provider")
	}
	reply, err := smoke.Generate(ctx, input)
	if err != nil || reply.Content != "1. Book a room." || calls != 1 {
		t.Fatalf("expected one complete synthetic reply, got error %v with %d calls", err, calls)
	}
	if _, err = smoke.Generate(ctx, input); err == nil || calls != 1 {
		t.Fatal("second attempt reached the provider")
	}
}

func TestCerebrasSmokeRejectsTruncatedReplyAndProviderError(t *testing.T) {
	for _, response := range []struct {
		name   string
		status int
		body   string
	}{
		{"truncated", http.StatusOK, `{"choices":[{"finish_reason":"length","message":{"role":"assistant","content":"Partial"}}]}`},
		{"unauthorized", http.StatusUnauthorized, `{"error":"private diagnostic text"}`},
	} {
		t.Run(response.name, func(t *testing.T) {
			smoke := NewCerebrasSmoke("mock-key")
			smoke.client.Transport = mockTransport(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: response.status, Body: io.NopCloser(strings.NewReader(response.body))}, nil
			})
			result, err := smoke.Generate(context.Background(), []*schema.Message{schema.UserMessage(CerebrasSmokePrompt)})
			if result != nil || err == nil || strings.Contains(err.Error(), "private diagnostic text") {
				t.Fatal("provider failure exposed output or diagnostic text")
			}
		})
	}
}
