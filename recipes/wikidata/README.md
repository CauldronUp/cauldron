# Wikidata

Emulates the Wikidata API (wikibase/v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

The English label is not there and the name is under
a code that is not a language. Douglas Adams's item carries 75 labels and `en` is
not one of them -- nor `de`, `fr`, `es`, `it`, `nl`, `pt`, `pl` or `sv`, because
every Latin-script language was consolidated into a single entry under `mul`,
which is ISO 639-3 for "multiple languages". Asking the label endpoint directly
confirms it: `/labels/en` is a 404 and `/labels/mul` is `"Douglas Adams"`. So
`item.labels[userLang] ?? item.labels.en`, which is what every client writes,
shows nothing at all -- while `item.descriptions.en` beside it works, because the
descriptions were not consolidated. Two sibling objects, keyed the same way,
disagreeing about which languages exist.

**And every key in `statements` is an opaque number** -- 312 of them on this
item, and nothing in the response says what any means. A value is
`{"type": "value", "content": ...}` where `content` is a string for an item
reference and an object for a date; a date is `time` with a leading `+` on the
year, an integer `precision` nothing decodes, and a `calendarmodel` that is a URL
you must resolve to learn the calendar is Gregorian. A statement identifier is
`Q42$F078E5B3-...` -- the item, a dollar sign and a UUID -- and the item's own
case is not consistent across them: 285 begin `Q42$` and 27 begin `q42$`. A
missing item and a missing label are the same body with one nested field
different.

## Sources

- Documentation: https://www.wikidata.org/wiki/Wikidata:REST_API
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wikidata     # run it
cauldron verify wikidata -v # check every claim
```
