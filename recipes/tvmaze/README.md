# TVmaze

Emulates the TVmaze API (v1), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

The failure puts the reason in `name` and leaves
`message` empty. A missing show answers
`{"name": "Not Found", "message": "", "code": 0, "status": 404}` -- four fields,
and the one a client would read is blank. `name` carries the reason phrase, and
on every successful response `name` is the title of the show, so the same key
holds "Under the Dome" on a 200 and "Not Found" on a 404. `code` is 0, which is
neither an HTTP status nor an error code, and `status` is the HTTP status
restated inside the body.

**And "where does it air" is two mutually exclusive fields.** Every show carries
`network` and `webChannel` and exactly one is an object, so `show.network.name`
throws on half the catalogue and the check has to be written both ways round --
and inside them the shapes disagree, the network carrying a country object where
the web channel carries `country: null`. Around that: `runtime` and
`averageRuntime` are two fields where one is often null and nothing says which to
prefer; `schedule.time` is the empty string on a show that has a broadcast day;
`externals` holds three identifiers in three types, `null`, a number and the
string `"tt9140554"`; `summary` is HTML in a JSON string with no plain-text
companion; and `updated` is epoch seconds.

## Sources

- Documentation: https://www.tvmaze.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve tvmaze     # run it
cauldron verify tvmaze -v # check every claim
```
