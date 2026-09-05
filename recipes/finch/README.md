# finch

Emulates the finch API (v1), for local development and tests.

**17 conformance cases, 6 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real connection; the refusal cases were struck live, unauthenticated, against api.tryfinch.com.

## What this Recipe found

Finch tells you what it can't do, which is what separates it from a plain unified API. Every connection exposes a `/provider` endpoint listing the fields that specific payroll or HR system actually supports -- the only way to tell "this employee has no manager" from "this payroll system has never had a concept of managers," both of which otherwise arrive as the identical null. Almost nobody calls it, so most integrations treat an unsupported field as an empty one until a customer connects the wrong system.

The live probe found `/providers` genuinely public: it answers byte-identically with no credential and with an invented one, rather than merely tolerating an absent one. It also found a missing credential and a wrong one are different failures on every other route -- different names, different sentences -- neither matching what this file declared before the probe, and that an unrouted path and a wrong method on a real path both answer identically before authentication is ever consulted.

Batch endpoints answer 200 with per-item failures folded inside: asking about four employees can return two employees, a 404 and a 500 in one array under one successful status code, and code that checks the status before iterating treats the error objects as people. Some Finch connections are also "assisted" rather than API-driven -- a human logs into the provider on a schedule -- so the data can be days old and writes aren't possible, with `authentication_method` on the connection the only visible sign.

## Sources

- Documentation: https://developer.tryfinch.com/api-reference
- Machine-readable description: https://api.tryfinch.com/openapi.json, last checked 2026-09-05
  `cauldron drift finch` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve finch     # run it
cauldron verify finch -v # check every claim
```
