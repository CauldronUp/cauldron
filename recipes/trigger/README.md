# trigger

Emulates the Trigger.dev run listing for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-14.**

Read from the management API reference at [`trigger.dev/docs`](https://trigger.dev/docs/management/runs/list), and struck live on 2026-09-14 with no credential, with a wrong secret key, with the key and no scheme, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**The problem document's `type` URI points at MDN.**

```json
{
  "title": "Unauthorized",
  "status": 401,
  "type": "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401",
  "detail": "No authorization header provided",
  "error": "No authorization header provided"
}
```

RFC 9457 asks the `type` URI to identify the *problem type*. This one identifies the status code — already in `status`, and on the status line — and it resolves to a third-party documentation site the API does not control, rather than to the specification or to Trigger's own docs. Anyone who follows it learns what 401 means, which they knew.

**And the same sentence is in the body twice.** `detail` and `error` carry identical strings, so a client reading either gets the same thing and a client reading both cannot tell which is authoritative.

**The response is pretty-printed.** Newlines and two-space indentation, on the wire, on every problem document this API sends.

**The well-formed envelope is for the failure that is not about your key.**

```
(no header)             401  application/problem+json   five fields, indented
Bearer <wrong secret>   401  application/json           {"error":"Invalid API key"}
<secret, no prefix>     401  application/json           {"error":"Invalid API key"}
```

The case a caller can fix by reading — a missing header — gets a title, a type, a detail and a link. The case that actually happens when a key is rotated gets three words in one field, in a different media type.

**Routing never happens.** An unrouted path and a method the path does not take both answer the "No authorization header provided" document. So this API cannot tell a caller their path is wrong until their key is right — and a caller *with* a wrong key is told the key is wrong even when the path does not exist.

**Eleven statuses, five of which are ways of failing.** `PENDING_VERSION`, `QUEUED`, `EXECUTING`, `REATTEMPTING`, `FROZEN`, `COMPLETED`, `CANCELED`, `FAILED`, `CRASHED`, `INTERRUPTED`, `SYSTEM_FAILURE`. What separates FAILED from CRASHED from INTERRUPTED from SYSTEM_FAILURE is not said in the field, and `CANCELED` is spelled with one L beside `REATTEMPTING`.

**`ttl` is a string or a number.** The reference gives its type as both, so one field holds either a duration expression (`"10m"`) or a count (`3600`), and the reader finds out at runtime.

**Money is a floating-point number of cents.** `costInCents` and `baseCostInCents` are `number`, not integer — the unit is in the name, the precision is binary, and there are two of them because the invocation is billed separately from the compute.

Also pinned: paging and filtering are nested objects flattened into the query string — `page[size]`, `page[after]`, `filter[status]`, `filter[createdAt][period]`; `pagination` carries `next` and `previous` that are run ids rather than cursors or URLs; `page[size]` is documented as 10 to 100 with a default of 25, so the smallest page anyone may ask for is ten; and `tags` is capped at ten entries of 1 to 128 characters each.

## Sources

- [Trigger.dev — list runs](https://trigger.dev/docs/management/runs/list) — the envelope, the run fields, the eleven statuses, and the `page` / `filter` parameter objects.
- Live: `api.trigger.dev`, struck 2026-09-14 with no credential, a wrong secret key, a key with no scheme, an unrouted path, and a wrong method.

## Modelling limits

- **No description is published.** Trigger.dev serves no OpenAPI document at any of the usual paths, so there is no `spec` to fingerprint.
- **The success side is document-derived.** Listing runs needs a real secret key; the records here are the reference's own field list with values of the documented types, and every case reading them is marked documentation-only.
- **The pretty-printing is described, not reproduced.** This Recipe serves compact JSON; the live problem documents arrive indented, which is the finding above.
- **`unknown_route` and `method_not_allowed` are declared as the credential failure.** That is what live sends, and it is what this Recipe serves — but no request ever reaches a routing failure here either, so the declarations name the shape rather than change behaviour.
- **`filter[createdAt]`, `filter[version]`, `filter[bulkAction]`, `filter[schedule]`, `filter[tag]` and `filter[error]` are not modelled.** Two filters are served, on `status` and `taskIdentifier`.
- **One route of many.** The run listing. Triggering, retrieving, replaying, cancelling, schedules, env vars and the queue surface are the rest.
- **Nothing is mapped in detection.** Trigger.dev is reached through `@trigger.dev/sdk` or a plain bearer call, and neither resolves to this host through a dependency file. Checked 2026-09-14.
