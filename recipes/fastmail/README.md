# fastmail

Emulates Fastmail's JMAP API for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-10.**

JMAP — [RFC 8620](https://www.rfc-editor.org/rfc/rfc8620) — rather than REST or GraphQL. Written against the RFC and Fastmail's own published metadata, and struck live against `api.fastmail.com` on 2026-09-10 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The challenge header does the whole job, and it is the only one in this collection that does.**

```
GET /jmap/session     (no credential)
401 Content-Type: text/plain
WWW-Authenticate: Bearer resource_metadata="https://api.fastmail.com/.well-known/oauth-protected-resource/jmap/session"
```

That is [RFC 9728](https://www.rfc-editor.org/rfc/rfc9728), OAuth 2.0 Protected Resource Metadata. The URL resolves, and it answers:

```json
{"resource": "https://api.fastmail.com/jmap/session",
 "authorization_servers": ["https://api.fastmail.com"],
 "scopes_supported": ["urn:ietf:params:oauth:scope:mail",
                      "urn:ietf:params:oauth:scope:contacts",
                      "urn:ietf:params:oauth:scope:calendars",
                      "offline_access"],
 "bearer_methods_supported": ["header"],
 "resource_name": "Fastmail JMAP API",
 "resource_documentation": "https://www.fastmail.com/dev/"}
```

So an unauthenticated caller is told, mechanically: which authorisation server to go to, which scopes exist, that the token belongs in a header, what this resource is called, and where the documentation is.

Three Recipes shipped in the last two days get exactly this wrong, and between the four the whole space is covered:

| Recipe | `WWW-Authenticate` on a 401 |
| --- | --- |
| [cloudamqp](../cloudamqp) | not sent at all to the caller with no credential |
| [airship](../airship) | sent, filled with whichever scheme the caller already tried |
| [cloudsmith](../cloudsmith) | sent once, naming a scheme that does not work |
| fastmail | sent, and points at a document that answers every question |

It is worth recording that the correct answer exists, is standardised, and costs one header.

**And then the body beside it is a bare sentence in `text/plain`.**

```
(no credential)   No Authorization header
Bearer notreal    Invalid Authorization bearer parameters, not valid format
```

No JSON, no code, no field name, no envelope. **The machine-readable half of this failure is the best in the collection and the human-readable half is a string** — on the same response.

Both sentences are true, and the second is unusually precise: it says the *format* is wrong rather than the key. `notreal` is a syntactically valid bearer token by RFC 6750's grammar, so this is Fastmail saying its own tokens have a shape and that was not one — a distinction almost nothing else here draws.

**The standard's own example does not parse.** RFC 8620 §2.1 prints a Session object containing:

```
    "urn:ietf:params:jmap:mail": {}
    "urn:ietf:params:jmap:contacts": {},
```

There is no comma after the first `{}`. The example JSON in a Standards Track RFC, describing a JSON protocol, is not valid JSON.

Anyone copying it into a fixture — which is exactly what this Recipe was about to do — gets a parse error, and the reason is in the specification rather than in their editor.

**`accounts` is a map, and finding your own is two lookups.** The session carries `accounts`, keyed by opaque account id, and `primaryAccounts`, keyed by capability URN and valued by account id. So "the account I read mail with" is:

```js
session.accounts[session.primaryAccounts["urn:ietf:params:jmap:mail"]]
```

A lookup through one map to get a key for another. And somebody else's shared mailbox sits in the same `accounts` map, distinguishable only by `isPersonal` and `isReadOnly`.

**The object keys are URNs, and one of them is a URL.** `capabilities` is keyed by `urn:ietf:params:jmap:core` and friends, and a vendor extension is keyed by a full `https://` address. Namespaced identifiers as JSON object keys — correct, required by the specification, and something no JSON-to-struct mapping enjoys.

**`downloadUrl` is not a URL.** It is an [RFC 6570](https://www.rfc-editor.org/rfc/rfc6570) URI Template: `.../download/{accountId}/{blobId}/{name}?accept={type}`. A field ending `Url` holding something that has to be expanded before it can be fetched.

That is the same shape [matrix](../matrix) has with `mxc://`, arrived at differently: Matrix's is a scheme no client implements, this one is an address with holes in it. Two protocol-defined APIs, two ways of putting something un-fetchable in a field named for a URL.

**An unknown path is a 302.** Not a 404 — a mistyped JMAP path redirects, so a client following redirects lands on an HTML page and parses it as JSON.

## Modelling limits

- **One route and one well-known document.** The session resource and the metadata it points at. The JMAP API itself is a single `POST /jmap/api/` carrying an array of method invocations — `Email/query`, `Mailbox/get`, back-references between calls — and that is a protocol this format does not describe: one path, one method, and the shape of the answer determined entirely by the body.
- **No `spec:`.** JMAP is defined by RFC 8620 and RFC 8621 as prose and examples rather than as a machine-readable description, and Fastmail publishes no OpenAPI for it. That is what a protocol looks like when it is a protocol rather than one vendor's API.
- **The account identifiers and state token are fixtures.** They are shaped like Fastmail's own but were never read from a real session, because that needs an account.
- **Nothing is mapped in detection.** A project speaking JMAP holds a JMAP client pointed at whichever server it discovered, the same reason [matrix](../matrix) ships unmapped: the dependency says which protocol, never which host.
