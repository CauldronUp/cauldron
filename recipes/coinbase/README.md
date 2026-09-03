# Coinbase

Emulates the Coinbase API (rest), for local development and tests.

**14 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

Its **two live APIs disagree about staleness** -- one
carries a nanosecond timestamp and a sequence, the other carries neither and
puts empty strings where prices go.

## Sources

- Documentation: https://docs.cdp.coinbase.com/exchange/reference/exchangerestapi_getproductticker
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve coinbase     # run it
cauldron verify coinbase -v # check every claim
```
