# cloudsmith

Emulates the Cloudsmith package-management API for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-10.**

Written against Cloudsmith's API reference at `api.cloudsmith.io` and struck live on 2026-09-10 with no credential, with a credential in the documented scheme, with one in the wrong scheme, and with a bare token carrying no scheme at all.

## What this Recipe found

**A credential in the wrong scheme is discarded and the request proceeds as anonymous.**

```
(no header)                     200 {"authenticated":false,"email":null, …}
Authorization: Bearer notreal   200 {"authenticated":false,"email":null, …}
Authorization: notreal          200 {"authenticated":false,"email":null, …}
Authorization: token notreal    401 {"detail":"Invalid authentication credentials
                                      provided in Authorization header."}
```

Cloudsmith's documented scheme is the word `token` — its own worked example is `curl -H "Authorization: token $key"`. Send the key under `Bearer`, which is what every HTTP library's built-in helper does, and the header is thrown away without comment.

The request succeeds. The body says `"authenticated": false`. And `/v1/repos/` answers `[]`.

So a developer who wires up the standard bearer helper gets **200 and an empty repository list**, and concludes their account is empty. There is no signal anywhere in the response that a credential was sent and ignored — no header, no field, no status. The only thing separating "I have no repositories" from "my authentication is misconfigured" is a boolean five fields into an endpoint the developer had no reason to call.

**And the one time it does challenge, it names a scheme it does not take.** That 401 carries:

```
www-authenticate: Basic realm="Cloudsmith API"
```

`Basic`, on an API documented to take `token`, and one that answers **200 anonymous** to a Basic-shaped header. A caller following the challenge authenticates the wrong way and lands back at silently anonymous.

That is the third distinct challenge-header fault in this collection, and between them the three cover the whole space:

| Recipe | What it does with `WWW-Authenticate` |
| --- | --- |
| [cloudamqp](../cloudamqp) | withholds it from the caller with no credential, sends it to the one with a wrong credential |
| [airship](../airship) | always sends it, filled with whichever scheme the caller already tried |
| cloudsmith | sends it once, naming a scheme that does not work |

**The error body is Django REST Framework's `detail`.** The sentence gives it away by naming the *header* rather than the key — accurate, and not written by Cloudsmith. [insightly](../insightly) publishes ASP.NET's `Message` for the same reason: the most precise failure sentences in this collection tend to be the ones a framework wrote.

**Four custom headers duplicate what the Link header carries.** `X-Pagination-Count`, `X-Pagination-Page`, `X-Pagination-PageTotal` and `X-Pagination-PageSize`, beside a `Link` header the documentation says carries `first`, `prev`, `next` and `last`. Two complete pagination systems on one response.

**And the reference cites a specification that does not exist.** The documentation labels that header **"Link (based on H+C6988)"**. RFC 5988 is Web Linking. H+C6988 is nothing at all — a mangled reference, published, in the very row that tells a reader where the format comes from.

**Empty datasets are never paginated**, says the same page. So all five headers disappear exactly when a client checking `X-Pagination-Count` for zero would want to read one.

## Modelling limits

- **The wrong-scheme behaviour is recorded and not served.** This is the headline finding and the format cannot express it: `public: when-absent` exempts a request that presented *no* credential, and there is no mode for "a credential this server could not parse is discarded and the request proceeds anonymously". So the emulator answers 401 to `Authorization: Bearer …` where Cloudsmith answers 200 — stricter than the provider, which is the direction a fake should not err in, and is why it is written down here in full.

  One provider showing a shape is a note; two make it a mechanism. If a second turns up, `public: when-unreadable` is the field it wants.
- **The repository listing is only ever served empty.** That is what an anonymous request gets live, and no repository record shape could be read without an account — so no fields are invented. `cauldron verify` reports one declared-and-unserved name, which is the correct complaint to leave standing rather than silence with a made-up record.
- **One resource and one listing.** Packages, entitlements, webhooks, vulnerability policies, quotas and the whole upload surface each want their own evidence — the reference lists several hundred operations.
- **No `spec:`.** Cloudsmith's reference is a rendered application, and no OpenAPI document answered at any of the usual addresses on 2026-09-10.
