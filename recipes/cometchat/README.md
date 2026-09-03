# CometChat

Emulates the CometChat API (v3), for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-01.**

## What this Recipe found

Its conversation **has no create route at all** -- it is a
computed view of a group, and exists the moment the group does.

## Sources

- Documentation: https://www.cometchat.com/docs
- Machine-readable description: https://www.cometchat.com/docs/chat-apis.json, last checked 2026-09-01
  `cauldron drift cometchat` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cometchat     # run it
cauldron verify cometchat -v # check every claim
```
