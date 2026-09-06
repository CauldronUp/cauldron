# freshbooks

Emulates the FreshBooks accounting API, for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against FreshBooks' own reference — `freshbooks.com/api/start` and `/api/errors` — read on 2026-09-06.

## What this Recipe found

**Two wrappers before the payload, and the paging lives inside the second one.** A list answers `{"response": {"result": {"clients": […], "page": 1, "pages": 1, "per_page": 15, "total": 2}}}`. The records are at `response.result.clients`, two envelopes deep, and the four paging numbers are *siblings of the array* rather than a meta object beside it — so anything that walks `result` looking for "the data" has to know which key is not an integer.

**The key changes with the operation as well as the depth.** A single record is `{"response": {"result": {"client": {…}}}}` — singular. Plural on the listing, singular on the fetch, both two levels down.

**The account id is a path segment on every accounting call.** `/accounting/account/{accountId}/users/clients` — not a header, not a query parameter, so an integration cannot hold one base URL and append endpoints to it.

**FreshBooks has its own error numbers, and they are not HTTP statuses.** 1001 is `RequiredField`, 1002 `InvalidValue`, 1003 `AccessDenied`, and the published table runs to 1108. They travel with a much coarser status — 401, 403, 404, 409, 500 — so the status says what kind of refusal it was and the number says which one. Only the number is specific enough to branch on.

## Modelling limits

- **The failure envelope is not modelled from evidence.** FreshBooks publishes the error *numbers* and their names; no page reachable without an account shows the JSON a failure arrives in. What this Recipe serves is its own convention — a flat object with a message — and the numbers in it are FreshBooks'. If the real envelope nests failures under `response` the way successes are nested, this is wrong in the shape while right in the codes, and it is written down so whoever finds out has somewhere to correct.
- **Clients only.** Invoices, expenses, payments, projects and time entries are separate resources, and the projects half of the API lives under a different path root with a different identifier — worth its own Recipe rather than a guess appended to this one.
- **`vis_state` is not modelled.** The field that marks a record deleted rather than removing it is real, and exactly the kind of thing this project exists for; it wants an observation rather than an assumption.
