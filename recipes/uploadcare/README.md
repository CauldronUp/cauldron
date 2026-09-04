# uploadcare

Emulates the uploadcare API (v0.7), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A file that exists today isn't a file that will exist tomorrow. Uploading succeeds, the file gets a UUID, and the CDN serves it -- and unless something calls `store` afterward, it's deleted after twenty-four hours. Whether that happens automatically is a per-project setting invisible from the code, so identical client code works on one project and silently loses files on another. `datetime_stored` reads `null` both before storing and while a store request is still in flight, so checking it says "safe" right up until it doesn't.

A deleted file also still answers: `datetime_removed` gets set, but the rest of the metadata stays intact, so a client checking for a 404 gets a 200 and a working-looking object instead.

## Sources

- Documentation: https://uploadcare.com/api-refs/rest-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve uploadcare     # run it
cauldron verify uploadcare -v # check every claim
```
