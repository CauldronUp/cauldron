# SuperTokens

Emulates the SuperTokens API (5.4), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-09-01.**

## What this Recipe found

It **answers 200 to a session that does not exist**:
`{"status": "UNAUTHORISED"}` under a status line that says the request
succeeded, byte-identical with a garbage handle, a wrong `rid`, or no `rid`.

## Sources

- Documentation: https://github.com/supertokens/core-driver-interface
- Machine-readable description: https://raw.githubusercontent.com/supertokens/core-driver-interface/master/api_spec.yaml, last checked 2026-09-01
  `cauldron drift supertokens` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve supertokens     # run it
cauldron verify supertokens -v # check every claim
```
