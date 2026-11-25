package personal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignedIntent(t *testing.T) {
	handler := Handler(testKey, Fake{})
	input := validRequest()
	input.Version = 3
	input.Prompt = "Which organizations do I manage?"
	payload, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	send := func(signature string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, IntentPath, bytes.NewReader(payload))
		request.Header.Set(Signature, signature)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		return response
	}
	response := send(Sign(testKey, payload))
	var result IntentReply
	if response.Code != 200 || json.Unmarshal(response.Body.Bytes(), &result) != nil ||
		result.Intent != "managed_organizations" || result.Scope != input.Scope || result.Version != 3 {
		t.Fatalf("unexpected result: %d %s", response.Code, response.Body)
	}
	if invalid := send("wrong"); invalid.Code != 401 {
		t.Fatalf("invalid signature returned %d", invalid.Code)
	}
}

func TestActivityDetailsIntent(t *testing.T) {
	for _, tc := range []struct {
		prompt   string
		intent   string
		activity string
	}{
		{"When and where is the next workshop?", "activity_details", ""},
		{"Show details for Sample Workshop", "activity_details", "Sample Workshop"},
		{"How many registered for Sample Workshop?", "activity_statistics", "Sample Workshop"},
		{"How many people registered for the next workshop?", "activity_statistics", ""},
		{"Give me a workshop checklist", "chat", ""},
	} {
		input := validRequest()
		input.Prompt = tc.prompt
		result := resolveIntent(input)
		if result.Intent != tc.intent || result.Activity != tc.activity {
			t.Errorf("%q yielded %+v", tc.prompt, result)
		}
	}
}

func TestActivityDraftProposal(t *testing.T) {
	input := validRequest()
	input.Prompt = "Draft an activity: Intro to Python | Practice loops | 30"
	result := resolveIntent(input)
	if result.Intent != "activity_draft" || result.Draft == nil ||
		*result.Draft != (ActivityDraft{"activity", "Intro to Python", "Practice loops", "WORKSHOP", "30"}) {
		t.Fatalf("unexpected activity proposal: %+v", result.Draft)
	}
	input.Prompt = "Draft a workshop"
	if draft := resolveIntent(input).Draft; draft == nil || draft.Title != "" || draft.Capacity != "" {
		t.Fatalf("missing fields should remain editable: %+v", draft)
	}
	input.Prompt = "Draft an activity called Intro to Python for 30 people"
	if draft := resolveIntent(input).Draft; draft == nil || draft.Title != "Intro to Python" || draft.Capacity != "30" {
		t.Fatalf("unexpected natural-language proposal: %+v", draft)
	}
}
