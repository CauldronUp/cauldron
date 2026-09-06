# Datadog

Emulates the Datadog API (v1), for local development and tests.

**20 conformance cases, 1 checked against the live API.**

Struck live against api.datadoghq.com/api/v1/monitor on 2026-09-05, and found the opposite of what this file claimed: "Datadog answers 403 where most providers answer 401" was backwards for this route. A bad or missing API key gets a plain 401 "Unauthorized", every time. Datadog does answer 403 to a bad credential -- but on /api/v1/validate, a different endpoint this Recipe does not model, which is exactly the kind of generalisation that does not survive being checked against the specific route it was claimed of.

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
