# openFDA

Emulates the openFDA API (openfda), for local development and tests.

**9 conformance cases, 7 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

**Every value begins with its own field name.**
`"purpose"` is `["Purpose Pain reliever"]`, `"do_not_use"` is `["Do not use if
you are allergic to aspirin ..."]`, `"questions"` is `["Questions or comments?
Call 1-877-753-3935 ..."]`. These are scanned drug labels and the section
heading came along with the section, so the field name is both the key and the
first words of the value, and rendering "Purpose: Purpose Pain reliever" is the
natural thing to build.

**And two of them begin with a different field's name.**
`indications_and_usage` starts "Uses for the temporary relief" and
`storage_and_handling` starts "Other information store between", so stripping
the prefix is not a fix: the prefix is the heading printed on the box and the
key is what a database called it. Every field is an array of one, identifiers
included. Finding nothing is a 404 -- `{"code": "NOT_FOUND", "message": "No
matches found!"}`, with an exclamation mark -- and a search naming a field that
does not exist answers exactly that, byte for byte. A path that does not exist
is Express's own page in `text/html`. And the temperature is written with
U+00BA MASCULINE ORDINAL INDICATOR rather than U+00B0 DEGREE SIGN: a character
that looks right at a glance and matches no search for degrees.

## Sources

- Documentation: https://open.fda.gov/apis/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve openfda     # run it
cauldron verify openfda -v # check every claim
```
