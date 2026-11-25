# Review the Campusly-aware assistant locally

Phase B uses synthetic presidents, activities and discussions. The
branches in AI, CAMS Core and H5 have matching names and build on the
previous slice in each repository. Start with B1 for individual review
or use B8 to try the complete sequence in one lab.

| Slice | Branch in each repository | Owner scenario and evidence |
| --- | --- | --- |
| B1 | `feature/campusly-assistant-b1` | [Managed organizations](B1.md) |
| B2 | `feature/campusly-assistant-b2` | [Activity schedule and venue](B2.md) |
| B3 | `feature/campusly-assistant-b3` | [Registration and waitlist counts](B3.md) |
| B4 | `feature/campusly-assistant-b4` | [Editable activity draft](B4.md) |
| B5 | `feature/campusly-assistant-b5` | [Confirmed activity creation](B5.md) |
| B6 | `feature/campusly-assistant-b6` | [Grounded notice draft](B6.md) |
| B7 | `feature/campusly-assistant-b7` | [Confirmed notice publication](B7.md) |
| B8 | `feature/campusly-assistant-b8` | [Confirmed group message](B8.md) |

The original committed tips were AI `3c049ca`, Core `c8ae90c`
and H5 `566c88a`. Use matching branches when reviewing a slice. Each
new branch advances all three sides, even when a service's executable
code can carry forward.

## Run the synthetic lab

From the selected AI worktree, build the assistant and its fixture
with Go 1.25.14. Then run the Core lab with the matching H5 worktree,
a local Python environment with the CAMS dependencies, H5 dependencies
installed, and Docker's PostgreSQL 15 test image available.

```sh
cd service
go build -buildvcs=false -o /tmp/campusly-assistant-bin/ ./cmd/...
```

```sh
env -u TEST_POSTGRES_ADMIN_URL -u DATABASE_URL CAMS_ENV_FILE=/dev/null \
  python -m scripts.m2.assistant_lab /tmp/campusly-assistant-bin /path/to/matching/h5 --review
```

The lab prints a loopback sign-in URL and synthetic accounts. Sign
in as the synthetic president, select **Sample Workshop Society** under
**Manage**, open **Ask Campusly**, and create a chat. Use the linked
scenario for the slice. Ctrl+C stops the processes and removes the
disposable database. For the automated mobile test, use a local
Puppeteer Core module and Chromium executable as described in the
[A1 lab instructions](A1.md#reproduce).

The Go intent resolver recognizes the sample prompts. The discussion
checks open committed messages over REST. Hosted language coverage,
realtime delivery in this lab and Cloud measurements are pending
separate evaluation. Owner accept-or-fix is pending for each slice.
The [living ML/AI draft](../reports/ML_REPORT_DRAFT.md) and
[blog drafts](../blog-drafts/) record the engineering evidence.
