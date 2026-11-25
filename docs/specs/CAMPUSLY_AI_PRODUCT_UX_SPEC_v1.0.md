# Campusly AI — Product & UX Specification v1.0

**Document type:** Product and user-experience specification  
**Status:** **Product direction agreed for specification; not implementation approval**  
**Date:** 2026-11-23  
**Product:** Campusly — organizer-first activity and organization management  
**Intended readers:** Product owner, designers, engineering agents, backend/frontend engineers, ML evaluators  
**Suggested repository location:** `docs/specs/CAMPUSLY_AI_PRODUCT_UX_SPEC_v1.0.md`

> **Scope of approval.** This document consolidates the product and UX decisions reached in the November 2026 brainstorming conversation. “Agreed” means the intended **future user experience**, not that the feature has been implemented, tested, deployed, or approved as an engineering milestone. Exact schemas, AI providers, job frequency, permissions for proposed organization/team groups, and rollout phases require separate engineering review.
>
> **Precedence.** Existing approved Campusly activity-messaging behavior remains binding. This AI spec extends it; it must not silently rewrite the activity-group contract. Where implementation evidence or the original approved business spec disagrees with an assumption here, investigate and resolve explicitly before coding.

---

## 1. Executive summary

Campusly AI should turn Campusly into a **unified organization-management workspace** rather than adding an isolated chatbot. The organizer remains in control; the assistant helps them understand activities, communicate across groups, prepare authorized actions, and learn from the organization's past experience.

**Core interaction model:**

1. **Organizer home:** one place to notice important updates, recent conversations, and enter the personal assistant.
2. **Personal assistant:** a private, permission-bound interface to ask questions and coordinate work across the selected organization, activities, and accessible groups.
3. **Human groups:** the real place people collaborate. Normal chat continues to work without AI.
4. **Shared group AI:** a visible participant that can answer direct mentions and selectively answer question-mark messages.
5. **Institutional memory:** private conversation, group, and organization-scoped knowledge; periodic extraction from authorized group history, source references, and reviewed lessons from completed activities.
6. **Human approval:** all consequential operational writes, including messages sent as an organizer, require a clear review/approval step. Automatic, reversible *working-memory extraction* is a distinct opt-in background behavior; it does not grant permission to modify official data.

### 1.1 Product promise

> **Coordinate people, activities, and organizational knowledge from one place without losing the ability to work directly in groups.**

### 1.2 Goals

- Reduce repeated cross-group navigation, copying, and manual follow-up.
- Make organizer tasks simpler, not dependent on an AI conversation.
- Preserve real group communication while providing an AI-assisted path to it.
- Help new organizers benefit from verified lessons learned by previous teams.
- Provide useful, grounded answers with references when evidence is available.
- Maintain strict identity, permission, data-source, and action boundaries.
- Keep AI optional and failure-isolated from ordinary messaging.

### 1.3 Non-goals for the initial AI experience

- Replacing WeChat as a general-purpose social network, or replicating every messaging feature.
- Making AI the mandatory interface for any existing manual operation.
- Unattended administrative actions or messages sent without an organizer's approval.
- Continuous LLM analysis of every incoming message for proactive replies.
- Automatic publication of organization policies from casual group chat.
- Automatic daily AI briefings or unsolicited cross-group summaries.
- MCP/OpenCode integration as a prerequisite for the in-product experience. MCP may reuse the same authorized capability layer later.
- A full task/project-management suite, complex knowledge graph, or mandatory rigid group hierarchy in V1.

---

## 2. Status and terminology

| Term | Meaning |
| --- | --- |
| **Agreed** | Product behavior accepted in brainstorming; not proof of implementation |
| **Existing baseline** | Behavior already established by the approved activity messaging contract or reported implementation; verify actual code before relying on new integration points |
| **Proposed extension** | A new Campusly feature requiring separate design/implementation scope |
| **Working memory** | Reversible, attributed information derived from conversations; may be incomplete or unverified |
| **Confirmed organizational knowledge** | Reviewed guidance or lessons approved by a president or delegated authorized organizer |
| **Authoritative data** | Current Campusly records and business rules returned by authorized FastAPI APIs; not inferred from memory |
| **Action receipt** | Auditable record linking the requesting user, approved action, result, and target resource |

### 2.1 Existing platform boundary

The source documentation establishes an **activity discussion** model, not arbitrary team rooms:

