package personal

import (
	"regexp"
	"strings"
)

type ActivityDraft struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"activity_type"`
	Capacity    string `json:"capacity"`
}

var seatsSuffix = regexp.MustCompile(`(?i) (?:for|with) ([0-9]+) (?:people|students|seats)$`)

func proposeActivity(prompt string) *ActivityDraft {
	prompt = strings.TrimSpace(prompt)
	lower := strings.ToLower(prompt)
	draft := &ActivityDraft{Kind: "activity", Type: "WORKSHOP"}
	for _, prefix := range []string{"draft an activity:", "draft a workshop:"} {
		if strings.HasPrefix(lower, prefix) {
			fields := strings.SplitN(prompt[len(prefix):], "|", 3)
			values := []*string{&draft.Title, &draft.Description, &draft.Capacity}
			for i, value := range fields {
				*values[i] = strings.TrimSpace(value)
			}
			return draft
		}
	}
	for _, prefix := range []string{"draft an activity called ", "draft a workshop called "} {
		if strings.HasPrefix(lower, prefix) {
			draft.Title = strings.TrimSpace(prompt[len(prefix):])
			if match := seatsSuffix.FindStringSubmatchIndex(draft.Title); match != nil {
				draft.Capacity = draft.Title[match[2]:match[3]]
				draft.Title = draft.Title[:match[0]]
			}
			return draft
		}
	}
	if lower == "draft an activity" || lower == "draft a workshop" {
		return draft
	}
	return nil
}
