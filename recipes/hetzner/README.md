# hetzner

Emulates the hetzner API (v1), for local development and tests.

**20 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real project, because a token creates real machines and bills for them. The refusal cases were struck live, unauthenticated, against api.hetzner.cloud.

## What this Recipe found

Powering off a server does not answer with the server. Every mutation in this API -- power off, attach a volume, change a type, rebuild -- answers 201 with an action: a job with `status: running`, a progress percentage, and a null finished timestamp. The machine is still on. Whether it actually goes off is a separate request to a different endpoint about a different object, and nothing here waits for you. A test that powers a server off and immediately asserts it's off passes against a mock and is a race against the real thing.

An action can fail well after its 201 succeeded -- `status` moves from `running` to `success` or `error` on its own timeline, and the error, when it comes, is an object (`{code, message}`), not a string, so logging it directly prints `[object Object]`. A server itself has nine possible statuses, and "not running" covers eight of them (initializing, starting, stopping, off, deleting, migrating, rebuilding, unknown), so code that only branches on `running` versus `off` mishandles seven cases, two of which are just ordinary booting.

And the one that costs an actual machine: creating a server without an SSH key returns a `root_password` exactly once, alongside the action, and no other endpoint ever surfaces it again. A client that reads the action out of that response and discards the rest ends up with a server it can never reach.

The live probe found this file's authentication error wrong: the real sentence is "the token you have provided is invalid," not "unable to authenticate," and a missing token entirely gets a third sentence, "token is required," under the same code. An unrouted path answers a distinct 404 before authentication is ever consulted -- though a wrong method on a real path does not follow the same rule, and answers the credential refusal instead, a narrower distinction this file's format cannot carry alongside the first.

## Sources

- Documentation: https://docs.hetzner.cloud/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hetzner     # run it
cauldron verify hetzner -v # check every claim
```
