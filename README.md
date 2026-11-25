# Campusly AI Agent

Campus organizers need to check activity details, draft messages and follow up in the discussions they already use. This repository owns the proposed assistant and its synthetic integration. Campusly core owns permissions, official actions and the human timeline.

## Getting started

Start with the [product specification](docs/specs/CAMPUSLY_AI_PRODUCT_UX_SPEC_v1.0.md), [architecture](docs/engineering/ARCHITECTURE.md) and [active milestones](docs/engineering/MILESTONES.md). Walk through the [synthetic scenario](examples/workshop.md) without credentials or a core checkout. Only the engineering foundation exists today. Build, container, local run commands and the first president-only assistant UI are deliverables of A1, once approved after AI01. Each user-facing slice is manually testable and delivers complete, non-streamed responses. Personal memory is planned before shared group AI.

This is an intentionally public, product-first engineering portfolio. The initial prototype targets one host with independently runnable AI and core boundaries; later scale-out is a separate phase. Technical blogs and living course reports capture the decisions, failures and verified results along the way. No production rollout or measured AI capability is claimed here.

## Technical direction

Go is the selected production AI service language, with standard-library HTTP initially and complete, bounded personal requests. Eino is the leading candidate pending AI01 synthetic validation; no framework dependency or hosted provider is selected yet. Python remains useful for offline ML evaluation. Shared group processing will be asynchronous and independent of normal messaging. See [ADR 0003](docs/decisions/0003-go-bounded-requests.md) for rationale, validation plan and guardrails.
