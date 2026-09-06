# Sumsub

Emulates the Sumsub API (unversioned), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A rejection here isn't necessarily final. Once `reviewStatus` reaches `completed`, `reviewAnswer` is a blunt GREEN or RED with no middle value -- but `reviewRejectType` sits beside it as FINAL or RETRY, and RETRY means the same applicant id can walk the failed steps again and come back completed a second time with a different answer. Nothing else on the resource says a decision might still move; a client that reads `reviewAnswer` alone can't tell a closed case from one resubmission away from flipping to GREEN.

Sumsub is also the most forthcoming of the identity-verification providers in this collection about why a decision went the way it did: rejection reasons, a moderator's comment, and a client-facing comment all live directly on the same `reviewResult` object as the answer, rather than requiring a second request the way comparable providers do.

## Sources

- Documentation: https://docs.sumsub.com/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sumsub     # run it
cauldron verify sumsub -v # check every claim
```
