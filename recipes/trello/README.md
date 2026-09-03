# Trello

Emulates the Trello API (1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

This is the clearest case in the collection for reproducing something rather
than improving on it. Trello's credentials travel in the query string, which is
a bad idea, and its failures are plain text rather than JSON.

A fake that accepted a header instead would hide the fact that the key ends up
in access logs. One that answered in JSON would hide that calling `.json()` on a
real error throws. Both would be more pleasant and less true.

## Sources

- Documentation: https://developer.atlassian.com/cloud/trello/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve trello     # run it
cauldron verify trello -v # check every claim
```
