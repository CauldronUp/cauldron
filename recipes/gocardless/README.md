# gocardless

Emulates the GoCardless API (2015-07-06), for local development and tests.

**14 conformance cases, none checked against a live API.**

Written against GoCardless's own API reference — `docs.gocardless.com/docs/api-reference` — read on 2026-09-06. The documentation is GoCardless's; the responses are this emulator's, and no case here claims otherwise.

## What this Recipe found

**The envelope key is plural even when it holds one object.** GoCardless's own example of fetching a single mandate is `{"mandates": {"id": "MD123", …}}` — not `mandate`, and not an array. The same key that holds a list on the collection endpoint holds one object here, so `body.mandates[0].id` works on a list and reads `undefined` on a get, while `body.mandate` is `undefined` on both. Code written against one endpoint and reused against the other fails on a key that is present, correctly spelled, and the wrong shape.

**The error array changes its keys depending on why the request failed.** A failure is `{"error": {type, code, message, documentation_url, request_id, errors: […]}}`, and what sits inside `errors` depends on `type`: a `validation_failed` entry carries `field` and `request_pointer`, everything else carries `reason`. So `error.errors[0].field` is a string on a 422 and `undefined` on a 409, and a client rendering it against a form has nothing to highlight.

**`type` is the discriminator and the status is not.** There are four types, and `invalid_api_usage` alone covers 400, 401, 403, 404, 405, 406, 415, 426 and 429. One status does not identify a failure; one type spans nine of them.

**A mandate has ten states and one is "working".** `pending_customer_approval`, `pending_submission`, `submitted`, `active`, `suspended_by_payer`, `failed`, `cancelled`, `expired`, `consumed`, `blocked`. Three are on their way to working, one works, six are not coming back. Branching `active`-versus-everything-else treats a mandate two days from collecting the same as one the payer has blocked — and `suspended_by_payer` is neither `cancelled` (the merchant) nor `failed` (the bank).

**The version header is required and is not a version of your client.** Every request must carry `GoCardless-Version: 2015-07-06`, a fixed date unchanged since it was minted, so it reads like boilerplate someone forgot to update. Omitting it is refused, not defaulted to the newest.

## Modelling limits

- **Mandates only.** Payments, refunds, payouts, subscriptions and the mandate-import flow have their own lifecycles and want evidence of their own rather than a shape copied from this one.
- **`links` is served; its targets are not.** A mandate points at a customer, a creditor and a bank account. Following one of those here 404s, because this Recipe models the mandate.
- **Stated and not served:** `consent_parameters`, `authorisation_source`, `funds_settlement` and `consent_type`. They belong to schemes this Recipe does not model, and inventing values would be describing a scheme rather than reproducing one.