- One primary discussion room per activity, automatically provisioned with the activity.
- Authenticated students join via the activity's discussion link; participation/registration status does **not** determine chat membership.
- Durable messages and permissions belong to FastAPI/PostgreSQL; Redis handles live fanout, and the Go WebSocket gateway is a separately deployable delivery boundary.
- Notices and pinning have an established message-domain contract; actual client UX should be verified in the current branch.
- The approved messaging baseline has president-only group management and hard-delete semantics. It does **not** include internal team groups, independent moderation, or general organizer management privileges.

**Reported later UX state:** The owner has described working activity discussions, a “My Discussions” list, student display names, and private student–organizer inquiries marked `/private`. These are useful contextual reports but are not independently verified by the provided engineering pack. The new AI design **excludes private inquiry content by default** irrespective of implementation details.

**Proposed:** organization groups, persistent internal team groups, an organizer home/inbox, personal AI UI, shared group AI, and organizational memory. None should be labeled as an implemented capability on the basis of this document.

---

## 3. Personas and authority

| Persona | Main jobs | AI experience |
| --- | --- | --- |
| Organization president | Lead teams, oversee all activities, coordinate communications and policies | Personal AI; authorized cross-group operations; organization-memory governance |
| Delegated organizer / team leader | Coordinate assigned work, manage relevant activities or team discussions | Personal AI and team AI within actual delegated privileges |
| Volunteer / organization member | Collaborate in authorized internal groups | Shared group AI; personal AI only to the extent their account is entitled to it |
| Student / activity participant | Ask questions and receive updates in joined activity discussions | Shared discussion AI; access only to information appropriate to that member |
| Future president | Take over organization leadership without losing validated institutional knowledge | Inherits reviewed organization knowledge, not private predecessor conversations |

**Authorization rule:** Titles are descriptive, not permission grants. Authorization comes from Campusly's actual capabilities and scoped membership. An organization president does not thereby receive access to another user's personal memory or to a conversation they are not otherwise authorized to inspect. A group AI never adopts the privileges of the most powerful group member.

---

## 4. Information architecture and navigation

### 4.1 One combined organizer home — agreed direction

The organizer's home has **three primary content areas**, not a competing collection of dashboards:

1. **Ask Campusly:** prominent, clear entry to the personal assistant; optionally a voice input affordance later.
2. **Needs attention:** small, meaningful, permission-appropriate items such as replies to an organizer's request, direct mentions, and pending approvals.
3. **Conversations:** recent, joined activity discussions and, once implemented, internal organization/team groups.

Keep the existing manual activities, chats, announcements, and management screens reachable. Do **not** assume the home requires a separate permanent AI bottom tab; navigation needs mobile usability validation.

### 4.2 Three complementary destinations

| Destination | Primary question | User action |
| --- | --- | --- |
| **Home / inbox** | “What needs my attention?” | Review update; open source conversation; ask AI for an on-demand summary |
| **Personal AI** | “What do I want to know or get done?” | Ask, plan, prepare, inspect and approve an action |
| **Group chat** | “Where do we collaborate directly?” | Read, reply, discuss, mention shared AI |

**Non-duplication invariant:** Sending via the assistant creates a message in the actual target group; receiving replies happens in that same canonical conversation. The inbox and assistant may reference it, but must not create parallel copies of group chats.

### 4.3 Organization context

- A visible **organization selector** anchors the personal assistant's current workspace.
- Users belonging to multiple organizations must be able to distinguish which is active before an operation.
- Activity and group context should ordinarily be inferred from the request and conversation, **not** require navigation into that activity/group first.
- If the target is ambiguous (“technical team” exists in two organizations, or multiple competitions match), ask a concise disambiguation question.
- A read-only question about another authorized activity need not permanently switch the conversation focus.
- A deliberate instruction such as “Let's work on the Python workshop instead” may change the active conversational focus, visibly.
- Before executing any write, the target organization/resource is explicit in the approval preview, regardless of inferred context.

### 4.4 Do not require opening personal AI from inside a group

An organizer should be able to ask “What happened in our technical team today?” directly from personal AI. Opening a group to activate a group-specific private assistant is **not** the primary navigation pattern. Optional contextual shortcuts may be considered later only if usability research supports them.

---

## 5. Personal AI — conversation experience

### 5.1 One assistant, optional multiple conversations

