# Who gets to ask the Campusly assistant?

A president already signs into Campusly to manage activities. A1 gives
that dashboard an **Ask Campusly** entry. When someone sends a message, Core
looks up their current role for the organization and calls a Go service on
loopback. Core releases its database transaction during generation and checks
the current role again before showing the reply.

I wanted the Go service to have just enough information to answer. Sending
the Campusly access token would make the service responsible for a broader
credential, so Core signs the exact bytes of a short-lived request instead.
Changing its actor or organization breaks the signature. Eight slots put an
upper limit on simultaneous generation, and the request deadline bounds each
call. The model is deterministic, which makes a synthetic workshop question
repeatable without a provider account.

The first Core checks denied a non-president, a president from another
organization, and a request that tried to supply an actor ID in the browser
payload. See the [A1 evidence](../evidence/A1.md) for commands and the UI check.
Live provider behavior and latency still need separate work.
