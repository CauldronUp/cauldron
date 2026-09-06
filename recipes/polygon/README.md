# polygon

Emulates the polygon API (v2), for local development and tests.

**15 conformance cases, 4 checked against the live API.**

Everything about bars, counts, and adjustment still cites documentation rather than an observation, because reaching it needs a paid account. The credential and routing checks were verified directly against api.polygon.io, unauthenticated, on 2026-09-05.

## What this Recipe found

Checked live: a missing API key and a wrong one are different sentences -- `{"error":"API Key was not provided"}` versus `{"error":"Unknown API Key"}` -- and this file had only modelled the second, under a name nothing in the Recipe actually wired to a credential check, so no request the emulator served could ever have reached it. Routing also runs ahead of the credential entirely: an unrouted path and a wrong method both answer Go's own plain-text net/http defaults, "404 page not found" and "405 method not allowed", needing nothing sent at all.

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
