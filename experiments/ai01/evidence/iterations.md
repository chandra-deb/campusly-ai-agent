# AI01 iteration log

**Date:** 2026-11-25. **Branch:** `experiment/ai01-eino-validation`. **Revision at execution:** `8e8fc0e529c7a953ecc1fa2178b7da261813f21d`. All changes are local and uncommitted; nothing was pushed, deployed or published. Preconditions are in [PLAN.md](../PLAN.md).

## Environment

Go was not on PATH. I downloaded the pinned official tarball into `/tmp/opencode/go-sdk` and unpacked it there; no system package or sudo was used.

```
curl -fsSL https://go.dev/dl/go1.25.14.linux-amd64.tar.gz -o go1.25.14.linux-amd64.tar.gz
echo 'a21ae5633a269bcd7e90cf767e48225633795e99d831742cbf3397064fee7712  go1.25.14.linux-amd64.tar.gz' | sha256sum -c -
tar -xzf go1.25.14.linux-amd64.tar.gz -C /tmp/opencode/go-sdk
/tmp/opencode/go-sdk/go/bin/go version  # go version go1.25.14 linux/amd64
```

`GOTOOLCHAIN=local` for every Go command so the toolchain cannot silently change. `go mod init campusly.local/ai01`, `go mod edit -go=1.25`, `go get github.com/cloudwego/eino@v0.9.23`, then `go mod tidy`.

- `go.mod`: SHA-256 `736574b417d54c1f9e42f4e74f522021c074adac3cb9d688dded7c0b2a154bd7`
- `go.sum`: SHA-256 `fe540fa3dbc2fb9720fc9252447ada68d61ebff5312c7190bb06aee8c71683aa`. The first `go.sum` before `tidy` had hash `ea38ffc5f66970d0c2c2489dd65b7f3ecaac352437bcf3a0ad3aaf93acbe691e`; the final hash is the pinned record.
- `github.com/cloudwego/eino v0.9.23 h1:f369x1HXckTj6SlMQccgyvKWr+QMMr2yxunqTAOAF1E=`
- `github.com/cloudwego/eino v0.9.23/go.mod h1:OBD1mrkfkt/pJa4rkg1P0VnaMeOVl7l8IAdEqY//3IQ=`
- Transitive pins that affect schema/JSON: `github.com/eino-contrib/jsonschema v1.0.3`, `github.com/bytedance/sonic v1.15.0`. Full graph is in `go.mod`/`go.sum`.

## Commands and raw output

Run from `experiments/ai01/`. Raw files are kept next to this log.

| Command | Outcome | Evidence |
| --- | --- | --- |
| `gofmt -l .` | Three files needed formatting; `gofmt -w .` applied, second `gofmt -l .` empty | [gofmt.txt](gofmt.txt) |
| `go build ./...` | First attempt failed (see below); fixed build exit 0 | [build-first-failure.txt](build-first-failure.txt), [build.txt](build.txt) |
| `go vet ./...` | exit 0, no findings | [vet.txt](vet.txt) |
| `go test ./... -v` | First attempt failed to compile (see below); after fix all 36 test/subtest `RUN` lines pass, `ok campusly.local/ai01 0.050s` | [test-first-failure.txt](test-first-failure.txt), [test.txt](test.txt) |
| `go test -race ./... -v` | exit 0; no race reports | [test-race.txt](test-race.txt) |

No full repository suite was run: this spike is the only Go module in the checkout. No CI workflow, cloud execution, provider call, real data, Core/UI integration or deployment was involved.

## Failures and corrections

1. `model.BaseModel` is generic in v0.9.23. I had written `func NewAgent(..., m model.BaseModel, ...)` from the older non-generic docs. Build error: `cannot use generic type model.BaseModel[M model.messageType] without instantiation`. Checked `components/model/interface.go` at the pinned tag and changed the parameter to `model.BaseChatModel`, the alias for `BaseModel[*schema.Message]`.
2. `adk.CancelHandle.Wait()` returns one value (`error`), not two. Compile error at `e05_lifecycle_test.go:105: assignment mismatch: 2 variables but handle.Wait returns 1 value`. Checked `adk/cancel.go` at v0.9.23 and corrected the call.
3. The first `gofmt -l .` listed `e03_tool_test.go`, `e04_draft_test.go` and `internal/fixture/fixture.go`. Formatted and confirmed clean.

