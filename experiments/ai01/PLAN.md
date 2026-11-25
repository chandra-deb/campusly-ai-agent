# AI01 — Eino local validation: pre-execution plan

**Date:** 2026-11-25. **Status:** Planned before the spike code was written.
**Scope:** Isolated experiment under `experiments/ai01/`. Not production code, not A1, no hosted-provider calls, no real data, no CI and no deployment.

## Problem and expectation

Campusly AI needs a bounded Go assistant: one complete non-streamed response, deterministic fake-first tests, scoped read-only tools, structured drafts, safe cancellation/retries and concurrent independent requests. The decision to make is whether Eino earns a place in A1 or whether a small Go adapter is enough.

Expectations (engineering expectations, not measured results):

- Eino's component interfaces let the spike use real framework abstractions (`model.ToolCallingChatModel`, `tool.InvokableTool`, ADK `ChatModelAgent`/`Runner`) while all inference is scripted locally.
- A fake model can exercise the same `Generate`/tool-call path a hosted adapter will use, without credentials or network inference.
- The framework's retry and cancel options cover the request lifecycle; authorization, draft validation and safe errors remain application-owned.
- If the ReAct/tool loop or lifecycle needs more code than a plain Go function, that is evidence for narrowing Eino's use, not a reason to hide a failure.

Alternatives considered before execution: plain Go adapter with hand-written tool loop (kept as the fallback recommendation if Eino adds little), Pydantic AI/AgentScope in Python (not appropriate for the Go service; no parallel implementation approved).

## Toolchain and pins

| Item | Value |
| --- | --- |
| Go | `go1.25.14 linux/amd64`, downloaded from `https://go.dev/dl/go1.25.14.linux-amd64.tar.gz`, SHA-256 `a21ae5633a269bcd7e90cf767e48225633795e99d831742cbf3397064fee7712`, unpacked under `/tmp/opencode/go-sdk` (no system install, no sudo) |
| Eino | `github.com/cloudwego/eino v0.9.23`, latest stable release at planning time (tag commit `3548031`, released 2026-11-09). The `v0.10.0-alpha.*` line is deliberately excluded. |
| Module | `campusly.local/ai01`, `go 1.25`, isolated under `experiments/ai01/` |
| Dependency pin record | `experiments/ai01/go.mod` and `go.sum`; SHA-256 of `go.sum` recorded in [evidence/iterations.md](evidence/iterations.md) after resolution |
| Direct imports expected | `github.com/cloudwego/eino/components/model`, `components/tool/utils`, `adk`, `compose`, `schema` (no eino-ext runtime dependency needed for the fake; only its maintainer skills were read) |

No paid call, provider credential, real campus data, private core source or network inference is used. Public module downloads are allowed by the brief.

## Official references (checked before coding)

- Eino docs: <https://www.cloudwego.io/docs/eino/>; ChatModel <https://www.cloudwego.io/docs/eino/core_modules/components/chat_model_guide/>; tools <https://www.cloudwego.io/docs/eino/core_modules/components/tools_node_guide/how_to_create_a_tool/>; ReAct agent <https://www.cloudwego.io/docs/eino/core_modules/flow_integration_components/react_agent_manual/>; ADK agent <https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_implementation/chat_model/>; runner/events <https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_extension/>; cancel <https://www.cloudwego.io/docs/eino/core_modules/eino_adk/agent_cancel_and_turnloop_quickstart/>; v0.9 notes <https://www.cloudwego.io/docs/eino/release_notes_and_migration/eino_v0.9._agentic-runtime/eino_v0.9_migration_notes/>.
- Release and source: <https://github.com/cloudwego/eino/releases/tag/v0.9.23>, module source in the local module cache at `github.com/cloudwego/eino@v0.9.23`. `llms.txt` in that tag was also checked.
- Maintainer skills (eino-ext `main`): `skills/eino-agent/SKILL.md`, `skills/eino-agent/reference/chat-model-agent.md`, `skills/eino-agent/reference/runner-and-events.md`, `skills/eino-component/SKILL.md`, `skills/eino-component/reference/model/overview.md`, `skills/eino-component/reference/tool/overview.md`, `skills/eino-component/reference/tool/utility.md`. Local copies: [evidence/references](evidence/references/). The skills describe the v0.9 ADK surface; every API used below was additionally checked against the pinned v0.9.23 source.

The public docs were last edited in March 2026 and can drift from a pinned tag. Source at v0.9.23 is the tie-breaker; mismatches are recorded in the iteration log.

## Synthetic fixtures

From [examples/workshop.md](../../examples/workshop.md), kept entirely synthetic:

- Organizations `north` and `south`. Activity `north-workshop`: venue `Room 12`, title `Intro to Python`, owner org `north`. Activity `south-workshop`: same title, venue `Room 9`, owner org `south`.
- Actors/scopes: `mira` (organizer in `north`), `lee` (member in `north`), `sam` (organizer in `south`).

The read-only tool receives `activity_id`, `organization` and `actor`, matches all three against the in-memory fixture and returns venue/title/source. It never writes. `mira`/`north` + `north-workshop` succeeds; wrong tenant, unknown ID, malformed/missing arguments and a member acting outside an organizer read are rejected without a successful lookup. These negative checks exercise the spike's boundary only; they do **not** establish production Core authorization.

## Chosen limits (fixed before running)

| Limit | Value |
| --- | --- |
| Retry cap | `MaxRetries = 2` → at most 3 total model attempts; retry only on a sentinel transient error; permanent errors stop after 1 attempt. Backoff forced to 0 in tests for determinism. |
| ReAct/tool iterations | `MaxIterations = 3` for the agent; a scripted endless tool-call script is used to confirm it stops at the cap. |
| Request workload | 8 independent concurrent requests, each with a unique prompt key; 5 attempts each; no shared mutable per-request state. Not a capacity benchmark. |
| Deadline | 25 ms per synthetic request in the deadline case; fake delay 200 ms. Race check with `-race` on the focused packages only. |
| Draft bounds | Non-empty title, capacity integer in `[1, 200]`, unknown fields rejected, no `action`/write field accepted. |

No latency, throughput, memory or cost claim is made from these numbers.

## Exercise mapping

| ID | Check | Evidence location |
| --- | --- | --- |
| E01 | Complete response through `Generate`, no partial-token delivery | `complete_response_test.go` |
| E02 | Scripted outputs, transient/permanent injection, retry cap, no credentials/network | `fake_model_test.go`, `retry_test.go` |
| E03 | Real tool abstraction + bounded agent, schema, call/result correlation, negative IDs and scopes | `tool_lookup_test.go` |
| E04 | Structured draft; framework schema vs local Go validation | `draft_test.go` |
| E05 | 8 isolated concurrent requests, cancellation, deadline, bounded retries, race | `lifecycle_test.go` |
| E06 | Official refs, API corrections, complexity and unresolved issues | `evidence/iterations.md`, report/blog drafts |

Commands, raw output and any failures are retained under `evidence/`; living report/blog drafts are updated only from observed results.
