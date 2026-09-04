# braze

Emulates the braze API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Braze's whole API rests on one property, "message", carrying both the literal string "success" and the text of a failure -- there's no separate error object, code, or type, so a client that checks the HTTP status and reads message for logging produces a log full of the word success with an occasional real error nobody notices, because Braze answers some rejections with a 2xx.

A campaign in a listing is also not the campaign: /campaigns/list returns five properties, /campaigns/details returns a different and larger set, and the two don't even agree on the timestamp's name (last_edited versus updated_at), so reading channels or description off a list entry gets undefined. The identifier for the details call also arrives in the query string rather than the path, so a client built around path-shaped URLs has nowhere natural to put it.

The segment export is asynchronous in a way that's easy to miss: POST /users/export/segment answers 201 immediately with a success message and a storage prefix, and the actual data lands in cloud storage minutes later -- a test that reads users directly off that response reads nothing, forever.

## Sources

- Documentation: https://www.braze.com/docs/api/basics
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve braze     # run it
cauldron verify braze -v # check every claim
```
