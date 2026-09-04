# polygon

Emulates the polygon API (v2), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The field `t` means two different things on two endpoints of the same API. On the range/aggregates endpoint it is documented as the start of the window; on the grouped daily endpoint it is the end of the window -- same one-letter field name, same shape, same position in an otherwise identical object. Code that joins bars from the two endpoints on `t` is off by one window on half its data, and nothing about the join looks wrong, because the field is present, numeric, and in the right order.

A few more one-character or one-word traps sit right beside it: `T` is the ticker and `t` is the timestamp on the same grouped-endpoint object, and a case-insensitive language or careless destructure gets whichever one it grabs. `otc` is left off the object entirely when false, Polygon's own wording, so `bar.otc === false` is never true anywhere -- the only way to check is testing whether the key exists at all. `limit` also does not limit the results returned; it limits how many underlying minute bars are used to build the aggregate, so `limit=10` on a daily range does not return ten days.

Prices are also split-adjusted by default, so the same historical trading day has two legitimate, different prices depending on a query parameter most callers never set. A basic-plan account gets `status: DELAYED` at HTTP 200 with data up to fifteen minutes stale, which this Recipe describes rather than serves, since a format's envelope constants are fixed per Recipe.

## Sources

- Documentation: https://polygon.io/docs/rest/stocks/aggregates/custom-bars
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve polygon     # run it
cauldron verify polygon -v # check every claim
```
