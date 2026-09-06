# Basecamp

Emulates the Basecamp API (3), for local development and tests.

**22 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Basecamp refuses a request outright when it doesn't carry a User-Agent identifying the calling application with a contact address -- something nearly every HTTP client sends by default, which is why the rule only ever gets discovered from a client that sends nothing. Everything in Basecamp is also soft-deleted: a trashed record still exists, is still reachable by id, and comes back with status "trashed" rather than a 404, so code that treats a successful fetch as proof something is live shows deleted content.

One fidelity gap: Basecamp pages results with a Link header, which Cauldron doesn't model.

## Sources

- Documentation: https://github.com/basecamp/bc3-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve basecamp     # run it
cauldron verify basecamp -v # check every claim
```
