# Lightspeed

Emulates the Lightspeed API (API/V3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Lightspeed's REST API still looks like the XML service it used to be. The paging envelope sits under the key `@attributes`, the count inside it is the string `"2"` rather than the number 2, and the records live under the resource's singular capitalised name -- `Sale`, not `sales` -- so the shape is different for every endpoint and never the plural a client would guess. Money is a string with four decimal places, `"47.6000"`, so it neither adds nor compares like a number.

Relations are absent unless requested via a JSON array embedded in a query parameter, `?load_relations=["SaleLines"]`, and a sale without them has no lines and no error to say so. The Recipe also states its one deliberate inaccuracy: the real API returns a single matching sale as an object and a list as an array, so code that maps over the result works until exactly one sale matches. Cauldron cannot make a response shape depend on how many records matched, so this emulator always answers with an array, hiding the bug in the direction that is easiest to miss. K-Series, Lightspeed's restaurant product and a completely different API under the same brand, is not modelled at all.

## Sources

- Documentation: https://developers.lightspeedhq.com/retail/introduction/introduction/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lightspeed     # run it
cauldron verify lightspeed -v # check every claim
```
