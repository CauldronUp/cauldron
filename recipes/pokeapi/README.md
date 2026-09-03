# PokeAPI

Emulates the PokeAPI API (v2), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

The description is still formatted for a Game Boy.
`flavor_text` comes back as
`"A strange seed was\nplanted on its\nback at birth.\fThe plant sprouts\nand
grows with\nthis POKéMON."` -- three things in one string. The newlines are hard
wraps at the original 1996 screen's column width, so re-flowing them in a modern
layout breaks mid-phrase. The `\f` is a **form feed**, U+000C, the page break
from the handheld's text box, sitting in a JSON string today. And the name is
spelled with a lowercase é between two capitals, because the font in those games
had no capital É. None of it is escaped, flagged or duplicated in a cleaned-up
field.

**And there are twenty-eight English descriptions, not one.**
`flavor_text_entries` is a flat array of a hundred and two, one per game version
per language, each carrying `language` and `version` as references -- so "the
English description" is a filter that returns twenty-eight answers and picking
the first is picking Red from 1996. Around that: a miss is the bare plain text
`Not Found` with a 404, from an API whose every success is JSON; a listing row is
a `name` and a `url` and **no identifier at all**, so the id exists only as a
path segment inside that URL and anything more than the name costs a request per
row; `previous` and `next` are both always present and each is null at its own
end, so `"next" in body` is true on every page; the lookup is case-insensitive
and says nothing about it; a type is `{slot, type: {name, url}}`, so the name is
two levels down and the array index is not the slot; and four of the ten
top-level `sprites` keys are present and null.

## Sources

- Documentation: https://pokeapi.co/docs/v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pokeapi     # run it
cauldron verify pokeapi -v # check every claim
```
