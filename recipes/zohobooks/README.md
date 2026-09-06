# zohobooks

Emulates the Zoho Books API (v3), for local development and tests.

**14 conformance cases, none checked against a live API.**

Written against Zoho's own OpenAPI document — the `invoices.yml` linked from the reference, 93 paths — fetched on 2026-09-06.

## What this Recipe found

**`page_context` is an array.** A listing answers `{"code": 0, "message": "success", "invoices": […], "page_context": [{…}]}` — the paging object is wrapped in a list with one element. So `page_context.has_more_page` is `undefined` and `page_context[0].has_more_page` is the boolean, and the loop everybody writes reads `undefined`, treats it as false, and stops after one page of ten. Zoho's own schema declares it `"type": "array"`.

**Money is a string.** `total` and `balance` are `"type": "string"` on an accounting API, so `a.total + b.total` concatenates, `Number(balance)` is a conversion someone has to remember, and `balance > 100` compares lexicographically — `"9.00"` sorts above `"100.00"`. `invoice_id` is a string too, described by Zoho as "a set of numeric characters": a number-looking value that must not be parsed as one.

**`code: 0` lives inside a 200, and the same key carries failures.** Success is `{"code": 0, "message": "success", …}`; failure is `{"code": <non-zero>, "message": …}`. The discriminator is a *value*, not a key — `if ("code" in body)` and `if (body.code !== undefined)` are both true on success, and `if (body.code)` only works by accident, because zero is falsy.

**The credential prefix is not Bearer.** `Authorization: Zoho-oauthtoken <access_token>` — an OAuth2 token under a word no OAuth library sends by default, so a standard bearer helper is refused with a perfectly valid token.

**There are eight base URLs and an account lives on exactly one.** `zohoapis.com`, `.eu`, `.in`, `.com.au`, `.jp`, `.ca`, `.com.cn`, `.sa`. The token, the organisation and the data all belong to one data centre; the same request against another is not a routing detail but a different installation of Zoho.

**`organization_id` rides on every request** as a query parameter — including calls that already name the record by id.

## Modelling limits

- **`page_context` is recorded, not served.** This format cannot describe a paging block shaped as a one-element array, so a listing here carries `code`, `message` and the records and no `page_context` at all. The shape is in the Recipe header so nobody meets it for the first time in production.
- **Invoices only, and a small part of them.** That one document has 93 paths — status transitions, email, reminders, attachments, templates — and each wants evidence rather than a shape copied from a neighbour.
- **One data centre**, the US host.
- **Statuses are values, not a lifecycle.** An invoice here does not move from draft to sent to overdue on its own.
