# firebaseauth

Emulates Firebase Authentication — Google's Identity Toolkit v1 API — for local development and tests.

**13 conformance cases, all of them checked against the live API on 2026-09-06.**

Written against Google's own machine-readable description of the API — the Identity Toolkit v1 discovery document at `identitytoolkit.googleapis.com/$discovery/rest?version=v1`, 41 methods, 61 schemas, 446 properties, revision `20260813` — and struck live against `identitytoolkit.googleapis.com` on 2026-09-06. Every live observation here was made with no credential and then with a deliberately invalid one, which is all the failure path needs.

## What this Recipe found

**The document describes no failure.**

61 schemas, and not one of them is the error envelope. There is no `Status` schema, no `GoogleRpcStatus`, and the string `"error"` occurs **once** in 183KB. The two schemas with "Error" in the name — `ErrorInfo` and `BatchDeleteErrorInfo` — are per-row failures carried *inside* a 200 response from the batch endpoints, each an `{index, message}` pair.

So the only errors this document describes are the ones that arrive with a success status.

`EMAIL_NOT_FOUND` and `INVALID_PASSWORD` — the two strings every Firebase Auth client on earth switches on — appear nowhere in it.

**A missing credential is 403. A wrong one is 400. Neither is 401.**

Struck live, no key at all:

```
POST https://identitytoolkit.googleapis.com/v1/accounts:lookup
HTTP 403
{"error": {"code": 403,
           "message": "Method doesn't allow unregistered callers ...",
           "errors": [{"message": "...", "domain": "global", "reason": "forbidden"}],
           "status": "PERMISSION_DENIED"}}
```

Struck live, key present and invalid:

```
POST .../v1/accounts:lookup?key=not-a-real-key
HTTP 400
{"error": {"code": 400,
           "message": "API key not valid. Please pass a valid API key.",
           "errors": [{..., "reason": "badRequest"}],
           "status": "INVALID_ARGUMENT",
           "details": [{"@type": "type.googleapis.com/google.rpc.ErrorInfo",
                        "reason": "API_KEY_INVALID", ...},
                       {"@type": ".../google.rpc.LocalizedMessage", ...}]}}
```

Both credential failures, and neither is 401. A client written the ordinary way — `if (status === 401) refresh()` — fires on neither. And 400 is the status for malformed input, so **a wrong credential is reported in the same class as a malformed request body**: a caller that retries 400s never retries, and a caller that treats 400 as "my payload is broken" looks for the bug in the wrong place.

**The error format is a query parameter, and it is called `$.xgafv`.**

The document declares a global parameter named `$.xgafv`, enum `["1", "2"]`, described as "V1 error format". Struck live, it changes the body:

| | `error.code` | `error.message` | `error.errors[]` | `error.status` | `error.details[]` |
|---|---|---|---|---|---|
| `$.xgafv=1` (default) | ✓ | ✓ | ✓ | ✓ | ✓ |
| `$.xgafv=2` | ✓ | ✓ | — | ✓ | ✓ |

v2 drops the legacy `errors` array and keeps `details`. So one failure has two body shapes, a query parameter picks which, the parameter's name begins with `$.` and has to be percent-encoded as `%24.xgafv` to survive a URL builder — and **neither shape is written down anywhere in the document that declares the switch**.

The message itself arrives **three times** in one v1 body: at `error.message`, at `error.errors[0].message`, and at `error.details[1].message`.

**The machine-readable reason is three levels down, in an array keyed by type.** `API_KEY_INVALID` — the only part of the 400 a program should branch on — lives at `error.details[N].reason`, and `N` is not fixed: entries are identified by an `@type` URL, so a client has to scan the array for `type.googleapis.com/google.rpc.ErrorInfo` rather than index it.

**41 methods, 35 of them POST, and no DELETE anywhere.**

