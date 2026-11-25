# Engineering rules

The owner approved the product direction and a standalone public repository. Each implementation milestone needs an approved design and explicit scope. Local acceptance is the owner's decision. Real-data processing, recurring jobs, external rollout, spending and publication require separate permission.

Core checks current tenant, membership, activity and action capability on every tool read and write. A draft, memory entry, prompt or service credential cannot confer a user's authority. Bind approval to the exact actor, action, target and payload. An edited draft needs fresh approval. Post shared replies as the bot and approved drafts as the human with an internal receipt.

Keep private inquiries outside shared AI delivery, retrieval, caches, prompts, logs and telemetry. Group conversations live in core. Store only necessary AI-owned state, isolate scopes and define deletion before persistence. Synthetic fixtures must work without importing private source or contacting a provider.

Start from an organizer's problem. Record a decision and a reproducible check as the work happens. Register hypotheses before formal measurement and retain adverse findings. This repository is intentionally public as a portfolio and product demonstration for prospective customers, supervisors and employers. Review ownership, licensing, privacy, dependencies, documentation and changes before every public commit; keep private core source, secrets, identifiable data and sensitive operational details out of public history. Owner approval is required for product rollout, external blog/paper publication and real-data processing.

## Product-first development and deployment

Deliver the smallest useful, secure AI prototype on a single host alongside core. Keep independently runnable service/worker boundaries and versioned interfaces so HTTP, workers and realtime delivery can later scale separately; no distributed deployment is required for the prototype. Owner local review, relevant tests and a separate rollout decision precede shipping. Scale-out is a later engineering phase, not a blocker for the local AI slice. Do not let research ambitions delay academic deadlines.

## Living evidence

Engineering decisions, debugging, diagrams, tests and actual results feed incremental technical blog and Cloud/ML report drafts as work happens. The work is driven by product needs and approved engineering objectives, not by a predetermined report or paper conclusion. Every substantive task updates the relevant drafts or records an explicit N/A reason. Formal experiments freeze hypotheses, controls and scoring before measurement; papers follow only if evidence warrants them. Track owner understanding and report review independently of engineering verification.
