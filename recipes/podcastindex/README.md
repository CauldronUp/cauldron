# Podcast Index

Emulates the Podcast Index API (1.0), for local development and tests.

**11 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

Index's, **the signed API that actually answers** -- five
distinct sentences in a strict order where its two siblings collapse everything
into one.

## Sources

- Documentation: https://podcastindex-org.github.io/docs-api/
- Machine-readable description: https://podcastindex-org.github.io/docs-api/pi_api.json, last checked 2026-09-02
  `cauldron drift podcastindex` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve podcastindex     # run it
cauldron verify podcastindex -v # check every claim
```
