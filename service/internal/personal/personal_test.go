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
	return Request{1, Scope{11, 21}, time.Now().Add(time.Second).UnixMilli(), "Workshop checklist", nil}
}

func send(t *testing.T, handler http.Handler, input Request, signBefore func(*Request)) *httptest.ResponseRecorder {
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
	r := httptest.NewRequest(http.MethodPost, Path, bytes.NewReader(payload))
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

	sendContext := func(input Request) *httptest.ResponseRecorder {
		payload, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		r := httptest.NewRequest(http.MethodPost, ContextPath, bytes.NewReader(payload))
		r.Header.Set(Signature, Sign(testKey, payload))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	w := sendContext(request)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "assign a host") {
		t.Fatalf("history was not used: %d %s", w.Code, w.Body)
	}
	request.History = nil
	w = sendContext(request)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "I need a workshop checklist") {
		t.Fatalf("context leaked between calls: %d %s", w.Code, w.Body)
	}
	request.History = make([]Turn, maxHistory+1)
	if w = sendContext(request); w.Code != 400 {
		t.Fatalf("too many turns: %d", w.Code)
	}
}
