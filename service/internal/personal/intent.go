package personal

import "strings"

// IntentReply proposes a read or draft operation. Core authorizes it and
// supplies the data. The resolver has no Campusly credentials or tools.
type IntentReply struct {
	Version int `json:"version"`
	Scope
	Intent string `json:"intent"`
}

func resolveIntent(input Request) IntentReply {
	reply := IntentReply{Version: input.Version, Scope: input.Scope, Intent: "chat"}
	prompt := strings.ToLower(strings.TrimRight(strings.TrimSpace(input.Prompt), "?.!？"))
	switch prompt {
	case "which organizations do i manage", "what organizations do i manage",
		"show my organizations", "list the organizations i manage":
		reply.Intent = "managed_organizations"
	}
	return reply
}
