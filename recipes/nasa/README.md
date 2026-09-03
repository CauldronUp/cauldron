# NASA

Emulates the NASA API (v1), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

**One host and one key front three different APIs.**
The gateway, APOD and NEO answer three incompatible error envelopes --
`{"error": {"code", "message"}}`, `{"code", "msg", "service_version"}` and
`{"code", "http_error", "error_message", "request"}` -- so a client that
unwraps one receives the other two as `undefined`. APOD leaks a Python
`strptime` exception verbatim as its message; NEO answers from Java behind
Heroku and echoes back a request URL with its own path segment missing. A date
in the future and a date before the archive begins produce byte-identical
messages. And a reversed date range is not refused but silently repaired: ask
for the 5th to the 1st and `links.self` echoes the 1st to the 5th, so a caller
who transposed two variables gets a plausible answer to a question they did not
ask.

## Sources

- Documentation: https://api.nasa.gov/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nasa     # run it
cauldron verify nasa -v # check every claim
```
