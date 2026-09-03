# UPS

Emulates the UPS API (v1), for local development and tests.

**7 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **token endpoint tells the truth and whose API does not**:
two codes separate an absent credential from a wrong one before you have a
token, and one body covers both once you do.

## Sources

- Documentation: https://developer.ups.com/api/reference/tracking/product-info
- Machine-readable description: https://raw.githubusercontent.com/UPS-API/api-documentation/main/Tracking.yaml, last checked 2026-09-02
  `cauldron drift ups` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ups     # run it
cauldron verify ups -v # check every claim
```
