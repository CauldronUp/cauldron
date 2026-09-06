# incidentio

Emulates the incidentio API (v2), for local development and tests.

**20 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real workspace; the refusal cases were struck live, unauthenticated, against api.incident.io.

## What this Recipe found

Almost every status field in this collection is a fixed enum; incident.io's is not. Status and severity are configured per workspace, so the words are whatever someone at that company typed into a settings page, and code that hard-codes `"resolved"` works only against the workspace it was written for. The API's own answer is to give each value an id and a numeric rank alongside the name -- severity "Major" outranks "Minor" because of its rank, not its wording, so sorting on the name gives alphabetical order dressed up as urgency.

Status and severity also move completely independently -- an incident can be resolved and still marked critical, or still triaging and marked minor -- so inferring one from the other invents a relationship the product doesn't have. The one fixed thing is a status's `category` (`triage`, `active`, `learning`, `closed`), the same four words in every workspace underneath a status name that isn't; a client that wants to know whether an incident is actually over should read the category and never the display name.

The live probe found this file's authentication error wrong: a missing credential and a wrong one carry different codes and sentences (`missing_authorization_material` against `access_token_invalid`), neither the `unauthorized` this file declared. Both travel inside the same `{type, status, request_id, errors: [...]}` envelope, and the top-level `type`/`status`/`request_id` fields are not modelled -- this format's list-style envelope has no way to place a value that varies per error at the top level, only inside each entry. An unrouted path and a wrong method on a real path both answer an empty 404 before authentication is ever consulted.

## Sources

- Documentation: https://api-docs.incident.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve incidentio     # run it
cauldron verify incidentio -v # check every claim
```
