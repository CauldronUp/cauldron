# Webflow

Emulates the Webflow API (2.0.0), for local development and tests.

**18 conformance cases, 2 checked against the live API on 2026-09-05.**

Webflow has no sandbox, so the CMS and collection cases still cite documentation. The credential shape needed no site at all, and checking it live found this Recipe's own sentence and code wrong.

## What writing this Recipe changed

This Recipe found a bug that had already shipped. Every timestamp field was
being filled in automatically, so a site that had never been published still
carried a `lastPublished`. The emulator was claiming events that never happened,
and no test written against it could have caught that -- the value was always
there, so nothing ever looked wrong.

## What checking it live found

No token at all and a present, wrong one are byte-identical: `{"message":"Request not authorized","code":"not_authorized",...}`, not `"unauthorized"`/`"The access token is invalid or has been revoked"` this Recipe had claimed. A path nothing declares answers a completely different, unrelated shape -- a raw framework 404 with `"code"` as a bare number -- recorded in the header rather than modelled, since this Recipe's error envelope carries a constant pair of fields on every failure that the route-not-found response does not send.

## Sources

- Documentation: https://developers.webflow.com/data/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve webflow     # run it
cauldron verify webflow -v # check every claim
```
