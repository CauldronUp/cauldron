# Art Institute of Chicago

Emulates the Art Institute of Chicago API (artinstituteofchicago), for local development and tests.

**14 conformance cases, 11 checked against the live API on 2026-08-29.** The two unchecked ones are the paging pair: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Institute of Chicago's, where one field of ninety-eight is
licensed differently from the rest and the response says so in prose. `info.
license_text` reads "The `description` field in this response is licensed under
a Creative Commons Attribution 4.0 Generic License (CC-By) ... All other data in
this response is licensed under a Creative Commons Zero (CC0) 1.0 designation",
and that English sentence in a sibling key is the only statement of which field
is which. Nothing structural marks it: `description` sits in `data` beside the
other ninety-seven and looks the same, so a client that stores the response and
republishes it under one licence is wrong about exactly one column.

**And an image URL is three pieces from two levels of the response.** The
artwork carries `image_id` and no URL at all; the base lives in a top-level
`config` key beside `data` rather than in it; and the rest of the path --
`/full/843,/0/default.jpg` -- is in neither. `date_display` and `date_end`
disagree, so Seurat's Grande Jatte reads "1884-86, border added 1888-89" while
the number a query sorts on is 1886. The same coordinate is in the response
three times and the third is rounded: `latitude` and `longitude` carry thirteen
decimal places and `latlon` joins them as a string with twelve.
`colorfulness` is 0 on a pointillist painting made of coloured dots and 43.0038
on a Franz Marc, so nothing distinguishes "not computed" from "computed, and
grey". And `has_not_been_viewed_much` is a negated boolean, which makes the
ordinary question a double negative.

## Sources

- Documentation: https://api.artic.edu/docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve artinstituteofchicago     # run it
cauldron verify artinstituteofchicago -v # check every claim
```
