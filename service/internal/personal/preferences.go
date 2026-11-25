package personal

import (
	"strconv"
	"strings"

	"github.com/cloudwego/eino/schema"
)

type Preference struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

var preferenceText = map[Preference]string{
	{"format", "bullets"}:  "Use bullet points for lists.",
	{"format", "numbered"}: "Use numbered steps for lists.",
	{"detail", "concise"}:  "Keep the reply concise.",
	{"detail", "expanded"}: "Give extra practical detail.",
	{"tone", "direct"}:     "Use a direct tone.",
	{"tone", "warm"}:       "Use a warm tone.",
}

func preferenceInstructions(preferences []Preference) (string, bool) {
	if len(preferences) < 1 || len(preferences) > 3 {
		return "", false
	}
	instructions := []string{"These are saved reply style defaults. Follow the current request when it specifies a style."}
	seen := make(map[string]bool, len(preferences))
	for _, choice := range preferences {
		text, supported := preferenceText[choice]
		if !supported || seen[choice.Kind] {
			return "", false
		}
		seen[choice.Kind] = true
		instructions = append(instructions, text)
	}
	return strings.Join(instructions, "\n"), true
}

// workshopChecklist gives the synthetic model a predictable way to
// demonstrate style changes without interpreting user text as instructions.
func workshopChecklist(messages []*schema.Message) string {
	instructions := ""
	if messages[0].Role == schema.System {
		instructions = messages[0].Content
	}
	has := func(kind, value string) bool {
		return strings.Contains(instructions, preferenceText[Preference{kind, value}])
	}
	lines := []string{"Workshop checklist"}
	switch {
	case has("tone", "direct"):
		lines = append(lines, "Start with these steps.")
	case has("tone", "warm"):
		lines = append(lines, "Here’s a plan your team can work through together.")
	}
	steps := []string{
		"Confirm the room, time and capacity.",
		"Assign a host and a check-in lead.",
		"Test the projector and materials.",
		"Draft a reminder for participants.",
	}
	switch {
	case has("detail", "concise"):
		steps = steps[:2]
	case has("detail", "expanded"):
		steps = append(steps, "Plan who will help if participants need equipment.")
	}
	for i, step := range steps {
		prefix := strconv.Itoa(i+1) + ". "
		if has("format", "bullets") {
			prefix = "• "
		}
		lines = append(lines, prefix+step)
	}
	return strings.Join(lines, "\n")
}
