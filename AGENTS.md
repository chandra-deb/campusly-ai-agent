# Working in Campusly AI

Read the product spec, constitution, architecture, milestone and task queue before changing code. Inspect Git status first. Implement only the approved task. Use synthetic inputs, versioned interfaces and focused positive and negative tests. Keep real student data, private core source and secrets outside this repository. Production changes, recurring jobs and publication need separate approval.

## Documentation & Git Commit Style — Mandatory

Write documentation, ADRs, READMEs, engineering notes and commit messages in a natural, concise student-developer voice. Explain decisions, alternatives and trade-offs in plain language. Avoid marketing, boilerplate and repetitive summaries. Make one logical change per commit and describe what changed, for example `Document agent authorization boundaries`. Never invent personal experience, tests, findings or progress. Follow applicable disclosure requirements and omit unnecessary attribution. Keep a history the owner can maintain and explain.

## Engineering-First Research Philosophy — Mandatory

Start with a real problem. Design and implement a solution, verify and measure it, then document findings. Keep decisions, reproducible checks, failures and limitations as engineering evidence. Define hypotheses before formal experiments. Add complexity only for product needs or approved engineering objectives. Pursue a paper when findings justify it. Reports synthesize the evidence. Mandatory scaling, security and academic deliverables still apply.

## Continuous Engineering Evidence — Mandatory

Product value drives the work; blogs, course reports and research papers are outputs of engineering, not reasons to add unnecessary features. Before implementation, state the problem, alternatives and relevant evidence or hypothesis. During implementation, record decisions, failures, commands, revisions and limitations. At each task or milestone, update the relevant living course-report draft and technical blog draft (or explicitly record N/A and why), alongside tests and diagrams. Preserve raw results and reproducibility; never invent personal experience or measurements. Owner reviews and revises academic submissions and approves publication separately. Protect external course deadlines and master's preparation.

## Visible vertical slices — mandatory

Follow the active A–E roadmap and ADR 0002. Start with the president-only personal assistant. Every user-facing slice includes a locally usable UI, integrated backend, focused checks and explicit owner accept/fix decision. Capture sanitized screenshots and feedback. Return complete responses only; no token streaming. Personal memory precedes shared group AI in the recommended order. Do not implement future infrastructure prematurely. Preserve AI01/AI17 cross-repository references.

## Runtime and verified framework usage — mandatory

Read [ADR 0003](docs/decisions/0003-go-bounded-requests.md). Production AI is Go with standard-library HTTP initially; Eino is pending validation, not a dependency approval. Python is available for offline evaluation. Personal requests are bounded synchronous HTTP with complete responses, propagated cancellation and deadlines; no streaming or initial personal-job queue. Shared group inference later is durable and asynchronous, never on the normal messaging path. Before writing framework-dependent code, consult official release-specific docs/examples and relevant maintainer skills, verify APIs against pinned source, record versions and links, then run formatting, build, focused tests and vet (race checks where appropriate). Do not guess APIs or claim correctness merely from compiling. AI01 evaluation and A1 implementation require separate authorization; no paid provider calls or real data without approval.

## Role boundaries

The planning/review assistant researches, teaches, designs, prepares briefs and reviews evidence; it edits planning documents only when requested. The owner or an explicitly delegated coding agent writes code and runs local tests/builds in the owner's environment. The planning assistant must not create experimental code, install dependencies, run tests or create/trigger CI. Follow the constitution's responsibility boundaries. For AI01 use [the validation brief](docs/engineering/AI01_VALIDATION_BRIEF.md); no CI, A1 or provider calls are authorized by it.
