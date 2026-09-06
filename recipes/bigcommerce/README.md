# BigCommerce

Emulates the BigCommerce API (v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A BigCommerce integration talks to two APIs that disagree about almost everything, because catalogue and customers live on V3 while orders live on V2 -- nobody chose that, it's just where the endpoints ended up. V3 wraps every response in {data, meta}; V2 returns a bare array with no envelope, so a client with one response helper gets undefined from whichever half it wasn't written for, and placing an order needs both halves. Money is also a number in V3 and a string in V2 -- a catalogue price is 89.95, but the order that sells it carries total_inc_tax "89.9500" to four decimal places -- so adding two V2 totals in JavaScript concatenates them instead.

The credential header is X-Auth-Token, not Authorization and not a bearer scheme, so a client sending it the usual way gets a 401 that never mentions the header it actually wanted. V3 failures put the actual message under title rather than message, so code reading .message finds undefined on every failure BigCommerce sends. And the order status a store owner sees is a renameable string next to the real state, status_id -- branching on the string is branching on a setting.

One deliberately-wrong gap: V2 returns an empty collection as 204 No Content rather than an empty array, so calling .json() on it throws in production; Cauldron can't express a status conditioned on emptiness and answers 200 with an empty array instead -- wrong, but wrong in the forgiving direction.

## Sources

- Documentation: https://developer.bigcommerce.com/docs/rest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bigcommerce     # run it
cauldron verify bigcommerce -v # check every claim
```
