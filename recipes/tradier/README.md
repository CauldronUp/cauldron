# tradier

Emulates the tradier API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Whether `orders.order` is an object or an array depends on how many orders there are -- Tradier's own words: a single order comes back as a JSON object, multiple as an array, because the API is generated from XML where one child element and a repeated one look the same. A test account with two orders in it writes and passes the array-handling code; production, on the quiet morning after someone cancels the rest, has one order left and `orders.order.map` is not a function.

Tradier's own OpenAPI description contradicts its own prose here, declaring `order` always an array -- two official descriptions of the same endpoint disagreeing about the exact thing most likely to break an integration.

## Sources

- Documentation: https://docs.tradier.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tradier     # run it
cauldron verify tradier -v # check every claim
```
