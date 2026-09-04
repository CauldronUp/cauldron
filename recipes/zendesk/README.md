# Zendesk

Emulates the Zendesk API (v2), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The error envelope reuses `error` as a bare string code rather than an object -- a failure is `{"error": "RecordNotFound", "description": "..."}`, so code that reads `error.message` finds a string where it expected an object and typically throws on the property access instead of reporting the actual failure. A collection listing also reports its true total separately from the page length, so a pagination UI can't be built from the page alone.

## Sources

- Documentation: https://developer.zendesk.com/api-reference/ticketing/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zendesk     # run it
cauldron verify zendesk -v # check every claim
```
