# typesense

Emulates the typesense API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The relevance score is a number no browser can read correctly. `text_match` is an int64 with real values up around 578730123365711993 -- sixty-odd bits -- and JavaScript's float64 numbers round that to 578730123365712000, so two hits with genuinely different scores can come out equal after `JSON.parse`. A sort by `text_match` is stable only by accident, and nothing errors: the JSON is valid, the field is present, and the value is simply wrong before any of the caller's code runs. `num_documents` has the identical problem on a large collection's count.

`search_cutoff` is also a boolean, at 200, meaning the search stopped early because it was taking too long -- so `found` can be smaller than the true match count with nothing else in the response saying so.

## Sources

- Documentation: https://typesense.org/docs/latest/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve typesense     # run it
cauldron verify typesense -v # check every claim
```
