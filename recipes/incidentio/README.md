# incidentio

Emulates the incidentio API (v2), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Almost every status field in this collection is a fixed enum; incident.io's is not. Status and severity are configured per workspace, so the words are whatever someone at that company typed into a settings page, and code that hard-codes `"resolved"` works only against the workspace it was written for. The API's own answer is to give each value an id and a numeric rank alongside the name -- severity "Major" outranks "Minor" because of its rank, not its wording, so sorting on the name gives alphabetical order dressed up as urgency.

Status and severity also move completely independently -- an incident can be resolved and still marked critical, or still triaging and marked minor -- so inferring one from the other invents a relationship the product doesn't have. The one fixed thing is a status's `category` (`triage`, `active`, `learning`, `closed`), the same four words in every workspace underneath a status name that isn't; a client that wants to know whether an incident is actually over should read the category and never the display name.

## Sources

- Documentation: https://api-docs.incident.io/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve incidentio     # run it
cauldron verify incidentio -v # check every claim
```
