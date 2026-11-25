# Review personal memory locally

Phase C uses the three matching B8 branches as its base. Their tips
are AI `cebcf1c`, CAMS Core `09e60e2`, and H5 `8449f33`.
Each C branch builds on the previous slice in its own repository.

| Slice | Branch in AI, CAMS Core and H5 | Review scenario |
| --- | --- | --- |
| C1 | `feature/campusly-assistant-c1` | [Continue a longer conversation](C1.md) |
| C2 | `feature/campusly-assistant-c2` | [Save and remove a preference](C2.md) |
| C3 | `feature/campusly-assistant-c3` | [Apply a preference in a new chat](C3.md) |
| C4 | `feature/campusly-assistant-c4` | [Inspect and edit a preference](C4.md) |
| C5 | `feature/campusly-assistant-c5` | [Review a suggestion and erase it with its source](C5.md) |

For an individual review, select the same branch in all three
repositories. C5 includes the complete Phase C workflow. Use the
[synthetic lab setup](PHASE_B.md#run-the-synthetic-lab) and sign in
as the synthetic president. Open **Ask Campusly** in the Workshop
Society to try the linked scenario. The lab uses a disposable
PostgreSQL database on loopback. Press Ctrl+C to stop it. It does
not contact a model provider.

The Go reply contracts in C1 and C3 should be reviewed together
with their matching Core revisions. The C2 and C5 Core migrations
have separate downgrade behavior documented in their slice records.
Check the [decision](../decisions/0007-personal-memory.md), the
[living ML report](../reports/ML_REPORT_DRAFT.md) and
[technical blog drafts](../blog-drafts/) for the rationale, commands
and limitations. Each linked review scenario includes a synthetic
mobile capture.

Owner accept-or-fix is pending for each slice. Cloud measurements,
live data retention and hosted model behavior are outside the
local verification. The previous AI01 and AI17 checkpoints retain
their existing scope.
