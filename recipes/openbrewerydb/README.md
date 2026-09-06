# Open Brewery DB

Emulates the Open Brewery DB API (openbrewerydb), for local development and tests.

**15 conformance cases, 10 checked against the live API on 2026-08-29.** The unchecked one is the paging case: it sends the two parameter names this Recipe declares, which is a claim read from the provider's own description rather than struck against it.

## What this Recipe found

Brewery DB's, where `brewery_type` says `closed`. The field naming
what kind of brewery this is also answers whether it is one: its fourteen values
include `micro`, `brewpub` and `regional`, which are kinds, alongside `closed`,
which is a status, `planning`, which is a stage, and `location`, which is
neither. 642 are typed `closed` and 639 `planning` -- 1,281 of 11,848 -- so
filtering `by_type=micro` silently drops every closed micro, and showing all
types puts businesses that do not exist beside businesses that do.

**And two pairs of fields are the same field twice.** `state` and
`state_province` carry the same string on every record, and so do `street` and
`address_1` -- not similar, identical, `"Oklahoma"` and `"Oklahoma"`. Two are
current and two are kept from an older shape, and nothing says which pair is
which. Asking for one brewery too many is a redirect: `?per_page=201` answers
302 to the homepage rather than 400 or a clamp, so a client that follows
redirects gets the landing page and a 200. A brewery that does not exist answers
404 as an HTML page. The documented `autocomplete` path is a 301 to `search`.
`by_state` has 214 keys, most of which are not states -- ACT is Australian,
Argyll Scottish, Auckland a New Zealand region, Aveiro a Portuguese district,
all sorted in among Alabama and Arizona. And `phone` is two formats in one
column: `"+49 9261 628000"` and `"4058160490"`.

## Sources

- Documentation: https://www.openbrewerydb.org/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openbrewerydb     # run it
cauldron verify openbrewerydb -v # check every claim
```
