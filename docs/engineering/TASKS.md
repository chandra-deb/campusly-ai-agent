# AI work queue — small visible slices

Owner-agreed roadmap, 2026-11-25; real-model gate adopted 2026-11-26. A1–A4, B1–B8 and C1–C5 have deterministic local evidence; real-model and owner acceptance remain pending. AI01 is tested and owner-reviewed, with narrow Eino adoption decided. Dependencies are minimum technical prerequisites, not a rigid chain. Each slice needs explicit scope, usable local UI, focused checks, owner manual accept/fix decision and living report/blog evidence. Core/UI changes belong in their owning repositories.

| ID | Depends on | Deliverable | Manual test and focused checks |
| --- | --- | --- | --- |
| AI01 | Foundation | Validate Eino as the leading Go framework candidate using the delegated local [AI01 validation brief](AI01_VALIDATION_BRIEF.md) and [ADR 0003](../decisions/0003-go-bounded-requests.md); review model adapter, scoped identity, deadlines and A1 synthetic/UI path | Record versions, focused results, limitations and framework decision (narrow Eino adoption, 2026-11-25); A1 authorization, hosted-adapter check and paid calls remain separate approvals |
| A1 | AI01 | [A1 evidence](../evidence/A1.md) · President-only chat UI, core integration, standalone synthetic fixture, complete model response | Open assistant and chat; deny non-president and wrong actor; clean build |
| A2 | A1 | [A2 evidence](../evidence/A2.md) · Loading, timeout, error and safe retry improvements (basic handling already in A1) | Simulate failure; retry without duplicate turns |
| A3 | A1 | [A3 evidence](../evidence/A3.md) · Private conversation persistence | Reopen chat; verify owner isolation, retention and deletion |
| A4 | A3 | [A4 evidence](../evidence/A4.md) · Create, rename, switch and delete conversations | Manage separate activity chats |
| B1 | A1 | [B1 evidence](../evidence/B1.md) · Authorized organization awareness | Ask which organizations I manage; deny unrelated or revoked access |
| B2 | B1 | [B2 evidence](../evidence/B2.md) · Sourced activity details, schedule, venue and status | Ask about next workshop; deny wrong tenant, abstain if missing |
| B3 | B1 | [B3 evidence](../evidence/B3.md) · Authorized activity statistics | Compare counts with canonical data; deny unauthorized reads |
| B4 | B1 | [B4 evidence](../evidence/B4.md) · Structured editable activity drafts, no official write | Describe workshop and edit fields; validate incomplete drafts |
| B5 | B4 | [B5 evidence](../evidence/B5.md) · Confirmed activity creation with core receipt | Create locally; reject edited, expired or revoked approvals and duplicates |
| B6 | B2 | [B6 evidence](../evidence/B6.md) · Editable, grounded notice drafts | Review venue-change notice and source |
| B7 | B6 | [B7 evidence](../evidence/B7.md) · Confirmed notice publication with human attribution | Approve notice; test exact target, replay and revocation |
| B8 | B1 | [B8 evidence](../evidence/B8.md) · Confirmed personal group messaging with canonical receipt | Send as president; test exact payload, replay, expiry and revocation |
| C1 | A3 | [C1 evidence](../evidence/C1.md) · Long conversation context management | Continue long chat; test context budget and isolation |
| C2 | A3 | [C2 evidence](../evidence/C2.md) · Explicit personal memory with provenance and deletion | Remember a style; test isolation and removal |
| C3 | C2 | [C3 evidence](../evidence/C3.md) · Relevant cross-conversation personal recall | New chat uses relevant preferences, not deleted ones |
| C4 | C2 | [C4 evidence](../evidence/C4.md) · Personal memory management UI | Inspect, edit and delete; verify changed behavior |
| C5 | C2, C4 | [C5 evidence](../evidence/C5.md) · Reviewable personal memory suggestions | Accept or reject suggestions; exclude temporary/sensitive information |
| R1 | A1–C5 deterministic evidence | Configurable hosted model adapter; retain fake mode; validate provider compatibility and bounded complete responses | Real chat with paraphrases, failures, timeout/cancellation, history, provider/model and usage evidence; no silent fake fallback |
| R2 | R1 | Connect all existing Phase B capabilities to real orchestration: managed organizations, activity details/statistics, activity/notice/message drafts and confirmed creation/publication/send | Typed calls, real local dummy DB results, multi-tool follow-ups, missing/ambiguous targets, invalid arguments, wrong scope/revocation, exact UI confirmation and replay; no direct model-authorized writes |
| R3 | R2 | Connect A/C conversation and memory behavior to real inference: history, bounded context, scoped preference list/save/edit/delete, relevant recall, suggestions and acceptance | Natural paraphrases, new chat recall, removal and source deletion, isolation and fresh consent; no permissions from memory |
| R4 | R3 | Integrated real-model acceptance of all A1–C5 on compatible AI/Core/H5 revisions using owner's existing local dummy-data database | Cover every capability; compare to canonical data, inspect tool traces, record outcomes/limitations, run focused deterministic regressions, owner manual accept/fix; D1 blocked until accepted |
| D1 | B1, R4 owner accepted | Group AI identity and organizer enablement | Enable/disable bot in authorized room |
| D2 | D1 | Durable public human @mention replies; minimum retry, deduplication, limits and privacy safeguards | Mention bot; replay event; verify private/bot/system events excluded and chat survives outage |
| D3 | D2, B2 | Grounded scoped group answers | Ask venue; verify source, abstention and live authorization |
| D4 | D3 | Selective public question-mark participation (ASCII/fullwidth), reply/no-reply decision | Compare relevant questions, jokes and ordinary text; no extra inference for ineligible events |
| D5 | D2 | Strengthen group AI recovery, retry, concurrency and isolation | Crash after post before ack; verify one receipt, no private leakage and chat recovery |
| D6 | D3 | On-demand authorized group summaries without persistent group memory | Summarize public discussion; exclude private messages |
| E1 | D1 | Organization team groups, separate rules review | Open authorized team group; test membership |
| E2 | B8 | Needs-attention inbox with canonical links | Open item and verify source and access |
| E3 | D6 | On-demand cross-group retrieval | Ask across authorized groups; exclude inaccessible sources |
| E4 | D6 | Persistent attributed group memory | Inspect decision; check consent, revocation and deletion |
| E5 | E4 | Group memory management | Correct/delete entry with scoped rights |
| E6 | E5 | Reviewed organization knowledge | Approve promotion with source and visibility checks |
| E7 | E6 | Activity retrospectives | Review evidence-backed lessons |
| E8 | Separate approval | Scheduled intelligence and controlled automation; split before implementation | Review each workflow, permissions, budget, cancellation and schedule separately |
| AI17 | B8, D4, D5 | Integrated compatibility checkpoint for original personal/group workflow | Owner tests workshop scenario, permissions, private exclusion, worker failure and canonical receipts; separate Cloud authorization |

