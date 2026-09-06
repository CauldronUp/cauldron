# GrowthBook

Emulates the GrowthBook API (v1), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The field everybody actually branches on has no declared values, while the field beside it does. An Experiment's `type` is a proper three-value enum, but `status` -- the property a dashboard filters on, an alert watches, and a deploy gate reads -- is just `{"type": "string"}` with no enum, no example, no description, anywhere in a 2.7-megabyte spec. `resultSummary` even carries its own separate `status` field, also an unconstrained string, meaning something else entirely.

The listing envelope has six fields answering one question -- whether to fetch another page -- and they don't reduce to each other: `count`, `total`, `hasMore`, and `nextOffset` (present and explicitly null on the last page, rather than simply absent) are four independent ways to decide, and code written against any single one of them ignores the rest. And the identifier worth joining on isn't `id` -- it's `trackingKey`, the value that actually shows up in the analytics warehouse, which makes GrowthBook the third provider in this collection, after Grafana's `uid` and ConfigCat's `key`, where the field named `id` is the wrong one to keep.

## Sources

- Documentation: https://docs.growthbook.io/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve growthbook     # run it
cauldron verify growthbook -v # check every claim
```
