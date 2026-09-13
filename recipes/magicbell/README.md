# magicbell

Emulates the MagicBell notifications API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from MagicBell's own generated client — it is open — and struck live against `api.magicbell.com` on 2026-09-13 with no credential, with a wrong project key, with a wrong bearer token, and on a path that does not exist.

## What this Recipe found

**An unrouted path answers 405, and says the only method it takes is OPTIONS.**

```
GET /cauldron-nope   405 text/plain, zero bytes
                     allow: OPTIONS
```

Not 404. A path that does not exist is reported as a path whose GET is not allowed — and the `Allow` header offers exactly one method, the CORS preflight. A caller who mistypes an endpoint is told the endpoint exists and that the only thing it will do is answer a preflight.

**Every failure carries a suggestion and a help link.**

```json
{"errors":[{"code":"api_key_not_provided",
            "message":"API Key not provided",
            "suggestion":"Please provide a 'X-MAGICBELL-API-KEY' header containing your MagicBell project's API key",
            "help_link":"https://documenter.getpostman.com/view/2269098/2sAYdhLAjv"}]}
```

Four fields, two of them written for a person: a `suggestion` naming the exact header to send, and a `help_link`. Most APIs in this collection manage one sentence and no link.

**And the help link points at Postman.** `documenter.getpostman.com/view/2269098/…` — a third-party documentation host and an opaque numeric view id. The URL a failing caller is sent to does not belong to MagicBell, and nothing about it says which endpoint it describes.

**A correctly formatted bearer token is refused for being badly formatted.** Struck live with `Authorization: Bearer notarealtoken`:

```json
{"code":"invalid_jwt_format","message":"expected authorization header format: Bearer <token>"}
```

The header sent was exactly `Bearer <token>`. The message describes the failure that did not happen and says nothing about the one that did — and in passing it tells a caller who has proved nothing that the credential is a JWT.

**Two credentials, and which one it complains about depends on which header you touched.** No headers at all complains about the missing `X-MAGICBELL-API-KEY`; a wrong `Authorization` complains about the JWT. The API has two ways in, and the failure describes whichever one the caller reached for.

**A notification's state is in six places.** `status` is a string, and beside it are `sentAt`, `seenAt`, `readAt`, `archivedAt` and `discardedAt` — five nullable timestamps, each recording one transition. Nothing says what `status` holds when two of them are set, or which wins.

**`content` is capped at ten megabytes.** `z.string().max(10485760)` on a notification body — a field that arrives in a bell menu.

**`customAttributes` is `z.any()`** — no shape at all, on a field that ships in every record.

**And the collection makes both of its keys optional.** `{data?: Notification[], links?: Links}`, so `{}` is a conforming response to a request for notifications and a client cannot tell an empty page from a malformed one. The `links` object has `first`, `next` and `prev` and no `last`, so there is no way to jump to the end.

## Sources

- Live: `api.magicbell.com`, struck 2026-09-13.
- [`notification.ts`](https://github.com/magicbell/magicbell-js/blob/main/packages/magicbell-js/src/user-client/services/notifications/models/notification.ts) — the record, its five timestamps and its size caps.
- [`notification-collection.ts`](https://github.com/magicbell/magicbell-js/blob/main/packages/magicbell-js/src/user-client/services/notifications/models/notification-collection.ts) and [`links.ts`](https://github.com/magicbell/magicbell-js/blob/main/packages/magicbell-js/src/user-client/services/common/links.ts) — the envelope, and what it makes optional.

## Modelling limits

- **One route.** Listing a user's notifications. The project client, channels, integrations, broadcasts, categories, topics and the mark-as-read mutations each want their own evidence.
- **The JWT failure is recorded, not served.** MagicBell reads two credentials and answers about whichever header the caller touched; Cauldron declares one carrier per Recipe, so this serves the project-key failures and the `invalid_jwt_format` transcript is quoted above.
- **The generated client renames every field.** Its response schema reads `action_url`, `created_at`, `user_id` and so on, and transforms them to camelCase for application code. This Recipe serves the camelCase shape the client hands a developer, which is the one a test will assert against.
- **No `spec:`.** MagicBell generates its clients from an internal description and publishes its reference through Postman's documenter, so there is no document at an address to fingerprint.
- **Nothing is mapped in detection.** The `magicbell` packages are real clients, but a project holding one is pointed at whichever project its key belongs to — and most MagicBell integrations are the embeddable inbox component, which is markup rather than a dependency on this API. Checked 2026-09-13.
