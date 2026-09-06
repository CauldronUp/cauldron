# Lago

Emulates the Lago API (v1), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.getlago.com, no account and no key. The documented auth failure matched exactly: a missing Authorization header and a well-formed wrong one both answer `{"status":401,"error":"Unauthorized"}`, with no `code` field, byte for byte.

## What this Recipe found

An event with an unrecognized billing-metric code is silently ignored, not refused. Lago's own docs say so directly: if the `code` doesn't match any active billable metric, "it will be ignored during the process" -- the request is still accepted, the event is still stored, the API still echoes the event back with its bad code intact, and no fee is ever produced. A typo in a metric code, or a metric renamed on Lago's side, is a month of usage received and never billed, with every signal a client gets saying it worked.

The identifier's format also depends on which datastore your Lago instance runs on -- `lago_id` is documented as "not guaranteed to be a UUID," a composite string on installs using the ClickHouse events store -- so validating it as a UUID or storing it in a `uuid` column works against one deployment and fails against another, both genuinely Lago. And a field named `precise_total_amount_cents` holds a value with a decimal point in it (`'1234.56'`), so the `_cents` suffix that usually signals "no fractional part" means the opposite here.

URLs are also addressed by the caller's own identifiers, not Lago's: customers live at `/customers/{external_customer_id}`, so the primary key in the path is the one your system chose, while `lago_id` shows up in every payload but isn't what any route actually accepts.

## Sources

- Documentation: https://docs.getlago.com/api-reference
- Machine-readable description: https://getlago.com/openapi.json, last checked 2026-09-05
  `cauldron drift lago` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lago     # run it
cauldron verify lago -v # check every claim
```
