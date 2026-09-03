# Mercado Libre

Emulates the Mercado Libre API (v1), for local development and tests.

**15 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

Libre's, where **most of the API is a 403 before it is
anything.** `/sites`, `/sites/MLA` and `/currencies` all answer
`PA_UNAUTHORIZED_RESULT_FROM_POLICIES` with no credential sent and none
required -- a WAF verdict, not an authentication one, and the same shape
`/users/me` returns, so a client cannot tell a missing token from a blocked
route. On what is reachable, the same category in two countries disagrees about
a field's *type*: `MLA1051` carries `settings.vertical` of
`"consumer_electronics"` where `MLB1051` carries `null`.

## Sources

- Documentation: https://developers.mercadolibre.com
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mercadolibre     # run it
cauldron verify mercadolibre -v # check every claim
```
