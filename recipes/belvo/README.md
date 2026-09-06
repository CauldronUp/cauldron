# Belvo

Emulates the Belvo API (v2), for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**The one that names the re-authentication state** outright,
where its neighbours circle it.

## Sources

- Documentation: https://developers.belvo.com/reference/retrievelinks
- Machine-readable description: https://raw.githubusercontent.com/konfig-dev/belvo-sdks/main/api.yaml, last checked 2026-09-02
  `cauldron drift belvo` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve belvo     # run it
cauldron verify belvo -v # check every claim
```
