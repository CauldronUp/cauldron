# finnhub

Emulates the finnhub API (v1), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The same letter means two different types depending on the endpoint. Finnhub's own field descriptions say `c` is a "list of close prices" on a candle response and the "current price" on a quote -- one is an array, the other a scalar, and a helper written against one endpoint returns an array from the other. `o`, `h` and `l` do the same thing.

A candle response is columns, not rows: seven parallel arrays with no object representing a single bar, so every read is a zip that assumes the arrays are the same length. When there's no data for a range, the status still comes back 200 -- a string field `s` says `no_data`, and the price arrays aren't empty, they're absent entirely, so `response.c.length` throws on a request that technically succeeded. And the response never echoes back what was actually asked for -- no symbol, no resolution, no date range -- so a plan that silently shortens a requested range gives no way to notice short of comparing the first timestamp against what was sent.

## Sources

- Documentation: https://finnhub.io/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve finnhub     # run it
cauldron verify finnhub -v # check every claim
```
