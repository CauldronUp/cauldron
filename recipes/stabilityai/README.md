# Stability AI

Emulates the Stability AI API (v1), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-02.**

## What this Recipe found

AI's, which **hands the wrong key back masked** -- and the mask
is exactly as long as the value it hides.

## Sources

- Documentation: https://platform.stability.ai/docs/api-reference
- Machine-readable description: https://raw.githubusercontent.com/Stability-AI/rest-api-support/main/spec/v1.yaml, last checked 2026-09-02
  `cauldron drift stabilityai` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve stabilityai     # run it
cauldron verify stabilityai -v # check every claim
```
