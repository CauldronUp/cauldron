# Open Library

Emulates the Open Library API (openlibrary), for local development and tests.

**9 conformance cases, 8 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Library's, where the answer is keyed by the string you sent and a
miss is that key's absence. `GET /api/books?bibkeys=ISBN:0451526538` answers
`{"ISBN:0451526538": {...}}`, so reading the response means rebuilding the query
string -- prefix, colon and all -- to index it; and a lookup that matched nothing
answers `{}` with a 200, not a null under that key and not a 404. So the
difference between found and not found is whether a key you constructed yourself
is there, and `data["ISBN:" + isbn].title` throws a TypeError on the miss without
saying why. **And the same field name carries a reference on one endpoint and a
resolved object on the other**: the canonical document answers
`"authors": [{"key": "/authors/OL18319A"}]` while `/api/books`, for the same
book, answers `[{"url": ..., "name": "Mark Twain"}]`. One is a pointer that costs
another request and nothing says which you are holding.

The rest is the data model of a wiki on a JSON surface. `/isbn/{isbn}.json` answers a
**302 when it finds the book** and a 404 when it does not, so the redirect is the
happy path -- and that 404 is a full HTML page served as `text/html` from a path
ending in `.json`. Timestamps are typed objects,
`{"type": "/type/datetime", "value": "2008-04-01T03:28:50.625462"}`, with
microseconds and no timezone at all. `type` is a reference too. `revision` and
`latest_revision` are two fields carrying one number. Every identifier is an
array, including `"openlibrary": ["OL1017798M"]`, which cannot repeat. And there
are **two schemes in one document**: `url` and `authors[].url` are `http://`
while `subjects[].url` and every cover URL are `https://` -- the record links are
the insecure ones.

## Sources

- Documentation: https://openlibrary.org/developers/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openlibrary     # run it
cauldron verify openlibrary -v # check every claim
```
