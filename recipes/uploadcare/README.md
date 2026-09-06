# uploadcare

Emulates the uploadcare API (v0.7), for local development and tests.

**17 conformance cases, 3 checked against the live API on 2026-09-05.**

The storage and expiry cases still cite documentation, since a real file needs a real upload. The credential shape needed no account at all, and checking it live found three verdicts where this Recipe had claimed one sentence.

## What this Recipe found

A file that exists today isn't a file that will exist tomorrow. Uploading succeeds, the file gets a UUID, and the CDN serves it -- and unless something calls `store` afterward, it's deleted after twenty-four hours. Whether that happens automatically is a per-project setting invisible from the code, so identical client code works on one project and silently loses files on another. `datetime_stored` reads `null` both before storing and while a store request is still in flight, so checking it says "safe" right up until it doesn't.

A deleted file also still answers: `datetime_removed` gets set, but the rest of the metadata stays intact, so a client checking for a 404 gets a 200 and a working-looking object instead.

## What checking it live found

No credential at all answers `"Authentication credentials were not provided."`; a present, wrong `Uploadcare.Simple` pair names the public half instead, `"Public key {name} not found."`; and a scheme word Uploadcare does not recognise, such as `Basic`, is answered exactly like no credential at all -- not the third sentence an existing case had claimed and never run.

## Sources

- Documentation: https://uploadcare.com/api-refs/rest-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve uploadcare     # run it
cauldron verify uploadcare -v # check every claim
```
