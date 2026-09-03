# USPS

Emulates the USPS API (v3), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-09-02.**

## What this Recipe found

**A token is judged in three stages** -- absent and
non-JWT collapse together, and a JWT signed wrongly gets its own sentence.

## Sources

- Documentation: https://developers.usps.com/trackingv3
- Machine-readable description: https://developers.usps.com/sites/default/files/apidoc_specs/tracking_22.yaml, last checked 2026-09-02
  `cauldron drift usps` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve usps     # run it
cauldron verify usps -v # check every claim
```
