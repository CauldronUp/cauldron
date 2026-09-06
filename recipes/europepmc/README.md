# Europe PMC

Emulates the Europe PMC API (v6.9), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

PMC's, where **three failures have three shapes and the HTTP
status is right on exactly one of them.** A search with no criteria answers 200
carrying `{"errCode": 404, "errMsg": ...}`; a bad page size does the same; and a
path that does not exist answers a real 404 with `Content-Length: 0` and nothing
in it. So `response.ok` is true for both failures that explain themselves and
false for the one that does not.

**And asking for a format that does not exist gets you XML.** `?format=xxx`
answers 200 with `application/xml`, telling you -- in a document your JSON
client cannot parse -- that you should have asked for `"xml"` or `"json"` or
`"dc"`, with a space left before the closing tag. **The same document also sends
a real boolean in one object and letters in the next**: the echoed request
carries `"synonym": false` while eleven fields on the record beside it are `"Y"`
and `"N"`, so the API knows how to send a boolean and declines to, eleven times,
next door. `journalIssn` is `"0305-1048; 1362-4962; "` -- two identifiers, one
string, a trailing separator and a trailing space. `pubType` is joined the same
way with commas inside its items. And an array arrives wrapped in an object
under a singular key inside a field named for a list: `"fullTextIdList":
{"fullTextId": ["PMC3531190"]}`.

## Sources

- Documentation: https://europepmc.org/RestfulWebService
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve europepmc     # run it
cauldron verify europepmc -v # check every claim
```
