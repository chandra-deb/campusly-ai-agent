# AI01 — independent verification note

**Date:** 2026-11-25. **Branch:** `experiment/ai01-eino-validation`, merged to `main` in PR #1 (commit `7c51499`). Nothing was published or deployed.

This is a short verification note, not a second validation record. The canonical pre-execution plan is [experiments/ai01/PLAN.md](../../experiments/ai01/PLAN.md), the commands and failures are in [evidence/iterations.md](../../experiments/ai01/evidence/iterations.md), and the narrative drafts are the [ML report](../reports/ML_REPORT_DRAFT.md) and [engineering blog](../blog-drafts/2026-11-25-ai01-eino-validation.md).

## What I found and what I did

When I started reviewing, the isolated spike and its evidence were already present in the working tree (files created between 15:35 and 15:48 on 2026-11-25 by another local execution in the same checkout). I did not overwrite any of that work. I reviewed the implementation and re-ran its checks independently so the handback is backed by commands I actually ran.

I reviewed `internal/fakemodel`, `internal/lookup`, `internal/draft`, `internal/agentrun` and the E01–E05 tests. The fake model implements the real `model.ToolCallingChatModel` interface, treats `Stream` as a test failure, and keys scripted replies by prompt. The lookup tool implements `tool.InvokableTool` with strict JSON decoding and checks all three scope fields before returning a venue. The draft validator rejects unknown fields and out-of-range capacity in application code. The negative tests assert real errors and check that no venue leaks; they are not vacuous.

## Independent re-run

Run from `experiments/ai01/` with `PATH=/tmp/opencode/go-sdk/go/bin:$PATH` and `GOTOOLCHAIN=local`:

| Command | Outcome | Raw output |
| --- | --- | --- |
| `go version` | `go1.25.14 linux/amd64` | — |
| `gofmt -l .` | no files listed, exit 0 | [verify-gofmt.txt](../../experiments/ai01/evidence/verify-gofmt.txt) |
| `go build ./...` | exit 0 | [verify-build.txt](../../experiments/ai01/evidence/verify-build.txt) |
| `go vet ./...` | exit 0, no findings | [verify-vet.txt](../../experiments/ai01/evidence/verify-vet.txt) |
| `go test -count=1 ./... -v` | exit 0; 36 `RUN` and 36 `PASS` lines | [verify-test-uncached.txt](../../experiments/ai01/evidence/verify-test-uncached.txt) |
| `go test -race -count=1 ./... -v` | exit 0; no race reports | [verify-race-uncached.txt](../../experiments/ai01/evidence/verify-race-uncached.txt) |
| `go mod verify` | `all modules verified` | [verify-mod.txt](../../experiments/ai01/evidence/verify-mod.txt) |

Pins re-checked: `go.mod` SHA-256 `736574b4…a154bd7`, `go.sum` SHA-256 `fe540fa3…83aa`, matching the values recorded in `iterations.md`. Eino `v0.9.23` checksums match `go mod download -json`. I did not run the full repository suite; this spike is the only Go module in the checkout.

A second Go toolchain exists at `/tmp/cams-local-toolchain/go` (`go1.24.10`). It was not used because the spike module declares `go 1.25`; using it would have silently changed the toolchain and made the re-run non-reproducible.

## Corrected limits and status

The limits actually used differ from an earlier draft note of mine: retry cap is `MaxRetries = 2` (three total attempts), agent `MaxIterations = 3`, eight concurrent requests, a 25 ms deadline against a 200 ms fake delay, and a 5 ms default fake delay in the isolation test. [PLAN.md](../../experiments/ai01/PLAN.md) is the authoritative record of those choices.

E01–E06 were **TESTED locally** as recorded in the report and iteration log. The owner reviewed the evidence on 2026-11-25 and chose **narrow adoption**: keep Eino for the model/tool/agent boundary if A1 is authorized, keep the dependency scoped to the AI service, and defer a hard dependency lock until a separately approved adapter check. This is not A1 authorization or a full dependency lock. No provider call, production permissions check, performance measurement or publication has happened.
