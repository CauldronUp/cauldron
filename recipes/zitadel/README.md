# zitadel

Emulates the Zitadel Management API for local development and tests.

**10 conformance cases, 5 checked against the live API on 2026-09-11.**

gRPC behind a grpc-gateway. Struck live against `zitadel.cloud` on 2026-09-11 with no credential, with a deliberately invalid one, and on a path that does not exist. Record shapes come from Zitadel's own protobuf, which is open.

## What this Recipe found

**The body's `code` is not an HTTP status.**

```
(no header)      401 {"code":16,"message":"auth header missing"}
Bearer notreal   401 {"code":16,"message":"Errors.Token.Invalid (AUTH-7fs1e)", …}
/cauldron-nope   404 {"code":5,"message":"Not Found"}
```

16 is gRPC's `UNAUTHENTICATED`; 5 is its `NOT_FOUND`. The numbering comes from a different protocol, arriving in an HTTP body beside an HTTP status that already said the same thing in its own numbering — and a client switching on `code` needs the gRPC status table to read it.

**And `WWW-Authenticate` is not a challenge.**

```
www-authenticate: auth header missing
www-authenticate: Errors.Token.Invalid (AUTH-7fs1e)
```

Neither begins with an auth-scheme token, which RFC 9110 requires: a challenge is `<scheme> <params>`, so a parser reads "auth" as the scheme and "header missing" as malformed parameters after it. The field is being used as a free-text error channel.

That is the fifth challenge-header fault in this collection, and the most thorough:

| Recipe | What it does with `WWW-Authenticate` |
| --- | --- |
| [cloudamqp](../cloudamqp) | omits it from the caller who needs it |
| [airship](../airship) | fills it with whichever scheme the caller already tried |
| [cloudsmith](../cloudsmith) | names a scheme that does not work |
| [fastmail](../fastmail) | points at an RFC 9728 metadata document |
| zitadel | puts a sentence where the scheme belongs |

**The sentence is an untranslated resource key.** `Errors.Token.Invalid` is the *name* of a string in a message bundle, not a string. `(AUTH-7fs1e)` is an internal error id.

So the human-readable field contains the identifier of the human-readable text, and a caller who looks the key up finds it in Zitadel's i18n files rather than in its API reference.

**The failure carries a protobuf type URL.** The wrong-token body has `details: [{"@type":"type.googleapis.com/zitadel.v1.ErrorDetail", …}]` — `google.protobuf.Any`, serialised into JSON, repeating the same id and message one level down. A JSON client is handed protobuf's own reflection machinery.

**`sequence` means two different things depending on the direction of the call.** Zitadel's own comment on it, in the proto:

> on read: the sequence of the last event reduced by the projection
> on manipulation: the timestamp of the event(s) added by the manipulation

One field, one type, and its meaning chosen by whether you were reading or writing — an event counter coming back from a GET and a **timestamp** coming back from a POST. Nothing on the wire distinguishes them.

It is also a `uint64`, which protobuf's JSON mapping requires to be emitted as a **string**, so the field that is sometimes a counter and sometimes a clock is always quoted.

**The proto says `primary_domain` and the wire says `primaryDomain`.** The gateway lower-camels every field name, so the schema a reader consults and the keys a client receives are spelled differently throughout.

**And every enum has a zero value meaning "unset".** `ORG_STATE_UNSPECIFIED = 0` sits before `ACTIVE`, `INACTIVE` and `REMOVED`, because protobuf requires a zero. So `state` can come back saying nothing, and a client switching on it needs a branch for a value that means the field was never set.

## Sources

- Live: `zitadel.cloud`, struck 2026-09-11.
- [`proto/zitadel/org.proto`](https://github.com/zitadel/zitadel/blob/main/proto/zitadel/org.proto) — the `Org` message.
- [`proto/zitadel/object.proto`](https://github.com/zitadel/zitadel/blob/main/proto/zitadel/object.proto) — `ObjectDetails`, and the comment on `sequence`.
- [`proto/zitadel/management.proto`](https://github.com/zitadel/zitadel/blob/main/proto/zitadel/management.proto) — `OrgState`.

## Modelling limits

- **One route.** The caller's own organisation. Users, projects, applications, grants, policies, actions and the whole IAM surface each want their own evidence — `management.proto` is 595KB.
- **No `spec:`.** Zitadel's API is defined in protobuf, and its OpenAPI is generated per service at build time rather than published at an address. The proto is the source that document is made from, so it is what this Recipe cites.
- **The instance host is not modelled.** Every Zitadel customer gets its own domain; `zitadel.cloud` answers the failure paths above and nothing else.
- **Nothing is mapped in detection.** Zitadel's published clients are generated gRPC stubs per language, and a project holding one is talking to whichever instance it was configured with — the dependency names the protocol, never the host.
