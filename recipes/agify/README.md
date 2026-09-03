# Agify

Emulates the Agify API (agify), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**Three APIs share one budget**. api.agify.io,
api.genderize.io and api.nationalize.io are presented as separate services with
separate documentation, and interleaved calls decrement a single counter: 15,
14, 13, 12, 11 across the three hosts, one allowance of twenty-five a day. So
asking all three about a name -- the obvious thing to build -- costs three of
them, and a page that does it for eight visitors is finished until midnight.

**And `X-Rate-Limit-Reset` is seconds until midnight UTC, not a timestamp.** It
read `60485` while there were `60501` seconds left in the UTC day, so a client
passing it to `new Date(reset * 1000)` lands in January 1971. A name nobody has
heard of is a 200 with `"age": null`, not a 404. Present-but-empty and absent
are different: `?name=` answers 200, and omitting the parameter answers 422.
And the 422 is the only response carrying no rate-limit headers at all, on a
content type without the `charset=utf-8` every success declares -- the response
that says a client is doing something wrong is the one that will not say how
much budget is left.

## Sources

- Documentation: https://agify.io/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve agify     # run it
cauldron verify agify -v # check every claim
```
