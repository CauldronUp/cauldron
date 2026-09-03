# OpenAlex

Emulates the OpenAlex API (openalex), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The error message is the schema.** A filter that
does not exist answers 400 with `"nosuchfilter is not a valid field. Valid
fields are underscore or hyphenated versions of: abstract.search,
abstract.search.exact, ... updated_date, version"` -- 4,779 characters and 207
field names, comma-separated, inside one sentence in one string. The failure is
larger than most successful responses, the entire queryable surface of the API
arrives as prose, and finding out which field was wrong means splitting an
English sentence on a colon and then on commas.

**And a free API quotes a price:** every response carries `"cost_usd": 0.0001`
in its meta, a tenth of a hundredth of a cent for a request nobody is billed
for, beside an `x_query` that leaks the internal query language. Three failure
shapes share one host: Flask's default HTML page for a missing work, a lone
`message` for a missing entity, and `{error, message}` for a bad filter.
`title` and `display_name` are byte-identical. And the identifiers are URLs, so
getting the bare W-number or the bare DOI means string-splitting an address.

## Sources

- Documentation: https://docs.openalex.org/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openalex     # run it
cauldron verify openalex -v # check every claim
```
