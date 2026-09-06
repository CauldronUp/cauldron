# PDFMonkey

Emulates the PDFMonkey API (v1), for local development and tests.

**12 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

It **states its download URL's lifetime in writing** --
"a temporary signed link valid for 1 hour" -- so an integration cannot store the
link and nothing in the document record says so.

## Sources

- Documentation: https://www.pdfmonkey.io/docs/api/documents/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pdfmonkey     # run it
cauldron verify pdfmonkey -v # check every claim
```
