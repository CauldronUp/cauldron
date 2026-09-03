# CoinGecko

Emulates the CoinGecko API (coingecko), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The map is keyed by currency and one of the keys
is not a currency.** `/simple/price` answers `{"bitcoin": {"usd": 78148, "eur":
67468, "last_updated_at": 1788075830}}`, so the field that says when the prices
were taken sits at the same depth, under the same kind of key, as the prices
themselves. A client that iterates the object to list them -- the only way to
read a map whose keys are whatever the caller asked for -- gets a currency
called `last_updated_at` worth 1,788,075,830.

**And anything the API does not recognise is dropped in silence.** An unknown
coin is not an error and not a null: `?ids=bitcoin,zzzznotacoin` answers
`{"bitcoin": {...}}` with a 200 and nothing to say one was rejected, and asking
only for unknown coins answers `{}`. One price is an integer and the next a
float in the same response, because the value is truncated by magnitude rather
than typed by field. There are two error shapes -- a flat `{"error": "coin not
found"}` and a nested `{"status": {"error_code": 429, ...}}` -- so one client
needs two parsers. And `/ping` answers `{"gecko_says": "(V3) To the Moon!"}`,
where the only statement of which API version replied is the `(V3)` in the
middle of a sentence.

## Sources

- Documentation: https://docs.coingecko.com/reference/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve coingecko     # run it
cauldron verify coingecko -v # check every claim
```
