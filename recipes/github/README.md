# GitHub

Emulates the GitHub API (2022-11-28), for local development and tests.

**17 conformance cases, 15 checked against the live API on 2026-08-23.**

## Sources

- Documentation: https://docs.github.com/rest
- Machine-readable description: https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json, last checked 2026-08-30
  `cauldron drift github` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve github     # run it
cauldron verify github -v # check every claim
```
