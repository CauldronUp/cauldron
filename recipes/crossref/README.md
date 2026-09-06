# Crossref

Emulates the Crossref API (v1), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

A date is an array inside an array and it is not
always the same length. `"issued": {"date-parts": [[2013, 7, 31]]}` on one work
and `[[1970, 6]]` on another: the outer array is there because the field can hold
a range, and the inner one stops wherever the metadata stopped. Nothing says
which precision arrived, so `date-parts[0][2]` is a number on one work and
`undefined` on the next, and building a Date from three positional arguments
silently reads the month as the day. **And one work has three publication dates
at two precisions** -- `published` and `published-online` are `[[2013, 7, 31]]`
while `published-print` is `[[2013, 8]]`, so "when was this published" has three
answers, one of them a month, and the print date is after the online one.

`title` is an array: a paper has one title and it arrives in a list of one, as do
`container-title`, `short-container-title` and `ISSN`. `created` carries the same
instant three ways -- `date-parts`, `date-time` and a `timestamp` in epoch
milliseconds -- where `published` beside it carries only the first. The envelope
is `status`, `message-type`, `message-version` and `message`, with a version of
the envelope rather than the API that has read `1.0.0` for years and a `status`
of the string `"ok"` beside an HTTP 200. A single-work lookup carries a
relevance `score` for a query it did not make. And a DOI that does not resolve is
the bare plain text `Resource not found.` with a 404 and no charset -- no
envelope, no status field, from an API whose every success is that four-field
object.

## Sources

- Documentation: https://api.crossref.org/swagger-ui/index.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve crossref     # run it
cauldron verify crossref -v # check every claim
```
