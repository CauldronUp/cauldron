# UniProt

Emulates the UniProt API (uniprot), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**The same acronym is capitalised two ways in two
top-level keys of one object.** `uniProtkbId` has a lowercase `kb` and
`uniProtKBCrossReferences` an uppercase `KB`, four fields apart -- so anything
generating types from the response gets two spellings of one name, and anything
guessing the other finds nothing.

**And the failure echoes your URL back with the scheme downgraded.** A request
to `https://rest.uniprot.org/uniprotkb/NOSUCHACC` answers `{"url":
"http://rest.uniprot.org/uniprotkb/NOSUCHACC", "messages": [...]}` -- the
request was over TLS and the URL in the reply is not, so a client that logs it
or retries it drops to plain HTTP without being told. That failure is a 400
rather than a 404, while a path that does not exist is nginx's own HTML page
naming the version it runs. `proteinExistence` packs a number and its label into
one string, `"1: Evidence at protein level"`. And `entryType` is a sentence with
a parenthetical that names the database twice.

## Sources

- Documentation: https://www.uniprot.org/help/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve uniprot     # run it
cauldron verify uniprot -v # check every claim
```
