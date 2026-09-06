# PostHog

Emulates the PostHog API (v1), for local development and tests.

**17 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a project this Recipe cannot fabricate. The credential check itself was verified directly against us.posthog.com on 2026-09-05.

## What this Recipe found

Checked live: an absent credential and a fictitious one are two different DRF error bodies, not one message repeated. No Authorization header at all answers `{"code":"not_authenticated","detail":"Authentication credentials were not provided."}`; a syntactically fine but fake personal key answers `{"code":"authentication_failed","detail":"Personal API key found in request Authorization header is invalid."}` -- a different code and a different sentence. This file had modelled a single generic message under `authentication_failed` for both.

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
