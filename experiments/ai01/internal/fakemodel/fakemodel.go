// Package fakemodel implements Eino's real model interfaces with scripted,
// in-process replies. It never opens a network connection or reads credentials.
package fakemodel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

var (
	ErrTransient     = errors.New("transient model failure")
	ErrPermanent     = errors.New("permanent model failure")
	ErrStreamOff     = errors.New("streaming is disabled in the AI01 fake")
	ErrNoScriptedKey = errors.New("no scripted reply for prompt")
)

type Step struct {
	Content   string
	ToolCalls []schema.ToolCall
	Err       error
	Delay     time.Duration
}

// Script maps the most recent user message to a fixed list of steps.
// A missing key is a test error, so an unexpected prompt cannot silently pass.
type Script map[string][]Step

type Model struct {
	mu          sync.Mutex
	script      Script
	attempts    map[string]int
	calls       int
	streamCalls int
	toolNames   []string
}

var _ model.ToolCallingChatModel = (*Model)(nil)

func New(script Script) *Model {
	return &Model{
		script:   script,
		attempts: make(map[string]int),
	}
}

func lastUserText(input []*schema.Message) (string, error) {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] != nil && input[i].Role == schema.User {
			return input[i].Content, nil
		}
	}
	return "", errors.New("fake model received no user message")
}

func (m *Model) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	key, err := lastUserText(input)
	if err != nil {
		return nil, err
	}
	common := model.GetCommonOptions(&model.Options{}, opts...)

	m.mu.Lock()
	m.calls++
	index := m.attempts[key]
	m.attempts[key] = index + 1
	m.toolNames = make([]string, 0, len(common.Tools))
	for _, info := range common.Tools {
		if info != nil {
			m.toolNames = append(m.toolNames, info.Name)
		}
	}
	steps, ok := m.script[key]
	m.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrNoScriptedKey, key)
	}
	step := steps[len(steps)-1]
	if index < len(steps) {
		step = steps[index]
	}
	if step.Delay > 0 {
		timer := time.NewTimer(step.Delay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if step.Err != nil {
		return nil, step.Err
	}
	return &schema.Message{
		Role:      schema.Assistant,
		Content:   step.Content,
		ToolCalls: step.ToolCalls,
	}, nil
}

// Stream exists only to satisfy the Eino interface. The spike must not use it:
// Campusly returns complete responses, so a stream call is a test failure.
func (m *Model) Stream(_ context.Context, _ []*schema.Message, _ ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	m.mu.Lock()
	m.streamCalls++
	m.mu.Unlock()
	return nil, ErrStreamOff
}

// WithTools returns an immutable derived model, matching the contract in the
// pinned v0.9.23 source. The ADK passes tools via model.WithOptions, but a
// hosted adapter would use this method too.
func (m *Model) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	derived := &Model{
		script:      m.script,
		attempts:    m.attempts,
		calls:       m.calls,
		streamCalls: m.streamCalls,
	}
	derived.toolNames = make([]string, 0, len(tools))
	for _, info := range tools {
		if info != nil {
			derived.toolNames = append(derived.toolNames, info.Name)
		}
	}
	return derived, nil
}

func (m *Model) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

func (m *Model) StreamCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.streamCalls
}

func (m *Model) Attempts(key string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.attempts[key]
}

func (m *Model) ToolNames() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.toolNames...)
}

func ToolCall(name, id, arguments string) schema.ToolCall {
	return schema.ToolCall{
		ID:   id,
		Type: "function",
		Function: schema.FunctionCall{
			Name:      name,
			Arguments: arguments,
		},
	}
}

func ToolCallStep(name, id, arguments string) Step {
	return Step{ToolCalls: []schema.ToolCall{ToolCall(name, id, arguments)}}
}
