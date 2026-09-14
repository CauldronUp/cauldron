# astronomer

Emulates the Astro Platform API organisation listing for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Astronomer serves without a credential at [`api.astronomer.io/spec/platform/v1beta1`](https://api.astronomer.io/spec/platform/v1beta1), and struck live on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**Every failure is declared as JSON and arrives as plain text.** The document gives `ListOrganizations` a 400, a 401, a 403 and a 500, each pointing at an `Error` schema that `required`s `message`, `requestId` and `statusCode`. Live:

```
(no header)      403 text/plain  RBAC: access denied
Bearer notreal   401 text/plain  Jwt is not in the form of Header.Payload.Signature with two dots and 3 sections
```

Eleven bytes, and a sentence, with none of the three required fields and no JSON anywhere. A client built from the document parses both and throws.

**The 401 is a lesson in JWT structure.** "two dots and 3 sections" — the same kind of number spelled twice in one sentence, two ways — said to a caller who sent a string that was never meant to be a JWT.

**And a request carrying nothing is Forbidden.** `RBAC: access denied`, to a caller with no identity to deny, with the same answer for an unrouted path and a wrong method, because the check runs before the router.

**The Error schema permits a status of 600.** `statusCode` is `{minimum: 400, maximum: 600}`, and there is no 600.

**Two query parameters filter by product and their vocabularies do not overlap.** `product` is "Filters the Organization list by product" and `astronomerProduct` is "filter by astronomer product, should be one of ASTRO or OBSERVE" — lower case, no full stop, beside a page of proper sentences. Neither is deprecated. And the `product` *field* on the record is `["HOSTED", "HYBRID"]`, so the word means one thing in the query and another in the response.

**The support plan has twenty-one values.** `INACTIVE`, `INTERNAL`, `POV`, `TRIAL`, `BASIC`, `BASIC_PAYGO`, `TEAM_PAYGO`, `STANDARD`, `PREMIUM`, `BUSINESS_CRITICAL`, `BUSINESS`, `TEAM`, `ENTERPRISE`, `DEVELOPER`, `DEVELOPER_PAYGO`, `TEAM_V2`, `BUSINESS_V2`, `ENTERPRISE_V2`, `TRIAL_V2`, `ENTERPRISE_BUSINESS_CRITICAL`, `IBM_ENTERPRISE`. Five are `_V2` reissues sitting beside the originals, one names a partner, and two — `INACTIVE` and `INTERNAL` — are not plans at all.

**An internal abbreviation decides who can read your data.** `allowEnhancedSupportAccess` is "Whether the organization allows CRE to have view access to…" — CRE, undefined anywhere in the document, in the field that grants it.

**And every write is a POST to the thing itself.** `/organizations/{id}`, `/clusters/{id}` and `/deployments/{id}` each take `get`, `post` and sometimes `delete`, and no path in the document has a `put` or a `patch` — so updating is posting to the resource.

**The security scheme is named `JWT` and declared `{type: http, scheme: bearer}`**; `productPlans` is a plural array on the record beside a singular `productPlan` query filter; and the envelope requires `limit`, `offset`, `organizations` and `totalCount`, so a page always says how large the whole set is.

## Sources

- [`api.astronomer.io/spec/platform/v1beta1`](https://api.astronomer.io/spec/platform/v1beta1) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.astronomer.io`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route of twenty-three.** The organisation listing. Clusters, deployments, deploys, workspaces, teams, users, tokens, alerts, audit logs and the options endpoints are the rest.
- **The plain-text failures carry a `Content-Type` and nothing else.** Live they arrive with no `requestId` anywhere, which is the finding; the sandbox serves the same bytes.
- **The success fixture is document-derived.** Listing organisations needs a real token; the records here are `Organization`'s own properties with values from its declared enums, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Astronomer is reached through the `astro` CLI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
