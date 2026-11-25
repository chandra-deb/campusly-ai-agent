# 0005 — Keep private conversations in Core

The president should be able to reopen a chat on another browser session.
The owner chose to retain conversations until deletion in the synthetic
prototype. Core stores them by student and organization and hard-deletes
the turns with their conversation. Storing them in the browser would make
access and deletion depend on the device. Storing them in Go would require
a second authorization and database boundary.

Signed inferencing requests use a v2 endpoint for up to two preceding
turns, in order. The service still accepts v1 for the standalone fixture
and older Core chat. A v2 request is limited to 96 KiB of JSON, enough for
worst-case JSON escaping at the same 2,000-byte prompt and 4,000-byte
reply limits. The deadline and concurrency limits stay the same. Go builds
typed user and assistant messages from this context and holds them only
while the call is in flight.

The Core conversation row orders saves. Core checks for an existing turn
before inference and again while saving under the conversation lock. Two
attempts may generate twice, but the matching retry receives the original
saved reply. Different turns prepared from the same history cannot both
append silently. A hard delete during inference causes a late save to
fail. These are consistency and retention checks for a disposable local
database. Live-data consent, backup retention and provider transfer need
separate decisions for a real rollout.
