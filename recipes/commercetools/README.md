# commercetools

Emulates the commercetools API (v1), for local development and tests.

**19 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

commercetools refuses any write that doesn't quote the current version of the resource -- every update is a POST with a body of {version, actions}, and a stale version is refused with the current one handed back so the retry can be scripted. That's optimistic concurrency, and it's a contract a test suite is structurally unable to exercise: a suite is the one place where nothing else is writing to the record, so code that quotes a stale version or none at all passes every local test, and the failure only shows up as a silent overwrite in production -- the later write wins, the earlier one vanishes, nothing is logged.

A product also carries two complete copies of itself: masterData.current is published, masterData.staged is edited, and there's no product.name at all -- the name lives at masterData.current.name.en, three levels down behind a locale key, and every translatable string is an object keyed by locale rather than a plain string. Money is centAmount in the smallest indivisible unit, which for JPY is the yen itself ("5 JPY is specified as 5" in commercetools' own words), so a helper that divides every currency by a hundred is wrong by a factor of a hundred for currencies with no minor unit.

A listing's total is explicitly documented as "an estimation that is not strongly consistent," and offset paging has a hard, documented maximum of 10000 -- deep paging isn't slow here, it becomes unavailable.

## Sources

- Documentation: https://docs.commercetools.com/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve commercetools     # run it
cauldron verify commercetools -v # check every claim
```
