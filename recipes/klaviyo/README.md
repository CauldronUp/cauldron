# Klaviyo

Emulates the Klaviyo API (2024-10-15), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Klaviyo speaks JSON:API, so the data most integrations want sits one level down at `data.attributes`, with the object's own type sitting as a sibling of `id` rather than a field inside it -- a client that reads `profile.email` directly finds nothing.

The revision header is not optional and has no default. Klaviyo pins its behaviour to a dated revision string, and omitting it is a hard failure rather than a fall-through to the latest version -- usually the first wall a new integration hits. Authentication also isn't a bearer token: the scheme is `Klaviyo-API-Key`, and sending `Bearer` fails in a way that reads exactly like a bad key rather than a malformed header. Klaviyo has no sandbox and no hard delete for profiles through the API either, so test data accumulates permanently against the real plan.

## Sources

- Documentation: https://developers.klaviyo.com/en/reference/api_overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve klaviyo     # run it
cauldron verify klaviyo -v # check every claim
```
