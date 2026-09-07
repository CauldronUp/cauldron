# gong

Emulates the Gong API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Gong's reference at `gong.app.gong.io/settings/api` and struck live against `api.gong.io` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The refusal discloses the running build, the branch and the number.**

Struck live, on a request carrying no credential at all:

```
x-gong-service-version: PublicApiServer:17.main.20701
```

The service name, a major version, the branch it was cut from, and the build number. `main` is in there. So an anonymous 401 says which of Gong's services answered, which line of development it came from, and how many builds have been cut.

The same class of disclosure [imgix](../imgix) makes with `meta.server.version` — and here it is a **header** rather than a body, so nothing that logs response bodies will ever capture it and nothing that reviews response schemas will ever see it.

**Two correlation ids, neither labelled as the one to quote.** `x-traceid` in the headers and `requestId` in the body, different values on the same response. imgix does this too. A support ticket gets whichever one the client's error handler happened to keep.

**`errors` is an array of strings.**

```json
{"requestId" : "4cr7gi74j5xsi8t8orw",
 "errors" : ["Validate credentials failed. Please check your credentials and try again."]}
```

Not objects — so there is no code, no field name, no pointer, only prose. The array exists because one request can fail several ways at once, and each way is a sentence. [honeybadger](../honeybadger) puts a bare string in a field with the same name; Gong puts strings in an array, which sits between "one message" and "structured failures" and has the drawbacks of both.

**The JSON is pretty-printed with a space before every colon.** `"requestId" : "…"` — Jackson's `DefaultPrettyPrinter`, the Java default nobody turns off. Harmless on the wire, and it means every byte comparison, every hand-recorded fixture and every snapshot test written against this API carries the spacing.

**An unknown path answers 401, not 404**, with the credentials sentence — so route existence is not discoverable without a working key, and a typo is reported as an authentication problem.

**A private call is in the listing without its contents.** The record comes back and the transcript does not, so a client counting calls counts it and a client fetching content does not — two different totals from one listing.

**A duration is seconds with no unit anywhere.** `2714`, which is forty-five minutes. Nothing in the name or the record says seconds rather than milliseconds, and both readings produce a plausible call length.

`WWW-Authenticate: Basic realm="Gong.io API"` is correct, which is rarer in this collection than it should be.

## Modelling limits

- **One route.** Calls. Transcripts, users, workspaces, CRM objects, trackers, scorecards and the whole data-privacy surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **The pretty-printing is described, not served.** The emulator sends compact JSON; reproducing Jackson's spacing would be asserting a serialiser setting rather than a contract, and no conformance case here can see whitespace.
- **No `spec:`.** Gong publishes its reference behind an account.
