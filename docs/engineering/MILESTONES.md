# AI milestones — owner-agreed roadmap

**Decision: 2026-11-25.** Visible vertical slices, president-only personal assistant first, complete non-streamed responses. Product direction is agreed. Go and the initial execution direction are selected in [ADR 0003](../decisions/0003-go-bounded-requests.md); Eino validation and implementation scope remain separate.

| Phase | Slices | User-visible goal |
| --- | --- | --- |
| A — Personal assistant foundation | A1–A4 | Chat privately, receive complete responses, manage conversations |
| B — Campusly-aware assistant | B1–B8 | Authorized information, activity/notice drafts and confirmed actions |
| C — Personal memory | C1–C5 | Owner-controlled recall and management |
| D — Shared group AI | D1–D6 | Identity, mentions, grounding, selective replies, reliability, summaries |
| E — Organization intelligence | E1–E8 | Later provisional increments, refined before implementation |

See [TASKS.md](TASKS.md) for each slice and test. AI01 is tested and owner-reviewed; narrow Eino adoption is decided. A1–A4 are implemented and locally tested. The owner accept-or-fix decisions are pending. [Review Phase A](../evidence/A4.md). [B1–B8](../evidence/PHASE_B.md) are locally tested and ready for owner review. [C1](../evidence/C1.md) is locally tested and ready for owner review. C2–C5, Phases D–E and AI17 are NOT_STARTED.

Each user-facing slice delivers UI, core/API and AI together, focused checks and owner manual accept/fix before completion. Small infrastructure tasks may support a named slice. Use loading and error states, not token streaming; asynchronous jobs may return complete results later. Recommended order A → B → C → D → E; personal memory precedes shared AI, but may start before all B tools are done. Technical prerequisites govern independent work. Initial internal prototype after A; a useful president assistant can be separately released after selected B slices. Owner acceptance and separate rollout authorization required for every release. Single host with independent service boundaries first.

AI17 retains the existing cross-repository integrated acceptance checkpoint after B8, D4, D5 and relevant security/integration checks. It is not the first manual test or a requirement for a separately approved personal release. CAMS infrastructure scaling still requires AI17 and separate authorization unless its plan is revised. Cloud submission is independent of these gates. Keep living reports and blogs current, capture sanitized screenshots, and protect academic deadlines. No further dependency/provider selection, implementation, deployment or recurring work approved here.

AI01 complete: the owner reviewed the local Eino evidence and chose [narrow adoption](../decisions/0003-go-bounded-requests.md) — Eino for model/tool/agent plumbing if A1 is authorized, dependency scoped to the AI service, hard lock deferred until a separately approved adapter check. A1 uses the deterministic adapter and scoped contract in [ADR 0004](../decisions/0004-personal-assistant-boundary.md). No hosted-provider call, publication or rollout approved. Personal assistant starts with bounded synchronous HTTP and hosted inference behind a fake-testable adapter; group jobs remain asynchronous later.