Six GETs, thirty-five POSTs, no PUT, no PATCH, no DELETE. Deleting an account is `POST v1/accounts:delete`. Reading accounts is `POST v1/accounts:lookup`, with the filter in the request body. 35 of the 41 paths carry a colon action — `accounts:signUp`, `accounts:signInWithPassword`, `accounts:sendOobCode` — Google's AIP convention, the same shape [exoscale](../exoscale) and [boundary](../boundary) use and the opposite of [cilium](../cilium)'s, which puts a namespace *before* the colon rather than a verb after it.

The practical effect: **the HTTP method carries no information here.** Every cache, every proxy rule and every "is this safe to retry" heuristic that keys on the verb sees POST for reads, writes and deletions alike.

**One record, five timestamps, four encodings.**

From `UserInfo`, the account object this Recipe serves:

```
createdAt          {"type": "string", "format": "int64"}            milliseconds
lastLoginAt        {"type": "string", "format": "int64"}            milliseconds
validSince         {"type": "string", "format": "int64"}            SECONDS
lastRefreshAt      {"type": "string", "format": "google-datetime"}  RFC 3339
passwordUpdatedAt  {"type": "number", "format": "double"}           milliseconds
```

`createdAt` and `validSince` have **byte-identical type declarations** and differ by a factor of a thousand. Nothing in the schema distinguishes them; only the English descriptions do, and only if read. A generator produces the same Go `string` and the same TypeScript `string` for both, and a client that treats them alike lands in 1970 or in the year 58000.

`passwordUpdatedAt` is the same quantity as `createdAt` in a different JSON type, on the same object — so `typeof user.createdAt === typeof user.passwordUpdatedAt` is `false` for two fields that both mean "a moment".

**13 response fields say "Always present" in prose.**

Discovery has no `required` list for responses, so the document says it in English, thirteen times: "Always present in the response." Every generator makes all thirteen optional. The strongest is `SignInWithPasswordResponse.registered` — *"Whether the email is for an existing account. Always true."* — a boolean that is **simultaneously deprecated and documented as constant**, still emitted, carrying no information at all.

Nearby, `SignInWithPasswordRequest.returnSecureToken` is documented in full as *"Should always be true."* — an instruction where a description belongs, on an optional field that decides whether the response contains a token.

**49 properties are deprecated, 13 of them named `kind`.** `kind` is the JSON-API-era discriminator, marked deprecated on thirteen separate response schemas and still sent by all of them.

**`writeOnly` appears 0 times in 446 properties** — in a document whose subject is credentials, which declares `password` on `SignUpRequest` and `SignInWithPasswordRequest`, and `passwordHash` and `salt` on `UserInfo`. After [coolify](../coolify) at 0 of 192, [casdoor](../casdoor) at 0 of 234 and [exoscale](../exoscale) at 1 of 261, this is the fourth and the sharpest: the fields that most need the marker are here and unmarked. `readOnly` is used 7 times, so the vocabulary was available.

**Field visibility depends on the kind of credential, stated in prose.** `UserInfo.passwordHash` and `UserInfo.salt`: *"Only accessible by requests bearing a Google OAuth2 credential."* `GetProjectConfigResponse.allowPasswordUser`: *"only returned for authenticated calls from a developer."* One schema, two audiences, and which fields arrive is not a function of anything in the request a client can see.

**`{+targetProjectId}` is decorated with an operator the document then turns off.** 31 path placeholders use RFC 6570's `+` — reserved expansion, meaning slashes in the value are *not* escaped — and the top level declares `fullyEncodeReservedExpansion: true`, which says to escape them anyway. The `+` is left in the templates regardless.

**`basePath` and `servicePath` are both the empty string** while `baseUrl` is `https://identitytoolkit.googleapis.com/`, so every path in the document is written with no leading slash: `v1/accounts:signUp`. Joining `baseUrl` to `path` works. Joining an origin to `path` yields `https://hostv1/accounts:signUp`.

**`GET v1/projects` returns one object.** The method id is `identitytoolkit.getProjects` — plural — and there is no id parameter and no array in the response: it is `GetProjectConfig`, the configuration of the project the API key belongs to. A path that reads like a collection, a name that reads like a collection, and a single record.

