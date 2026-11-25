# 0004 — Keep personal chat behind Core

The goal of A1 is to let a president open the assistant from Campusly and see a complete answer. The president's identity and current organization role already live in Core. Passing the browser's campus token to Go would duplicate that authorization and expose a credential with broader purpose. Core instead sends only the actor ID, organization ID, expiry and prompt. Both sides sign or verify the same v1 payload bytes using a dedicated HMAC key. The response repeats the scope so Core can reject a mismatch.

The prototype stays on loopback. Core checks active presidency before forwarding and again before returning the answer. It releases the authentication transaction during generation. A disconnect cancels the Core call, and signed expiry plus a three-second Go cap bound generation even if the connection stays open. Eight inference slots reject excess work immediately. Core allows five seconds for its upstream exchange; the prompt and reply caps are 2,000 and 4,000 UTF-8 bytes. Requests are capped at 16 KiB. The first model is a deterministic implementation of Eino's complete-response interface at v0.9.23. There are no tools or hosted calls in this slice.

For the local check, a synthetic president should reach the H5 chat, send a workshop prompt and see the full checklist. A non-president or president of the wrong organization must be denied. A modified envelope, an expired request and a bad model response must produce safe errors. The standalone fixture exercises the same Go HTTP handler without Core. These are functional checks. Provider quality, latency and cost are unmeasured.

The service uses only `Generate` and verifies the reply's role and text before
returning it. The Eino API signatures were checked in the pinned module
source and the repo-owned skill. See [model interface at v0.9.23](https://github.com/cloudwego/eino/blob/v0.9.23/components/model/interface.go).
