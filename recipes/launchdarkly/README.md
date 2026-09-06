# LaunchDarkly

Emulates the LaunchDarkly API (v2), for local development and tests.

**18 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against app.launchdarkly.com, no account and no key, on three different paths and three different bad-credential shapes. This file declared the 401 body as `{"code":"unauthorized","message":"Invalid access token"}`, and the real, consistent answer is `{"code":"unauthorized","message":"Invalid account ID header"}` -- a sentence that does not mention a token at all. Fixed below.

## What this Recipe found

The Authorization header carries the API key with no scheme at all -- not `Bearer`, not `Token`, just the key. Every other provider in this collection wants a prefix, so the habit of adding one is exactly wrong here, and produces a 401 that reads like a bad key rather than a malformed header.

A flag also has no single on/off state: it has one per environment, nested under an `environments` map keyed by environment key, so "is this flag on" isn't answerable without saying where, and code reading `flag.on` directly finds `undefined`. A variation is referenced by index rather than by value -- `fallthrough.variation` is `0` or `1`, meaning true or false in whatever order the variations were originally defined -- so reading the index as the value can give `0` for true. And an archived flag still exists and still serves: it just disappears from the default listing, which makes it look deleted while the SDK keeps evaluating it everywhere it's still referenced.

## Sources

- Documentation: https://apidocs.launchdarkly.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve launchdarkly     # run it
cauldron verify launchdarkly -v # check every claim
```
