# Completion evidence

Record task ID, revisions on each side, changed contract version, exact focused commands, environment, outcomes, untested cases and rollback. Exercise relevant denial, replay and failure paths. A clean standalone build and synthetic run are required for the first executable release. Log identifiers and outcomes without prompts, credentials or conversations. Distinguish PROPOSED, APPROVED, IMPLEMENTED, TESTED and OWNER_ACCEPTED. Run the full repository suite only with explicit owner approval.

## Documentation and learning gate

- Record the real problem, chosen alternative, material decisions, failures, trade-offs and reproducible evidence. Update diagrams and ADRs when affected.
- Update the technical blog draft and relevant living course-report section at each substantive task/milestone; if not applicable, record why. Keep measured claims linked to revisions, commands, configuration and raw evidence. Mark unknowns PENDING or NOT MEASURED.
- Record owner review/teach-back separately from engineering verification. Owner verifies and revises final academic work; public publication, real-data pilot and deployment are separately authorized.

## Vertical-slice acceptance

- User-facing work includes a usable local UI and integrated core/API/AI behavior. Record a manual owner scenario, observed issues, fixes and explicit accept/fix status. Enabling infrastructure tasks link to their intended slice.
- Test complete-response behavior, loading, timeout/failure handling and current-scope authorization when relevant. No token streaming; R1 onward requires visible tool-event progress via authenticated SSE. Required privacy, deletion and idempotency safeguards precede the capabilities that need them.
- Capture sanitized screenshots and feedback for reports/blogs. Agent tests do not imply owner acceptance.

## Mandatory AI phase gate (2026-11-26)

- Deterministic focused tests and scoped end-to-end integration pass; fake mode remains reproducible without provider credentials.
- Real LLM and native tool orchestration are exercised against the owner's existing local Campusly dummy-data database. Cover all phase capabilities, including non-tool context/draft/confirmation behavior; verify answers against canonical Core records.
- Record real provider/model, revisions, case, tool selections/arguments/results, latency/usage if available, failures and omissions; check paraphrases, multi-tool reasoning, missing/ambiguous inputs, denial, retry and exact human confirmation when relevant.
- Keep deterministic verification, real-model verification and owner acceptance separate. A failed/blocked live call is not success. Owner accepts or requests fixes before advancing to the next AI phase; R4 gates D1. No production, credentials in logs, or unauthorized spend/deployment.

## Tool activity acceptance (R1–R4)

- Actual started/completed/failed events arrive while processing, with safe labels, stable call identities and measured durations. A plain chat shows no invented tool calls.
- Final result is one complete response; per-turn activity collapses and can be expanded/reopened. History is scoped and deleted with the turn.
- Test multiple calls, failed calls, cancellation/disconnect, duplicate/out-of-order events and conversation switching; reject unauthorized/stale events. Never display raw reasoning or sensitive payloads. Confirmed writes still require explicit user approval and canonical receipts.
- Capture sanitized H5 evidence and separate tool timing from overall turn timing. R4 owner acceptance remains mandatory.

## Framework and runtime verification

For framework-dependent changes, record official documentation/source references for the pinned version, exact dependency versions, formatting, build, focused Go test and vet outcomes (race checks where concurrency is exercised and supported). Verify behavior and permissions, not just compilation. For personal requests test deadline/cancellation, safe error and relevant concurrency/retry limits. No framework validation result may be claimed before AI01 runs. Retain errors and failed attempts as evidence.
