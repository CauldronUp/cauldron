# api.video

Emulates the api.video API (v1), for local development and tests.

**14 conformance cases, 7 checked against the live API on 2026-09-03.**

## What this Recipe found

It has **no single answer to whether a video is
ready** -- three vocabularies, one of which has no failure value at all.

## Sources

- Documentation: https://docs.api.video/reference
- Machine-readable description: https://raw.githubusercontent.com/apivideo/api.video-api-specification/master/oas_apivideo.yaml, last checked 2026-09-03
  `cauldron drift apivideo` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve apivideo     # run it
cauldron verify apivideo -v # check every claim
```
