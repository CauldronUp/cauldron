# Lithic

Emulates the Lithic API (v1), for local development and tests.

**19 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

The clearest trap is a field that answers two different questions depending on when you ask it. Lithic's deprecated `amount` property holds the authorization amount while a `card_transaction` is PENDING, then silently starts holding the settled amount the moment it flips to SETTLED -- same key, same identifier, two meanings, and the two amounts are not required to agree (a restaurant authorizes the bill and settles the bill plus a tip). Code that reads `amount` once and caches it is reading a number that goes stale without any change in the shape of the response.

Authentication was checked live rather than assumed, and it distinguishes three failures with three different sentences: no header, a header that is not shaped like a Lithic key, and a well-formed key that does not exist. Only two of the three are reachable here -- the middle one cannot be told apart from a wrong scheme entirely, because Lithic's key is a raw value with no declarable prefix or pattern to check against, and declaring one would break key matching outright. Simulation endpoints, account holders, KYC, disputes and the financial-account/ledger surface are not modelled; no transaction here is ever created by a real card event, only by a fixture.

## Sources

- Documentation: https://docs.lithic.com/reference/authentication
- Machine-readable description: https://raw.githubusercontent.com/lithic-com/lithic-openapi/main/lithic-openapi.yml, last checked 2026-09-01
  `cauldron drift lithic` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lithic     # run it
cauldron verify lithic -v # check every claim
```
