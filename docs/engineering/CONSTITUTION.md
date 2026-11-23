# Engineering rules

The owner approved the product direction and a standalone repository. Each implementation milestone needs an approved design and explicit scope. Local acceptance is the owner's decision. Real-data processing, recurring jobs, external rollout, spending and publication require separate permission.

Core checks current tenant, membership, activity and action capability on every tool read and write. A draft, memory entry, prompt or service credential cannot confer a user's authority. Bind approval to the exact actor, action, target and payload. An edited draft needs fresh approval. Post shared replies as the bot and approved drafts as the human with an internal receipt.

Keep private inquiries outside shared AI delivery, retrieval, caches, prompts, logs and telemetry. Group conversations live in core. Store only necessary AI-owned state, isolate scopes and define deletion before persistence. Synthetic fixtures must work without importing private source or contacting a provider.

Start from an organizer's problem. Record a decision and a reproducible check as the work happens. Register hypotheses before formal measurement and retain adverse findings. Before public release review ownership, licensing, privacy, dependencies, documentation and every commit, including deleted files. The owner decides whether to publish.