- Each user has **one conceptual personal Campusly assistant**.
- The assistant offers one familiar entry point and a history of optional private conversations.
- Conversations can cover one activity, a team, several related topics, or general organization management; no forced one-chat-per-activity rule.
- The organizer can start new conversations, resume older ones, and rename them.
- The assistant should maintain conversational references (“it”, “them”, “the competition”) while checking authorization and ambiguity.
- Opening a saved conversation restores its history and prior focus **subject to current permissions**.
- Conversation history should not become a separate source of truth for activity data.

**Proposed simple UI:** assistant title; current organization; conversation history/new chat; conversational canvas; typed input; contextual action previews; source links. Avoid a permanent chain of mandatory activity/team selectors.

### 5.2 Personal assistant use cases

| Intent | Example | Expected output |
| --- | --- | --- |
| Activity Q&A | “How many registered for our competition?” | Current authorized registration summary |
| Cross-group summary | “What did the technical team discuss today?” | On-demand sourced summary of accessible messages |
| Drafting | “Write an announcement about tomorrow's workshop.” | Editable draft; no send yet |
| Coordinated messaging | “Ask the technical and publicity teams for progress.” | Separate target-specific previews and approvals |
| Activity operations | “Change capacity to 100.” | Proposed change, permission check and confirmation; FastAPI enforces rules |
| Review | “What should we check before publishing this competition?” | Current details plus appropriately sourced previous lessons |
| Follow-up | “Did the technical team reply to my request?” | Relevant replies, unresolved state, and links to the actual group |

### 5.3 Read versus write behavior

**Read:** Answer directly when the user is authorized, data is sufficiently current, and the request is unambiguous. Cite or link original evidence when useful; clearly describe uncertainty, staleness, and missing records.

**Draft:** Create proposed messages/announcements/changes without committing them. Allow editing and cancellation.

**Write:** Show the exact content, sender, action type, organization, destination, and material impact. Require the organizer's explicit approval before performing a normal user-visible write. Re-check authorization and target state at execution; confirmation is **not** an authorization bypass.

**Group message attribution:** When an organizer approves a message sent via personal AI, the actual chat message displays the **organizer as sender**, not `CampuslyAI`. The system retains internally that this was AI-assisted, who approved it, which action it fulfilled, and where it was posted.

**Shared group AI attribution:** Replies by the shared group AI display as `CampuslyAI`, not as a human organizer.

**Action receipt:** After success, show what happened and a link to the canonical target (message/activity/notice). Do not claim success until the backend confirms it. Retries must not duplicate a previously successful action.

---

## 6. Human groups and internal organization structure

### 6.1 Agreed product direction

Campusly should support both:

- **Activity discussion groups:** student-facing, tied to an activity (existing approved model).
- **Internal organization/team groups:** persistent, organization-scoped conversations for leadership, technical, publicity, volunteer and other teams (**proposed extension**).

A group is not an AI conversation. Human discussion continues whether or not AI is available. Internal teams may work on multiple activities; a team group's lifetime is not automatically limited to one activity.

### 6.2 Preferred group structure — provisional

**Hybrid structure under consideration:** persistent general organization/team groups plus one discussion group per activity, with optional private activity-coordination groups *only when truly useful*. Avoid auto-creating dozens of groups. Internal-group membership, creation rights, naming, moderation and activity associations need a separate approved business spec; this document does **not** alter existing activity-group semantics.

### 6.3 Shared group AI

- A visibly identified AI participant can help brainstorm, prepare checklists, answer activity questions and summarize accessible discussion.
- In team groups, the AI may use authorized team context and reviewed relevant organizational knowledge, without impersonating any member or inheriting their role.
- In activity discussions, student-facing answers may use public/current activity details and permitted discussion context; do **not** disclose organization-private material or private student inquiries.
- The group AI can advise or draft, but consequential writes still obey the permissions of the **requesting user** and approval rules. Group AI should not silently become an administrator.

---

## 7. Group AI reply trigger — agreed V1 rule

**Candidate trigger (simple and deterministic):** A newly created, permitted human group message is considered for a group-AI reply **only when**:

1. It explicitly mentions `@CampuslyAI` (or the canonical structured AI mention), **OR**
2. Its text contains `?` or `？`.

If neither condition holds, **do not invoke the reply-decision model and do not generate an automatic group reply**.

If the condition holds, the AI then decides whether answering is useful, relevant, sufficiently grounded, and permitted. It may stay silent on low-value questions, jokes, conversation already answered, or requests that would violate privacy/authorization. For an explicit mention, it should ordinarily acknowledge the request; it may decline or explain a limitation instead of fulfilling it.

### 7.1 Examples

