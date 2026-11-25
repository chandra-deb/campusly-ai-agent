# Living ML/AI engineering report — draft

**Status:** Draft for owner review. AI01 evidence was reviewed by the owner on 2026-11-25; narrow Eino adoption decided. Nothing here is published.
**Date:** 2026-11-25.
**Scope:** This repository owns the AI engineering evidence and this living report. Campusly core owns its own Cloud report and real-data work. See [LEARNING_WORKFLOW](../engineering/LEARNING_WORKFLOW.md) and [DEFINITION_OF_DONE](../engineering/DEFINITION_OF_DONE.md).

## Why this report exists

Campusly needs a small, permission-bound assistant, not a general autonomous agent. The question for AI01 was practical: can Eino help build the Go service without adding machinery we do not need, or is a plain Go adapter enough for the first slice? Product direction and the A–E roadmap are agreed; narrow Eino adoption is decided, while A1 scope is not.

This report records evidence as work happens. Claims are labelled PROPOSED, TESTED, OWNER_DECIDED or NOT MEASURED. Owner review and revision of any academic submission happen separately.

## AI01 — local Eino validation (TESTED, owner-reviewed)

### Problem and expectation

A personal request must return one complete answer, use read-only scoped lookups, produce editable drafts, survive cancellation and bounded retries, and keep independent users isolated. Before writing framework code I expected Eino's component interfaces to cover the model/tool/agent plumbing while authorization and output validation stay in our code. That is an engineering expectation, not a measured superiority claim. The alternative was a hand-written Go adapter, kept as the fallback.

Hypotheses were fixed in [PLAN.md](../../experiments/ai01/PLAN.md) before execution: real Eino interfaces with a scripted model would exercise the same path a hosted adapter uses; framework retry/cancel would cover the lifecycle; anything beyond that would be application-owned.

### Method and pins

- Go `1.25.14` (local tarball, checksum verified), Eino `v0.9.23` (latest stable, not the 0.10 alpha), module `campusly.local/ai01`, isolated at [experiments/ai01](../../experiments/ai01/).
- Synthetic North/South workshop fixtures from the existing [workshop scenario](../../examples/workshop.md). No real student data, provider credentials, network inference, Core integration, CI or deployment.
- Limits fixed in advance: retry cap 2 (three attempts total), at most 3 ReAct iterations, eight concurrent requests, 25 ms deadline against a 200 ms fake delay, draft capacity 1–200. These are correctness bounds, not performance measurements.
- Official docs and v0.9.23 source were checked before coding; maintainer skills are copied under [evidence/references](../../experiments/ai01/evidence/references/). Full commands, pins, failures and raw output are in [iterations.md](../../experiments/ai01/evidence/iterations.md).

### Results

| Exercise | Result | Observed evidence |
| --- | --- | --- |
| E01 complete response | PASS | One `Generate` call returned the exact synthetic answer, no tool calls, zero `Stream` calls. |
| E02 deterministic fake | PASS | Three calls gave the three scripted strings; two transient failures then success used three attempts; a permanent error used one attempt. No credentials or inference network. |
| E03 typed read-only tool | PASS | Real `tool.InvokableTool` with required string schema; ADK emitted assistant tool call → correlated tool result → final answer; nine bad/unauthorized inputs rejected; endless loop stopped at three attempts with `exceeds max iterations`. |
| E04 structured draft | PASS | Eino derived `title`/`capacity` schema with no `action` field; local validation accepted one valid draft and rejected ten malformed/unsafe cases. No create or publish call exists. |
| E05 request lifecycle | PASS | Eight concurrent requests stayed isolated; deadline gave `context.DeadlineExceeded` with no leaked output; explicit `adk.WithCancel()` stopped the run; four concurrent retrying requests each used three attempts; `-race` clean. |
| E06 implementation reliability | RECORDED | Two API assumptions were wrong and fixed; all checks pass; unresolved items listed below. |

`gofmt` clean, `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./... -v` pass (36 test/subtest RUN lines), `go test -race ./... -v` exit 0. Raw logs: [test.txt](../../experiments/ai01/evidence/test.txt), [test-race.txt](../../experiments/ai01/evidence/test-race.txt).

### First-principles notes for review

- **Model vs agent vs tool.** The model turns messages into one message: text, or a request to call a named function with JSON arguments. A tool is that named function plus a schema and an execution result. The agent is the loop around them: ask the model, if it asked for a tool run the tool, append the result, ask the model again until it answers or a limit is reached. Eino supplies interfaces for all three; it does not decide what our data means or who may read it.
- **Tool execution loop observed here.** The trace in [test.txt](../../experiments/ai01/evidence/test.txt) has three events: assistant message with tool call `call-lookup-1`, tool event for `lookup_activity` carrying the same call ID and `Room 12`, then the final assistant answer. Correlation is the call ID, so a result cannot be silently attached to the wrong call. The loop ended after three model attempts when the script never produced a final answer.
- **Validation ownership.** Eino's schema helpers describe shapes (`title`, `capacity`, required fields). They do not enforce a free-form model answer. Strict JSON decoding, field checks and capacity bounds are our code, and authorization stays in the tool/Core boundary. A schema is documentation for the model, not a permission check.
- **Cancellation and deadlines.** A request carries a context. Every blocking step should select on that context: the fake model does, and ADK's cancel option stops the agent at a safe point. On timeout the caller gets a safe error, never partial output and never a claim that a write succeeded. Repeating a consequential write after an uncertain timeout needs a receipt from Core, not a blind retry.

### Engineering decision

TESTED locally: Eino is practical for the narrow A1 path. It gives typed model/tool/agent abstractions, a bounded ReAct loop and configurable retry/cancel with little glue. What remains application-owned is scoped identity, tool authorization, output validation, safe error contracts and all writes.

Owner decision (2026-11-25, OWNER_DECIDED): **narrow adoption** — keep Eino for the model/tool/agent boundary if A1 is authorized, keep the dependency scoped to the AI service, and do not lock it in until a credential-free hosted-adapter smoke and live tool-call compatibility are separately approved. A plain Go adapter remains the fallback for a minimal A1. This decision does not authorize A1. No comparison framework was built; no superiority claim is made.

### Limitations and next questions

- Live provider compatibility, latency, cost, model quality, token limits and permissions are NOT MEASURED. No paid call was made.
- Deadline/cancel tests use in-process delays; propagation through Core and a real provider is NOT TESTED.
- The North/South authorization fixture proves boundary handling in the spike only; Core remains authoritative at execution time.
- Retry tests cover model errors, not uncertain writes or idempotency. Eino's wider ADK surface (middleware, checkpoints, multi-agent, TurnLoop) was not exercised.
- Public docs can drift from a pinned tag; two compile errors came from exactly that. Pin source was the tie-breaker.

### Proposed A1 boundary (not authorized)

President-only personal chat with one complete non-streamed response; deterministic fake adapter for tests; hosted adapter behind the same interface only after approval; scoped actor/organization on every request; explicit deadline, retry cap, tool iteration cap and safe error text; no persistence, no tool writes, no group AI and no queue. Open questions: exact model adapter and vendor, contract versioning, identity propagation from Core, error/timeout wording, and final numeric limits. See [TASKS.md](../engineering/TASKS.md) A1; AI01 does not authorize it.

**Owner review status:** evidence reviewed on 2026-11-25; narrow Eino adoption decided. A1 authorization, adapter/vendor, contracts and limits remain open.
