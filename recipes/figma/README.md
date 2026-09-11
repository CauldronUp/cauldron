# figma

Emulates the Figma REST API for local development and tests.

**8 conformance cases, 5 checked against the live API on 2026-09-11.**

Written against Figma's published OpenAPI document at `github.com/figma/rest-api-spec` and struck live against `api.figma.com` on 2026-09-11 with no credential, with a deliberately invalid one in each of the two schemes it accepts, and on a path that does not exist.

## What this Recipe found

**The document declares 403 where the API answers 401.**

`GET /v1/me` lists exactly four responses — 200, 403, 429 and 500. A request without a credential answers:

```json
401 {"status":401,"err":"Missing credentials"}
```

So the status every new integration meets first is the one status the operation does not declare, and the refusal it *does* declare is a different one. A client generated from this document has a branch for a refusal it will rarely see and none for the one it will see immediately.

`cauldron drift` says the same thing from the other direction, without being told:

```
not backed: no operation the Recipe routes to answers 401, which missing_credentials declares
not backed: no operation the Recipe routes to answers 404, which unknown_route declares
```

That is the fourth document in a week to omit the failure everyone meets, and each does it differently:

| Recipe | What its document says about failures |
| --- | --- |
| [middesk](../middesk) | declares only the success |
| [close](../close) | declares three, types every one as `"Any type"` |
| [attentive](../attentive) | declares five, gives none of them a body |
| figma | declares a **different status** instead of the real one |

**The failure field is called `err`.** Three letters, where the rest of this collection says `error`, `message`, `detail` or `reason`. It is required by the schema, alongside `status` — which is the HTTP status again, in the body, as a number.

**Both credential schemes are read, and both refusals are identical.** `X-Figma-Token: notreal` and `Authorization: Bearer notreal` each answer `{"status":401,"err":"Invalid token"}`.

The document declares four schemes — `PersonalAccessToken`, `PlanAccessToken`, `OAuth2` and `OrgOAuth2` — and **applies them properly, per operation**, which is worth saying after three providers this week that defined schemes and referenced none of them.

**The three sentences are accurate and distinct.** "Missing credentials" for a request carrying none, "Invalid token" for one carrying a wrong one, and "Not found" for a path that does not exist. Three failures, three sentences, none describing the wrong one — rarer here than it should be.

**The user is an `allOf` of two schemas.** `{id, handle, img_url}`, all required, composed with `{email}`, also required.

So the shape a client reads is assembled from two pieces and every field in it is guaranteed — including the avatar URL. A user who has never set a picture still has an `img_url`, because the schema does not permit otherwise.

**The identifier is a long number sent as text.** Nineteen digits, quoted, so a client that parses it to a number loses the tail.

## Modelling limits

- **One route.** The authenticated user. Files, nodes, images, comments, components, styles, webhooks, projects and the whole variables surface each want their own evidence — the document has 47 paths.
- **Nothing is mapped in detection.** A project that talks to Figma holds either a plugin manifest, which is not a dependency, or a generic HTTP client and a personal token — so no package name says this API is in use. Checked 2026-09-11.
