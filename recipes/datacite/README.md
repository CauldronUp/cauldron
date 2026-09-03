# DataCite

Emulates the DataCite API (rest), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The same key holds a list in one failure and an
object in another.** A missing DOI answers `{"errors": [{"status": "404",
"title": "..."}]}` and an unparseable query answers `{"errors": {"title":
"parse_exception ..."}}`. A client reading `errors[0].title` gets the sentence
on the first and undefined on the second, and one reading `errors.title` gets it
the other way round -- and the key is spelled the same, plural, in both. The
object one is 318 characters of the query parser's own Lucene grammar, newlines
and token names included, forwarded to whoever typed the URL.

**And the listing says there are 133,915,241 records and 10,000 pages of one.**
`totalPages` is not the total divided by the page size; it is ten thousand
records' worth, whatever the page size. A client paging until `page >
totalPages` reads 10,000 of 133 million and stops, told both numbers by the same
object. Asking for page 10,001 answers 200 with `"page": 10000`, a different
number from the one asked for, in the field named after it. A `page[size]` that
is not a number answers 200 with `"data": []` and `"totalPages": 0` beside that
unchanged total. Absence is spelled three ways in one record -- `container: {}`,
`formats: []`, `reason: null` -- the same year is sent twice in two types, and
the record's type is sent six times in six spellings, five saying dataset and
the sixth saying `misc`.

## Sources

- Documentation: https://support.datacite.org/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve datacite     # run it
cauldron verify datacite -v # check every claim
```
