# GitHub

Emulates the GitHub API (2022-11-28), for local development and tests.

**21 conformance cases, 15 checked against the live API on 2026-08-23.**

## What this Recipe found

This Recipe exists as much to prove the emulator format generalises as to be useful on its own. GitHub disagrees with Stripe, this collection's original template, on two basic shapes: a list is a bare JSON array rather than an object with a `data` key, and identifiers are plain integers rather than prefixed strings. A client written by analogy with a Stripe-shaped API unmarshals a list into the wrong type and treats an id as a string when it's actually a number.

Issues and labels are also scoped to a repository -- `/repos/{owner}/{repo}/issues` only ever returns that repository's records, with no cross-repo listing modelled here.

## Sources

- Documentation: https://docs.github.com/rest
- Machine-readable description: https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json, last checked 2026-08-30
  `cauldron drift github` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve github     # run it
cauldron verify github -v # check every claim
```
