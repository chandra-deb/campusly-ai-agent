# What the first hosted assistant check verified

The deterministic Phase A adapter made it possible to test the same chat
and retry flow repeatedly. The first hosted check asked a narrower
question. Would one Cerebras completion return through the same signed
Core request and appear as one saved answer in H5?

The Go adapter accepted one fixed workshop prompt with no history and
used a 256-token completion cap. The checks with a mock transport
verified that request and rejected truncated replies. Running the
browser scenario sent the synthetic prompt to Cerebras once. The API
returned HTTP 403, which Core presented as an upstream failure. The
models endpoint returned a geographic-access diagnostic on the same
connection.

The outcome leaves one open engineering check. A hosted answer still
needs to pass through the three-second Go budget and reopen in H5.
The [test record](../evidence/CEREBRAS_SMOKE.md) keeps the rejected
attempt alongside the credential-free evidence.
