# Living ML/AI engineering report — draft

**Status:** Draft for owner review. AI01 evidence was reviewed by the owner on 2026-11-25; narrow Eino adoption decided. Phase A engineering evidence is ready for owner review.
**Date:** 2026-11-25.
**Scope:** This repository owns the AI engineering evidence and this living report. Campusly core owns its own Cloud report and real-data work. See [LEARNING_WORKFLOW](../engineering/LEARNING_WORKFLOW.md) and [DEFINITION_OF_DONE](../engineering/DEFINITION_OF_DONE.md).

## Why this report exists

Campusly needs a small, permission-bound assistant, not a general autonomous agent. The question for AI01 was practical: can Eino help build the Go service without adding machinery we do not need, or is a plain Go adapter enough for the first slice? The owner approved narrow Eino adoption after the isolated check, then approved local implementation of Phase A.

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

### A1 questions at the close of AI01

President-only personal chat with one complete non-streamed response; deterministic fake adapter for tests; hosted adapter behind the same interface only after approval; scoped actor/organization on every request; explicit deadline, retry cap, tool iteration cap and safe error text; no persistence, no tool writes, no group AI and no queue. Open questions: exact model adapter and vendor, contract versioning, identity propagation from Core, error/timeout wording, and final numeric limits. See [TASKS.md](../engineering/TASKS.md) A1; AI01 does not authorize it.

**Owner review status at that point:** evidence reviewed on 2026-11-25 and narrow Eino adoption chosen. The following sections record the subsequent Phase A implementation.

## A1 — A complete assistant reply from the president dashboard

A1 uses Go HTTP and the provisional Eino v0.9.23 model interface. The adapter
takes a typed user message and returns one typed assistant message. H5, Core
and Go use a synthetic workshop prompt to check request routing, scoped
authorization and complete-response handling. Core's focused checks passed
for positive access, a forged actor, missing presidency, wrong organization
and the disabled feature. See the [decision](../decisions/0004-personal-assistant-boundary.md)
and [reproduction commands](../evidence/A1.md).

This is a transport and authorization trial with a deterministic fixture.
Provider accuracy, latency and cost remain **NOT MEASURED**. The owner's
accept-or-fix decision and teach-back are **PENDING**.

## A2 — Timeout and retry in the personal chat

The retry test checked one pending message, a synthetic HTTP 504, reuse of
the same turn ID, and one completed response. Core validated the UUID and
H5 checked the echoed ID before completing the turn. This is a recovery
check for the UI and request contract. The delay was injected deliberately,
so it cannot estimate model latency or service availability. See the
[A2 evidence](../evidence/A2.md). Owner review and teach-back are **PENDING**.

## A3 — Private persistence and bounded context

The president can reopen a transcript. Core is responsible for the
stored history, and Go receives the two latest turns as typed model
messages. The fixture's answer to “Can you explain the second step?”
depends on a preceding workshop checklist. That is a controlled
demonstration of context transport and isolation. Hosted-model reasoning
quality remains unmeasured. See the
[retention and context decision](../decisions/0005-personal-conversation-retention.md).

The synthetic prototype retains conversations until the owner deletes
them. Core's tests exercise its owner predicate, turn-ID reuse and
cascade deletion. A real-data policy for backups and third-party model
transfer is **PENDING**, as are the owner's review and teach-back.

## A4 — Managing and switching conversations

The personal assistant can keep several private transcripts in one
organization. Renaming changes metadata, and deleting affects only the
selected conversation. The H5 organization selector replaces the
transcript instead of mixing activity planning across workspaces.
Core's tests check owner isolation on list, rename and delete, plus the
conversation cap. The UI scenario and evidence are linked in
[A4 local evidence](../evidence/A4.md).

Phase A establishes the product path and permission boundaries that
later evaluations can use. Grounding, abstention, quality differences
between models, latency and cost remain unmeasured. Owner acceptance
and review of the draft remain **PENDING**.

## Hosted adapter smoke — Cerebras

The owner chose `gpt-oss-120b` for a one-prompt A4 check. The
credential-free tests verified Chat Completions request shape, the fixed
synthetic input and the one-attempt bound. One real request from H5
reached the provider and received HTTP 403. The read-only models
endpoint returned a geographic-access diagnostic. Hosted completion,
key validity, model quality and model latency remain unverified.
See the [Cerebras smoke evidence](../evidence/CEREBRAS_SMOKE.md).

## B1 — Organization answers from current authorization

B1 distinguishes recognizing a question from supplying its answer. The
local Go service recognizes the managed-organizations intent. Core checks
active president assignments and sends the matching organization names
straight to H5. The answer is transient, so a previously visible role is
outside the next prompt's history. The synthetic cases exercise two
managed organizations, a revoked role, and an unrelated organization.

The deterministic recognizer exercises transport and permissions on the
supported prompts. Language coverage and hosted-model grounding remain
**NOT MEASURED**. The Cloud report is outside this single-host local
functional check. See the [decision](../decisions/0006-live-context-and-confirmed-actions.md)
and [B1 evidence](../evidence/B1.md). Owner review is **PENDING**.

## B2 — Sourced activity details and abstention

The local route separates activity intent from data selection. Core
selects the next scheduled workshop or an exact title within the current
organization. The synthetic checks found the correct venue despite an
earlier foreign workshop and a cancelled session. Missing venues and
ambiguous titles did not yield guessed facts. The mobile answer showed
the schedule and activity source. See [B2 evidence](../evidence/B2.md).

This is a controlled grounding and authorization check. We have no
estimate of open-ended intent accuracy or performance at load. Cloud
analysis is outside the local test. Owner review is **PENDING**.

## B3 — Authorized activity statistics

The assistant reports committed registration and waitlist counts from
Core. Its snapshot uses existing individual status predicates and the
separate team-capacity unit. The synthetic individual fixture excludes a
withdrawal and a promoted waitlist entry. The mobile path displays an
activity source and the same active totals. See [B3 evidence](../evidence/B3.md).

The result is an aggregate grounding check. Model accuracy and Cloud
throughput remain **NOT MEASURED**. Owner review is **PENDING**.
