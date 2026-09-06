# Segment

Emulates the Segment API (v1beta), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real workspace. The credential check itself was verified directly against api.segmentapis.com on 2026-09-05.

## What this Recipe found

Checked live: an absent credential and a wrong one are not just different sentences, they are different HTTP statuses. No Authorization header answers 401, `{"type":"unauthorized","message":"Authorization header is required"}`; a syntactically fine but fictitious bearer answers 403, `{"type":"forbidden","message":"Not authorized to perform this operation"}` -- a key that is definitely wrong reads as forbidden rather than unauthenticated. This file had one message under one 401 for both.

The source collection sits two levels down, under `data.sources`, rather than at the top level or under a single obvious key -- code written for either of the usual shapes finds nothing. A disabled source also still exists, still has its write key, and still appears in every listing; it simply drops the events sent to it, silently, which is the failure mode that takes longest to notice.

## Sources

- Documentation: https://docs.segmentapis.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve segment     # run it
cauldron verify segment -v # check every claim
```
