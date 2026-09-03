# Chess.com

Emulates the Chess.com API (chesscom), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The country is a URL.** A player's `country` is
`"https://api.chess.com/pub/country/US"` -- not `"US"` and not `"United
States"` -- so a client rendering `player.country` puts an API address on the
page, and getting the two letters costs a second request to a resource that is
`{"@id", "code", "name"}` and nothing else.

**And the player has two URLs that are not the same URL:** `"@id"` is
`https://api.chess.com/pub/player/hikaru` and `"url"` is
`https://www.chess.com/member/Hikaru`. Two hosts, two paths, two
capitalisations of one username, and nothing in either name saying which is
which -- and `@id` begins with a sigil, so it needs a subscript in most
languages. `twitch_url` is sent twice, once at the top level and once as
`streaming_platforms[0].channel_url`, byte-identical, so the field and the
structure meant to replace it both ship. `status` is a subscription tier,
`"premium"`, in a response whose HTTP status is 200. And the failure code is
zero: a player who does not exist answers `{"code": 0, ...}`, and so does a path
that does not exist, whose message is `Data provider not found for key
"/pub/player/hikaru/notanendpoint".` -- internal vocabulary, with the path
echoed back inside escaped quotes.

## Sources

- Documentation: https://www.chess.com/news/view/published-data-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve chesscom     # run it
cauldron verify chesscom -v # check every claim
```
