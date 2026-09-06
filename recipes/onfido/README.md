# Onfido

Emulates the Onfido API (v3.6), for local development and tests.

**21 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.onfido.com, no account and no key -- and this file's declared authentication failure matched exactly, word for word, whether no Authorization header was sent at all or the wrong scheme (Bearer) was.

## What this Recipe found

An Onfido check being `complete` says nothing about whether the person passed -- `status` and `result` are separate fields answering separate questions, and only once status reaches `complete` does result mean anything at all. There are three results, not two: `clear`, `consider`, and `unidentified`, and `consider`, the most common non-clear outcome, means a human has to look, which is neither a pass nor a failure. Treating it as either one is the kind of wrong somebody eventually complains about.

A check's own result is not simply the worst result among its reports -- a check can land on `consider` because exactly one report did, while every other report came back clear, and the reason for that lives on the individual report, not on the check as a whole. `sub_result`, which only appears on document reports, is where the real story is: `rejected` and `suspected` both look like failure but mean very different things, a bad photograph versus a suspected forgery.

Nothing is decided yet when a check is created -- the fields that will eventually carry a result are absent entirely rather than null, so code reading them right after creation gets `undefined`, and code that specifically checks for null finds nothing wrong at all.

## Sources

- Documentation: https://documentation.onfido.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve onfido     # run it
cauldron verify onfido -v # check every claim
```
