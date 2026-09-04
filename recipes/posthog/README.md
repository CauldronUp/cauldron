# PostHog

Emulates the PostHog API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

PostHog's feature-flag API serves the rules, active, a rollout percentage, conditions, not the answer for any particular user; evaluating those rules against a person's properties happens locally in the SDK, so reading `active` and treating it as "this user has the flag on" is wrong for everybody outside the rollout. A rollout percentage of zero is not the same as the flag being off, either: the flag is still live and its conditions still evaluated, so anyone matching an override still gets it, which makes turning a flag off and setting it to 0% two different acts with very different blast radii.

Capture answers the same way no matter what is sent -- a bare `{status: 1}`, no identifier, no validation -- so an event with a malformed property is accepted here and silently dropped later, somewhere nobody is watching. A cohort is computed on a schedule rather than live: `is_calculating` is true while a recalculation runs and the count shown meanwhile is the previous one, so adding a property to a person does not put them in a cohort immediately and nothing in the response says when it will update.

After merging two person records, the survivor carries both original ids, so a lookup by the older id still resolves and the count of distinct people drops without anyone actually being deleted.

## Sources

- Documentation: https://posthog.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve posthog     # run it
cauldron verify posthog -v # check every claim
```
