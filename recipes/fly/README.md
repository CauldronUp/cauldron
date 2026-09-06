# fly

Emulates the fly API (v1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real organisation; the refusal cases were struck live, unauthenticated, against api.machines.dev.

## What this Recipe found

`started` does not mean the application is up. A Fly machine has a `state` field that everyone reads, and three other fields on the same object can independently contradict it: `host_status` can be `unreachable` while the machine is still `started`, `cordoned` can be true while it's started and deliberately taking no traffic, and `checks` can be `critical` while the process runs and fails its own health check. A deploy script that waits for `state == "started"` and calls it done has waited for the least informative of four possible answers.

The identifier worth tracking is also not the one most integrations pick: `instance_id` changes on every version of the machine while the machine id stays stable across updates, so metrics or logs keyed on the instance id lose their history at every deploy. And `nonce`, the token needed to update a machine later, is returned exactly once, only on creation, and only if a lease duration was requested -- ask for a machine without one and there's no way to write to it afterward.

Waiting for a state change is a blocking endpoint (`GET .../wait`) rather than something a client polls for, and this Recipe deliberately does not answer it immediately: doing so would teach a client that waiting is free and instant, which is the opposite of what the real endpoint exists to demonstrate.

The live probe found the declared authentication error had never actually been reachable: it was named `unauthorized`, which nothing in this file wired to a credential failure, so it sat unused while every real refusal fell through to a generic default. It is now named correctly and carries the real sentence, "Authenticate: token validation error". An unrouted path and a wrong method on a real path both answer Go's own plain-text 404 before authentication is ever consulted.

## Sources

- Documentation: https://docs.machines.dev/
- Machine-readable description: https://docs.machines.dev/openapi.json, last checked 2026-08-31
  `cauldron drift fly` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fly     # run it
cauldron verify fly -v # check every claim
```
