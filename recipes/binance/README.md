# Binance

Emulates the Binance API (v3), for local development and tests.

**14 conformance cases, 10 checked against the live API on 2026-09-02.** The 3 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**Three price endpoints prove three different things**
and only one of them lets a caller detect a gap. Its book never names the symbol
it was asked for.

## Sources

- Documentation: https://developers.binance.com/docs/binance-spot-api-docs/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve binance     # run it
cauldron verify binance -v # check every claim
```
