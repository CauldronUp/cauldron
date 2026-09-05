# Frankfurter

Emulates the Frankfurter API (v1), for local development and tests.

**9 conformance cases, all of them checked against the live API on 2026-08-28.**

## What this Recipe found

Asking for a Sunday gets you Friday with a 200 and
no word about it. `GET /v1/2026-08-23` answers `{"date": "2026-08-21", ...}`,
because the European Central Bank publishes on working days and this service
falls back to the most recent fixing. The only thing that says the date moved is
the `date` field, which a client that already knows what it asked for has no
reason to read -- so a two-day-old rate arrives labelled as the rate for the day
you wanted, and over Easter that is four days. **And the base is never among the
rates**: with `base=USD` the object holds EUR, GBP and JPY and no USD, so
`rates[base]` is undefined, and `symbols=USD,EUR` answers with EUR alone -- the
currency removed from your own list without a word.

A range changes the shape under the same key names: `date` becomes `start_date`
and `end_date`, and `rates` goes one level deeper, so `rates.EUR` works on one
and is undefined on the other. The range you get is the range that has data --
asking to the 31st of December answers to the 27th of August, silently clamped
-- and the days between are simply absent, so iterating dates and indexing
`rates[date]` finds nothing on two days in seven. A currency that does not exist
and a date before the series began are the same 404 with the same
`{"message": "not found"}`, naming neither. But a **path** that does not exist
answers `{"status": 404, "message": "not found"}` -- the same words with a status
added, so code branching on `body.status` finds a number when it mistyped the
URL and undefined when it mistyped the currency.

## Sources

- Documentation: https://frankfurter.dev/
- Machine-readable description: https://api.frankfurter.dev/v1/openapi.json, last checked 2026-09-05
  `cauldron drift frankfurter` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve frankfurter     # run it
cauldron verify frankfurter -v # check every claim
```