Both compiler failures are kept as evidence rather than edited out of the history. No test assertion failed after compilation succeeded.

## Exercise observations

- **E01 PASS.** `TestE01CompleteResponse` calls `Generate` once with the fixed prompt and receives the exact string `Room 12 — Intro to Python (north).`, no tool calls, `StreamCalls() == 0`. The fake's `Stream` deliberately returns an error, so any streaming path would fail loudly.
- **E02 PASS.** Three repeated calls return `reply-1`, `reply-2`, `reply-3`. With `MaxRetries = 2` and `IsRetryAble` restricted to `ErrTransient`, two transient failures then success take exactly three attempts; a permanent failure takes exactly one attempt and surfaces `ErrPermanent`. `TestE02NoCredentialsNeeded` uses an empty fake and fails if no error appears. Nothing reads credentials or opens a network connection.
- **E03 PASS.** `Info()` exposes `lookup_activity` with `activity_id`, `organization`, `actor` all required strings; `ToJSONSchema()` was inspected in the test. The ADK run produced three events in order: assistant message with tool call `call-lookup-1`, tool result for `lookup_activity` with matching `ToolCallID` and venue `Room 12`, then the final assistant answer. The model received the tool schema. Nine negative invocations were rejected (malformed JSON, missing actor, unknown field, wrong type, unknown ID, malformed ID, wrong tenant, unauthorized actor, unauthorized organization) with no venue disclosed. A scripted endless tool loop stopped at three model attempts with `exceeds max iterations`.
- **E04 PASS.** Eino's `utils.GoStruct2ParamsOneOf[Draft]()` + `ToJSONSchema()` produce required `title`/`capacity` properties and no `action` property. Application-owned `Validate` accepted one valid draft and rejected ten cases: malformed JSON, missing fields, wrong type, zero/negative/above-bound capacity, blank title, `action` field and trailing JSON. No create or publish call exists in the module.
- **E05 PASS.** Eight concurrent requests through one shared runner returned their own `answer-N` with no cross-talk; `StreamCalls() == 0`. A 25 ms deadline against a 200 ms scripted delay returned `context.DeadlineExceeded` with an empty answer and no synthetic content in the error. `adk.WithCancel()` produced a cancellation error and no answer. Four concurrent retrying requests each used three attempts. `-race` found no data race.
- **E06 (recorded here).** Two API assumptions were wrong and fixed; see above. `gofmt`, build, `go vet`, focused `go test` and `go test -race` all pass. About 1,048 lines of Go including tests: `fakemodel` 177, `lookup` 69, `draft` 67, `fixture` 68, `agentrun` 88, tests 575, `doc.go` 4. The spike has no adapter for a real provider, no HTTP handler, no persistence and no Core integration.

## Limitations and unresolved questions

- The fake model proves interface and loop mechanics. It says nothing about a hosted model's tool-call reliability, latency, cost, token limits, schema adherence or rate-limit behavior. No live provider compatibility check was authorized or performed.
- Cancellation and deadline behavior was tested with an in-process delay. Real cancellation across Core → AI → provider is not tested here.
- The authorization fixture is synthetic. It does not establish Core's production permissions; Core still decides at execution time.
- Retries cover model-call errors only. Uncertain consequential writes and idempotency remain application/Core concerns outside AI01.
- Eino v0.9.23 exposes a much larger ADK surface (middleware, checkpoints, multi-agent, agentic messages, TurnLoop). Only the small path needed by A1 was exercised. Complexity of the untested surface is unknown.
- Public docs were last edited in March 2026 and are ahead of or beside the pinned tag in places. The pinned module source was the tie-breaker; the two compiler errors above came from that drift.
- No performance or capacity claim: eight requests with 5 ms fake delays is a correctness/race check, not a benchmark.

Status: **TESTED locally, AWAITING OWNER REVIEW**. Eino is not approved as a dependency, A1 is not authorized, and no provider call, publication or rollout has happened.
