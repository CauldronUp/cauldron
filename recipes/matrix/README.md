# matrix

Emulates the Matrix Client-Server API for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-10.**

Written against the Matrix specification's own OpenAPI sources at `github.com/matrix-org/matrix-spec` and struck live against `matrix.org` on 2026-09-10 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**Three failures, three codes, one vocabulary — and the specification defines all of them.**

```
(no token)        401 {"errcode":"M_MISSING_TOKEN","error":"Missing access token"}
Bearer notreal    401 {"errcode":"M_UNKNOWN_TOKEN","error":"Token is not active","soft_logout":false}
/cauldron_nope    404 {"errcode":"M_UNRECOGNIZED","error":"Unrecognized request"}
```

Every one of those `M_` codes is in the protocol specification rather than invented by this server, so a client can branch on them across Synapse, Dendrite and Conduit alike instead of against one vendor's choices.

That is the thing almost nothing else in this collection has: **an error vocabulary that is part of the protocol** rather than part of a product. Everywhere else, a code is a promise one company makes and can quietly change.

**And the 404 is JSON in the same shape as the 401.** Most providers here fall through to a web server on an unrouted path and answer HTML — [affirm](../affirm) sends a styled marketing page, [middesk](../middesk) sends zero bytes typed `text/html`, [airship](../airship) sends "Error 404". Matrix answers its own envelope, so one parser reads every failure.

**`soft_logout` is a third field and it is the useful one.** It tells the client whether the token was invalidated in a way that should discard local state — encryption keys, device identity — or merely expired and can be refreshed.

No other provider in this collection distinguishes those. Getting it wrong means either throwing away a device's keys unnecessarily, or keeping keys the server has already repudiated.

**One path is anonymous for GET and authenticated for POST.** The specification gives `GET /publicRooms` `security: null` — explicitly no authentication, not merely omitted — and gives `POST /publicRooms` `accessTokenBearer`. The POST is the filtered version, and carrying a filter is what costs you a token. Struck live: the GET answers 200 with real rooms and no credential at all.

**The published server URL names no host, and that is correct.**

```yaml
servers: ['{protocol}://{hostname}{basePath}']
```

Every component templated. Matrix is federated — there is no `api.matrix.com`, a client discovers its homeserver through `.well-known/matrix/client`, and any concrete host in the document would be a lie about somebody's deployment.

It is the honest version of a shape [adobesign](../adobesign) has by accident, whose `servers` entry is the bare path `/api/rest/v6` on an API that silently shards by region. Same missing host, opposite reasons.

**`avatar_url` is not a URL any client can fetch.** It is an `mxc://` URI — `mxc://t2l.io/0e92cc03f71f...` — which has to be rewritten into a `/_matrix/media/` request against your own homeserver before anything can load it. A field whose name ends `_url`, holding a scheme no HTTP client implements.

**The total is named an estimate.** `total_room_count_estimate` — optional, and the specification's own description says "an estimate on the total number of published rooms, if the server has one". After a collection full of `total_count` fields that quietly mean something narrower, a provider putting the approximation in the field name is worth recording.

**The array is called `chunk`.** Named for the streaming shape rather than for what is in it, so two endpoints returning different things share a key — the same trade [ironclad](../ironclad) makes with `list`.

**Three fields describe who may read a room**, and they are not redundant: `world_readable`, `guest_can_join` and `join_rule`. A room can be world-readable and closed to guests. A client showing a padlock has to decide which of the three it means.

**And the room directory contains things that are not rooms.** `room_type: "m.space"` marks a space — a container of rooms — returned in the same array as the rooms themselves. A client rendering a room list shows folders among the files unless it filters.

## Modelling limits

- **Two routes and a version list.** The public room directory, `joined_rooms`, and `/versions`. Sync, messages, state events, media, devices, encryption and the whole federation surface each want their own evidence — and sync in particular is a long-polling protocol this format does not describe.
- **`joined_rooms` is served raw.** It answers a bare array of room ID strings rather than records, so it exists here to give the credential failures a gated route rather than to model a resource.
- **No `spec:` is pinned.** Matrix publishes OpenAPI one file per endpoint, and every schema is a relative `$ref` into a sibling `definitions/` file — so no single published file resolves on its own, and fingerprinting one would compare a document whose fields cannot be read. That is the position [basiq](../basiq) and [customerio](../customerio) are in; here it is entered knowingly and written down rather than discovered later.
- **Nothing is mapped in detection.** The Matrix clients on npm — `matrix-js-sdk` and friends — speak this protocol to whichever homeserver a deployment configures, so the dependency says the project speaks Matrix and nothing about which server, which is the one thing an emulator would need to stand in for.
