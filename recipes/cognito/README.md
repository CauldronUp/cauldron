# AWS Cognito

Emulates the AWS Cognito API (2016-04-18), for local development and tests.

**14 conformance cases, 7 checked against the live API on 2026-09-02.**

## What this Recipe found

**The wrong region is the same as no pool** -- the
identical sentence, differing only in the id echoed back, for an id whose first
half names the region.

## Sources

- Documentation: https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/Welcome.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cognito     # run it
cauldron verify cognito -v # check every claim
```