| Incoming message | Candidate? | Typical outcome |
| --- | --- | --- |
| “Hello everyone!” | No | No AI reply |
| “Where is tomorrow's workshop?” | Yes | Answer if venue can be verified |
| “Are we meeting at 3 PM？” | Yes | Verify, then answer/clarify |
| “Haha, really?” | Yes | Usually remain silent |
| “@CampuslyAI, summarize our preparation.” | Yes | Respond within authorized scope |
| “We should test the laptops tomorrow.” | No | No instant AI reply; eligible for scheduled memory processing |
| AI's own posted reply | **No** | Prevent loops |

### 7.2 Mandatory safeguards

- Don't recursively trigger from AI messages, action receipts, system events or duplicate delivery.
- Deduplicate by the canonical message/event identity; don't generate multiple replies to one message during retries.
- Check group existence, visibility and access when reading context and posting.
- Limit response frequency and resource usage so bursts do not flood a group.
- Respect a clear product-level AI availability/opt-out setting where supported (exact controls to define).
- Retain an on-demand human escalation path when sources do not support a reliable answer.
- No external LLM disclosure of unapproved real participant conversations; use synthetic/approved anonymized data during evaluation and establish a privacy/consent basis before any real-data AI processing.

**Important:** This trigger applies to **interactive replies only**. Scheduled memory maintenance operates independently and can consider ordinary messages that did not trigger a reply.

---

## 8. Combined inbox and follow-up UX

### 8.1 What the inbox is

A **personalized view** over authorized, existing conversations and action outcomes. It is **not** a new chat room or separate message store.

Potential items:

- Replies to requests the organizer sent manually or through AI.
- Direct mentions in groups.
- Pending action approvals (when the user is present and able to approve them).
- Recent important updates based on reliable, initially conservative signals.
- Canonical links to the original messages/groups.

Avoid a noisy summary for every unread message. Keep ordinary unread chats visible in the conversations area.

### 8.2 Reply-follow-up journey (core end-to-end UX)

1. President tells personal AI: “Ask our technical team to confirm readiness.”
2. AI resolves the authorized group and presents an editable message preview.
3. President approves; FastAPI writes the message as that president.
4. UI shows success, an action receipt and a link to the real group message.
5. Teammates reply normally in that group; the canonical conversation stays intact.
6. The home surfaces relevant replies; the president may read the group directly.
7. Alternatively, the president asks personal AI: “What did they say?” and receives a sourced, permission-checked summary.
8. If requested, AI prepares a further reply for explicit approval.

### 8.3 Lightweight coordination tracking — later option

Link approved requests to subsequent relevant replies and show an *unconfirmed* status such as “awaiting response.” Do **not** claim a task was completed solely because someone replied, and do not build a full task manager into the first release without validating this demand.

---

## 9. Memory model — agreed conceptual design

Memory is **not** identical to chat history. Chat history is retained according to its own domain rules; memory is a selective, scoped representation of useful knowledge, usually with provenance.

### 9.1 Memory scopes

| Scope | What it can contain | Access and governance |
| --- | --- | --- |
| **Personal** | Private preferences, recurring work style, summaries from user's private assistant interactions | User-owned; other organizers/presidents do not inherit it |
| **Conversation** | Topic continuity and action references within an individual private assistant conversation | Its owner, under current authorization |
| **Group** | Working decisions, unanswered questions, recurring problems and useful discussion facts | AI may retrieve only for authorized requesters; inspection/correction restricted to authorized group managers |
| **Organization** | Durable, shared procedures, validated lessons, announcement guidelines and institutional knowledge | Reviewed/confirmed by president or expressly delegated administrators |

**Interpretation:** Personal/conversation/group memory are the **working level**; organization memory is the **reviewed institutional level**. These are scoped stores/views, **not** an unrestricted inheritance chain. Group knowledge is not silently visible to every group or student. An activity retrospective is a structured view over evidence, not necessarily a fifth storage silo.

### 9.2 Automatic working-memory extraction

- Process new, permitted messages since the last successful processing checkpoint.
- Extract useful recurring questions, attributed commitments, tentative plans, disagreements, and unresolved issues.
- Separate **observation**, **inference**, **proposal**, and **verified fact**.
- Keep references to canonical messages/activities and timestamp/status where available.
- Deduplicate and update existing related facts rather than accumulate contradictory copies.
- Support correction, dismissal, supersession, expiry and policy-controlled deletion.
- Do not automatically convert a suggestion into an organizational rule or an official activity-data change.
- Avoid collecting unnecessary student personal information, especially into durable organization-wide knowledge.

