# Royal Mail

Emulates the Royal Mail API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Deleting an order here can cost real money. Royal Mail's own documentation for `DELETE /orders/{orderIdentifiers}` warns that a label generated on a deleted order is no longer valid, and if it still reaches a Royal Mail van the account is charged plus a handling fee -- the call is reversible in a database and not reversible on a parcel already moving.

The same path parameter takes a semicolon-separated list of up to a hundred identifiers, with string references quoted and percent-encoded because a reference can itself contain a semicolon. This Recipe matches one path segment as one identifier rather than parsing that list, so the encoding trap it documents is prose here, not behaviour -- and the array-shaped response Royal Mail always returns from that endpoint, one or a hundred, comes back here as a single object instead.

## Sources

- Documentation: https://api.parcel.royalmail.com/
- Machine-readable description: https://api.parcel.royalmail.com/swagger/v1/swagger.json, last checked 2026-08-31
  `cauldron drift royalmail` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve royalmail     # run it
cauldron verify royalmail -v # check every claim
```
