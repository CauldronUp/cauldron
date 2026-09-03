# Scryfall

Emulates the Scryfall API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The prices are five nulls and one string.** Six
fields about money, five holding nothing and the sixth holding `"46.91"` -- so
summing a basket means parsing every value that is not null, and finding out
which those are means checking all six. Scryfall is a well-built API, which is
what makes it worth describing: the traps are decisions rather than
carelessness, and that is the harder kind to notice because everything around
them is right.

**Four failures, four typographies, and one of them is minified.** A missing
card gets curly double quotes around what was asked for; a missing parameter
gets Markdown backticks; an empty search gets `You didn‘t enter anything` --
where the apostrophe is U+2018, a *left* single quotation mark, which is the
wrong one for a contraction; and a missing path gets plain prose with no marks
at all. Three arrive pretty-printed with newlines and two-space indentation and
the fourth arrives minified, and that fourth is also the only one carrying
`"warnings": null`, a key its three siblings do not have. Every response says
what it is -- `object` is `"card"`, `"error"` or `"related_card"` -- and the
card's own `scryfall_uri` ends `?utm_source=api`, so an application rendering
"view on Scryfall" passes on their attribution without deciding to.

## Sources

- Documentation: https://scryfall.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve scryfall     # run it
cauldron verify scryfall -v # check every claim
```
