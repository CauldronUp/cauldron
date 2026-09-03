# what3words

Emulates the what3words API (v3), for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its **two credential failures come from two machines**.
`MissingKey` is answered by CloudFront's edge with no CORS headers; `InvalidKey`
comes from the origin with them, so a browser sending no key gets a response its
own CORS check rejects.

## Sources

- Documentation: https://developer.what3words.com/public-api/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve what3words     # run it
cauldron verify what3words -v # check every claim
```
