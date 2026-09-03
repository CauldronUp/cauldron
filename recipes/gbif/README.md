# GBIF

Emulates the GBIF API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

A name it could not find comes back at a hundred per cent
confidence. `?name=Zzzznotaspeciesxyz` answers
`{"confidence": 100, "matchType": "NONE", "synonym": false}` with a 200 --
`confidence` is the highest the scale goes, on the one answer that found nothing,
so code ranking results by confidence puts the failures first. Every real match
scores lower: an exact one is 99, a fuzzy one 95, a genus standing in for a
species 94. **And `synonym` is present only when nothing was found**: the
no-match carries `synonym: false`, and the response whose `status` actually is
`SYNONYM` does not carry the field at all.

A deprecated name matches **EXACTly**. `Felis concolor` is what cougars were
called until 1993, and asking for it answers `matchType: "EXACT"` at confidence
98 with `status: "SYNONYM"` -- exact because the string exactly matched a name
nobody should use. That response carries two keys for two different taxa:
`usageKey` is the synonym and `speciesKey` and `acceptedUsageKey` are the
accepted species, so the field with the plainest name is the one not to store. A
misspelled species silently becomes a genus: `Puma notaspecies` answers
`matchType: "HIGHERRANK"` at 94 with `scientificName: "Puma Jardine, 1834"`,
which a client will print as though it were the species asked for. A typo is
corrected without saying what changed. And `/v1/species/notarealendpoint` answers
400 with the bare text `For input string: "notarealendpoint"` -- Java's
`NumberFormatException`, verbatim, because `/v1/species/{key}` takes an integer
and the segment was handed to a number parser.

## Sources

- Documentation: https://techdocs.gbif.org/en/openapi/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gbif     # run it
cauldron verify gbif -v # check every claim
```
