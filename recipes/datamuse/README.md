# Datamuse

Emulates the Datamuse API (datamuse), for local development and tests.

**10 conformance cases, all of them checked against the live API on 2026-08-29.**

## What this Recipe found

A definition is a part of speech, a tab character
and a sentence, in one string. `"defs": ["n\tA large amount. "]` packs two
fields into a JSON string with a control character between them, inside a format
that has arrays and objects for exactly this -- and every entry ends in a space,
so a client that splits on the tab and renders index 1 shows trailing whitespace
it never asked for.

**And `tags` is three vocabularies in one array of strings.** `["n", "v",
"pron:S L UW1 "]` holds bare part-of-speech codes beside a colon-prefixed
key-value pair, so `tags.includes("n")` works and getting the pronunciation
means scanning the same array for a prefix and splitting on a colon -- and the
pronunciation carries a trailing space too. Three queries on one endpoint answer
three different field sets, and `md=` explains only one of the differences:
`?ml=ocean` gets tags and no syllable count, `?rel_rhy=blue` gets a syllable
count nobody asked for and no tags at all, and adding `md=dpsr` to the rhyme
query gets all five. An unknown parameter answers `[]` with 200, so a typo in a
parameter name is indistinguishable from a query that matched nothing. And
`score` is an integer with no scale anywhere -- 40,041,792 for a synonym and
10,051 for a spelling suggestion, three orders of magnitude apart on one host,
with nothing saying what either means.

## Sources

- Documentation: https://www.datamuse.com/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve datamuse     # run it
cauldron verify datamuse -v # check every claim
```
