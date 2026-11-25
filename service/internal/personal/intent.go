package personal

import "strings"

// IntentReply proposes a read or draft operation. Core authorizes it and
// supplies the data. The resolver has no Campusly credentials or tools.
type IntentReply struct {
	Version int `json:"version"`
	Scope
	Intent   string         `json:"intent"`
	Activity string         `json:"activity,omitempty"`
	Draft    *ActivityDraft `json:"draft,omitempty"`
	Message  string         `json:"message,omitempty"`
}

func resolveIntent(input Request) IntentReply {
	reply := IntentReply{Version: input.Version, Scope: input.Scope, Intent: "chat"}
	if draft := proposeActivity(input.Prompt); draft != nil {
		reply.Intent, reply.Draft = "activity_draft", draft
		return reply
	}
	prompt := strings.TrimSpace(input.Prompt)
	for _, prefix := range []string{"draft a message for ", "send a group message to "} {
		if strings.HasPrefix(strings.ToLower(prompt), prefix) {
			target, body, _ := strings.Cut(prompt[len(prefix):], ":")
			if target = strings.TrimSpace(target); target != "" {
				reply.Intent, reply.Activity, reply.Message = "message_draft", target, strings.TrimSpace(body)
				return reply
			}
		}
	}
	prompt = strings.TrimSpace(strings.TrimRight(prompt, "?？"))
	lower := strings.ToLower(prompt)
	switch strings.TrimRight(lower, ".!") {
	case "which organizations do i manage", "what organizations do i manage",
		"show my organizations", "list the organizations i manage":
		reply.Intent = "managed_organizations"
	case "when and where is the next workshop", "where is the next workshop",
		"when is the next workshop", "show the next workshop", "what is the next workshop":
		reply.Intent = "activity_details"
	default:
		for _, operation := range []struct{ prefix, intent string }{
			{"show details for ", "activity_details"},
			{"how many registered for ", "activity_statistics"},
			{"how many people registered for ", "activity_statistics"},
			{"draft a venue-change notice for ", "notice_draft"},
			{"draft a notice for ", "notice_draft"},
		} {
			if strings.HasPrefix(lower, operation.prefix) {
				reply.Intent = operation.intent
				reply.Activity = strings.TrimSpace(prompt[len(operation.prefix):])
				if strings.EqualFold(reply.Activity, "the next workshop") {
					reply.Activity = ""
				}
				break
			}
		}
	}
	return reply
}