**`GET v1/publicKeys` has no declared response at all.** Every other method in the document names a response schema. That one has none, because it answers a bare map of key id to PEM certificate, which this schema language cannot express. So the document describes 41 methods and declines to describe one of their answers.

**The paging token ends empty, not null.** `nextPageToken`: *"If there are more accounts to be downloaded, a token ... Otherwise, this is blank."* The second document in this collection to end a pagination with `""` rather than `null`, after [waypoint](../waypoint). Empty is falsy, so `while (token)` terminates by accident and `token !== null` never does.

**And the product is named twice.** Everyone calls it Firebase Auth. The document calls itself "Identity Toolkit API", `canonicalName` is "Identity Toolkit", the host is `identitytoolkit.googleapis.com`, and the word *Firebase* appears in the title of nothing. A developer searching the description of the API they are using for the name of the product they are using does not find it.

## Detection

Six client packages across three ecosystems, each checked by fetching its published archive and grepping for the host this Recipe serves:

| Package | Ecosystem | Names `identitytoolkit.googleapis.com` |
|---|---|---|
| `@firebase/auth` | npm | yes, and `accounts:signInWithPassword` outright |
| `firebase` | npm | yes (the meta-package; it also depends on `@firebase/auth`) |
| `firebase-admin` | npm | yes |
| `firebase-auth-lite` | npm | yes — 9KB, no dependencies, straight to the REST API |
| `kreait/firebase-php` | Packagist | yes |
| `firebase.google.com/go/v4` | Go | yes, in 7 files under `auth/` |

**`firebase-auth` — the obvious name — is left out, and that is the finding.** It is version 0.1.2, last published 2016-03-19, "a simple wrapper around firebase token authentication", and it depends on `firebase-token-generator`, which belongs to the *pre-Google* Firebase. Ten years stale, and the name a developer would try first. The same shape as [hackernews](../hackernews)'s `hacker-news-api` and [wikipedia](../wikipedia)'s: the obvious name belongs to something else — here an abandoned wrapper rather than a competitor's product.

## Modelling limits

- **No `spec:` URL, deliberately.** Google publishes a complete machine-readable description and it is a **Google Discovery Document**, not OpenAPI — `kind: discovery#restDescription`, no `openapi` key, no `swagger` key. `cauldron drift` reads it as a format it cannot parse, and there is no fingerprint to take because the fingerprint is taken over parsed paths. Naming a URL nothing can check is what the shipping rule exists to prevent. The same call [codecov](../codecov) makes for a different reason: there the document is OpenAPI and invalid, here it is valid and not OpenAPI.
- **Seconds-as-a-string has no type here.** `validSince` is declared as a plain `string`, with the unit carried by the literal in the fixture. This format has `timestamp_ms_string` — which [clickup](../clickup) and `createdAt` both need — and no seconds counterpart. One provider needing a type is not enough to add one; the case is recorded so the next provider that needs it finds it already made.
- **`accounts:lookup` serves the whole fixture set.** The real endpoint filters on `localId`, `email`, `phoneNumber` or `federatedUserId` **in the request body**, and this format matches filters from the query string only. The filter names are recorded and not routed.
- **The `$.xgafv` switch is recorded and not served.** Selecting between two whole error envelopes on a query parameter is not something a Recipe can declare, and inventing a half of it would be worse than the note.
- **`USER_NOT_FOUND` and its siblings are recorded and not served.** Firebase's not-found is a bare screaming-snake token where every other provider modelled here sends a sentence — but every read in this API is a POST whose subject arrives in the body, so there is no route on which this format can raise it. Serving it from the listing would put a 400 where the real API sends a 200 and an empty array.
- **Credential-dependent fields are recorded and not served.** `passwordHash`, `salt` and `allowPasswordUser` would mean modelling two classes of credential, and the finding is that one schema covers both.
- **Sign-in is not modelled.** `accounts:signUp` and `accounts:signInWithPassword` are where the `returnSecureToken`, `expiresIn`-as-a-string and `registered` findings come from, and they are read out of the document rather than served: an emulator that mints tokens is a different and larger thing than one that serves records.
