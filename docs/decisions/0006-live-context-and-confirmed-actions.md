# 0006 — Resolve assistant questions against live Core data

A president should be able to ask what they manage from the same chat where
they plan an activity. Sending an organization's records to the model as
conversation context would make revocation difficult to enforce across
saved turns. Separate dashboard controls would avoid that problem but also
remove the natural-language entry point.

B1 adds a signed, versioned intent request from Core to the local Go
service. Go receives the prompt and request scope. It returns a bounded
operation, while Core resolves the operation with the authenticated
actor's current permissions. Core builds the reply from the authorized
rows. Campusly answers are transient in the open chat, since they may
contain facts whose permissions or values change. Regular planning
replies still use the existing private history path.

Expect the same question to reflect a changed role on the next request.
Verify with two managed organizations, an unrelated organization, and
a revoked role. The deterministic local adapter tests product behavior,
with language coverage limited to the supported sample requests.
