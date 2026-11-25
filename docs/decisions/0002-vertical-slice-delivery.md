# 0002 — Build visible vertical slices, personal assistant first

**Status:** Accepted product and delivery decision, 2026-11-25. No implementation or technology choice authorized by this ADR.

## Context

Local messaging is owner-accepted. Building the entire AI backend and shared group integration before exposing a UI would delay hands-on product feedback. The owner wants to use each small capability, learn from it, correct the design and document the engineering as it happens.

## Decision

Start with a president-only personal assistant and a complete response in a usable Campusly UI. Deliver UI and backend together in small slices with focused tests and manual owner acceptance. Add private conversation continuity, authorized information, activity and notice drafts, separately confirmed actions and personal memory. Personal memory comes before shared group AI in the planned order. Group identity, mentions, grounded answers, selective question-mark replies, reliability and summaries are separate slices. Organization knowledge and scheduled intelligence remain later, individually scoped work.

Return the whole response after it is ready, without token streaming, partial-token reconstruction or a streaming-specific milestone. Show loading and safe error/retry states. Later asynchronous jobs may deliver complete results. Keep independent services initially deployable on one host. Preserve core authority, private inquiry exclusion, exact write approval and deletion safeguards. Do not postpone minimum correctness until later reliability milestones.

The active milestone/task files define slices and actual prerequisites. AI17 retains the existing cross-repository integrated acceptance meaning, not the first manual test. An independently accepted personal release may precede it.

## Trade-offs and evidence

Earlier feedback and smaller change sets reduce integration surprises, but require coordinated UI/API checks each time. Complete responses simplify transport and structured drafts, but users wait until inference completes; loading, timeouts and latency measurements matter. Personal memory lets us learn retrieval and privacy in a simpler scope. Runtime, provider, storage and queue choices remain open.

This docs-only decision has no application-test or benchmark result. Carry the rationale into the first relevant living ML report/blog draft; capture sanitized screenshots, actual owner feedback and measured behavior during implementation. Next: AI01 design review, then explicit A1 authorization.
