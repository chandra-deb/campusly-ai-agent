// Package agentrun wires a scripted model into Eino's ADK agent and runner.
package agentrun

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"campusly.local/ai01/internal/fakemodel"
)

const (
	MaxIterations = 3
	MaxRetries    = 2
)

func NewAgent(ctx context.Context, m model.BaseChatModel, tools []tool.BaseTool) (*adk.ChatModelAgent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ai01-agent",
		Instruction: "Answer briefly from the supplied synthetic data. Use the read-only lookup tool when asked for a venue.",
		Model:       m,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{Tools: tools},
		},
		MaxIterations: MaxIterations,
		ModelRetryConfig: &adk.ModelRetryConfig{
			MaxRetries: MaxRetries,
			IsRetryAble: func(_ context.Context, err error) bool {
				return errors.Is(err, fakemodel.ErrTransient)
			},
			BackoffFunc: func(context.Context, int) time.Duration { return 0 },
		},
	})
}

func NewRunner(agent *adk.ChatModelAgent) *adk.Runner {
	return adk.NewRunner(context.Background(), adk.RunnerConfig{Agent: agent})
}

func Collect(iter *adk.AsyncIterator[*adk.AgentEvent]) []*adk.AgentEvent {
	var events []*adk.AgentEvent
	for {
		event, ok := iter.Next()
		if !ok {
			return events
		}
		events = append(events, event)
	}
}

// Outcome returns the last complete assistant answer and the first real
// error. Eino emits WillRetryError events during a retry; those are progress
// signals, not the final result, and are ignored when a final answer exists.
func Outcome(events []*adk.AgentEvent) (string, error) {
	var text string
	var finalErr error
	for _, event := range events {
		if event == nil {
			continue
		}
		if event.Err != nil {
			var retrying *adk.WillRetryError
			if !errors.As(event.Err, &retrying) && finalErr == nil {
				finalErr = event.Err
			}
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		message, err := event.Output.MessageOutput.GetMessage()
		if err != nil {
			if finalErr == nil {
				finalErr = err
			}
			continue
		}
		if event.Output.MessageOutput.Role == schema.Assistant && len(message.ToolCalls) == 0 {
			text = message.Content
		}
	}
	return text, finalErr
}
