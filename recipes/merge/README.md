# merge

Emulates the merge API (v1), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Merge normalizes dozens of providers into one shape, and the finding is everything that normalization cannot fix. A null field means two different things behind an identical response: the underlying provider does not support it at all, or this particular account just never filled it in. Code that treats null as "not set" is right half the time and wrong the other half, with no way to tell which case it is in.

A linked account exists from the moment someone starts the link flow, with status `INCOMPLETE`, and reading data from it answers 400 rather than an empty list -- so counting linked accounts counts abandoned attempts, not real connections. The first sync after linking is asynchronous and can take minutes to hours; every list is empty and correct during that window, which is indistinguishable from a customer who genuinely has no data. `remote_data`, the escape hatch back to the original provider's shape, is null by default and only populated if the linked account was explicitly configured to include it -- an integration that never sees it in testing can still meet it for the first time in production.

Only the HRIS category is modelled; the sync-status machinery is fixture-driven rather than actually progressing, so a client handling the not-yet-synced state against this Recipe finds whatever the fixture holds and nothing moves it forward.

## Sources

- Documentation: https://docs.merge.dev/hris/overview/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve merge     # run it
cauldron verify merge -v # check every claim
```
