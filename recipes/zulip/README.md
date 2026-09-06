# zulip

Emulates the Zulip chat API (v1), for local development and tests.

**11 conformance cases, 6 checked against the live API on 2026-09-06.**

The live cases were struck against `chat.zulip.org`. `/server_settings` needs no account, and the refusal was recorded by asking for messages without one. Written against Zulip's own OpenAPI document, published in its own repository — `zulip/zulip`, `zerver/openapi/zulip.yaml`, 119 paths.

## What this Recipe found

**The path says v1 and the real version is an integer in the body.** Every Zulip endpoint lives under `/api/v1` and has for the life of the product. The version that decides what a client may do is `zulip_feature_level`, a monotonically increasing integer answered by `/server_settings`. Struck live:

```json
"zulip_version": "12.0-919-gb1bd7d7d3e",
"zulip_feature_level": 509
```

The documentation's own words: a value of N means the server supports all API features introduced before feature level N. So a client cannot know what the server does from the URL it is calling, and cannot know it from `zulip_version` either — that string is a Git reference on any server tracking main. It has to call `/server_settings` first and branch on an integer. For self-hosted Zulip this is the whole game: two servers on the same `/api/v1` can be hundreds of feature levels apart.

**A success tells you which of your parameters it ignored.** Since Zulip 7.0 (feature level 167) a successful response may carry `ignored_parameters_unsupported`, an array naming the parameters the endpoint did not understand.

That is the same hazard the [gitea](../gitea) Recipe in this collection records from the other side — Gitea silently ignores `per_page` and then echoes it back in the `Link` header — except Zulip says so out loud, in the success body. A client has to look.

It is not universal. Struck live: `/server_settings?notaparam=1&another=2` answered `{"result":"success","msg":""}` with no ignored array at all. So the field's absence does not mean everything was understood; it means either that, or that this endpoint does not report.

**Every response carries `result` and `msg`, successes included.**

```
success:  {"result":"success","msg":"", ...}
failure:  {"result":"error",
           "msg":"Not logged in: API authentication or user session required",
           "code":"UNAUTHORIZED"}
```

`msg` is on both — the empty string on success, so truthiness happens to work and `"msg" in body` does not discriminate at all. `code` is on the failure only, which makes its presence the reliable test. It is also the one a client is least likely to reach for, because `result` reads like the obvious discriminator and is a string that has to be compared rather than a status that can be branched on.

**The credential is HTTP Basic, and the 401 says so.** Struck live: `WWW-Authenticate: Basic realm="zulip"`. The username is the user's email address and the password is the API key — so the credential is two values joined by a colon and base64'd, not a token, and a client library built around bearer tokens has nowhere to put the email.

**`realm_url` and `realm_uri` are both sent, with the same value.** Struck live, both `https://chat.zulip.org`. The spec says `realm_uri` was deprecated in Zulip 9.0 (feature level 257) because the term URI is deprecated in web standards. Nothing in the response marks which is which.

**`is_incompatible` is the server's opinion of your client** — a boolean in `/server_settings` meaning the client that sent this request is deemed incompatible with this server, decided from the request. A field about the caller, in a response about the server.

**Paging is bidirectional around an anchor, and the anchor is a string.** `GET /messages` takes `anchor`, `num_before` and `num_after`: a position and how far to read in each direction. No page number, no cursor. The anchor's declared type is `string`, and its values are message ids — which are integers — plus the words `newest`, `oldest` and `first_unread`. A field that is a string and usually a number, where arithmetic is a type error exactly when it is not.

## Modelling limits

- **`/server_settings` and a messages listing.** 119 paths is a chat platform: streams, topics, subscriptions, users, drafts, scheduled messages, reactions, the events queue and the real-time protocol each want their own evidence.
- **Paging is not modelled, deliberately.** Zulip's anchor-plus-two-counts is neither offset paging nor cursor paging, and this format has neither shape. Declaring it as offset paging made the emulator read the anchor — a message id like `1904489` — as an offset, and answer an empty page for every real anchor a caller could hold. That is worse than not paging: it teaches a client its correct request is empty. So this route answers the whole fixture with `found_newest` and `found_oldest` both true, which is what Zulip answers when the window covers everything.
- **Only the password half of the credential is checked.** Zulip is unusual in that the username half varies too — it is the user's own email, not a constant the way Mailgun's `api` is — so a request with the right key and somebody else's email is accepted here and refused by Zulip.
- **`narrow` is not modelled.** It is a JSON-encoded array inside a query parameter, which is its own escaping problem.
- **`ignored_parameters_unsupported` is served empty** on the messages listing, because the conformance requests send nothing unsupported. The point is that the key exists to be read.
