# Listen Notes

Emulates the Listen Notes API (v2), for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-02.**

## What this Recipe found

Notes', whose **mock host ignores the identifier entirely**:
every podcast id returns the same record, so an integration can pass every test
without once addressing a specific thing.

## Sources

- Documentation: https://www.listennotes.com/api/docs/
- Machine-readable description: https://listen-api.listennotes.com/api/v2/openapi.json, last checked 2026-09-02
  `cauldron drift listennotes` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve listennotes     # run it
cauldron verify listennotes -v # check every claim
```
