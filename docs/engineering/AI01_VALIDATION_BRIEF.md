# AI01 — Eino practical validation brief

**Date:** 2026-11-25
**Status:** Ready for delegation; execution and owner acceptance pending.
**Scope:** A small, isolated framework evaluation in the owner's development environment. This brief is not production code or an A1 implementation authorization.

## Read first

- [Agent instructions](../../AGENTS.md), [constitution](CONSTITUTION.md), [architecture](ARCHITECTURE.md), [milestones](MILESTONES.md), [tasks](TASKS.md), [definition of done](DEFINITION_OF_DONE.md), and [learning workflow](LEARNING_WORKFLOW.md).
- [ADR 0003](../decisions/0003-go-bounded-requests.md) is the authoritative language, execution and validation decision; [ADR 0002](../decisions/0002-vertical-slice-delivery.md) defines delivery order.
- Official [Eino documentation](https://www.cloudwego.io/docs/eino/), [repository](https://github.com/cloudwego/eino), [examples](https://github.com/cloudwego/eino-examples), and [extensions and maintainer skills](https://github.com/cloudwego/eino-ext). Inspect the release-specific source and applicable skills before coding. Confirm current released versions and compatibility rather than assuming the newest tag is suitable.

## Objective and boundaries

Determine whether Eino is a practical fit for Campusly AI's bounded personal assistant. Go, standard-library HTTP initially, hosted inference and complete non-streamed responses are agreed. Eino is the leading candidate; its production adoption remains pending. Python remains available for offline ML evaluation. Core retains all authorization and business writes. Later shared group AI will have independent durable asynchronous processing; it is not part of this exercise.

The owner or explicitly delegated coding agent performs the implementation and tests in the owner's development environment. The architecture/review assistant researches, specifies, explains and reviews supplied evidence; it does not create experimental code, run tests, configure environments or launch CI. This brief does not authorize any GitHub Actions workflow, cloud execution, deployment, provider call, spending, real data transfer, Core/UI integration or A1 implementation. Inspect the current checkout before work; keep the spike isolated from production code. Use local execution only and obtain separate approval before pushing or merging any code. The previous experimental branch is superseded; do not depend on or reuse its unverified code or workflow.

## Validation plan

Before execution, document toolchain and compatible pinned Eino release, dependencies and checksums, official documentation references, synthetic inputs, expected outcomes and chosen test limits. A deterministic fake model must exercise real Eino interfaces; do not mock away the feature being evaluated. No provider credentials or external inference are required. A network download of public dependencies in the owner's environment is not a model call.

| ID | Exercise | Pass criteria |
| --- | --- | --- |
| E01 | Complete response | Invoke Eino's non-streaming model interface with a fixed synthetic prompt; receive the exact expected complete response. No partial-token delivery. |
| E02 | Deterministic fake | Repeated inputs produce expected scripted outputs; inject transient and permanent failures. Tests run without model credentials or inference network access. Document the selected retry cap; transient read attempts stay within it, permanent errors do not retry. |
| E03 | Typed read-only tool | Use an actual Eino tool abstraction and bounded agent execution to look up a synthetic workshop. Verify argument schema, tool-call/result correlation and expected venue. Reject malformed, missing, unknown and unauthorized activity IDs without successful lookup. Repeated tool calls stop within documented limits. No real business write. |
| E04 | Structured draft | Produce a synthetic activity draft with title and positive bounded capacity. Reject malformed JSON, missing fields, invalid types/ranges and unexpected action fields. Distinguish framework/provider schema features from local Go validation. No creation or publication. |
| E05 | Request lifecycle | Run independent synthetic requests concurrently, verify per-request isolation, cancellation, deadline, bounded retries and safe errors. Include race-enabled focused checks where supported. Choose and record a small reproducible workload and budget before running; do not claim production capacity. |
| E06 | Implementation reliability | Record official source/skill references, any API assumptions disproven by compilation or tests, corrections, code complexity and unresolved issues. Build, format, test and vet the actual spike. Do not infer model training familiarity or claim superiority from a single coding attempt. |

Use a synthetic workshop and two distinct synthetic organization/actor scopes where relevant. Include negative checks on authorization boundaries in the fixture; these do not establish production Core authorization. Keep experiments small: no actual Campusly API, database, memory, queues, group bot or multi-agent orchestration. Do not build parallel framework implementations unless the owner later approves a comparison prompted by a meaningful limitation.

## Expected handoff evidence

Return a concise validation record containing:
- Repository revision and local experiment path; Go/Eino versions and dependency pins/checksums; official documentation and skill links.
- Exact commands and outcomes for formatting, build or test compilation, focused tests, vet and race checks where supported; retain failure output and fixes. Do not run the full repository suite without approval.
- E01–E06 individually marked PASS, FAIL or NOT_RUN with supporting observations, including negative cases, actual limitations and any environment blockers.
- The minimal architecture and code paths, why Eino is useful, what remains application-owned, and whether a simpler Go adapter would suffice for A1.
- A recommendation to adopt Eino, narrow its use, or investigate a fallback, with concrete rationale. Live provider compatibility, latency, cost, model quality, production permissions and scaling remain unverified by fake tests. No unmeasured performance claims.
- Proposed A1 boundary and remaining design questions (model adapter, scoped identity, complete response/error contracts, limits); do not implement A1 or lock dependencies on behalf of the owner.

Update the living ML report and engineering blog draft from actual work and observed evidence, or record why an update is not applicable. Preserve failures; redact credentials, private data and sensitive operational details. UI screenshots are not required for this non-UI evaluation. Explain model versus agent versus tool, the tool execution loop, validation ownership and context cancellation in first-principles terms for the owner's review.

## Exit gate

AI01 engineering evaluation is complete only when the agreed focused evidence is available and limitations are explicit. The owner and architecture/review assistant review the results together. The owner separately approves or rejects Eino and the scope of A1. If a blocker prevents testing, report it and leave the relevant checks NOT_RUN. Do not create a CI workaround. Nothing in this brief claims Eino has passed, any implementation is complete, or any production action is authorized.
