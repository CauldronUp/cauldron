# PDBe

Emulates the PDBe API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The response is keyed by what you asked for, and not as
you asked it.** `GET /pdbe/api/pdb/entry/summary/4HHB` answers `{"4hhb": [{...}]}`
-- a map with one key, and that key is your own query lowered. `res["4HHB"]` is
undefined and `res["4hhb"][0]` is the answer, and no field anywhere in the
response holds the identifier, so recovering it means `Object.keys(res)[0]`.

**The same record from RCSB does the opposite.** That Recipe pins `rcsb_id:
"4HHB"` for a request that said `4hhb`; this one pins a key of `"4hhb"` for a
request that said `4HHB`. Two organisations, one deposition, normalisation in
opposite directions -- and they do not agree on the date either: PDBe sends
`"19840307"`, eight digits and no separators, where RCSB sends
`"1984-03-07T00:00:00.000+00:00"`. Three failures have three shapes and one is a
405 for a GET: a missing entry is 404 under `message`, a missing path is 404
under `detail`, and the same path without an identifier is `405 Method Not
Allowed` on an endpoint that takes GET and only GET. `number_of_entities`
carries a key with a slash in it, `"dna/rna"`, beside the two it joins. And the
experimental method is two arrays in two cases, `["x-ray"]` and `["X-ray
diffraction"]`.

## Sources

- Documentation: https://www.ebi.ac.uk/pdbe/api/doc/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pdbe     # run it
cauldron verify pdbe -v # check every claim
```
