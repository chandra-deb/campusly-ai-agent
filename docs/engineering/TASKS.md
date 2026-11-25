# AI work queue

The rows below are the recommended MVP execution order, not blanket technical dependencies. Approve scope before implementation. AI01 precedes AI10; personal read/approval work (AI11–AI13) and shared event/reply work (AI14–AI15) require the applicable reviewed contracts and AI10 foundation, not necessarily each other's UI. AI16 integrates both paths; AI17 follows successful integration. Keep the MVP focused and reassess independent later increments after local acceptance; parallel work requires explicit scope and compatible interfaces. Each substantive task updates the relevant living report/blog drafts or records N/A. Integration tasks include corresponding core or UI changes reviewed in their owning repositories.

| ID | Deliverable | Focused check and exit |
| --- | --- | --- |
| AI01 | Review ADR 0001 and draft v1 API/event examples, approval lifecycle, privacy filters and worker ownership | Walk through denial, replay, edited confirmation and crash-after-post. Record owner decisions and authorized AI10 scope. |
| AI10 | Standalone mock Campusly API, typed client, deterministic model fixture and local container entry point | Build and run from a clean checkout without secrets. Exercise allowed, wrong-tenant, revoked and unsupported-version cases. Document exact commands under scripts and put contract checks under tests. |
| AI11 | Core scoped read adapter and organizer's private assistant entry | Select organization explicitly when ambiguous. Two actors cannot open each other's conversation. |
| AI12 | Read one authorized activity fact with source reference | Verify correct venue, inaccessible source and abstention on missing facts. |
| AI13 | Preview, confirm and send one message via core | Verify exact target/text and human identity, internal receipt, edit/expiry, revocation and duplicate confirmation. |
| AI14 | Durable delivery of eligible public human messages | Mention or ASCII/fullwidth question mark qualifies. Ordinary text, private events, bot/system messages and duplicate deliveries produce zero extra inference calls. Restart before ack and recover. |
| AI15 | Scoped group tool reads and reply/no-reply policy | Answer sourced mentions as CampuslyAI, abstain on unsupported questions, stay silent on low-value candidates. Enforce shared limits and recheck deletion before posting. |
| AI16 | Integrated slice in owner’s local environment | Walk workshop scenario with mock then authorized provider if approved. Stop AI mid-request, verify chat send/history/reconnect, one receipt after retry and zero private metadata exposure. |
| AI17 | Owner personally tests the §14.1 slice | Record revision, issues and explicit accept/fix decision. Only acceptance opens eligibility for separate Cloud authorization. |

Later increments require individual decisions: combined attention inbox and reply links, multiple conversations, internal team-group rules before implementation, authorized cross-group summaries, reviewed scoped memory and retrospective promotion. Resolve consent, retention/export, hard-delete propagation, review rights and budget before any memory job. No recurring schedule is approved. External MCP/OpenCode clients can follow stable tools. Track remaining entitlement, rollout, moderation and interaction questions against spec §16.
