package personal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	cerebrasCompletions = "https://api.cerebras.ai/v1/chat/completions"
	CerebrasSmokePrompt = "Give three short numbered checklist items for a fictional campus workshop. Start each with its number."
	cerebrasTokenCap    = 256
	CerebrasSmokeModel  = "gpt-oss-120b"
	cerebrasBodyCap     = 64 << 10
)

type CerebrasSmoke struct {
	key       string
	client    *http.Client
	attempted atomic.Bool
}

func NewCerebrasSmoke(key string) *CerebrasSmoke {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	return &CerebrasSmoke{
		key: key,
		client: &http.Client{
			Transport:     transport,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

type completionMessage struct {
	Role      string            `json:"role"`
	Content   string            `json:"content"`
	ToolCalls []json.RawMessage `json:"tool_calls,omitempty"`
}

type completionRequest struct {
	Model           string              `json:"model"`
	Messages        []completionMessage `json:"messages"`
	Stream          bool                `json:"stream"`
	MaxTokens       int                 `json:"max_completion_tokens"`
	ReasoningEffort string              `json:"reasoning_effort"`
}

type completionResponse struct {
	Choices []struct {
		Finish  string            `json:"finish_reason"`
		Message completionMessage `json:"message"`
	} `json:"choices"`
	Usage struct {
		Prompt     int `json:"prompt_tokens"`
		Completion int `json:"completion_tokens"`
	} `json:"usage"`
}

func (c *CerebrasSmoke) Generate(ctx context.Context, messages []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	if len(messages) != 1 || messages[0] == nil || messages[0].Role != schema.User ||
		messages[0].Content != CerebrasSmokePrompt || len(messages[0].ToolCalls) != 0 {
		return nil, errors.New("Cerebras smoke requires the fixed synthetic prompt in an empty conversation")
	}

	payload := completionRequest{
		Model: CerebrasSmokeModel, Messages: []completionMessage{{Role: "user", Content: CerebrasSmokePrompt}},
		MaxTokens: cerebrasTokenCap, ReasoningEffort: "low",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cerebrasCompletions, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if !c.attempted.CompareAndSwap(false, true) {
		return nil, errors.New("Cerebras smoke attempt already used")
	}
	start := time.Now()
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Cerebras smoke transport failed after %d ms", time.Since(start).Milliseconds())
		return nil, errors.New("Cerebras transport failed")
	}
	defer resp.Body.Close()
	log.Printf("Cerebras smoke returned HTTP %d in %d ms", resp.StatusCode, time.Since(start).Milliseconds())
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cerebras HTTP %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, cerebrasBodyCap+1))
	if err != nil || len(data) > cerebrasBodyCap {
		return nil, errors.New("Cerebras response exceeded the read limit or could not be read")
	}
	var reply completionResponse
	if err := json.Unmarshal(data, &reply); err != nil || len(reply.Choices) != 1 {
		return nil, errors.New("invalid Cerebras completion")
	}
	choice := reply.Choices[0]
	if choice.Finish != "stop" || choice.Message.Role != "assistant" ||
		len(choice.Message.ToolCalls) != 0 || !validText(choice.Message.Content, maxReply) {
		log.Print("Cerebras smoke rejected an incomplete or unsupported completion")
		return nil, errors.New("Cerebras did not return a complete assistant reply")
	}
	log.Printf("Cerebras smoke completed with %d input tokens and %d completion tokens",
		reply.Usage.Prompt, reply.Usage.Completion)
	return schema.AssistantMessage(strings.TrimSpace(choice.Message.Content), nil), nil
}

func (*CerebrasSmoke) Stream(context.Context, []*schema.Message, ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("streaming is unavailable")
}

var _ model.BaseChatModel = (*CerebrasSmoke)(nil)