## Guardrails

Every future AI phase follows deterministic checks → real LLM/tool calls → local integrated tests → owner acceptance. Test all introduced capabilities, not only reads. Native tool calls are typed model API requests; drafts, context and UI confirmation may use structured output or application logic rather than separate registered tools. Local dummy-data use is approved for manual tests; keep production, secrets and unapproved real data out. R1–R4 are the next work, not completed tasks. See [constitution](CONSTITUTION.md).

Server-side authorization, version rejection and safe logging start at A1. Before operational writes, bind approval to actor/action/target/payload, expiry and current permissions; edits need fresh approval. Before persistence, define retention and deletion; no memory grants permissions or replaces current core facts. Personal memory never enters shared AI. Private inquiries remain excluded. D2 includes minimum reliable delivery from day one; D5 strengthens it. No token streaming. No real data, provider calls, deployment, spending or recurring jobs without separate approval. Record actual tests and sanitized screenshots; update report/blog drafts or explain N/A. Split oversized later slices before implementing.

For an approved data tool, add a scoped Core API when the current endpoints cannot supply the needed Campusly data. The API must enforce live authorization for its own resource and actions.

## Legacy IDs

AI10–AI16 are superseded, not completed: AI10 → A1; AI11 → A1/B1; AI12 → B2; AI13 → B8; AI14 → D2/D4/D5; AI15 → D3/D4; AI16 → per-slice integration and AI17. AI01 and AI17 remain for CAMS compatibility.
