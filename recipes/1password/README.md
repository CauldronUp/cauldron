# 1Password

Emulates the 1Password API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

1Password's own official Go SDK disagrees with itself about whether a listing includes secrets. GET /v1/vaults/{v}/items returns items without their fields; GET .../items/{i} returns the same item complete with values attached. The SDK's GetItems() decodes the listing as though every item were whole, while GetItemsByTitle() forty lines later re-fetches each one individually to get the real thing. Code written against GetItems() reads item.fields on a response that never carried any and finds nil, silently, using 1Password's own client.

The id 1Password calls a UUID is not one -- twenty-six lowercase letters and digits, no hyphens, matched by a hand-written regexp rather than any UUID library -- while a section id one level down really is RFC 4122 shaped. And a missing token and a rejected one produce the identical message, so absent and wrong collapse to one sentence; not-found is a different status and sentence entirely.

This models Connect (self-hosted, no shared host to observe) rather than the newer Service Accounts API (which needs an account to create one, ruled out by this project's sourcing rules), and every case is read from 1Password's hand-written Go SDK rather than observed live.

## Sources

- Documentation: https://developer.1password.com/docs/connect/connect-api-reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve 1password     # run it
cauldron verify 1password -v # check every claim
```
