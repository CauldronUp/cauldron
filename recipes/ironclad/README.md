# ironclad

Emulates the Ironclad API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Ironclad's reference at `developer.ironcladapp.com` and struck live against `ironcladapp.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**One response to three mistakes, padded with spaces.**

```
GET /public/api/v1/records          (no credential)
GET /public/api/v1/records          Authorization: Bearer notreal
GET /public/api/v1/cauldron-nope    Authorization: Bearer notreal

401 { "code": "UNAUTHORIZED", "message": "invalid authentication token" }
```

Byte-identical all three times, **including the whitespace** — a space after the opening brace and a space before the closing one. A serialiser configured for readability rather than for bytes, so every fixture anybody records by hand carries the padding. The same class of tell as [gong](../gong)'s space *before* each colon and [tinybird](../tinybird)'s space *after*.

**"invalid authentication token" is the answer to sending no token** — the seventh provider in this collection to describe the wrong failure. It is also lower case throughout, beside a `code` in screaming snake, so one object has a shouted constant and an uncapitalised sentence.

**An unknown path is a 401.** The route is judged after the credential, so a path cannot be shown not to exist — and for a contracts API whose whole surface is behind an account, that means nothing about it is discoverable from outside.

**The path says `public`.** `/public/api/v1/…` — a segment distinguishing this surface from Ironclad's internal one. A reasonable thing for a company to need, and an odd thing for a customer to type: every URL carries a word that is only meaningful inside the company.

**Every property is wrapped in an object.** `record.properties.counterpartyName.value` — two hops even for a string, because a property carries its own metadata. A client that forgets `.value` gets an object rather than a value, and an unset property is **absent** rather than a wrapper holding null, so two records in one listing have different property sets.

**The envelope key is `list`** — named for the shape rather than the contents, so two endpoints returning different things share a key.

## Modelling limits

- **One route.** Records. Workflows, approvals, signatures, launch forms and the whole webhook surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Ironclad publishes a rendered reference site.
