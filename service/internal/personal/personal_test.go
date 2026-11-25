package personal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

var testKey = bytes.Repeat([]byte{'t'}, 32)

func validRequest() Request {
	return Request{
		Version: 1, Scope: Scope{11, 21}, Expires: time.Now().Add(time.Second).UnixMilli(), Prompt: "Workshop checklist",
	}
}

func send(t *testing.T, handler http.Handler, input Request, signBefore func(*Request)) *httptest.ResponseRecorder {
	t.Helper()
	return sendTo(t, handler, Path, input, signBefore)
}

func sendTo(t *testing.T, handler http.Handler, path string, input Request, signBefore func(*Request)) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	signature := Sign(testKey, payload)
	if signBefore != nil {
		signBefore(&input)
		payload, err = json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
	}
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	r.Header.Set(Signature, signature)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func TestCompleteReply(t *testing.T) {
	input := validRequest()
	w := send(t, Handler(testKey, Fake{}), input, nil)
	var reply Reply
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &reply) != nil ||
		reply.Scope != input.Scope || reply.Version != 1 || !strings.Contains(reply.Text, "4. Draft") ||
		w.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body)
	}
}

func TestRejectsTamperingExpiryAndWrongVersion(t *testing.T) {
	handler := Handler(testKey, Fake{})
	if w := send(t, handler, validRequest(), func(r *Request) { r.Actor++ }); w.Code != 401 {
		t.Fatalf("tampered actor returned %d", w.Code)
	}
	if w := send(t, handler, validRequest(), func(r *Request) { r.Organization++ }); w.Code != 401 {
		t.Fatalf("tampered organization returned %d", w.Code)
	}
	expired := validRequest()
	expired.Expires = time.Now().Add(-time.Second).UnixMilli()
	if w := send(t, handler, expired, nil); w.Code != 401 {
		t.Fatalf("expired scope returned %d", w.Code)
	}
	future := validRequest()
	future.Version = 2
	if w := send(t, handler, future, nil); w.Code != 400 {
		t.Fatalf("unknown version returned %d", w.Code)
	}
}

type testModel struct {
	Fake
	generate func(context.Context) (*schema.Message, error)
}

func (m testModel) Generate(ctx context.Context, _ []*schema.Message, _ ...model.Option) (*schema.Message, error) {
	return m.generate(ctx)
}

func TestDeadlineAndMalformedOutput(t *testing.T) {
	handler := Handler(testKey, testModel{generate: func(ctx context.Context) (*schema.Message, error) {
		<-ctx.Done()
		return schema.AssistantMessage("partial text", nil), ctx.Err()
	}})
	input := validRequest()
	input.Expires = time.Now().Add(50 * time.Millisecond).UnixMilli()
	w := send(t, handler, input, nil)
	if w.Code != 504 || strings.Contains(w.Body.String(), "partial") {
		t.Fatalf("deadline leaked output: %d %s", w.Code, w.Body)
	}
	bad := Handler(testKey, testModel{generate: func(context.Context) (*schema.Message, error) {
		return schema.UserMessage("malformed"), nil
	}})
	w = send(t, bad, validRequest(), nil)
	if w.Code != 502 || strings.Contains(w.Body.String(), "malformed") {
		t.Fatalf("malformed result leaked: %d %s", w.Code, w.Body)
	}
}

func TestContextIsBoundedAndUsesSignedConversation(t *testing.T) {
	handler := Handler(testKey, Fake{})
	request := validRequest()
	request.Version = 2
	request.Prompt = "Explain the second step"
	request.History = []Turn{{"Workshop checklist", "Workshop checklist: host and check-in lead"}}

	w := sendTo(t, handler, ContextPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "assign a host") {
		t.Fatalf("history was not used: %d %s", w.Code, w.Body)
	}
	request.History = nil
	w = sendTo(t, handler, ContextPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "I need a workshop checklist") {
		t.Fatalf("context leaked between calls: %d %s", w.Code, w.Body)
	}
	request.History = make([]Turn, maxHistory+1)
	if w = sendTo(t, handler, ContextPath, request, nil); w.Code != 400 {
		t.Fatalf("too many turns: %d", w.Code)
	}
}

func TestRecentContextUsesWholeTurnsWithinByteBudget(t *testing.T) {
	handler := Handler(testKey, Fake{})
	request := validRequest()
	request.Version, request.Prompt = 4, "Explain the second step"
	request.History = []Turn{{"Workshop checklist", "Workshop checklist: host and check-in lead"}}
	for range maxRecentTurns - 1 {
		request.History = append(request.History, Turn{"General planning question", "General reply"})
	}
	w := sendTo(t, handler, RecentPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "assign a host") {
		t.Fatalf("window lost an eligible older turn: %d %s", w.Code, w.Body)
	}

	request.History = append(request.History, Turn{"One more question", "One more reply"})
	if w = sendTo(t, handler, RecentPath, request, nil); w.Code != 400 {
		t.Fatalf("too many recent turns returned %d", w.Code)
	}

	wide := Turn{strings.Repeat("界", 666), strings.Repeat("答", 1333)}
	request.History = []Turn{wide, wide}
	if w = sendTo(t, handler, RecentPath, request, nil); w.Code != 400 {
		t.Fatalf("over-budget UTF-8 text returned %d", w.Code)
	}

	request.History = nil
	w = sendTo(t, handler, RecentPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "I need a workshop checklist") {
		t.Fatalf("empty context reused a previous request: %d %s", w.Code, w.Body)
	}
}

func TestPreferredReplyUsesOnlySupportedChoices(t *testing.T) {
	handler := Handler(testKey, Fake{})
	request := validRequest()
	request.Version = 5
	request.Preferences = []Preference{{"format", "bullets"}, {"tone", "warm"}, {"detail", "concise"}}
	w := sendTo(t, handler, PreferredPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "• Confirm") ||
		!strings.Contains(w.Body.String(), "work through together") || strings.Contains(w.Body.String(), "3. Test") ||
		strings.Contains(w.Body.String(), "• Test") {
		t.Fatalf("preference reply: %d %s", w.Code, w.Body)
	}

	if w = sendTo(t, handler, PreferredPath, request, func(r *Request) { r.Preferences[0].Value = "numbered" }); w.Code != 401 {
		t.Fatalf("tampered preference returned %d", w.Code)
	}
	for _, choices := range [][]Preference{
		nil,
		{{"tone", "student details"}},
		{{"format", "bullets"}, {"format", "numbered"}},
	} {
		request.Preferences = choices
		if w = sendTo(t, handler, PreferredPath, request, nil); w.Code != 400 {
			t.Fatalf("invalid preferences returned %d", w.Code)
		}
	}
	request.Version = 4
	request.Preferences = []Preference{{"format", "bullets"}}
	if w = sendTo(t, handler, RecentPath, request, nil); w.Code != 400 {
		t.Fatalf("obsolete version accepted new fields: %d", w.Code)
	}
	request.Preferences = nil
	w = sendTo(t, handler, RecentPath, request, nil)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "1. Confirm") {
		t.Fatalf("preference leaked into another reply: %d %s", w.Code, w.Body)
	}
}
