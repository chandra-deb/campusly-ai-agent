# Proposed local architecture

```mermaid
flowchart LR
  UI[Organizer workspace] --> API[Campusly FastAPI]
  API --> DB[(Canonical activity, chat and outbox)]
  DB --> Live[Existing dispatcher and Redis]
  Live --> Go[WebSocket gateway]
  API -. private assistant request .-> AI[AI service]
  DB -. eligible event via durable handoff .-> AI
  AI --> Model[Replaceable model adapter]
  AI -->|Scoped tool calls and proposals| API
```

Solid core messaging edges exist. The personal request path below is implemented as a local prototype. The other AI links are proposals. Core owns the integration endpoints and organizer UI. This repo owns orchestration, a versioned client, synthetic counterpart, evaluation and worker lifecycle. The gateway stays in core until the Cloud phase.

Propose versioned request, proposal, receipt and event schemas. Include event ID, version, organization, group and canonical source reference. Reject unknown versions. Core filters eligibility before handoff. Check live visibility again before retrieval and posting. Keep context scoped to actor and audience; never copy a group's timeline into personal history. A personal, conversation, group or organization memory has its own owner, readers, provenance and deletion policy. Promotion between scopes needs review.

The existing dispatcher claims outbox rows with SKIP LOCKED and records Redis publication. A durable AI consumer needs separate progress. Discuss a core-owned job record committed with the eligible event and lease/ack APIs. On timeout reclaim the lease. Retry with bounded backoff, quarantine exhausted jobs, key bot posts by source identity and return the original receipt after uncertain success. Never hold a business transaction open during inference.

For later scaling, keep request handlers stateless and partition worker capacity independently of HTTP and sockets. Persist claims and receipts outside process memory, bound concurrency and share provider quotas across replicas. Compare one, two and four workers on identical synthetic arrivals, recording queue age, completion latency, quota errors, duplicates and recovery after worker loss. Formal measurement needs separate authorization.

## Incremental delivery (ADR 0002)

The diagram describes eventual boundaries, not A1 requirements. A1 needs only president-scoped UI, Core, AI service and model adapter with one complete non-streamed response and basic failure/loading states. Use a synthetic counterpart for standalone tests. Persistence, reads, confirmed actions and personal memory follow individual slices. Durable group handoff/worker design is reviewed before D2, not required for A1. Independent service boundaries and single-host deployment remain. [ADR 0003](../decisions/0003-go-bounded-requests.md) selects Go and bounded synchronous personal requests; Eino is narrowly adopted for the model/tool/agent boundary, while provider, concrete state, queue and exact contract choices remain open.

## A1 local personal request boundary

```mermaid
sequenceDiagram
  participant U as President UI
  participant C as Campusly Core
  participant A as Go AI service
  participant M as Model adapter
  U->>C: Request (loading state)
  C->>C: Authenticate and authorize
  C->>A: Signed v1 request with actor, organization and expiry
  A->>M: Generate complete response
  M-->>A: Complete output or error
  A-->>C: Complete scoped result or safe error
  C->>C: Verify reply scope and recheck presidency
  C-->>U: Complete response
```

A1 uses net/http and context, with a deterministic fake adapter for credential-free tests. Hosted inference requires separate approval. Propagate cancellation and deadlines, bound concurrent work and retries, and validate output. Do not hold a core database transaction during inference. No personal job queue or streaming. Later group jobs use a separately designed durable asynchronous handoff; ordinary messaging never waits for AI. See [ADR 0004](../decisions/0004-personal-assistant-boundary.md) for the contract and limits.


## A3 stored turns

```mermaid
sequenceDiagram
  participant H as H5 assistant
  participant C as Core
  participant D as PostgreSQL
  participant A as Go assistant
  H->>C: Send in owned conversation with turn UUID
  C->>D: Authorize and read recent turns
  C->>D: End read transaction
  C->>A: Signed v2 request with up to two turns
  A-->>C: Complete scoped reply
  C->>D: Recheck presidency and lock conversation
  alt Matching turn already saved
    D-->>C: Existing answer
  else Same conversation revision
    C->>D: Insert prompt and reply, then commit
  else Deleted or advanced
    C-->>H: Safe error
  end
  C-->>H: Canonical turn on success
```

The conversation is owned by the student within the organization. A
database cascade removes its turns on hard delete. The Go service
processes the signed context in memory for the duration of the call.
See [ADR 0005](../decisions/0005-personal-conversation-retention.md).

## B1 live organization context

```mermaid
sequenceDiagram
  participant H as H5 personal chat
  participant C as Core
  participant A as Local Go service
  participant D as Campusly database
  H->>C: Ask in the current organization
  C->>C: Check presidency and conversation owner
  C->>A: Signed v3 intent request
  A-->>C: Scoped managed-organizations intent
  C->>D: Recheck presidency and read current roles
  C-->>H: Complete answer with organization sources
```

The v3 intent exchange has a one-second Core budget. A plain chat intent
uses the existing v2 model exchange. Campusly source cards disappear when
the user switches chats or reloads the transcript. Core rebuilds the answer
on the next question using the then-current authorization. Go receives no
organization rows and has no credentials to query Core.

## B5 confirmed activity creation

```mermaid
sequenceDiagram
  participant H as H5 president chat
  participant C as Core
  participant D as PostgreSQL
  H->>C: Review action ID and exact draft fields
  C->>C: Check current presidency and chat ownership
  C-->>H: Exact preview and expiring signed approval
  H->>C: Confirm the same action ID, draft and approval
  C->>C: Recheck scope and verify approval
  C->>D: Check for a receipt under the president lock
  alt New and unexpired
    C->>D: Commit draft activity, discussion and receipt together
  else Same action already committed
    C->>D: Read the original receipt
  end
  C-->>H: Canonical activity receipt
```

The confirmation path does not call Go. Receipts contain the actor,
organization, action ID, target and keyed payload fingerprint. Deleting
a private chat removes its transcript. The operational audit stays with
Core and is scoped to the actor and organization.

B7 uses the same confirmation path for official notices. The payload
names the canonical discussion and Core checks membership and president
capability before committing the public notice, its outbox events and
the receipt. Manual discussion paths use the default transaction
behavior. H5 opens the confirmed post by its message ID.

B8 uses that same room and receipt path for a public text message.
It sends the president's approved text through the existing chat
service and creates only its ordinary message event. Private
inquiries stay in the discussion's existing private composer.

## C1 recent context

Core reads a bounded suffix of whole turns from the owned transcript
and releases its read transaction before inference. The signed v4
reply contract checks the number of turns and UTF-8 text size again in
Go. The full saved transcript stays in Core. See
[ADR 0007](../decisions/0007-personal-memory.md).

## C2–C3 personal preferences

```mermaid
sequenceDiagram
  participant H as H5 assistant
  participant C as Core
  participant D as PostgreSQL
  participant A as Go assistant
  H->>C: Save a supported reply preference
  C->>D: Recheck presidency and save for owner and organization
  H->>C: Ask for a workshop checklist in another chat
  C->>D: Read current preferences in the selected scope
  C->>A: Signed v5 reply request with approved style values
  A-->>C: Complete reply
  C->>D: Recheck presidency, conversation and preferences
  alt Preferences still match
    C-->>H: Saved reply
  else Preference removed or changed
    C-->>H: Retry with current choices
  end
```

The intent request for live Campusly reads carries the prompt only.
Using or editing a preference does not confer access to any activity or
discussion. See [ADR 0007](../decisions/0007-personal-memory.md).
