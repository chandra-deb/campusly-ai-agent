# Campusly AI Agent

Campus organizers need to check activity details, draft messages and follow up in the discussions they already use. This repository owns the proposed assistant and its synthetic integration. Campusly core owns permissions, official actions and the human timeline.

## Getting started

Start with the [product specification](docs/specs/CAMPUSLY_AI_PRODUCT_UX_SPEC_v1.0.md), [architecture](docs/engineering/ARCHITECTURE.md) and [active milestones](docs/engineering/MILESTONES.md). Walk through the [synthetic scenario](examples/workshop.md) without credentials or a core checkout. The engineering foundation and the isolated AI01 spike exist today. Build, container, local run commands and the first president-only assistant UI are deliverables of A1, once separately approved. Each user-facing slice is manually testable and delivers complete, non-streamed responses. Personal memory is planned before shared group AI.

## Technical direction

Go is the selected production AI service language, with standard-library HTTP initially and complete, bounded personal requests. AI01 validated Eino locally and the owner chose narrow adoption for the model/tool/agent boundary; the dependency stays scoped to the AI service and a hard lock awaits a separately approved adapter check. A1 is not authorized and no hosted provider is selected. Python remains useful for offline ML evaluation. Shared group processing will be asynchronous and independent of normal messaging. See [ADR 0003](docs/decisions/0003-go-bounded-requests.md) for rationale, validation plan and guardrails.

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
