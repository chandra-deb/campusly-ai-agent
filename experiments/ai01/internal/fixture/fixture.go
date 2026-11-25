// Package fixture provides the synthetic workshop data from the AI01 plan.
// It is immutable in memory and contains no real campus data.
package fixture

import (
	"errors"
	"regexp"
)

var (
	ErrInvalidArguments = errors.New("invalid activity arguments")
	ErrNotFound         = errors.New("activity not found in scope")
	ErrDenied           = errors.New("activity lookup denied")
)

var activityIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

type Activity struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Venue  string `json:"venue"`
	Source string `json:"source"`
	Org    string `json:"organization"`
}

type Repository struct {
	activities map[string]Activity
	organizers map[string]string
	members    map[string]string
}

func New() *Repository {
	return &Repository{
		activities: map[string]Activity{
			"north-workshop": {ID: "north-workshop", Title: "Intro to Python", Venue: "Room 12", Source: "synthetic:north:activity-1", Org: "north"},
			"south-workshop": {ID: "south-workshop", Title: "Intro to Python", Venue: "Room 9", Source: "synthetic:south:activity-1", Org: "south"},
		},
		organizers: map[string]string{"mira": "north", "sam": "south"},
		members:    map[string]string{"lee": "north"},
	}
}

// Lookup is a read-only, scope-checked lookup. A member cannot use an
// organizer read, and an activity from another organization is reported as
// not found so the venue is not disclosed.
func (r *Repository) Lookup(id, org, actor string) (Activity, error) {
	if id == "" || org == "" || actor == "" || !activityIDPattern.MatchString(id) {
		return Activity{}, ErrInvalidArguments
	}
	if organizerOrg, ok := r.organizers[actor]; !ok || organizerOrg != org {
		return Activity{}, ErrDenied
	}
	activity, ok := r.activities[id]
	if !ok || activity.Org != org {
		return Activity{}, ErrNotFound
	}
	return activity, nil
}

func (r *Repository) IsOrganizer(actor string) bool {
	_, ok := r.organizers[actor]
	return ok
}

func (r *Repository) IsMember(actor string) bool {
	_, ok := r.members[actor]
	return ok
}
