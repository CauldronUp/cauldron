# Deck of Cards

Emulates the Deck of Cards API (deckofcards), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

Of Cards', where a failed draw empties the deck and
hands you the cards anyway. Asking a fifty-two card deck for 999 cards does not
refuse: it draws all fifty-two, sets `remaining` to 0, says `"success": false`,
returns them, and answers HTTP 200. There is no safe way to read that. A client
that checks `success` and discards the body has thrown away fifty-two real cards
and left the deck empty; a client that ignores `success` has fifty-two it did not
ask for the right number of. Either way the deck is spent, and the request that
spent it is the one the API called a failure.

**And the ten is a zero.** Card codes are two characters, so the ten of spades is
`"0S"` -- while its `value` field says the string `"10"`. Anything parsing a code
by taking its first character reads a nought where a ten should be. `value` is a
string for numbers and a word for faces, so sorting a hand by it puts `"10"`
before `"6"` and both before `"ACE"`. Every image is in the response twice, as
`image` and `images.png`, character for character. And a count that is not a
number answers 500 with Django's default error page in HTML, from a JSON API --
the framework's page, not anything this API wrote.

## Sources

- Documentation: https://deckofcardsapi.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve deckofcards     # run it
cauldron verify deckofcards -v # check every claim
```
