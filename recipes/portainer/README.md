# portainer

Emulates the Portainer container-management API, for local development and tests.

**12 conformance cases, none checked against a live API.**

Written against Portainer's own generated Swagger 2.0 document, published in its own repository — `portainer/portainer`, `api/docs/swagger.yaml`, 224 paths, version 2.45.0 — read on 2026-09-06. Portainer is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**A 2xx that means some of it did not happen.** The bulk delete of environments answers **207**, described in the document as "Partial success. Some environments were deleted successfully, while others failed."

207 is in the 2xx range. `response.ok` is true. `if (status < 300)` passes. `raise_for_status()` does not raise. Every idiom for "did that work" says yes — and the answer is that some of the environments a caller asked to remove are still there.

**And the report of what failed is a list of integers:**

```json
{"deleted": [3], "errors": [4]}
```

`errors` is an array of environment *identifiers*. Not errors — ids. The response says which environments failed and nothing whatever about why: no code, no message, no per-item status. A caller wanting a reason has to retry them one at a time.

**The destructive verb is the deprecated one.** `DELETE /endpoints` carries, in its own description, "Deprecated: use the `POST` endpoint instead." So the supported way to remove many environments is `POST /endpoints/delete`, and a reviewer scanning a codebase for destructive calls by HTTP verb finds nothing. The safe-looking verb is the one being retired.

**The field names are Go struct names.** `Id`, `Name`, `URL`, `Type`, `Status`, `TLSConfig`, `PublicURL`, `TagIds`, `UserAccessPolicies` — PascalCase on the wire, because the structs carry no JSON tags. A client that lowercases the first letter, or maps camelCase by convention, reads `undefined` from every one. And `URL` is three capitals where `Url` would be the camelCase guess.

**Two integers carry the two things you most want to branch on:**

| field | values |
|---|---|
| `Status` | 1 up, 2 down, 3 provisioning, 4 error |
| `Type` | 1 Docker, 2 agent, and on |

Bare integers with the meaning in a description rather than a name. `endpoint.Status === 1` is "up", there is no string anywhere to compare against, and every value is truthy — so `if (endpoint.Status)` is true for a dead environment.

**Access control travels inside each environment.** `UserAccessPolicies` and `TeamAccessPolicies` are embedded in the record, keyed by user or team id **as a string** — a Go map with an integer key, and JSON has no integer keys. So a listing of environments is also a dump of who may reach each one, with identifiers that are strings on the wire and integers everywhere else.

**Paging is `start`, not `offset` or `page`,** and the listing is a bare array with no total anywhere — so a caller pages by asking for more until it gets fewer. There are twelve query parameters on that listing (`search`, `groupIds`, `status`, `types`, `platformTypes`, `outdated`, `excludeGroupIds`, `tagIds` and more), so the same endpoint is a filter language and still cannot say how many things there are.

**A 404 surfaces the database's own words:** `"object not found inside the database"`.

## Modelling limits

- **Nothing here is verified against a live API.** Portainer is self-hosted and there is no public instance.
- **Environments, listed and fetched.** 224 paths is a container manager: stacks, containers, images, volumes, networks, registries, teams, the whole Kubernetes proxy surface and the edge agent protocol each want their own evidence.
- **The 207 is served raw rather than through the error table.** Portainer's own schema for that body is `{deleted, errors}` and nothing else, so routing it through the errors machinery would add a `message` the provider never sends.
- **`Status` and `Type` are served as the integers the document declares.** Nothing here maps them to names, because the provider does not.
