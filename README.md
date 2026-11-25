# Campusly AI Agent

Campus organizers need to check activity details, draft messages and follow up in the discussions they already use. This repository owns the Go assistant and its synthetic integration. Campusly core owns permissions, official actions and the human timeline.

## Getting started

Start with the [product specification](docs/specs/CAMPUSLY_AI_PRODUCT_UX_SPEC_v1.0.md), [architecture](docs/engineering/ARCHITECTURE.md) and [active milestones](docs/engineering/MILESTONES.md). Walk through the [synthetic scenario](examples/workshop.md) without credentials or a core checkout. The Phase A branches include the first president-only H5 assistant, with complete replies, scoped history and conversation controls. The disposable local lab uses a PostgreSQL container and runs H5, Core and Go as loopback processes. The [Phase C review guide](docs/evidence/PHASE_C.md) walks through the locally tested personal memory slices.

## Technical direction

Go is the selected production AI service language, with standard-library HTTP initially and complete, bounded personal requests. AI01 validated Eino locally and the owner chose narrow adoption for the model/tool/agent boundary; the dependency stays scoped to the AI service and a hard lock awaits a separately approved adapter check. The Phase A lab uses a deterministic model adapter. Python remains useful for offline ML evaluation. Shared group processing will be asynchronous and independent of normal messaging. See [ADR 0003](docs/decisions/0003-go-bounded-requests.md) for rationale, validation plan and guardrails.

## Current next milestone

A–C have deterministic local evidence. The real-model gate R1–R4 is next, before D1: configurable hosted inference, all existing business capabilities, conversation/memory and integrated owner testing on the existing local dummy-data Campusly database. Keep fake mode for regression. The restricted [Cerebras smoke](docs/evidence/CEREBRAS_SMOKE.md) was blocked and is not a general real-model adapter. Follow [milestones](docs/engineering/MILESTONES.md) and [tasks](docs/engineering/TASKS.md) for the active sequence.

## Run the local assistant

The Go service and its fixture use the provisional Eino v0.9.23 adapter. From
`service/`, set a fresh local key and start the loopback server.

```sh
export CAMPUSLY_AI_KEY="$(openssl rand -hex 32)"
go run ./cmd/assistant
```

In another terminal with the same key, run `go run ./cmd/fixture`. It sends a
signed synthetic workshop question through the HTTP endpoint. For the full H5
path, use the synthetic lab documented in [the A1 evidence](docs/evidence/A1.md).

### Review Phase A locally

Check out the same `feature/personal-assistant-a1` through
`feature/personal-assistant-a4` branch in the AI, CAMS and H5 repositories
for the slice you want to review. A4 includes a review mode for the whole
foundation. It starts the services and disposable PostgreSQL, prints the
loopback sign-in URL and synthetic credentials, and cleans up on Ctrl+C.

Build the two Go binaries from `service/` and then start the lab from CAMS.
Use the corresponding H5 worktree with its dependencies installed.

```sh
go build -buildvcs=false -o /tmp/campusly-assistant-bin/ ./cmd/...
```

```sh
CAMS_ENV_FILE=/dev/null python -m scripts.m2.assistant_lab \
  /tmp/campusly-assistant-bin /absolute/path/to/cams-h5 --review
```

Use [the slice evidence](docs/evidence/A4.md) for owner review scenarios
and automated verification commands.

### Review Phase B locally

Use the matching local branches in AI, CAMS Core and H5. The
[Phase B review guide](docs/evidence/PHASE_B.md) lists the branches,
synthetic lab command, owner scenarios and captured verification for B1–B8.
