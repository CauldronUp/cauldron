# Typeform

Emulates the Typeform API (1), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

Typeform has no sandbox, so the response and form cases still cite documentation. The credential and routing shapes needed no account at all, and checking them live found this Recipe's own status code wrong.

## What writing this Recipe changed

It shared in a bug this collection found through Webflow: every timestamp field
was being filled in automatically, so a response somebody abandoned still
carried a `submitted_at`. An emulator that dates an event which never happened
teaches an integration to trust a field it should be checking for absence.

## What checking it live found

No credential at all is a 401, and a present, wrong bearer token is a 403 -- different statuses under the same `AUTHENTICATION_FAILED` code, not the single 401 this Recipe had for both. A path nothing declares is a 404, resolved before the credential is judged at all.

## Sources

- Documentation: https://www.typeform.com/developers/create/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve typeform     # run it
cauldron verify typeform -v # check every claim
```
