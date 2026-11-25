# 0001 — Repository boundary and first design discussion

Accepted boundary: one public, independent AI repository (intentional portfolio and product demonstration; private core source and sensitive information stay outside), versioned contracts and a reproducible synthetic integration. No dependency on private backend imports. The options below are proposed.

| Choice | Starting proposal | Trade-off and decision test |
| --- | --- | --- |
| Runtime | Python service with an explicit sequence of typed tool calls | Reuses team skills and simplifies evaluation. Go offers a single binary and a familiar concurrency model. Compare packaging effort and cancellation handling before optimizing runtime overhead. |
| Orchestration | Small adapter around a bounded tool loop | Custom code keeps the approval path inspectable but owns validation and retries. Pydantic AI offers typed tools and outputs. LangGraph adds durable graph execution at the cost of another state model. Revisit for branching, resumable workflows. |
| Provider | Deterministic fake first, then one owner-approved adapter | Compare a hosted tool-capable model with a local model using the same synthetic cases. Evaluate tool schema compatibility, grounded abstention, latency, quotas, retention and cost. No vendor chosen. |
| Delivery | Core-owned durable job with lease, ack and retry | Simplifies atomic enqueue against PostgreSQL but adds polling and pool load. Redis Streams offers consumer groups but requires reliable outbox bridging, pending-entry recovery and retention. Pub/Sub has no durable acknowledgement. |
| State | Start with source references and minimal scoped proposal/receipt state | An AI-owned relational store could support conversations and reviewed memory. A separate database adds backup, retention and deletion coordination. Defer vector retrieval until known examples show a need. |
| Packaging | One service image with independently runnable worker and a mock core profile | Easier to reproduce than several services per capability. Decide local commands and dependency pins with the runtime. |

First review a synthetic send with a crash after posting but before acknowledgement. Agree who issues the scoped credential, persists consent, checks expiry/revocation and returns the original receipt. Approve the interface and AI10 scope before selecting libraries.

Primary references: [Pydantic AI](https://ai.pydantic.dev/), [LangGraph](https://docs.langchain.com/oss/python/langgraph/overview), [Redis Streams](https://redis.io/docs/latest/develop/data-types/streams/).