### 9.3 Schedule and processing

**Initial experiment:** approximately **two extraction runs per day** per eligible group, with one/day and three/day as alternative settings to evaluate. This is a tuning proposal, **not** a deployed schedule or a commitment to create any cron job now.

- Scheduled memory work is **batch/asynchronous**; it must not interrupt group chat.
- Process by incremental checkpoints/cursors and support retry/deduplication.
- An activity-end retrospective can be generated on demand or at a defined lifecycle checkpoint once permitted.
- Memory maintenance does not require posting an AI message to the group.
- Interactive mention/question-mark replies stay event-triggered and do not wait for the next memory run.
- The product must explain to participants when and under what policy their group messages may be processed for AI memory; a real-data consent/privacy review precedes rollout.

### 9.4 Memory lifecycle and reliability

A working-memory entry may progress through:

`Observed → Suggested insight → Reviewed/Confirmed → Superseded/Archived`

These are knowledge states, not chat-room states. Every entry should, where feasible, retain: content, scope, type, source resource IDs, observed date, processing/review date, author/approver (if any), confidence or uncertainty, expiration/review rule, and revision history.

A deleted group or loss of membership must invalidate access to its associated source material and relevant memory according to approved retention/erasure policy; old memory must not become a bypass around room deletion. Exact deletion/retention behavior requires separate design before implementation.

### 9.5 Memory management UI

Authorized presidents and delegated organizers should be able to:

- Review and search relevant memory entries.
- See where a claim originated when its source remains accessible.
- Correct, dismiss, approve, supersede or archive entries within their authority.
- Distinguish working observations from official organization guidelines.
- Control or disable eligible automatic memory collection consistent with the privacy policy.

Ordinary members can still receive appropriate answers from the group AI without having administrative memory-edit access. Personal memory remains inspectable/manageable by its owner.

### 9.6 Privacy and evidence boundaries

- Existing `/private` student–organizer inquiries are **excluded by default** from group AI and scheduled shared memory extraction.
- A memory entry does not grant a user permission to read its source message or reveal otherwise private content.
- Public/student-facing group AI does not access private team discussion or organization-private memories merely because the activity belongs to the organization.
- If evidence is inaccessible, deleted or insufficient, do not fabricate a reference or disclose prohibited contents.
- Avoid sending identifiable real conversations to unreviewed AI providers. Data minimization, retention and user notice are required for production.

---

## 10. Learning from past activities

### 10.1 Why this matters

Organizations lose operational knowledge as presidents and volunteers graduate. Campusly should help preserve useful lessons **without treating every chat opinion as a confirmed failure or policy**.

### 10.2 Activity retrospective

For a completed activity, a reviewable retrospective may capture:

- What worked well (supported by data or authorized participant feedback).
- Repeated student questions and communication gaps.
- Coordination mistakes or unresolved issues, with the evidence and uncertainty.
- Proposed preventive improvements for similar future activities.
- Links to relevant source discussions, current activity records and approved corrections.

**Example (illustrative):** Many students ask whether laptops are provided. The AI records the repeated question. It may *suggest* that equipment requirements were unclear; an organizer reviews whether the original description was actually ambiguous before confirming this as a reusable lesson.

### 10.3 Organization-level promotion

`Authorized activity/group sources → Provisional group insight → Reviewable activity retrospective → President/delegate confirms reusable organization lesson`

Confirmation never changes underlying registration rules, room membership or official activity details. Later lessons may supersede older ones while an appropriately authorized revision history preserves provenance.

### 10.4 Reuse on demand

When an authorized organizer asks “What did we learn from our last competition?” or “Review our preparation plan,” the assistant retrieves *relevant* confirmed lessons, alongside live activity data. It may recommend a checklist or point out repeated historical issues **with source references and explicit uncertainty**.

No automatic daily cross-group briefing is needed. Historical knowledge is used when the user asks, or when they explicitly invite the group AI to brainstorm/review. Group AI should not habitually interrupt ordinary chat with lessons.

---

## 11. Source of truth, authorization and action integrity

