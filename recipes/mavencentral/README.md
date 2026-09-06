# Maven Central

Emulates the Maven Central API (solrsearch), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-08-27.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Central's, where the search API is Apache Solr's response handed
over unedited. A successful search reports `"status": 0` -- Solr's convention,
not HTTP's -- inside a body that arrived with a 200, beside `QTime` and a
`params` object echoing the caller's query along with the internal field list
the service asked Solr for and the sort expression it used. The results
themselves are single letters: `g`, `a`, `v`, `p` and `ec` for groupId,
artifactId, version, packaging and the file extensions that exist, and nothing
in the response says so. `start` appears twice in one document in two types --
`0` where the service computed it, `""` where the caller's value is read back as
text. A search matching nothing is a 200 that says it succeeded, so a mistyped
query and a missing artifact are the same response. And a query Solr will not
parse is plain text that stops mid-sentence: `Solr returned 400, msg:` with
nothing after the colon, an error relayed by a service that did not read it.

## Sources

- Documentation: https://central.sonatype.org/search/rest-api-guide/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mavencentral     # run it
cauldron verify mavencentral -v # check every claim
```
