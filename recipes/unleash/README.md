# Unleash

Emulates the Unleash API (v8), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The same field, `enabled`, means something different on each of Unleash's three APIs, and none of the three meanings is "this flag is on." On the Frontend API it's a constant, always `true`, because the endpoint only returns flags that already evaluated on -- code checking `.enabled` there is asking a question whose answer is always yes when it doesn't throw. On the Client API it's one input ANDed with a separately-evaluated strategies list, so a flag is only truly on when both agree, and nothing states that rule in one place. On the Admin API it's a single boolean summarising an `environments` array whose members can disagree with it -- a flag can be on in development and off in production while the top-level field says just one thing.

The two SDK-facing APIs don't even agree on what to call the collection: one answers `{"toggles": [...]}` and the other `{"features": [...]}`, for the same underlying flags.

## Sources

- Documentation: https://docs.getunleash.io/api-overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve unleash     # run it
cauldron verify unleash -v # check every claim
```
