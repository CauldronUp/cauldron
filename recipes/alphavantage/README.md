# Alpha Vantage

Emulates the Alpha Vantage API (alphavantage), for local development and tests.

**12 conformance cases, 9 checked against the live API on 2026-08-31.** The 3 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Vantage's, where **every failure is an HTTP 200.**
Wrong symbol, missing function, bogus function, missing parameter, throttled --
every response observed was 200, with no `Retry-After` and no `X-RateLimit`
header anywhere, so a client branching on `response.ok` takes the success path
for all of them. The two failure shapes then disagree about which key to read:
leaving off `apikey` answers `{"Error Message": ...}` and everything else wrong
answers `{"Information": ...}`, so code checking one finds `undefined` on the
other. And `Information` means two unrelated things -- a request to fix and a
request to retry -- in the same key, with the same status, and nothing to
separate them.

**The keys are unreachable with a dot.** `"Meta Data"` has a space,
`"Time Series (Daily)"` has parentheses, and a level down they are `"1. open"`
and `"5. volume"` -- beginning with a digit, containing a full stop. The same
figure is even prefixed to different widths by different functions:
`TIME_SERIES_DAILY` calls it `"1. open"` and `GLOBAL_QUOTE` calls it
`"02. open"`.

## Sources

- Documentation: https://www.alphavantage.co/documentation/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve alphavantage     # run it
cauldron verify alphavantage -v # check every claim
```
