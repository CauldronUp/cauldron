# replicate

Emulates the replicate API (v1), for local development and tests.

**19 conformance cases, 3 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real account. The credential check itself was verified directly against api.replicate.com on 2026-09-05.

## What this Recipe found

Checked live: a missing credential and a present-but-wrong one (whether the wrong scheme word or the right one with a fictitious value) answer sentences one word apart, "You did not pass an authentication token" versus "You did not pass a valid authentication token" -- both under the title "Unauthenticated", which this file had not modelled at all, and neither the "You did not provide a valid API token" this file had assumed. An unrouted path and a wrong method both land on the missing-credential sentence too, so the credential is checked first.

A newly created Replicate prediction returns immediately with `status: "starting"` and no `output` property at all, not null, absent, and only once it moves to `processing` does output start to appear, incrementally for streaming models. Code reading `prediction.output` right after creation gets undefined, and code waiting for output to become non-null is waiting for a key that does not exist yet.

`succeeded` does not mean the model produced what you wanted, either -- a model can succeed and return an empty array, a refusal, or a safety-filtered placeholder, and the status says succeeded for all of them identically. A cold model, Replicate scales to zero, costs a full minute of container boot time before it costs anything else, with no signal anywhere in the response distinguishing a cold run from a warm one, and billing is by the second of compute reported after the fact in `metrics.predict_time` rather than by request count, so two calls to the same model can differ tenfold in cost.

Output is also a link, not a value, and the link expires -- a succeeded prediction's output URL points to a file Replicate deletes after an hour, so a client that stores the URL instead of the bytes stores something that will 404 the next day.

## Sources

- Documentation: https://replicate.com/docs/reference/http
- Machine-readable description: https://api.replicate.com/openapi.json, last checked 2026-09-05
  `cauldron drift replicate` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve replicate     # run it
cauldron verify replicate -v # check every claim
```
