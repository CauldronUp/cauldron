# Segment

Emulates the Segment API (v1beta), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The source collection sits two levels down, under `data.sources`, rather than at the top level or under a single obvious key -- code written for either of the usual shapes finds nothing. A disabled source also still exists, still has its write key, and still appears in every listing; it simply drops the events sent to it, silently, which is the failure mode that takes longest to notice.

## Sources

- Documentation: https://docs.segmentapis.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve segment     # run it
cauldron verify segment -v # check every claim
```
