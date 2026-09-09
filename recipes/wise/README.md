# Wise

Emulates the Wise API (v1), for local development and tests.

**20 conformance cases, 7 checked against the live API on 2026-09-09.**

## What this Recipe found

**The two keys arrive in either order, and which one you get is a coin flip.**

Six identical requests to `https://api.wise.com/2026Q3/profiles` with no credential, seconds apart, on 2026-09-09:

```
{"error":"missing_token","error_description":"Missing token"}
{"error_description":"Missing token","error":"missing_token"}
{"error":"missing_token","error_description":"Missing token"}
{"error_description":"Missing token","error":"missing_token"}
{"error_description":"Missing token","error":"missing_token"}
{"error":"missing_token","error_description":"Missing token"}
```

Three of each. Same URL, same headers. Four more calls against the older `/v1/` path split the same way, so this is not a property of one route or one release — it is a serialiser walking an unordered map.

The values never vary and the shape never varies, so nothing that *parses* JSON is affected. What is affected is everything that treats the body as bytes:

- a recorded HTTP fixture matched on an exact body — VCR, nock, WireMock in strict mode — matches or misses depending on which order happened to be recorded
- a snapshot test of the failure flakes about half the time
- any signature, checksum or ETag over the body is unstable

It is the fourth key-order instability in this collection, after Zuora, Browserbase and Postman, and the first proven by **repeating one request** rather than by comparing two different ones.

**It is recorded rather than served.** An emulator that returned a random key order would be faithful and would make every test written against it flake, which is the one thing a fake must not do. This Recipe sends one order and writes the other down.

**Two credential failures, two codes, both true.**

```
(no header)              401 {"error":"missing_token","error_description":"Missing token"}
Authorization: Bearer …  401 {"error":"invalid_token","error_description":"Invalid token"}
```

Both accurate about the request that earned them, which is rarer in this collection than it should be — seven providers here answer "invalid key" to a request carrying none.

When this Recipe was first written the format had one credential branch, and the errors table said in as many words that the invalid-token wording "is recorded in the header rather than modelled, since there is no second branch here to carry it". The format grew `rejected_error` since. It is served now, and the stale note is gone.

**A header with no scheme is reported as no header at all.** `Authorization: notreal` — present, non-empty, missing only the word `Bearer` — answers `missing_token`, word for word, on both paths. So the one mistake where the caller demonstrably sent something is the one this API describes as sending nothing.

**One endpoint answers without a credential and the rest do not.** `/v1/quotes` computes a real quote for anybody. It also changes the order it checks routing and authentication depending on *which* credential problem you have.

**The published document does not declare the failure everybody meets.** `GET /profiles` in Wise's own OpenAPI declares exactly two responses: `200` and `429`. There is no 401 in it, on an API whose 401 is reachable four different ways. It does, unlike [close](../close) and [middesk](../middesk), apply its security schemes properly — `security` is set per operation rather than defined and forgotten.

**The version is a calendar quarter.** The document's server is `https://api.wise.com/2026Q3`, with `preview`, `latest` and `legacy` in the documentation's version picker. A URL a client hardcodes names a quarter, and quarters end.

## Modelling limits

- **The random key order is not reproduced**, for the reason above.
- **No `spec:` is recorded.** Wise publishes a real OpenAPI document at `docs.wise.com/_spec/api-reference/@latest/index.yaml`, and it describes the `2026Q3` surface while this Recipe models `v1` paths. Pinning it would fingerprint a document that declares none of these routes — the position [basiq](../basiq) and [customerio](../customerio) are in, and not one worth entering on purpose.

## Sources

- Documentation: https://docs.wise.com/api-reference
- Live: `api.wise.com`, struck 2026-08-31 and again 2026-09-09.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.
