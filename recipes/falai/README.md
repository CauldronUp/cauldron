# Fal.ai

Emulates the Fal.ai API (v1), for local development and tests.

**9 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

Its **queue position vanishes rather than reaching zero**,
so a caller treating a missing number as zero reports next-up for something
already running.

## Sources

- Documentation: https://docs.fal.ai/model-apis/quick-start/queue-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve falai     # run it
cauldron verify falai -v # check every claim
```
