# Packagist

Emulates the Packagist API (p2), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-08-26.**

## What this Recipe found

It is the sharpest of the registry findings: the
versions array is **a chain of diffs**. `GET /p2/monolog/monolog.json` returns
eighty-seven entries where the first has twenty-one keys and the second has
eight, because each entry after the first carries only the fields that differ
from the one before it. So `packages["monolog/monolog"][1].license` is undefined
-- not because 3.9.0 has no licence, but because it has the same one as 3.10.0.
The chain runs **newest first**, so the deltas apply backwards in version order.
A field that goes away is the string `"__unset"`. And the only signal that any
of it is happening is a sibling key at the top of the document, `"minified":
"composer/2.0"` -- nothing inside an entry says it is a patch.

## Sources

- Documentation: https://packagist.org/apidoc
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve packagist     # run it
cauldron verify packagist -v # check every claim
```
