# snowflake

Emulates the snowflake API (v2), for local development and tests.

**14 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

One endpoint answers with two entirely different bodies, and which one you get depends on how fast the query ran. A quick query gets 200 and the results inline; a slow one gets 202, a statement handle, and no data at all, and the caller is expected to poll. Nothing in the request decides which -- the same query against the same warehouse answers one way when it's warm and the other when it's cold, so the code path that handles 202 is the one nobody exercises locally and everyone meets in production.

Every value in a result row is also a string, including numbers and booleans, with the real types living separately in `resultSetMetaData.rowType` -- and a row is a positional array, not an object, so reordering a `SELECT` silently reorders everything downstream.

## Sources

- Documentation: https://docs.snowflake.com/en/developer-guide/sql-api/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve snowflake     # run it
cauldron verify snowflake -v # check every claim
```
