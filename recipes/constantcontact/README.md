# Constant Contact

Emulates the Constant Contact API (v3), for local development and tests.

**18 conformance cases, 10 checked against the live API on 2026-09-01.**

## What this Recipe found

Contact's, whose **second create route sets consent**: one door
refuses a duplicate with a 409 and the other marks permission explicit with no
exception for somebody who unsubscribed.

## Sources

- Documentation: https://developer.constantcontact.com/api_reference/index.html
- Machine-readable description: https://developer.constantcontact.com/api_reference/bundledWithSamples.yaml, last checked 2026-09-01
  `cauldron drift constantcontact` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve constantcontact     # run it
cauldron verify constantcontact -v # check every claim
```
