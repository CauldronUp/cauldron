# Kraken

Emulates the Kraken API (v0), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**Every success carries an empty error array, and an
empty array is true.** `{"error": [], "result": {...}}` -- so `if (res.error)
throw` throws on every success and never on a failure, in JavaScript, in Ruby
and in PHP alike. The only reading that works is `res.error.length`, and nothing
in the response says so.

**And a failure answers 200 with no `result` key at all** -- not null, not
empty, absent -- so `res.result.XXBTZUSD` is a TypeError rather than a message
anyone could show. One failure is 200 and another is 404 in the same shape, and
both messages are a machine code and a sentence joined by a colon. **You ask
about `XBTUSD` and are answered about `XXBTZUSD`**, Kraken's own legacy naming,
and no field anywhere holds either name -- it is only the key. The ticker itself
is nine single-letter fields in three shapes: `a` and `b` are three-entry
tuples, six more are two-entry tuples, and `o` is a bare string. Every number is
a string with its trailing zeros kept, `"78311.50000"`, except `t`, which holds
two bare integers.

## Sources

- Documentation: https://docs.kraken.com/api/docs/rest-api/get-ticker-information
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kraken     # run it
cauldron verify kraken -v # check every claim
```
