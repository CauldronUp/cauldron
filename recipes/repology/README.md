# Repology

Emulates the Repology API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-25.**

## What this Recipe found

It is the one registry here that does not answer a
missing thing with a 404: asking for a project nobody packages is a **200 with
an empty array**, so `if (!response.ok) throw` never fires and a typo is
indistinguishable from a project packaged nowhere. One project's response is a
bare array with an entry per packaging -- 806 of them for `curl` -- and only
five of the fourteen fields are on every entry. `srcname` is missing from
twelve, `vulnerable` from a hundred and seventy-four, and `binname` and
`binnames` are the same idea singular and plural with no entry carrying both.
Six statuses appear at once: 377 outdated, 255 legacy, 163 newest, 7 rolling, 3
devel and 1 `incorrect`, which is the index's own judgement that a version
string could not be read as a version, sitting in an ordinary entry beside the
rest.

## Sources

- Documentation: https://repology.org/api/v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve repology     # run it
cauldron verify repology -v # check every claim
```
