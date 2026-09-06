# Shortcut

Emulates the Shortcut API (v3), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

Shortcut has no sandbox, so the story/epic/iteration cases still cite documentation rather than a workspace. The three that do not need one -- what an unauthenticated request gets back -- were struck live and found the Recipe's own claim wrong.

## What this Recipe found

Whether a story is "done" isn't a field you can read directly -- workflow states are ids a workspace admin defines, so "Ready for Review" and "Done" are opaque identifiers rather than names, and Shortcut's own `completed`/`started` booleans, derived from the state's type, are the only stable thing to branch on. An archived story is also still returned by a plain search: archiving isn't deleting, so an unfiltered count includes work somebody deliberately set aside.

**And the credential case this Recipe already had was wrong.** It claimed sending `Authorization` instead of `Shortcut-Token` gets a plain `{"message":"Unauthorized"}`. Live, Shortcut reads only its own header, so that request is not a wrong credential -- it's no credential -- and gets `{"message":"Sorry, the organization context for this request is missing...","tag":"organization2_missing"}`, a different sentence with a `tag` field this Recipe never modelled. A garbage-but-present `Shortcut-Token` gets a third answer: `{"message":"Unauthorized","tag":"unauthorized"}`. Both verdicts arrive identically for a path that does not exist, which says the credential is judged before the request is routed at all.

## Sources

- Documentation: https://developer.shortcut.com/api/rest/v3
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shortcut     # run it
cauldron verify shortcut -v # check every claim
```
