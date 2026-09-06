# Dropbox

Emulates the Dropbox API (2), for local development and tests.

**14 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource shapes cite documentation rather than an observation on a real account; the refusal cases were struck live, unauthenticated, against api.dropboxapi.com.

## What writing this Recipe changed

This Recipe found a limitation in the conformance checker rather than in the
emulator. Dropbox names a field `.tag`, where the leading dot is part of the
name, and the checker split every path on dots -- so there was no way to assert
on that field at all.

The emulator had been sending it correctly the whole time, while every case that
mentioned it failed.

The live probe found that a missing Authorization header and a wrong one fail completely differently: no header at all gets plain text naming the function that was called, while a well-formed invented token gets the JSON `invalid_access_token` shape this file already declared. An unrouted path answers with Dropbox's own branded marketing-site 404 page, a different backend entirely from the API. A wrong method on a real path turned out not to be modellable at all: Dropbox resolves the function from the path alone and checks the credential before ever looking at the HTTP method, which this format cannot reproduce alongside an unrouted path correctly bypassing authentication -- stated in the header rather than asserted around.

## Sources

- Documentation: https://www.dropbox.com/developers/documentation/http/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve dropbox     # run it
cauldron verify dropbox -v # check every claim
```
