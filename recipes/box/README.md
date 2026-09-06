# Box

Emulates the Box API (2.0), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Two were struck live against api.box.com on 2026-09-05, and corrected a case that had asserted a flat JSON body for an unauthorized request. Box sends no body at all either way -- Content-Length 0 -- and puts the failure entirely in a WWW-Authenticate header, which names an absent token and a present, wrong one differently: "invalid_request" / "The access token was not found." for one, "invalid_token" / "The access token provided is invalid." for the other.

## What this Recipe found

Every Box object declares a type, which is the only way to tell a file from a folder in a listing -- they share an identifier space in appearance but not in fact, so file 12345 and folder 12345 are different objects, and code that ignores the type field happily fetches the wrong one. Errors also carry an HTTP status, a code, and a separate request id -- support asks for the request id specifically, so an integration that only logs the message has nothing useful to hand over when something goes wrong.

## Sources

- Documentation: https://developer.box.com/reference
- Machine-readable description: https://raw.githubusercontent.com/box/box-openapi/main/openapi.json, last checked 2026-08-30
  `cauldron drift box` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve box     # run it
cauldron verify box -v # check every claim
```
