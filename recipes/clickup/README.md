# ClickUp

Emulates the ClickUp API (v2), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Two were struck live against api.clickup.com on 2026-09-05, and found a case this file had never modelled: an entirely absent Authorization header is not even the same status code as a present, wrong one. A wrong token is 401 OAUTH_019; no header at all is 400 OAUTH_017 -- the one failure in this collection where the status code itself, not just the message, depends on whether a credential was sent.

## What this Recipe found

A ClickUp task's status is an object, not a string -- it carries a name, a colour and a type, and the type is what actually says whether the status counts as done. Code that compares task.status to "complete" is comparing an object to a string and is false forever. Errors are also their own shape: "err" carries the message and "ECODE" the machine-readable code, with nothing named message or error, so code written for any other provider finds nothing at all.

## Sources

- Documentation: https://developer.clickup.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve clickup     # run it
cauldron verify clickup -v # check every claim
```
