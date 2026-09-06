# BoC Valet

Emulates the BoC Valet API (bocvalet), for local development and tests.

**7 conformance cases, 6 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Valet's, where **the response documents its own
key names inside itself.** `seriesDetail.FXUSDCAD.dimension.key` is `"d"` and
`.name` is `"Date"`, and the observations carry the date under `"d"` -- so a
client that wants it has to read the dimension first and subscript by whatever
it finds, or hardcode one letter and hope.

**And the number is two levels down, under a key the caller supplied.** A rate
is `obs["FXUSDCAD"].v`: the series name is a JSON key rather than a value, and
what sits under it is an object with one field, holding a string. The same key
differs by one letter between endpoints -- `/observations` returns
`seriesDetail`, a map keyed by series name, and `/series/{name}` returns
`seriesDetails`, flat, with the name as a field. The terms of use are the first
key of every successful body. And a bad parameter answers with a link to a
generated OpenAPI operation id, underscores and all, in a URL fragment.

## Sources

- Documentation: https://www.bankofcanada.ca/valet/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve bocvalet     # run it
cauldron verify bocvalet -v # check every claim
```
