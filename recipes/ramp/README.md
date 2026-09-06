# Ramp

Emulates the Ramp API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Ramp amounts are integers in cents with nothing in the field name to say so -- a transaction of `4250` is forty-two dollars fifty, and code that renders it directly shows a bill four thousand times too large. A pending transaction has no settlement yet and can still change: its amount is only an authorization hold, and a restaurant tip or hotel adjustment means the eventual settled amount can differ from what has already been written to a ledger.

A terminated card keeps producing transactions for a while after termination, so filtering a spending report by active cards loses real spending that the card's own state gives no hint about. A declined transaction is still a transaction too -- it shows up in the list with `state: DECLINED` and a real amount, so summing the list without filtering on state counts money nobody actually spent.

Nothing here is actually spent; which state a transaction holds, and whether it settled for its authorized amount, is entirely what a fixture puts there, since a live Ramp sandbox cannot be made to produce a declined or pending-then-settled transaction on demand.

## Sources

- Documentation: https://docs.ramp.com/developer-api/v1/reference
- Machine-readable description: https://docs.ramp.com/openapi/developer-api.json, last checked 2026-09-05
  `cauldron drift ramp` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ramp     # run it
cauldron verify ramp -v # check every claim
```
