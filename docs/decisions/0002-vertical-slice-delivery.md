# 0002 — Build visible vertical slices, personal assistant first

**Status:** Accepted product and delivery decision, 2026-11-25. No implementation or technology choice authorized by this ADR.

## Context

Local messaging is owner-accepted. Building the entire AI backend and shared group integration before exposing a UI would delay hands-on product feedback. The owner wants to use each small capability, learn from it, correct the design and document the engineering as it happens.

## Decision

Start with a president-only personal assistant and a complete response in a usable Campusly UI. Deliver UI and backend together in small slices with focused tests and manual owner acceptance. Add private conversation continuity, authorized information, activity and notice drafts, separately confirmed actions and personal memory. Personal memory comes before shared group AI in the planned order. Group identity, mentions, grounded answers, selective question-mark replies, reliability and summaries are separate slices. Organization knowledge and scheduled intelligence remain later, individually scoped work.

Return the whole response after it is ready, without token streaming, partial-token reconstruction or a streaming-specific milestone. Show loading and safe error/retry states. Later asynchronous jobs may deliver complete results. Keep independent services initially deployable on one host. Preserve core authority, private inquiry exclusion, exact write approval and deletion safeguards. Do not postpone minimum correctness until later reliability milestones.

The active milestone/task files define slices and actual prerequisites. AI17 retains the existing cross-repository integrated acceptance meaning, not the first manual test. An independently accepted personal release may precede it.

## Trade-offs and evidence

Earlier feedback and smaller change sets reduce integration surprises, but require coordinated UI/API checks each time. Complete responses simplify transport and structured drafts, but users wait until inference completes; loading, timeouts and latency measurements matter. Personal memory lets us learn retrieval and privacy in a simpler scope. At the time of this decision, runtime, provider, storage and queue choices were open. [ADR 0003](0003-go-bounded-requests.md) subsequently selected Go and bounded personal requests; narrow Eino adoption is decided, while provider, storage, group queue and exact limits remain open.

This docs-only decision has no application-test or benchmark result. Carry the rationale into the first relevant living ML report/blog draft; capture sanitized screenshots, actual owner feedback and measured behavior during implementation. AI01 is tested and owner-reviewed with narrow Eino adoption decided; next is explicit A1 authorization.

## Amendment — 2026-11-26

Owner adopted a mandatory phase-level gate: deterministic integration first, then real LLM/native tool orchestration and testing of all capabilities against the owner's existing local dummy-data Campusly database, followed by explicit owner acceptance. Keep the deterministic fake for regression; a provider smoke alone cannot complete the gate. A–C have deterministic evidence; R1–R4 close their open real-model gate before D1. Subsequent phases repeat the same cycle. See [constitution](../engineering/CONSTITUTION.md), [milestones](../engineering/MILESTONES.md) and [task queue](../engineering/TASKS.md). Core authorization and exact human write confirmations remain mandatory; no deployment or provider spending is implied.

## Amendment — visible tool activity (2026-11-26)

In the personal assistant, show actual tool activity as it happens: a human-readable operation label, running/succeeded/failed state, and measured duration. After the final answer, collapse the activity into an expandable per-turn summary (tool-call count and elapsed time). Include repeated calls as separate executions. No fabricated progress, raw model reasoning, secrets, or unfiltered arguments/results; optional developer diagnostics must be separately gated and sanitized. Tool-event streaming over authenticated SSE is distinct from token streaming: deliver the final assistant response as one complete message. Preserve Core authorization, scoped conversation ownership, cancellation, deadlines and exact human confirmation for writes. R1 introduces the event contract and a harmless fixture; R2–R3 connect the existing business and memory tool calls; R4 verifies the integrated H5 experience with the local dummy database.

This permits streaming lifecycle events, not model text tokens; the complete-response and local-first delivery decisions remain intact. See [architecture](../engineering/ARCHITECTURE.md) for the proposed cross-service path.
