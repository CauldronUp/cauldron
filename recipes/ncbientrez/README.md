# NCBI Entrez

Emulates the NCBI Entrez API (entrez-eutils), for local development and tests.

**17 conformance cases, 16 checked against the live API on 2026-08-31.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Entrez's, where **a `Content-Type` of
`application/json` is simply false.** `efetch` sends that header and answers
plain text: ask for an abstract with `retmode=json` and you get a word-wrapped
citation; ask without a `rettype` and you get the literal bytes `42673563\n`,
which parses only because a bare number is valid JSON. `esearch` and `esummary`
honour `retmode=json` properly and `efetch` never does, so the header tells the
truth on two endpoints of three.

**No two failures agree.** `esearch` on an invalid database answers 200 with
`ERROR` in capitals -- the only uppercase key in the API -- sitting exactly
where `idlist` sits on success. `esummary` on the same mistake renames its
envelope, `result` becoming `esummaryresult`, and turns it from an object into
a bare array of one string: the failure is not a field of the response, it is a
different response. `efetch`'s failures are `text/plain` and never URL-decoded,
`Database%3A+notarealdb+-+is+not+supported` with literal plus signs. And the
quietest is 200 with three bytes -- `1. ` -- a citation number, a full stop and
a space, with nothing anywhere saying the record does not exist.

## Sources

- Documentation: https://www.ncbi.nlm.nih.gov/books/NBK25501/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ncbientrez     # run it
cauldron verify ncbientrez -v # check every claim
```
