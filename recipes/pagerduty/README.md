# PagerDuty

Emulates the PagerDuty API (2), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It authenticates with `Token token=` rather than `Bearer`, which is the kind of
detail a fake accepting anything would hide until production.

It also exposed a real inconsistency in the runtime later on. When the
conformance checker learned to compare a scalar's kind as well as its value, the
nested error style turned out to be sending PagerDuty's numeric codes as
text.

## Sources

- Documentation: https://developer.pagerduty.com/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pagerduty     # run it
cauldron verify pagerduty -v # check every claim
```
