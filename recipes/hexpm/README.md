# Hex.pm

Emulates the Hex.pm API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-26.**

## What this Recipe found

The API hands you the dependency line to paste and the
three it hands you disagree. For `phoenix` 1.8.13 the `configs` object is
`mix.exs: {:phoenix, "~> 1.8"}`, `rebar.config: {phoenix, "1.8.13"}` and
`erlang.mk: dep_phoenix = hex 1.8.13` -- the first a range that will pick up
1.8.14 on the next resolve, the other two exact pins. So the same request tells
an Elixir project to float within a minor line and an Erlang project to freeze
on a patch, and the only thing choosing between those policies is which key you
copied. Also: four download counters where `all`, `day` and `week` say what they
measure and `recent` does not, with no `month` beside it to infer from; a
`meta.maintainers` that is always empty while `owners` holds the actual people;
and an advisory carrying three identifiers for one finding.

## Sources

- Documentation: https://github.com/hexpm/specifications/blob/main/apiary.apib
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve hexpm     # run it
cauldron verify hexpm -v # check every claim
```
