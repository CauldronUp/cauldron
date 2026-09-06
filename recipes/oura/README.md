# Oura

Emulates the Oura API (v2), for local development and tests.

**12 conformance cases, 10 checked against the live API on 2026-09-02.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**Three credential sentences and a typo worth keeping**: the
absent-token message ends mid-word, and it is reproduced as sent.

## Sources

- Documentation: https://cloud.ouraring.com/v2/docs
- Machine-readable description: https://cloud.ouraring.com/v2/static/json/openapi-1.37.json, last checked 2026-09-02
  `cauldron drift oura` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve oura     # run it
cauldron verify oura -v # check every claim
```
