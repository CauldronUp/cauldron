# Taddy

Emulates the Taddy API (graphql), for local development and tests.

**13 conformance cases, 5 checked against the live API on 2026-09-02.**

## What this Recipe found

**Sending no key is a 200 and sending a bad one is a
500**, both carrying the identical sentence.

## Sources

- Documentation: https://taddy.org/developers/intro-to-taddy-graphql-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve taddy     # run it
cauldron verify taddy -v # check every claim
```
