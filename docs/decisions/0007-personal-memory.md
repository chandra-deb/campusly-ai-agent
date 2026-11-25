# 0007 — Keep personal context and preferences within one workspace

The president should be able to continue a long planning chat and carry
deliberately saved preferences into another chat in the same workspace.
On November 26, 2026, the owner approved local C1–C5 implementation
from the three Phase B tips. The owner also approved the retention and
suggestion rules below.

## Conversation context (C1)

Keep the existing transcript and choose a contiguous suffix for each
request. Starting at the latest saved turn, Core takes up to eight
whole prompt/reply pairs. The pairs and the new prompt together fit
within 12,000 UTF-8 bytes. Go independently enforces those limits on the
signed v4 route. The v1 standalone reply, v2 legacy reply and v3
Campusly intent retain their original contracts. Transient reads and
drafts do not enter the saved context.

Passing the whole transcript would make request size grow with usage.
Summarization could preserve older intent but would add another model
call and introduce a second representation of the conversation to
correct or erase. The recent window deliberately forgets older text
at inference time while the owner can still read it in the transcript.
The limit is measured in bytes. Model token counts and the best window
size for real tasks are open questions.

## Preferences (C2–C5)

Core owns each personal preference under the current president's
student ID and selected organization. Saving takes a deliberate action.
Use a small supported set of style preferences with a single active
value per kind. This gives the owner control without allowing student
records, temporary activity facts or arbitrary instruction text to
become durable preferences.

- C2 adds explicit saving and removal with provenance.
- C3 sends relevant saved preferences only to the private planning
  reply path. Live reads and operational actions keep their own
  authorization and sources.
- C4 adds inspection, editing and deletion.
- C5 suggests a supported preference only from a saved user message
  the owner selects. Accepting saves a source-linked preference.
  Rejecting saves nothing. Deleting the source conversation erases
  preferences saved from its suggestions. Manually saved preferences
  remain until their owner deletes them.

Personal memory grants no authority. Keep it outside shared AI and
private inquiry processing. This policy covers the synthetic local
prototype. Live consent, backup erasure, export and provider transfer
remain rollout decisions.

C2 stores format, detail and tone as enumerated values in PostgreSQL.
The database rejects unknown combinations and enforces one value per
kind per owner and organization. Core authenticates and checks the
active presidency again on every memory API call. The role lock
serializes saves for a workspace. Retrying the same choice returns
its existing entry, and a different choice cannot overwrite it
through the create path. The deletion URL uses the saved entry's ID
so replay cannot erase a new entry created later under the same kind.

C3 applies saved preferences to the supported workshop checklist
requests that leave reply style open. Version 5 of the Go reply
contract carries only the supported kind/value pairs. Go maps
them to fixed style text for its model adapter. Conversation text,
sources and database identifiers never appear in the style fields.
Campusly intent handling continues to read current Core data without
personal memory.

The role lock is released before generation. Before saving a new
styled reply, Core reacquires it and compares the current choices
to those sent. A changed or deleted preference discards the late
result. Earlier replies remain in the private transcript and can
be removed by deleting their conversation.
