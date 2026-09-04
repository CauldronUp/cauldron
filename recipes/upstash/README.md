# Upstash

Emulates the Upstash API (v2), for local development and tests.

**6 conformance cases, 1 checked against the live API on 2026-08-31.**

## What this Recipe found

There are three different 401 bodies on this API, and which one a caller gets depends on how the credential was wrong, not that it was wrong. No `Authorization` header at all gets a 401 with an empty body and no `Content-Type`. A garbage bearer token gets a Go JWT library's raw parse error leaking onto the wire as plain text. A garbage Basic credential gets `{"error":"Unauthorized"}` as JSON. Three requests to the same endpoint, three different content types, and the difference is entirely which flavor of wrong credential was sent.

A field can also be typed like a boolean and answer in quotes: `db_acl_enabled` is a string enum of the literal words `"true"` and `"false"`, sitting beside several ordinary JSON booleans on the same object, so `if (db.db_acl_enabled)` is true for a database with ACLs off. And `region` is nearly useless -- its enum has exactly one value, `"global"` -- while the field that actually names where a database runs is `primary_region`.

## Sources

- Documentation: https://upstash.com/docs/devops/developer-api/introduction
- Machine-readable description: https://upstash.com/openapi.json, last checked 2026-08-31
  `cauldron drift upstash` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve upstash     # run it
cauldron verify upstash -v # check every claim
```
