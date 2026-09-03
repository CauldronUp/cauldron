# RCSB

Emulates the RCSB API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**Two of its fields are the letters
Y and N, and "N" is true.** `has_released_experimental_data` is `"N"` and
`pdb_format_compatible` is `"Y"` -- single-character strings where a boolean
belongs -- and in JavaScript the string `"N"` is truthy, so `if
(info.has_released_experimental_data)` runs the yes branch on an entry that said
no. The same code is correct for `"Y"`, which is why it survives review. There
is no value either field can hold that is falsy.

**And the identifier you get back is not the one you sent.** Ask for `4hhb` and
the record answers `"rcsb_id": "4HHB"`; ask for a nonsense one and the failure
says `No data found for entryId: NOTANID`. A client keying a cache on what it
requested misses every hit, and one matching the failure's text against its own
input never matches either. Two 404s share the host and share one key: the
application's is `{status, message, link}` and a path the application never sees
gets Spring Boot's `{timestamp, status, error, path}`, which has no `message` at
all. The experimental method is sent twice in two spellings, `"X-RAY
DIFFRACTION"` and `"X-ray"`. The crystallographic keys keep their CIF capitals
among snake_case siblings -- `cell.Z_PDB`, `symmetry.Int_Tables_number`,
`space_group_name_H_M` -- so any client that lower-cases its keys destroys
three. And a 1984 deposition brings its punched-card formatting with it, a block
of capitals wrapped at sixty columns inside a JSON string.

## Sources

- Documentation: https://data.rcsb.org/redoc/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rcsb     # run it
cauldron verify rcsb -v # check every claim
```
