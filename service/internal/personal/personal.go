package personal

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	Path        = "/internal/v1/personal/reply"
	Signature   = "X-Campusly-Signature"
	maxBody     = 16 << 10
	maxPrompt   = 2000
	maxReply    = 4000
	maxLifetime = 5 * time.Second
)

type Scope struct {
	Actor        int64 `json:"actor_id"`
	Organization int64 `json:"org_id"`
}

type Request struct {
	Version int `json:"version"`
	Scope
	Expires int64  `json:"expires_ms"`
	Prompt  string `json:"prompt"`
}

type Reply struct {
	Version int `json:"version"`
	Scope
	Text string `json:"reply"`
}

type service struct {
	key    []byte
	model  model.BaseChatModel
	slots  chan struct{}
	budget time.Duration
}

func Sign(key, payload []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func Handler(key []byte, model model.BaseChatModel) http.Handler {
	s := &service{bytes.Clone(key), model, make(chan struct{}, 8), 3 * time.Second}
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+Path, s.reply)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

func (s *service) reply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBody))
	if err != nil {
		fail(w, http.StatusBadRequest, "invalid_request")
		return
	}
	if len(s.key) < 32 || !hmac.Equal([]byte(r.Header.Get(Signature)), []byte(Sign(s.key, body))) {
		fail(w, http.StatusUnauthorized, "invalid_signature")
		return
	}
	var input Request
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	if !utf8.Valid(body) || dec.Decode(&input) != nil || dec.Decode(new(any)) != io.EOF ||
		input.Version != 1 || input.Actor <= 0 || input.Organization <= 0 ||
		strings.TrimSpace(input.Prompt) == "" || len(input.Prompt) > maxPrompt {
		fail(w, http.StatusBadRequest, "invalid_request")
		return
	}
	now := time.Now()
	expires := time.UnixMilli(input.Expires)
	if !expires.After(now) || expires.After(now.Add(maxLifetime)) {
		fail(w, http.StatusUnauthorized, "expired_scope")
		return
	}
	select {
	case s.slots <- struct{}{}:
		defer func() { <-s.slots }()
	default:
		fail(w, http.StatusServiceUnavailable, "busy")
		return
	}
	ctx, cancel := context.WithDeadline(r.Context(), expires)
	defer cancel()
	ctx, cancelBudget := context.WithTimeout(ctx, s.budget)
	defer cancelBudget()
	result, err := s.model.Generate(ctx, []*schema.Message{schema.UserMessage(input.Prompt)})
	if ctx.Err() != nil {
		fail(w, http.StatusGatewayTimeout, "timeout")
		return
	}
	if err != nil || result == nil || result.Role != schema.Assistant || len(result.ToolCalls) != 0 ||
		strings.TrimSpace(result.Content) == "" || len(result.Content) > maxReply || !utf8.ValidString(result.Content) {
		fail(w, http.StatusBadGateway, "generation_failed")
		return
	}
	write(w, http.StatusOK, Reply{1, input.Scope, result.Content})
}

func write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func fail(w http.ResponseWriter, status int, code string) {
	write(w, status, map[string]string{"error": code})
}

type Fake struct{}

func (Fake) Generate(ctx context.Context, messages []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	prompt := strings.ToLower(messages[len(messages)-1].Content)
	text := "This local demo can show how a complete reply arrives. Try asking for a workshop checklist."
	if strings.Contains(prompt, "workshop") {
		text = "Workshop checklist\n1. Confirm the room, time and capacity.\n2. Assign a host and a check-in lead.\n3. Test the projector and materials.\n4. Draft a reminder for participants."
	}
	return schema.AssistantMessage(text, nil), nil
}

func (Fake) Stream(context.Context, []*schema.Message, ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("streaming is unavailable")
}
