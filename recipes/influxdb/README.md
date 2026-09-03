# InfluxDB

Emulates the InfluxDB API (v2), for local development and tests.

**6 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

Its Cloud host **does not serve InfluxDB**.
`/api/v2/buckets` ignores the Authorization header and redirects to a login page,
and InfluxData's own documented write example 404s verbatim.
Five are Sendcloud's, where **buying a label is a boolean rather than a call**:
one POST creates a parcel, and a field on that same request decides whether it
also spends money -- a field that never appears on the record it produced.

## Sources

- Documentation: https://docs.influxdata.com/influxdb3/cloud-serverless/api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve influxdb     # run it
cauldron verify influxdb -v # check every claim
```
