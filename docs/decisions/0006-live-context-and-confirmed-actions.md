# 0006 — Resolve assistant questions against live Core data

A president should be able to ask what they manage from the same chat where
they plan an activity. Sending an organization's records to the model as
conversation context would make revocation difficult to enforce across
saved turns. Separate dashboard controls would avoid that problem but also
remove the natural-language entry point.

B1 adds a signed, versioned intent request from Core to the local Go
service. Go receives the prompt and request scope. It returns a bounded
operation, while Core resolves the operation with the authenticated
actor's current permissions. Core builds the reply from the authorized
rows. Campusly answers are transient in the open chat, since they may
contain facts whose permissions or values change. Regular planning
replies still use the existing private history path.

Expect the same question to reflect a changed role on the next request.
Verify with two managed organizations, an unrelated organization, and
a revoked role. The deterministic local adapter tests product behavior,
with language coverage limited to the supported sample requests.

B2 keeps the same boundary for activity questions. A named lookup matches
within the selected organization. For “next workshop,” Core uses the next
real scheduled session and falls back to the activity start only when
there are no sessions. The alternative of picking the first title match
would make ambiguous results look reliable. Expect missing or ambiguous
records to produce a clear answer without a guessed venue or date.

B3 uses the same activity lookup and source for statistics. Core reads
the committed registration and queue rows at a single statement snapshot.
For team activities, it gives team counts separately from registered
people. The alternative of sending participant lists to Go would disclose
unneeded fields. Confirm that status exclusions match the activity
service and that another organization's counts are absent.

B4 returns an editable activity proposal with a title, description,
type and capacity. Core validates the response shape, and incomplete
fields stay visible in H5 so the president can fill them in. The
alternative of encoding a creation command in free-form text would
mix drafting and authorization. Expect edits and incomplete drafts
to leave the activity database unchanged.

B5 adds a separate Core review and confirm path for activity creation.
The review returns the exact proposed action and a short-lived signed
token bound to the actor, workspace, conversation, action ID and payload.
A token can only be used after live authorization and ownership checks.
The activity, its discussion and the receipt share one transaction.
Core returns the original receipt when the same authorized request
is repeated after a lost response. A changed payload or an expired
unused token needs a fresh review.

An in-memory receipt would lose the action after a restart, while a
receipt written after the business commit leaves a gap where retry can
create another activity. Expect the database to contain one activity,
one discussion and one receipt after a replay, and none after a failed
transaction. The receipt stores identifiers and a keyed payload
fingerprint. It does not store the draft text or the private conversation.

B6 drafts a notice from the next scheduled session's recorded time and
venue. It describes the recorded venue without inferring that the venue
changed. When the source has no venue, Core asks the president to update
the activity. This is preferable to passing the activity's discussion
through the model, since the schedule owns the facts and the discussion
may include private inquiries. Expect the draft to link to the activity
and leave the timeline unchanged.

B7 sends the approved notice through the existing discussion service.
The review binds the real room ID as well as the activity ID and text.
Core checks membership and the president capability, and keeps the
discussion service's room-before-role lock order. The message,
announcement event and assistance receipt commit in one transaction.
The sender is the approving president. We expect a wrong room, an
edited notice or a revoked role to leave the public timeline unchanged.
