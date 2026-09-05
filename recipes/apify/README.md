# apify

Emulates the apify API (v2), for local development and tests.

**13 conformance cases, 2 checked against the live API.**

Two were struck live against api.apify.com on 2026-09-05, and they closed a real gap: this file had declared a token-not-provided error and never wired it to anything -- no route and no auth verdict pointed at it, so an actually-absent credential fell through to an "authentication_error" entry that did not exist in this file at all. Apify also distinguishes an absent token from a present, wrong one, with a different code and sentence for each.

## What this Recipe found

A run that SUCCEEDED tells you nothing about whether it produced anything. Status describes the process (it started, didn't crash, exited zero); what it scraped lives in a separate dataset under a different id, so a scraper that found nothing, got blocked, or hit a page whose markup changed still reports SUCCEEDED. And the dataset endpoint is the one response in the whole API that isn't wrapped in {"data": ...} -- it's a bare array -- so the code reading a run and the code reading its results unwrap differently against the same provider.

Of the eight statuses, three are the -ING half of a pair (TIMING-OUT vs TIMED-OUT, ABORTING vs ABORTED), so matching on a prefix or `includes("ABORT")` treats a still-running job as a finished one. TIMED-OUT is also not an empty result -- everything written before the cutoff is still in the dataset, so it's the status most likely to hold real, paid-for partial data and also the one most likely to get discarded. And duration appears twice in two different units (durationMillis, runTimeSecs) on the same object, which is a subtraction waiting to happen.

## Sources

- Documentation: https://docs.apify.com/api/v2
- Machine-readable description: https://docs.apify.com/api/openapi.json, last checked 2026-08-31
  `cauldron drift apify` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apify     # run it
cauldron verify apify -v # check every claim
```
