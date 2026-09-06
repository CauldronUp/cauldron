# unit

Emulates the Unit banking-as-a-service API, for local development and tests.

**13 conformance cases, none checked against a live API.**

Written against Unit's own OpenAPI document, published in its own GitHub organisation — `unit-finance/openapi-unit-sdk`, 117 paths — read on 2026-09-06. The server named in that document is Unit's sandbox, `api.s.unit.sh`.

## What this Recipe found

**The balance is not the money you can spend.** An account carries three integers — `balance`, `available` and `hold` — and they are *cents*. `available` is what a payment can draw on; `balance` includes money already held and on its way out. Code that shows `balance` to a customer overstates what they have by exactly the amount about to leave, and code that authorises against `balance` authorises money that is already spoken for.

**Nothing is where a client first looks for it.** The envelope is JSON:API: a record is `{"type": "depositAccount", "id": …, "attributes": {…}, "relationships": {…}}`, so the balance is at `data.attributes.balance` and `data.balance` is `undefined`. A listing is `{"data": […], "included": […], "meta": {"pagination": …}}`, where `included` holds side-loaded *customers* — so `body.data.concat(body.included)` produces a list of two different kinds of thing.

**The content type is not `application/json`.** It is `application/vnd.api+json; charset=utf-8`. A client that checks the header before parsing will not parse it, and a framework configured to decode only `application/json` hands the handler an empty body with no error anywhere.

**The paging parameters are bracketed.** `page[limit]` and `page[offset]` — a `deepObject` in the spec's own words, defaulting to 100 and capped at 1000. Not `limit`, not `per_page`. A query builder that does not escape brackets, or a proxy that normalises them, silently sends nothing and receives the first hundred every time.

**A failure carries `detail` *and* `details`, one letter apart**, on the same object, inside an array. Reading `error.details` when the sentence is in `error.detail` produces `undefined` and no complaint from anything. And `status` there is a **string** — `"401"`, not `401` — so comparing it to the HTTP status with `===` is false.

**Frozen is neither open nor closed.** A frozen account lists, holds money, and moves none of it.

## Modelling limits

- **Deposit accounts only.** `type` discriminates `depositAccount`, `creditAccount` and `walletAccount`, each with different attributes. Serving one and claiming the others would describe a shape nobody sent.
- **`included` is not served.** Side-loading is driven by an `include` parameter this Recipe does not model, and an empty array would claim that nothing *was* included rather than that nothing was asked for.
- **`relationships` is served; its targets are not.** An account points at a customer this Recipe does not serve.
- **Statuses are values, not a lifecycle.** Nothing here freezes an account over time.
