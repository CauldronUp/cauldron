# Klaviyo

Emulates the Klaviyo API (2024-10-15), for local development and tests.

**10 conformance cases, 3 checked against the live API.**

Struck live 2026-09-05 against a.klaviyo.com, no account and no key. This file declared one authentication_error, "The API key is missing or invalid.", for every failure; the real API sends three distinct verdicts through two sentences -- a missing credential gets `not_authenticated` / "Authentication credentials were not provided.", while a wrong prefix (Bearer, sent from habit) and the right prefix with a key nobody issued both get an identical `authentication_failed` / "Incorrect authentication credentials." Split and fixed below.

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