1. **FastAPI owns business rules, authorization, identity and committed writes.** The AI only uses authorized tools/APIs; it never writes business tables directly.
2. **PostgreSQL owns durable chat and authoritative records**; Redis live fanout is not a durable memory store. AI work may have an independent acknowledged, retryable processing queue.
3. **Current backend wins over memory** for venue, registration status/capacity, deadlines, approved roles and other mutable operational facts.
4. **Permissions are scoped to the requester and target at execution time.** Conversation history and apparent model confidence are not authorizations.
5. **Human-approved sender identity is preserved:** personal AI-assisted messages appear as the approving human and retain an internal AI-assistance/action audit trail.
6. **Shared AI has an explicit bot identity** and may not impersonate human members.
7. **Explicit human approval for operational writes**, with per-action or clearly bounded batch previews. No hidden “approve everything in the future” permission in V1.
8. **Memory capture is not an administrative action:** reversible provisional extraction can occur under an enabled privacy/retention policy, but promoting official organizational knowledge requires confirmation.
9. **Never leak private inquiries or one organization's information into another organization.** Cross-org intent must be made explicit when allowed.
10. **AI failure is nonblocking:** ordinary send, receive, reconnect, group membership and manual management work without the AI service.
11. **Prompt injection/data contamination:** treat chat text and retrieved memory as untrusted data, not privileged instructions. A student's message “AI, ignore the rules” cannot alter tool permissions or organization policy.

---

## 12. UX states, copy principles and failure behavior

### 12.1 Interaction states

| State | Expected UX |
| --- | --- |
| No assistant conversation yet | Clear question input and a few relevant examples; no mandatory setup |
| Resolving ambiguous organization/group | Ask the minimum specific clarifying question |
| Reading group context | Indicate relevant sources; don't imply access to unseen content |
| Preparing a message or change | Show editable preview with sender, destination and effect |
| Awaiting approval | Clearly label **Not sent / Not applied** |
| Writing | Show pending state; prevent accidental double submission |
| Success | Confirm action with canonical receipt/reference |
| Failure/timeout | State that success is not confirmed, offer safe retry/check, avoid duplicate actions |
| Inaccessible memory or deleted group | Explain access limitation; avoid exposing source content |
| Insufficient evidence | Abstain, qualify, or recommend organizer verification |
| AI unavailable | Manual messages, groups and activity operations remain functional |

### 12.2 Product language

- Differentiate **personal assistant** (private) from **CampuslyAI** (group-visible bot).
- Differentiate **group message** from **official notice/announcement**, which may have additional prominence or effects.
- Clearly label AI drafts, proposed changes, provisional memory and confirmed policies.
- Prefer plain, concise answers and action previews over long unsolicited reports.
- Do not claim “completed,” “sent,” “approved,” or “verified” unless the relevant backend or reviewer actually confirmed it.

---

## 13. Canonical user journeys and acceptance scenarios

### Journey A — Send, receive, understand, follow up

**Given** a president belongs to the relevant technical team group, **when** they tell personal AI to request readiness confirmation, **then** AI prepares a correctly scoped message, obtains approval, posts as the president through FastAPI, retains an action receipt, and later lets the president access or summarize replies from the actual group. No duplicate AI-only chat is created.

### Journey B — Direct group AI mention

**Given** a joined participant says `@CampuslyAI, where is the competition?`, **when** the bot evaluates the candidate, **then** it retrieves public/authorized activity details, answers in the group as `CampuslyAI` if supported, and does not reveal private inquiries or internal team notes.

### Journey C — Question mark without mention

**Given** a group member asks `Do I need my own laptop?`, **when** the new message is processed, **then** the candidate filter passes, the AI evaluates whether a useful answer is warranted, and either replies or remains silent. A non-question-mark ordinary message without a mention does not invoke the immediate reply decision.

### Journey D — Silent memory preservation

**Given** three organizers discuss an unresolved equipment check without mentioning AI or using a question mark, **when** the next authorized scheduled memory pass runs, **then** it can capture the pending check as attributed **working** memory without sending an unsolicited group reply or treating the check as completed.

### Journey E — Lessons for future organizers

**Given** previous students repeatedly asked about equipment, **when** a new president asks AI to review a similar competition plan, **then** the assistant may surface the relevant reviewed lesson with a source or retrospective link, rather than claiming a historical mistake without evidence.

### Journey F — Wrong permission / private boundary

**Given** a volunteer requests an organizer-only capacity change in group chat, **then** AI cannot execute it, regardless of whether the president is in the group. **Given** a student asks about another student's private inquiry, **then** shared AI does not reveal the inquiry or related restricted memory.

### Journey G — Multi-organization ambiguity

**Given** one person belongs to both CESS and DAD, **when** they say “send this to the technical team” and multiple targets match, **then** AI clarifies and previews the exact organization and group before any send.

