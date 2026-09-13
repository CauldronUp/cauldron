# gandi

Emulates the Gandi domain listing for local development and tests.

**9 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from the [`go-gandi`](https://github.com/go-gandi/go-gandi) client's published wire types, and struck live against `api.gandi.net` on 2026-09-13 with no header, with a wrong bearer, and on a path that does not exist.

## What this Recipe found

**Three failures, three formats, one hostname.**

```
(no header)        401  {"message":"You must provide an access token or an API Key. See https://api.gandi.net/docs/authentication/"}
Bearer notreal     403  {"object":"HTTPForbidden","cause":"Forbidden","code":403,"message":"Access was denied to this resource."}
/v5/cauldron-nope  404  <html>…<pre>Cannot GET /v5/cauldron-nope</pre>…</html>
```

A one-key JSON object, a four-key JSON object, and an HTML page. The status changes too: a credential that is missing is 401 and one that is wrong is 403, so the two halves of "you may not do this" land on opposite sides of the authentication boundary.

**And the four-key one names a Python class.** `"object": "HTTPForbidden"` is the exception type from the framework that raised it, serialised into the response. `"cause": "Forbidden"` beside it is the HTTP reason phrase, and `"code": 403` is the status a third time — four keys, three of which restate the status line.

**The HTML one is a different stack entirely.** `Cannot GET /v5/…` inside a minimal document is Express's default handler, so one host answers with a Python exception on one path and a Node router's error page on another.

**`autorenew` is a boolean or an object, depending on the endpoint.** In the listing it is `*bool`; on the detail record it is `{href, dates, duration, enabled, org_id}`. One key, two JSON types, and which one arrives is decided by whether you asked for one domain or for all of them.

**So is the nameserver field, and it changes number too.** The listing carries `nameserver`, an object; the detail carries `nameservers`, an array of strings. Singular and plural, object and array, for the same thing.

**Three fields say who owns the domain.** `owner`, `domain_owner` and `orga_owner` sit beside each other, with `sharing_id` making four — and the detail record replaces all of them with a `sharing_space` object.

**The domain name is on the record twice.** `fqdn` and `fqdn_unicode` — the punycode form and the readable one, which for an ASCII domain are the same string.

**The lifecycle is eleven timestamps in a flat bag.** `dates` carries `registry_created_at`, `updated_at`, `authinfo_expires_at`, `created_at`, `deletes_at`, `hold_begins_at`, `hold_ends_at`, `pending_delete_ends_at`, `registry_ends_at`, `renew_begins_at` and `renew_ends_at`, nine of them optional — so which of them is present is the actual state machine.

**And `status` is a list.** A domain is in several states at once, and nothing orders them.

## Sources

- [`go-gandi/go-gandi`](https://github.com/go-gandi/go-gandi) — `domain/types.go`, the wire types for both the listing and the detail record.
- Live: `api.gandi.net`, struck 2026-09-13 with no header, a wrong bearer, and an unrouted path.

## Modelling limits

- **One route.** The domain listing. The detail record, contacts, nameservers, LiveDNS, transfers, renewals, the billing surface and the organisation API are the rest — and the detail record is where `autorenew` becomes an object, which is recorded here rather than served.
- **The HTML 404 is served as its sentence.** Live, `Cannot GET /v5/…` arrives wrapped in a minimal HTML document from Express. This Recipe serves the sentence with the same `text/html` type and does not reproduce the surrounding markup.
- **No `spec:`.** `api.gandi.net/v5/domain/openapi.json` sits behind the credential check, and no machine-readable description is served at any address this Recipe could reach anonymously — so there is nothing for `cauldron drift` to record.
- **The success fixture is client-derived.** Listing domains needs a real key; the records here are `ListResponse`'s own fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Gandi is reached through `gandi.cli` on PyPI, `go-gandi`, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
