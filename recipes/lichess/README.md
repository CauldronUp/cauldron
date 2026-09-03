# Lichess

Emulates the Lichess API (v1), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The game counts do not add up.** `count` reports
23,723 games in total and 11,329 won, 11,345 lost and 1,044 drawn -- which is
23,718 -- on a record where `playing` is 0. Five games are unaccounted for,
nothing in the object names them, and there is no fifth field to blame. It is
the kind of thing that survives forever, because each number is right on its own
and only the sum is wrong, and nobody sums them.

**And a list of links is separated by a Windows line ending, inside a JSON
string**: `"github.com/ornicar\r\nmas.to/@thibault"`. Splitting on `\n` leaves a
carriage return on the first, which renders as nothing and compares as
something. Two time units share the record and neither is named -- `createdAt`
and `seenAt` are epoch milliseconds while `playTime.total` is seconds -- so a
client dividing the wrong one by a thousand is out by a factor of a million and
gets a plausible answer either way. Sibling objects disagree about their keys:
`prov` is sent only when it is true, so one rating has it and the next does not.
And two 404s share the host, one `{"error": "Not found"}` and one the whole
Lichess website in HTML.

## Sources

- Documentation: https://lichess.org/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lichess     # run it
cauldron verify lichess -v # check every claim
```
