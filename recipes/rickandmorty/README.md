# Rick and Morty

Emulates the Rick and Morty API (rickandmorty), for local development and tests.

**11 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Morty's, where asking for a character by a name
instead of a number is a 500. `GET /api/character/abc` answers
`{"error": "Hey! you must provide an id"}` with a server-fault status, so the
retry wrapper retries, the alert fires and the dashboard records an outage --
for a request that was never going to work. The 404 beside it, for a number that
is simply not a character, is the correct status, so the API can produce the
right answer and does not for the wrong input.

**And `"unknown"` is the only lowercase value in the enum.** `status` is
`"Alive"`, `"Dead"` or `"unknown"`, so the one value meaning "we have no answer"
is also the one that breaks a switch written from the other two. One object says
"no value" two ways at once: on a character whose origin is not known,
`origin.name` is the string `"unknown"` and `origin.url` is the empty string --
a word in one field and an absence in its neighbour, inside the same two-key
object. Related records are URLs rather than identifiers, so joining anything
means parsing an integer back out of a string. And three failures carry three
sentences and two statuses, only one of which ends in a full stop.

## Sources

- Documentation: https://rickandmortyapi.com/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rickandmorty     # run it
cauldron verify rickandmorty -v # check every claim
```
