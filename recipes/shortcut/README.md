# Shortcut

Emulates the Shortcut API (v3), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Whether a story is "done" isn't a field you can read directly -- workflow states are ids a workspace admin defines, so "Ready for Review" and "Done" are opaque identifiers rather than names, and Shortcut's own `completed`/`started` booleans, derived from the state's type, are the only stable thing to branch on. An archived story is also still returned by a plain search: archiving isn't deleting, so an unfiltered count includes work somebody deliberately set aside.

## Sources

- Documentation: https://developer.shortcut.com/api/rest/v3
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve shortcut     # run it
cauldron verify shortcut -v # check every claim
```
