# Carbon Intensity

Emulates the Carbon Intensity API (carbonintensity), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

Intensity's, where **a path that does not exist answers 200 and
the body says 400.** `GET /nosuchpath` returns HTTP 200 carrying `{"error":
{"code": "400 Bad Request", "message": "Please enter a valid path e.g.
/intensity/"}}` -- so `response.ok` is true for a URL the API does not have, and
the code is not a code but a status line as a string, with the number and the
reason phrase in one field. A date it cannot parse, meanwhile, answers a real
400 with the same body.

**And the same fuels are named two ways on two endpoints.** `/generation` calls
them `"gas"`, `"imports"`, `"nuclear"` -- lowercase, one word -- and
`/intensity/factors` calls them `"Gas (Combined Cycle)"`, `"Gas (Open Cycle)"`,
`"Dutch Imports"`, `"French Imports"`, `"Irish Imports"`: capitalised, spaced,
parenthesised, and split finer. Neither vocabulary is a key of the other, so
joining a fuel's share to its emissions factor needs a mapping the API does not
publish. The timestamps have no seconds -- `"2026-08-30T10:00Z"`, which is the
form its own error message asks for -- and every half hour carries a `forecast`,
an `actual` and an `index`, two numbers and a category derived from one of them
by a threshold nothing states.

## Sources

- Documentation: https://carbon-intensity.github.io/api-definitions/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve carbonintensity     # run it
cauldron verify carbonintensity -v # check every claim
```
