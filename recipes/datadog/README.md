# Datadog

Emulates the Datadog API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It sends its errors as an array of bare strings, so a client reading `.message`
from each entry finds `undefined` on every one.

## Sources

- Documentation: https://docs.datadoghq.com/api/latest
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve datadog     # run it
cauldron verify datadog -v # check every claim
```
