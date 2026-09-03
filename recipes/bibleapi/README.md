# Bible API

Emulates the Bible API API (bibleapi), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

API's, where **the same verse comes back with different
whitespace depending on the translation, and not the same kind of whitespace.**
The World English Bible pads `John 3:16` with a leading newline and two trailing
ones; the King James pads it with ten leading spaces and one trailing newline.
One query parameter apart, on the same path, for the same verse -- so a client
that renders the text raw gets a blank line for one and an indent for the other,
and a client comparing two translations compares their padding as well as their
words.

**And the text is sent twice, both copies untrimmed:** the top-level `text` and
`verses[0].text` are byte-identical, padding included. The two endpoints share
no keys at all -- a reference answers `{reference, verses, text, translation_id,
translation_name, translation_note}` and `/data/web/random` answers
`{translation, random_verse}`, not one field name in common. The same concept
carries two names, `book_name` on one and `book` on the other, and the string
`"Public Domain"` arrives as `translation_note` in one place and `license` in
the other. The two 404s do not come from the same place: an unknown book is the
application's `{"error": "not found"}`, and a chapter past the end of a book
never reaches the application at all, answering nginx's own HTML instead.

## Sources

- Documentation: https://bible-api.com/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bibleapi     # run it
cauldron verify bibleapi -v # check every claim
```
