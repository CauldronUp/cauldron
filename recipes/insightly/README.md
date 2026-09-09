# insightly

Emulates the Insightly CRM API (v3.1) for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-09.**

Written against Insightly's published Swagger 2.0 document at `api.insightly.com/v3.1/swagger/docs/v3.1` and struck live against `api.na1.insightly.com` on 2026-09-09 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**946 documented string fields say no valid value exists.**

The document carries `maxLength: 0` on **2,521** schema nodes, and **946** of those are typed `string`. `maxLength` constrains strings, and zero means the empty string is the longest permitted value — so on every one of those 946 fields, the document declares that nothing the API can send would validate.

The two that matter most on a Contact are its timestamps:

```
DATE_CREATED_UTC   type: string   format: date-time   maxLength: 0
DATE_UPDATED_UTC   type: string   format: date-time   maxLength: 0
```

A generator that emits validation from this document rejects every record the API returns. And the document is not even consistent about it: `DATE_OF_BIRTH`, three fields away, is also a `date-time` string and carries `maxLength: 255`.

The integers carry it too — `CONTACT_ID` and `OWNER_USER_ID` are `type: integer` with `maxLength: 0`, where the keyword does not apply and is ignored. That is the harmless half. The 946 strings are not.

**`FIRST_NAME` is required and `CONTACT_ID` is not.** The Contact schema's entire `required` list is one entry, and it is the given name. A contact must have a first name and need not have an identifier.

It is the third shape of that mistake in as many Recipes shipped this week:

| Recipe | Required | And yet |
| --- | --- | --- |
| [close](../close) | `name` | typed `["string", "null"]` |
| [adobesign](../adobesign) | `hidden` | `id` and `status` are not required |
| insightly | `FIRST_NAME` | `CONTACT_ID` is not on the list |

Three providers, three ways of guaranteeing the wrong field.

**The Authorization header is a hand-rolled parameter again.** `required: true`, `type: string`, `in: header`, on every operation — the same malpractice [close](../close) and [middesk](../middesk) publish. Three providers in one batch describe authentication as a string the caller assembles rather than as a security scheme, while a perfectly good `securityDefinitions` mechanism goes unused.

**One sentence answers everything, and ASP.NET wrote it.**

```json
401 {"Message":"Authorization has been denied for this request."}
```

Struck live for a missing credential, a wrong one, and a path that does not exist — the same body all three times.

`Message` is **PascalCase**, the only error field in this collection spelled that way, because it is ASP.NET Web API's own `HttpError` rather than anything Insightly composed. The passive voice names no actor and no reason. And because an unknown path answers it too, nothing about this API's surface is discoverable without a key that works.

**The host contains a pod.** `https://api.{pod}.insightly.com/v3.1/`, where the pod is `na1`, `eu1` and so on — the reference writes the hostname with a placeholder in it, and a caller has to know which shard its account is on before its first request.

[scalekit](../scalekit) refuses on exactly this axis and says so in four words. Insightly's apex answers normally, so a request to the wrong pod comes back as a credential problem rather than as a routing one — which is the harder failure to diagnose.

**The total is opt-in.** `count_total` is a boolean query parameter defaulting to off, so a listing reports how many records matched only if the caller asked in advance. `brief` is a second boolean that changes the shape of every record in the response.

**Paging borrows OData's names.** `top` and `skip`, on an endpoint that is not OData, defaulting to 100 and capped at 500.

**The field names are database columns.** `CONTACT_ID`, `DATE_CREATED_UTC`, `SOCIAL_LINKEDIN` — upper snake case on the wire, which is what a CRM built directly on its own schema ends up publishing.

**And the timestamps are not RFC 3339.** `2026-04-11 08:30:00` — a space where the `T` belongs and no zone at all, in a field whose name ends `_UTC`.

## Modelling limits

- **One route.** Listing contacts. Organisations, opportunities, projects, leads, tasks, events and custom objects account for most of the document's 303 paths, and each wants its own evidence.
- **`brief` and `count_total` are not modelled.** Both change the response, and this Recipe serves the full record with no total, which is the default either way.
- **The pod is not modelled.** The emulator serves the routes; reproducing "this shard is not yours" would mean modelling the hostname.
- **Nothing is mapped in detection.** No client for this API on npm, Packagist or the Go module proxy under an obvious name, checked 2026-09-09.