### Journey H — Failure isolation

**Given** the AI worker is down or a scheduled memory job fails, **then** human messaging persists normally; the job can retry without corrupting or duplicating messages. Memory processing delay is not misrepresented as real-time knowledge.

---

## 14. Prioritization: owner-agreed incremental delivery (2026-11-25)

The product vision remains intact. The active [milestones](../engineering/MILESTONES.md) and [task queue](../engineering/TASKS.md) supersede the earlier bundled delivery sequence. Each slice requires separate implementation authorization and owner manual acceptance.

### 14.1 First interactive AI vertical slice (A1)

- One president-only personal assistant entry and usable local chat UI.
- Scoped core integration and complete model response, without token streaming.
- Basic loading/failure handling, standalone synthetic fixture and positive/negative authorization checks.
- No tool calling, group AI, persistent memory or background event queue required for A1. Conversation persistence and management follow A3/A4.

### 14.2 Next product increments

- Personal assistant: private conversation continuity, authorized organization/activity information and statistics, editable activity and notice drafts, separately confirmed creation/publication and group-message actions with receipts.
- Personal memory: bounded context, explicit owner-controlled memory, cross-conversation retrieval, management UI and approved memory suggestions. Personal memory precedes shared group AI in the planned order.
- Shared group AI: visible identity and enablement, explicit mentions, grounded responses, selective question-mark participation, reliability and on-demand summaries. Minimum safety and retry checks start with the first shared reply. The original combined workflow is the later AI17 compatibility checkpoint, not the first UI milestone.

### 14.3 Later possibilities

- Independently scoped internal team groups, attention inbox, cross-group retrieval, group memory and review, institutional knowledge and retrospectives.
- Scheduled intelligence and controlled automation require separate policy, consent and owner approval. No recurring schedule is approved.
- Voice, external MCP/OpenCode clients, deeper coordination and advanced retrieval/orchestration when justified. No token streaming requirement.

---

## 15. Evaluation and research questions

The product should be judged by user outcomes and controlled evidence, not by how many AI tools are installed.

| Dimension | Possible measure |
| --- | --- |
| Organizer usefulness | Time and steps to complete a cross-group coordination task manually vs with AI |
| Action integrity | Correct target, sender, permission, confirmation and idempotency rate |
| Group reply behavior | Precision of reply/no-reply decisions; useful answer rate; unnecessary interruptions |
| Grounding | Supported answers with valid, accessible references; hallucination/abstention rates |
| Memory quality | Precision of captured facts, uncertainty labeling, duplicate/stale entry rate |
| Institutional learning | Relevance and factual support of retrieved historical lessons |
| Privacy/security | Cross-organization leakage, private-inquiry leakage, unauthorized writes (target zero) |
| Performance/cost | End-to-end response latency, inference calls per candidate, memory-processing cost |
| Reliability | Chat unaffected by AI outage; job retries and duplicate prevention |

**Research protocol:** start with synthetic or approved anonymized group transcripts, known ground truth, permission tests, and held-out examples. Explicitly test question-mark jokes, unsupported answers, ambiguous references, stale memories, deleted groups, and adversarial chat messages. Compare on-demand versus scheduled memory extraction as a later measurable experiment.

---

## 16. Explicit open decisions — not silently settled

1. **Internal team/group business rules:** creation, naming, membership, organization roles, visibility, lifetime, and optional activity-specific organizer groups.
2. **Personal assistant entitlement:** which non-president members receive one and what limited tasks it supports.
3. **Memory consent, retention, deletion and export:** student notices, opt-in/opt-out model, source retention, hard-delete propagation, third-party AI processing policy.
4. **Memory schedule:** one/two/three runs daily, activity-end trigger, workload thresholds and cost limits.
5. **Memory governance:** exact delegated review rights for group and organization scopes; audit/revision-history visibility.
6. **Inbox relevance:** initial reliable signals, notification boundaries, grouping and ranking without noisy/unverifiable AI claims.
7. **Conversation resumption:** whether activity shortcuts offer resume/new and how to keep the choice unobtrusive.
8. **Reply policy details:** rate limits, how long to wait for human replies, quote/mention parsing, attachment-only messages and moderation controls.
9. **Approval ergonomics:** safe multi-action batch confirmation, notice-specific wording, interrupted/expired approvals.
10. **Rollout and acceptance:** which organization pilots first, accessibility/localization, and how owner approval is captured.

