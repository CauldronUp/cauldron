# appwrite

Emulates the Appwrite health surface for local development and tests.

**5 conformance cases, all of them checked against the live API on 2026-09-13.**

Struck live against `cloud.appwrite.io` on 2026-09-13 with no credential, with a project header naming a project that does not exist, with an API key beside it, and on a path that does not exist.

## What this Recipe found

**Every failure carries the API's version number.**

```
/v1/account       401 {"message":"User (role: guests) missing scopes ([\"account\"])",
                       "code":401,"type":"general_unauthorized_scope","version":"2.0.0"}
wrong project     404 {"message":"Project with the requested ID could not be found. Please
                       check the value of the X-Appwrite-Project header…","code":404,
                       "type":"project_not_found","version":"2.0.0"}
/v1/cauldron-nope 404 <!DOCTYPE html>… (the Appwrite console)
```

`version` is on every JSON failure, and nowhere on a success except the one endpoint whose entire body it is. `GET /v1/health/version` answers `{"version":"2.0.0"}` — so the only endpoint that will talk to an anonymous caller returns exactly the field every failure already carries.

**The message contains an escaped JSON array.** `missing scopes ([\"account\"])` — a JSON array serialised into the human sentence, with its quotes escaped for the second encoding. The machine-readable part of the failure is embedded as text inside the human-readable part, so a client that wants the scope has to parse prose.

**And the anonymous caller is a user with a plural role.** `User (role: guests)` — one caller, described as a user, with the role `guests`.

**A project that does not exist is a 404, with or without a key.** Sending `X-Appwrite-Project: cauldron-nope` answers `project_not_found` whether an `X-Appwrite-Key` accompanies it or not. The project is resolved before the credential is read, so a mistyped project id and a missing resource are the same answer, and the key is never examined.

**The 404's message is instructions.** "Please check the value of the X-Appwrite-Project header to ensure the correct project ID is being used." — a sentence addressed to a person, in the field a client reads.

**`code` is the HTTP status and `type` is the code.** `"code":401` duplicates the status line, and the string a client would actually branch on — `general_unauthorized_scope`, `project_not_found` — is under `type`. The field named for a code carries a number already in the response, and the one named for a category carries the code.

**One health endpoint is public and its siblings are not.** `/v1/health/version` answers 200 to a caller with nothing; `/v1/health` and `/v1/health/time` answer `401 missing scopes (["health.read"])` to the same caller.

**And an unrouted path answers the console.** A full `<!DOCTYPE html>` page with a viewport meta tag and a link to Appwrite's own favicon, served from the API host where a JSON client expected an error.

## Sources

- Live: `cloud.appwrite.io`, struck 2026-09-13. Every case here is a live one.

## Modelling limits

- **Four routes, and three of them exist to fail.** The public version endpoint, its two neighbours that refuse, and a project lookup that cannot find one. Accounts, databases, documents, storage, functions, messaging and teams each want their own evidence — and all of them need a project that exists, which no anonymous probe can reach.
- **The project header is not modelled as a credential.** Appwrite resolves `X-Appwrite-Project` before it examines `X-Appwrite-Key`, which Cauldron has no way to express; the project failure is served by a route that always answers it.
- **The unrouted 404 serves the page's opening.** Live it is the whole Appwrite console; here it is the first tag, which is enough to reproduce what breaks.
- **No `spec:`.** Appwrite generates its OpenAPI per platform at build time and serves the reference as a documentation site, so there is no single document at an address to fingerprint.
- **Nothing is mapped in detection.** The `appwrite` and `node-appwrite` packages are real clients, but a project holding one is pointed at whichever endpoint and project id its configuration names — self-hosted or cloud — and neither is in the dependency. Checked 2026-09-13.
