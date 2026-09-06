# tradier

Emulates the tradier API (v1), for local development and tests.

**13 conformance cases, 2 checked against the live API on 2026-09-05.**

Tradier has a sandbox with its own separate host, so the order cases still cite documentation. The credential shape needed no account at all, and checking it live found this Recipe's own fault case had never been run against a real request.

## What this Recipe found

Whether `orders.order` is an object or an array depends on how many orders there are -- Tradier's own words: a single order comes back as a JSON object, multiple as an array, because the API is generated from XML where one child element and a repeated one look the same. A test account with two orders in it writes and passes the array-handling code; production, on the quiet morning after someone cancels the rest, has one order left and `orders.order.map` is not a function.

Tradier's own OpenAPI description contradicts its own prose here, declaring `order` always an array -- two official descriptions of the same endpoint disagreeing about the exact thing most likely to break an integration.

## What checking it live found

The fault-shaped failure this Recipe already modelled had never actually been produced by a real credential check: nothing wired `invalid_token` to the verdict a garbage token reaches, so its only case armed the fault directly rather than sending a request nobody told what to answer. Fixed to an unarmed garbage token, which answers the same shape -- and a second, real layer turned up beside it: no `Authorization` header at all is a different gateway entirely, plain text, `Content-Type: plain/text` (backwards, and Tradier's own), body `"Invalid access token"`, not the JSON fault envelope a present, wrong token gets.

## Sources

- Documentation: https://docs.tradier.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tradier     # run it
cauldron verify tradier -v # check every claim
```
