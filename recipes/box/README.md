# Box

Emulates the Box API (2.0), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

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
