# OpenAlex

Emulates the OpenAlex API (openalex), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

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
- Machine-readable description: https://developers.openalex.org/openapi.json, last checked 2026-09-05
  `cauldron drift openalex` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openalex     # run it
cauldron verify openalex -v # check every claim
```
