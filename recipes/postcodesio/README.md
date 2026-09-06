# Postcodes.io

Emulates the Postcodes.io API (postcodesio), for local development and tests.

**16 conformance cases, 15 checked against the live API on 2026-08-29.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Four of the fields are named after
organisations that no longer exist. `nhs_ha` is a Health Authority, abolished in
2002. `primary_care_trust` is a PCT, abolished in 2013. `nuts` is an EU region
the UK left in 2020. And `ccg` is a Clinical Commissioning Group, abolished in
2022 -- sitting next to `icb`, the body that replaced it, holding almost the
same string, with nothing marking either as historical.

**And in Wales the two abolished English fields hold the name of a Welsh body.**
CF10 3NP answers "Cardiff and Vale University Health Board" for both `nhs_ha`
and `primary_care_trust` -- two fields named after two different defunct English
structures, filled with a body that is neither -- while `icb`, which is an
English structure, says "Wales". The unversioned `lsoa` is an alias, and only a
postcode whose boundaries moved can say which census it follows: in Westminster
`lsoa`, `lsoa11` and `lsoa21` all agree, and in Cardiff they do not, which is
what shows `lsoa` tracking 2021. A null field's code is a row of nines --
`admin_county` is null and `codes.admin_county` is `"E99999999"` -- and Cardiff
carries the same sentinel under `codes.icb` while `icb` itself says "Wales".
`ccg_id` is a different shape on every record: `"W2U3Z"` on the English one,
`"7A4"` on the Welsh, against a letter and eight digits everywhere else. And
`date_of_introduction` is `"198001"`, six characters with no separator.

## Sources

- Documentation: https://postcodes.io/docs
- Machine-readable description: https://api.postcodes.io/openapi.json, last checked 2026-08-31
  `cauldron drift postcodesio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve postcodesio     # run it
cauldron verify postcodesio -v # check every claim
```
