# Neon

Emulates the Neon API (v2), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Almost nothing is finished when a Neon response arrives. Creating a branch answers 201 with the branch already in state `init`, which Neon's own docs define as "being created but not available for querying" -- the 201 is a receipt for having asked, not a promise the branch works, and code that connects to it on the next line connects to nothing. A branch also carries both `current_state` and `pending_state`, and they disagree exactly while something is in progress: reading either one alone answers a different, incomplete question about whether it is ready.

The create response's `operations` array is the only real way to know when a branch becomes usable, sitting in the same body a caller is likely to skim past looking for the branch id. An operation's status has eight values, not two -- `finished` is the only success, and `failed`, `error`, `cancelled`, and `skipped` are all just different ways of having stopped without working, so a poll written as `while status === 'running'` exits early on `scheduling` and calls the job done.

`default` and `primary` are both fields on a branch, with `primary` deprecated in favor of `default` -- the Recipe sends both so that code reading the wrong one does not silently pass by accident.

## Sources

- Documentation: https://neon.tech/api_spec/release/v2.json
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve neon     # run it
cauldron verify neon -v # check every claim
```
