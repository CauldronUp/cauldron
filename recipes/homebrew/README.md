# Homebrew

Emulates the Homebrew API (v1), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-25.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

An object called `versions` holds two strings and a
boolean -- `{"stable": "8.21.0", "head": "HEAD", "bottle": true}` -- so one of
its three values is a version, `head` is a git ref spelled the same way on every
formula that has one, and `bottle` says whether a prebuilt binary exists. Four
more fields describe a computer that is not involved: `installed` is `[]`,
`pinned` and `outdated` are `false` and `linked_keg` is `null` on everything the
API serves, because it is a static document on a CDN and the same schema is what
`brew info --json` prints locally. A missing formula answers a **full HTML
page** from a path ending in `.json`, which is the third way this collection has
seen that refused. And Ruby symbols arrive as strings carrying their colon:
`keg_only_reason.reason` is `":provided_by_macos"` and every bottle's `cellar`
is `":any"`.

## Sources

- Documentation: https://formulae.brew.sh/docs/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve homebrew     # run it
cauldron verify homebrew -v # check every claim
```
