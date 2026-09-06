# Semantic Scholar

Emulates the Semantic Scholar API (v1), for local development and tests.

**13 conformance cases, 11 checked against the live API on 2026-09-01.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Scholar's, where **a paper's own identifier cannot be routed
to.** A DOI is `10.<registrant>/<suffix>`, so `DOI:10.1145/3197026` carries a
literal slash, and this router matches on an exact count of slash-delimited
segments -- `/graph/v1/paper/{id}` compiles to four while that request splits
into six. The live API answers it perfectly; Cauldron falls through to its own
404. Its fields are also opt-in to a degree that makes the default nearly
empty: a paper fetched without a `fields` parameter carries `paperId` and
`title` and nothing else, so `.year` is `undefined` on every paper and no error
says why. Absence is then spelled two ways in one record -- a requested field
with no value is `"abstract": null`, an unrequested one is missing entirely.

## Sources

- Documentation: https://api.semanticscholar.org/api-docs/graph
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve semanticscholar     # run it
cauldron verify semanticscholar -v # check every claim
```
