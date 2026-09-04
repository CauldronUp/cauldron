# finch

Emulates the finch API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Finch tells you what it can't do, which is what separates it from a plain unified API. Every connection exposes a `/provider` endpoint listing the fields that specific payroll or HR system actually supports -- the only way to tell "this employee has no manager" from "this payroll system has never had a concept of managers," both of which otherwise arrive as the identical null. Almost nobody calls it, so most integrations treat an unsupported field as an empty one until a customer connects the wrong system.

Batch endpoints answer 200 with per-item failures folded inside: asking about four employees can return two employees, a 404 and a 500 in one array under one successful status code, and code that checks the status before iterating treats the error objects as people. Some Finch connections are also "assisted" rather than API-driven -- a human logs into the provider on a schedule -- so the data can be days old and writes aren't possible, with `authentication_method` on the connection the only visible sign.

## Sources

- Documentation: https://developer.tryfinch.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve finch     # run it
cauldron verify finch -v # check every claim
```