These are specification refinement items, not reasons to reopen the agreed organizer-first philosophy or the simple reply trigger.

---

## 17. Decisions log (November 2026 brainstorming)

| ID | Decision | Status |
| --- | --- | --- |
| UX-01 | Organizer-first unified Campusly workspace, not an isolated chatbot | **Agreed** |
| UX-02 | Combined home = Ask Campusly + Needs attention + Conversations | **Agreed direction** |
| UX-03 | One conceptual personal assistant with optional multiple private conversations | **Agreed** |
| UX-04 | Resolve activity/group by natural language; don't require entering the group to use personal AI | **Agreed** |
| UX-05 | Human groups remain canonical and usable manually | **Agreed** |
| UX-06 | AI-assisted outgoing messages appear under approving human's name; internally audited | **Agreed** |
| UX-07 | Every operational write requires explicit approval and backend authorization | **Agreed** |
| UX-08 | Shared bot replies as CampuslyAI; only mention or `?`/`？` triggers reply decision | **Agreed V1** |
| UX-09 | Cross-group summaries are on demand, not proactive digests | **Agreed** |
| UX-10 | Personal/conversation/group memory separate from president-controlled organization knowledge | **Agreed concept** |
| UX-11 | Scheduled group memory extraction includes ordinary messages; two/day is experimental | **Agreed approach; frequency proposed** |
| UX-12 | Reviewed activity retrospectives can inform later organizers with references | **Agreed direction** |
| UX-13 | Existing private inquiries excluded by default; no permission inheritance | **Required safeguard** |
| UX-14 | MCP is not part of the initial in-product UX | **Agreed** |
| UX-15 | Persistent internal team groups + optional activity coordination group | **Provisional hybrid structure** |
| UX-16 | President-only initial assistant; visible end-to-end slices with UI and manual acceptance; personal memory before shared AI | **Agreed 2026-11-25** |
| UX-17 | Complete responses only; loading/error UX, no token streaming | **Agreed 2026-11-25** |

---

## 18. Engineering handoff boundaries and references

This is a **product/UX contract**, not an approval to modify code, stage, production, data, scheduled jobs, or Notion. Engineering design should inspect the real repository before choosing frameworks, APIs, new tables, workers or migrations. Prefer the smallest credible implementation that demonstrates the agreed interactions.

**Existing project references from the provided handoff pack:**

- `AGENTS.md` — repository agent instructions and scope gates.
- `docs/engineering/CONSTITUTION.md` — engineering, security and academic rules.
- `docs/engineering/ARCHITECTURE.md` — messaging ownership, existing topology and events.
- `docs/engineering/MILESTONES.md` — current engineering milestones and unaccepted owner-review gates.
- `docs/engineering/M1_IMPLEMENTATION.md`, `M2_IMPLEMENTATION.md`, `M3_IMPLEMENTATION.md` — evidence for persistence, live gateway and benchmarks.
- `docs/engineering/M4_REPORT_EVIDENCE.md` — draft report preparation and remaining owner actions.
- `docs/specs/Campusly_Activity_Messaging_Business_Spec_DRAFT_v0.4.md` — approved existing activity-group business rules (referenced by engineering docs; the original spec was not in this attachment set).

**Reconciliation note:** Some older planning text describes the Go gateway or messaging as future work, whereas the later M1–M4 evidence shows implemented/tested local messaging components; owner acceptance and production rollout are separate questions. Use **current repository code and the approved business spec** to verify actual state before extending them. The latest user conversation defines this new AI product direction; it does not retroactively change the original messaging contract.

**Academic constraint:** The Cloud Computing report has a December 8, 2026 deadline. Any implementation scheduling must protect course submission, the owner's review/teach-back, graduation requirements, and Fall 2027 master's preparation. Drafts and experiments must distinguish proposed, implemented, tested, and verified claims.

---

## 19. Acceptance of this document

- [x] Brainstorming decisions consolidated into a product/UX specification.
- [x] Existing messaging foundation separated from proposed AI and internal-team capabilities.
- [x] Reply trigger, memory separation, action approval and cross-group experience recorded.
- [ ] Product owner reviews this document and requests any revisions.
- [ ] Exact privacy/permissions/memory policy is approved for a real-data pilot.
- [ ] Engineering implementation milestone(s) are separately authorized.
- [ ] Any feature is implemented, tested and accepted (none claimed here).

**Document footer:** Campusly AI Product & UX Specification v1.0 — discussion baseline, 2026-11-23.
