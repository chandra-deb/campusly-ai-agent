# Validating Eino without a model key — AI01 notes

**Draft, 2026-11-25.** For owner review. Not published. This is an engineering note from a small local spike, not a product announcement or a benchmark result.

## The question

Campusly's first assistant slice is small: a president asks something, gets one complete answer, and later gets read-only lookups and editable drafts. Go is the service language; standard-library HTTP and context handle the boundary. Eino looked like a useful framework, but "looks useful" is not evidence. AI01 asked whether we can build the narrow path with it while keeping authorization and validation in our own code — and where the framework stops being helpful.

I used a scripted model, not a hosted one. That is the point: no credentials, no paid calls, no real data, and a result I can reproduce offline.

## What the three pieces actually do

It is easy to blur these together, so here is the simple version I would explain to another student developer:

- The **model** is a function from messages to one message. It can answer with text or with a request to call a named function and some JSON arguments. Eino calls this `Generate` for a complete answer; `Stream` is the token-by-token version we deliberately never use.
- A **tool** is the named function plus a description of its arguments. Eino's `ToolInfo` carries that schema. Running the tool is still ordinary Go code — ours — and in this spike it can only read an in-memory workshop fixture.
- The **agent** is the loop: ask the model, run any tool it asked for, append the result as a message, and ask again until there is a final answer or a limit is hit. Eino's ADK `ChatModelAgent` and `Runner` implement that loop.

So Eino supplies plumbing and types. It does not decide whether Mira may read a North workshop, whether a draft is valid, or whether an action is allowed. Those are product rules and stay with us and Core.

## The loop I actually saw

The fake model was scripted to ask for `lookup_activity` once and then answer. The runner emitted three events, in this order:

1. assistant message containing call `call-lookup-1` with `{"activity_id":"north-workshop","organization":"north","actor":"mira"}`;
2. tool result for `lookup_activity`, with the same call ID and `Room 12`;
3. final assistant answer, `Room 12 — Intro to Python (north).`

The call ID is the thread that keeps results attached to the right request. When I scripted the model to keep asking for the tool forever, the agent stopped after three attempts with `exceeds max iterations`. That is useful: a bad model or a bad prompt cannot spin indefinitely inside a user request.

## Validation is still our job

Eino can derive a JSON schema from a Go struct. For the draft experiment it produced required `title` and `capacity` fields and no `action` field. That describes the shape to a model; it does not enforce anything on the output. A model can still return broken JSON, a negative capacity or a stray `action`. Our validator uses strict decoding, rejects unknown fields, and checks title and capacity bounds. Ten bad drafts were rejected while one good draft passed. Nothing was created or published.

The same split applies to authorization. Nine tool calls were rejected — malformed JSON, missing actor, unknown field, wrong type, unknown ID, malformed ID, wrong tenant, a member instead of an organizer, and a wrong organization scope. The South workshop's `Room 9` was never disclosed to a North-scoped caller. That is a synthetic fixture, not proof of Core's production permissions; it just shows where the boundary belongs.

## Retry, deadline and cancel

I set the retry cap before running: at most two retries, three model attempts total, and only for a sentinel transient error. A permanent error stopped after one attempt. Eino's ADK handled that through its model retry config.

For deadlines, a scripted 200 ms reply behind a 25 ms context gave `context.DeadlineExceeded` and no answer. An explicit `adk.WithCancel()` also stopped the run. Eight concurrent requests with one shared runner each got their own answer. `go test -race` found nothing. This is correctness evidence on a tiny workload, not a capacity claim.

## Two things I got wrong

The public docs are a bit ahead of the pinned release. I wrote `model.BaseModel` as a plain interface; in v0.9.23 it is generic, and the build told me so. The fix was `model.BaseChatModel`, the alias for `BaseModel[*schema.Message]`. I also assumed `CancelHandle.Wait()` returned two values; it returns one. Both compiler errors are kept in the [iteration log](../../experiments/ai01/evidence/iterations.md) rather than quietly cleaned up. Lesson: read the pinned source, not just the current docs page.

## What this does not prove

The fake proves that the interfaces, loop, schema metadata and lifecycle hooks work. It says nothing about a real provider's tool-call reliability, latency, cost or model quality. No live call was made and no performance number is claimed. Core authorization, real cancellation through a provider, uncertain-write idempotency and Eino's larger ADK surface remain untested.

## Where I landed

For the narrow path, Eino is practical and cheap to use: real interfaces, a bounded tool loop, configurable retry/cancel, and no streaming machinery in our code. For A1's very first "one complete answer" slice, a plain Go adapter would also be enough. My suggestion is to keep Eino as the leading candidate but narrow its role — model/tool/agent plumbing only — and defer any dependency lock until a separately approved, credential-free hosted-adapter check confirms live tool-call behavior. The owner decides, and A1 is still its own authorization.

Reproduce the checks from `experiments/ai01/`: `gofmt -l .`, `go build ./...`, `go vet ./...`, `go test ./... -v`, `go test -race ./... -v`. Versions, pins, fixtures and raw output are in [PLAN.md](../../experiments/ai01/PLAN.md) and [iterations.md](../../experiments/ai01/evidence/iterations.md).
