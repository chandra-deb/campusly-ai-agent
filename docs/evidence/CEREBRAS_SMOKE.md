# Cerebras smoke through Phase A

The synthetic Phase A flow verified permissions, complete replies and
persistence with a deterministic model. The question here is whether one
Cerebras completion can pass through the same signed request, the three-second
Go generation budget and the H5 assistant without changing the model contract.

Use one fixed workshop prompt in a fresh disposable conversation. The opt-in
Go adapter sends that prompt once to Cerebras Chat Completions, requests a
complete reply and caps the completion at 256 tokens. The owner chose `gpt-oss-120b` and
the request sets reasoning effort to `low`. The platform receives
only the synthetic prompt. The adapter checks the choice, assistant role,
finish reason and output bounds. Each run can attempt one request. A provider
error or timeout is evidence, so the same run does not retry or change models.

A direct curl request would answer whether the key works. The H5 scenario
also checks the actual integration and saves the complete reply in the
disposable database. The default local adapter remains deterministic.
The credential-free checks exercise the real HTTP request and response
contract with an injected transport. The hosted test does not exercise
tools or establish provider accuracy or performance. See the
[Cerebras Chat Completions reference](https://inference-docs.cerebras.ai/api-reference/chat-completions).

## Reproduction

The base branches are AI `7a530b9`, CAMS `ddbb254` and H5 `5f31928`.
Use the A4 test dependencies and a fresh, owner-readable local key file.
Build both Go binaries from `service/` and run the following from the AI
repository root with explicit paths to the A4 worktrees.

```sh
go test -race -buildvcs=false ./internal/personal
go vet -buildvcs=false ./cmd/... ./internal/personal
go build -buildvcs=false -o /tmp/campusly-cerebras-bin/ ./cmd/...
```

```sh
CAMS_PUPPETEER_MODULE=/absolute/path/to/puppeteer-core.js \
python3 scripts/run-cerebras-smoke.py --core /absolute/path/to/CAMS \
  --h5 /absolute/path/to/cams-h5 --binaries /tmp/campusly-cerebras-bin \
  --core-python /absolute/path/to/CAMS/.venv/bin/python \
  --key-file /path/to/owner-only-key --screenshot /tmp/cerebras-smoke.png
```

The launcher starts the disposable lab and passes the synthetic login to the
browser fixture over standard input. Only the Go process reads the API key.
The browser checks the original prompt, the complete three-item reply and
the same saved turn after reloading. Ctrl+C or the end of the check stops
the lab and removes its disposable database.

## Results on 2026-11-25

- Credential-free Go race tests, vet and build passed with Go 1.25.0 and
  Eino v0.9.23. They covered the signed integration's existing checks,
  the exact provider request, prompt and history rejection, the attempt
  limit, and safe handling of provider errors and truncated output.
- One browser-initiated synthetic completion reached Cerebras. The
  provider returned HTTP 403. CAMS mapped it to HTTP 502. There was no
  hosted reply to display or reopen.
- An authenticated, read-only GET to `/v1/models` also returned HTTP 403
  with a geographic-access diagnostic. Key validity and model access
  remain unverified from this connection. The 403 response time is
  not a model latency measurement.
- The provider's terms prohibit bypassing its access restrictions.
  A hosted rerun needs a Cerebras-authorized environment.
  [Cerebras terms of use](https://www.cerebras.ai/terms-of-service)

No Cloud report update applies because this check did not involve the
Cloud research tasks. The live result is **BLOCKED** and the adapter's
credential-free checks are **TESTED**. The evidence does not settle
broader Eino adoption or hosted tool calling.
