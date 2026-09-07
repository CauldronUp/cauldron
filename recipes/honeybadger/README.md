# honeybadger

Emulates the Honeybadger reporting API for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-07.**

Written against Honeybadger's reference at `docs.honeybadger.io` and struck live against `app.honeybadger.io` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**`errors` is a string.**

```
GET /v2/projects                (no credential)
403 {"errors":"Access denied"}
```

Plural, and holding two words.

Every other provider in this collection that names a field `errors` puts an array there — SendGrid, Terraform Cloud, beehiiv, Middesk, Gong, Confluent — because the point of the plural is that one request can fail several ways at once. Here it is a bare string, so `errors[0]` is `"A"` and `errors.length` is `12`, and **both of those are what a client written against the plural will compute**. Neither throws. Neither is right.

**A missing credential is 403, and so is a wrong one.** Struck live both ways, byte-identical. 403 means the caller was identified and is not permitted; sending nothing at all is neither of those. So the status that would tell a client to re-authenticate never arrives, and the two failures needing different fixes are one response.

**An unknown path is HTML.** 404 with a full application page, so the JSON above is reachable only on paths that exist — and a client cannot tell a typo from a permissions problem by parsing, only by reading.

**A project counts its faults twice.** `fault_count` is lifetime and `unresolved_fault_count` is now: 412 against 7 on the same record. A dashboard reading the first shows four hundred problems where there are seven, and the field names are one word apart.

**Environments are bare strings.** `["production", "staging"]` — no identifiers, so an environment cannot be addressed except by whatever name a client happens to be holding.

**Reporting and reading take different credentials in different headers.** Sending an error uses a per-project API key in `X-API-Key`; reading data uses a personal auth token as the Basic username with an empty password. Two credentials, two shapes, one product, and this Recipe serves the reading half.

## Detection

`@honeybadger-io/js` names **both** hosts in its published archive — `api.honeybadger.io` seven times for reporting and `app.honeybadger.io` six times for reading, which is the surface this Recipe serves. It is mapped, along with the PHP and Go notifiers.

## Modelling limits

- **One route.** Projects. Faults, notices, deploys, comments, teams and the check-in surface each want their own evidence.
- **The success path is documentation-derived.** Reading anything needs an account; the refusals come from the wire.
- **No `spec:`.** Honeybadger publishes a rendered documentation site.
