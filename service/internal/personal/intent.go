package personal

import "strings"

// IntentReply proposes a read or draft operation. Core authorizes it and
// supplies the data. The resolver has no Campusly credentials or tools.
type IntentReply struct {
	Version int `json:"version"`
	Scope
	Intent   string `json:"intent"`
	Activity string `json:"activity,omitempty"`
}

func resolveIntent(input Request) IntentReply {
	reply := IntentReply{Version: input.Version, Scope: input.Scope, Intent: "chat"}
	prompt := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(input.Prompt), "?？"))
	lower := strings.ToLower(prompt)
	switch strings.TrimRight(lower, ".!") {
	case "which organizations do i manage", "what organizations do i manage",
		"show my organizations", "list the organizations i manage":
		reply.Intent = "managed_organizations"
	case "when and where is the next workshop", "where is the next workshop",
		"when is the next workshop", "show the next workshop", "what is the next workshop":
		reply.Intent = "activity_details"
	default:
		const prefix = "show details for "
		if strings.HasPrefix(lower, prefix) {
			reply.Intent = "activity_details"
			reply.Activity = strings.TrimSpace(prompt[len(prefix):])
		}
	}
	return reply
}
