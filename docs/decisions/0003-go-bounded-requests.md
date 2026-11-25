# 0003 — Go for the AI service, bounded requests first

**Date:** 2026-11-25. **Status:** Accepted language and initial execution direction; Eino is the leading candidate, pending AI01 validation. No implementation, provider spending, or deployment authorized.

## Context and alternatives

Campusly AI is a task-oriented assistant, not a general autonomous research agent. Initial inference is remote through a hosted API. Personal requests (chat, reads, drafts and separately confirmed actions) are intended to finish within bounded request budgets. Normal group messaging must never wait for inference. The owner wants deeper Go/cloud-native experience and familiarity with Chinese engineering ecosystems while maintaining a manageable product and research workflow.

Considered Go with Eino or custom orchestration, Python with AgentScope or Pydantic AI, and heavier workflow frameworks. Python remains useful for ML evaluation, data analysis and inference research; Go is chosen for the independent production AI service. Remote inference may dominate latency, so no claim of lower end-to-end latency or memory use is made without measurements. Using Go does not itself teach GPU inference optimization.

## Decisions

- Go is the production AI service language. Start with standard-library net/http and context for the HTTP boundary; add router or other dependencies only for a concrete need. Core remains Python/FastAPI; the existing Go messaging gateway is unchanged.
- Eino (CloudWeGo) is the first framework to validate, not yet a selected dependency. Retain Pydantic AI as the reference fallback; AgentScope and custom Go remain alternatives if evidence warrants reconsideration. Do not build three parallel products.
- Personal assistant uses bounded synchronous HTTP request-response initially. The request can wait for the complete model result; Go handlers can concurrently serve independent requests. Propagate cancellation/deadlines through Core, AI and provider calls. Bound input/output, tool iterations, concurrency, retry budget and provider quota usage. Do not hold a business transaction or lock during inference. Return a safe error on timeout; never retry an uncertain consequential write blindly. Future long-running workflows require a separate design decision.
- No token streaming. Initial inference is through an approved hosted model API, behind a replaceable adapter. Synthetic deterministic model first; no paid call or real-data transfer without explicit approval. Model/vendor, Go/Eino versions, identity/contracts and concrete limits remain to be decided in AI01/A1.
- No personal-task queue, local GPU serving, multi-agent workflow engine or dedicated model gateway in A1. Later shared group replies are asynchronous, durable, idempotent and isolated from ordinary messaging; their queue/storage design remains open until D2. No implication that all future work is synchronous.
- Docker Compose, initially one host with independent service boundaries. Python is available for offline evaluation/analysis, not a required production AI runtime.

## AI01: narrow validation before locking Eino

Use version-pinned released dependencies and synthetic fixtures. Predefine pass criteria and record exact versions, commands, failures and limits. Review official documentation and source for the chosen release. Start with Eino alone:

| Check | Required evidence |
| --- | --- |
| E01 complete response | One deterministic non-streaming response through the adapter; an actual hosted API smoke is separate and only after approval. |
| E02 deterministic testing | Tests run without provider credentials/network inference; controllable errors and outputs. |
| E03 typed tool call | Synthetic activity lookup, invalid arguments rejected, bounded loop; no actual Campusly business writes. |
| E04 structured draft | Validate a small synthetic activity draft; handle malformed output without publishing anything. |
| E05 request lifecycle | Concurrent synthetic requests remain isolated; cancellation/deadline and bounded retry work; document observed local resource usage only if measured. |
| E06 coding reliability | Record actual incorrect API assumptions, compile/test failures and fixes; use official docs and pinned source. No unsupported claims about training data or comparative agent reliability. |

A failure is evidence, not an instruction to hide it or automatically migrate. Compare a reference alternative only if meaningful limitations persist. Any comparative experiment must use equivalent tasks, model/config, documentation access and controls. Keep the spike limited to AI01, outside production integration. Obtain explicit authorization before executing it, selecting dependencies, making paid requests or starting A1.

## Documentation-driven implementation

Before writing framework-dependent code, consult official release-compatible docs/examples and relevant maintainer-provided coding skills if available. Verify uncertain APIs against the pinned module source; do not invent methods. Record the documentation link, release/tag and dependency pins. For Go changes, use gofmt, focused go test, go vet and a build; use race-enabled tests for concurrency where supported. Full repository suite remains separately authorized. Review semantic correctness, authorization, privacy, cancellation and failure paths; compilation alone is not proof. Explain key mechanisms to the owner from first principles and record actual findings in living report/blog drafts.

Official starting points: [Eino](https://github.com/cloudwego/eino), [Eino docs](https://www.cloudwego.io/docs/eino/), [Eino extensions and skills](https://github.com/cloudwego/eino-ext), [Eino examples](https://github.com/cloudwego/eino-examples), [Go net/http](https://pkg.go.dev/net/http), [Go context](https://pkg.go.dev/context), [Pydantic AI](https://ai.pydantic.dev/), [AgentScope](https://github.com/agentscope-ai/agentscope). Recheck availability and version applicability when executing the validation.

## Delegation

Follow the [AI01 validation brief](../engineering/AI01_VALIDATION_BRIEF.md). The owner or designated coding agent executes locally; the architecture/review assistant prepares the brief and reviews returned evidence. No CI workflows or environment setup by the planning assistant.

## Current evidence

Planning rationale only. No framework spike, benchmark, implementation or provider call performed. AI01 remains pending; A1 and all later slices remain NOT_STARTED. Update the relevant ML report/blog with the rationale when those drafts are opened; no fabricated results.
