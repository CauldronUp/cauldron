# Lob

Emulates the Lob API (2020-02-11), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Lob's test keys are the usual trap: a test letter is always accepted, always renders, and never gets returned to sender, so an integration validated in test mode can still end up posting to addresses nobody lives at. Address verification in test mode reports a deliverability the real service would not agree with.

Deliverability itself is a graded verdict, not a boolean -- `deliverable`, `deliverable_unnecessary_unit`, `deliverable_incorrect_unit`, `deliverable_missing_unit` and `undeliverable` are the five answers, and three of the five are yes with a caveat. Code that checks for equality with `deliverable` rejects addresses that would actually arrive. A letter's tracking events are the only way to know it went anywhere: `processed_for_delivery` is not `delivered`, and both `re-routed` and `returned_to_sender` can appear after a letter already looked fine.

## Sources

- Documentation: https://docs.lob.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lob     # run it
cauldron verify lob -v # check every claim
```
