# ROR

Emulates the ROR API (v2), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**A path that does not exist says your authentication is
missing, on an API that has none.** `GET /v2/nosuchthing` answers `403
{"message": "Missing Authentication Token"}` -- the gateway's default for an
unmatched route, on a registry with no credential to get wrong and no header
that would make it go away. It sends whoever typed the path wrong looking for an
API key that does not exist, past the one place the fault actually is, and the
403 means a client branching on the status treats a typo as a permissions
problem.

**And the identifier is a whole URL, inside the path.** `GET
/v2/organizations/https://ror.org/02mhbdp94` returns a record whose `id` is
`"https://ror.org/02mhbdp94"` -- so fetching one means putting a URL, scheme and
slashes and all, into the path of another, and a client that percent-encodes its
path segments correctly gets a 404 while one that does not gets the record.
There are two failure shapes and the one that is not your fault answers 200: an
invalid id is 404 with `{"errors": ["..."]}`, an array of bare strings, and a
page outside the range is that same shape at 200. The record carries two schema
versions and says which is which, `"1.0"` for when it was created and `"2.1"`
for when it was last touched. And the display name is chosen by finding
`"ror_display"` inside an array inside an array element.

## Sources

- Documentation: https://ror.readme.io/v2/docs/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ror     # run it
cauldron verify ror -v # check every claim
```
