# umami

Emulates the Umami Cloud analytics API for local development and tests.

**12 conformance cases, 5 checked against the live API on 2026-09-11.**

Struck live against `api.umami.is` on 2026-09-11 with no credential, with a deliberately invalid one, with the self-hosted scheme's header, and on a path that does not exist. Record shapes are read from Umami's own source, which is open, rather than from its reference site — see the sources below.

## What this Recipe found

**A missing credential is a 400, not a 401.**

```
(no header)                400 {"error":{"message":"No API key specified.",
                                 "code":"bad-request","status":400}}

x-umami-api-key: notreal   401 {"error":{"message":"Invalid API key.",
                                 "code":"unauthorized","status":401}}
```

The absence of authentication is modelled as a malformed request, and 401 is reserved for a credential that was present and wrong.

There is an argument for it — nothing was attempted, so there is nothing to be *unauthorized* about — and it breaks the check every client writes. `if (res.status === 401)` is false for the commonest failure there is, so a token-refresh path keyed on 401 never fires, and a 400 handler that assumes a malformed body goes looking at the request it sent.

What it buys is a real distinction: **two failures, two statuses, two codes.** Most providers in this collection cannot tell the two apart at all.

**A header in the wrong scheme is an invalid key rather than a missing one.** `Authorization: Bearer notreal` answers "Invalid API key.", not "No API key specified."

The reason is that self-hosted Umami authenticates with `Authorization: Bearer <token>` and Umami Cloud with `x-umami-api-key`, and one codebase serves both. So Cloud reads the self-hosted scheme's header, finds a credential in it, and rejects it as a key.

That is the fourth provider in this collection on that axis, and the first to read it this way round:

| Recipe | `Authorization` in the wrong shape reads as |
| --- | --- |
| [plausible](../plausible) | missing |
| [wise](../wise) | missing |
| [todoist](../todoist) | missing |
| umami | **invalid** |

The same request, four providers, two opposite verdicts — and nothing a client can inspect to know which it will get.

**`code` is the status as a slug.** `bad-request` beside `status: 400`, `unauthorized` beside `status: 401` — three representations of one fact, twice in the body and once on the status line. Here they always agree, which is worth saying after [todoist](../todoist), whose `error_code` is `477` and matches nothing.

**And the failure is nested.** `error.message`, where most of this collection puts the sentence at the top level. A client reading `body.message` or `body.error` as a string gets an object.

**The listing says when its own count is a lie.** Umami's paged result carries `isCapped` beside `count` — the response telling a client that the total it just reported was truncated rather than counted.

Almost nothing else here admits that. [sourcegraph](../sourcegraph) faces the same cost and reports a flat `totalCount: 0` on a page holding records; [matrix](../matrix) puts the word "estimate" in the field name. Three providers, three ways of handling a count too expensive to compute, and only one of them lies.

**A website may have no domain.** The schema types `domain` nullable, on a product whose entire purpose is measuring traffic to a domain.

**And it may belong to nobody.** `userId`, `teamId` and `createdBy` are three separate nullable references — an owner, a team, and whoever made it — so a record can carry all three, one, or none.

**Deletes are soft.** `deletedAt` is a column, so a website that has been removed is still a row, and any endpoint that forgets to filter returns it. `resetAt` is a fourth timestamp beside it, for wiping a site's statistics without deleting the site — a distinct event with its own column.

## Sources

- Live: `api.umami.is`, struck 2026-09-11.
- [`umami.is/docs/api/authentication`](https://umami.is/docs/api/authentication) for the two schemes.
- [`prisma/schema.prisma`](https://github.com/umami-software/umami/blob/master/prisma/schema.prisma) — the `Website` model.
- [`src/lib/types.ts`](https://github.com/umami-software/umami/blob/master/src/lib/types.ts) — `PageResult` and `PageParams`.

The record and envelope shapes come from the product's own source rather than from a docs page, which is a better source and is only available because Umami is open.

## Modelling limits

- **One route.** Websites. Stats, events, sessions, reports, funnels, goals, segments, teams and the whole session-replay surface each want their own evidence.
- **`isCapped` is served as false.** It is a real field on every paged result; what makes a listing capped is a `maxResults` the caller sets, and modelling that would mean modelling the cap rather than the field.
- **Self-hosted Umami is not modelled.** It authenticates with a login token from `POST /api/auth/login` and serves from `/api/` rather than `/v1/`, so it is a second surface with its own evidence to gather.
- **Nothing is mapped in detection.** A project using Umami holds a script tag or the `@umami/node` tracker, which posts events; neither is a client of this reading API. Checked 2026-09-11.
