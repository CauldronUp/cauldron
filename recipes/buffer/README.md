# buffer

Emulates the Buffer legacy REST API for local development and tests.

**11 conformance cases, 8 checked against the live API on 2026-09-07.**

Written against Buffer's REST migration guide at `developers.buffer.com` and struck live against `api.bufferapp.com` on 2026-09-07 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**The two halves of one failure travel on two different channels.**

```
GET /1/profiles.json                        (no token)
401  Content-Type: text/html; charset=UTF-8
WWW-Authenticate: OAuth realm='Service', error='invalid_request',
  error_description='The request is missing a required parameter, includes an
  unsupported parameter or parameter value, repeats the same parameter, uses
  more than one method for including an access token, or is otherwise malformed.'
<empty body>

GET /1/profiles.json?access_token=notreal
401  Content-Type: application/json;charset=UTF-8
{"error":"The provided access token is invalid","code":401,"deprecation":{…}}
(no WWW-Authenticate)
```

A missing credential says everything in the header and nothing in the body. A wrong one says everything in the body and sends no header at all. So **a client that reads only the body sees a blank 401 for the commonest mistake there is**, and a client that reads only the header sees nothing for the second commonest. Neither request gets both.

It is the mirror image of [pipefy](../pipefy), which puts its OAuth error in the body and leaves `WWW-Authenticate` empty. Buffer does both, one each, on one endpoint.

**The header is single-quoted, which makes it unparseable.** RFC 9110 requires an auth-param value to be a token or a quoted-string, and a quoted-string uses double quotes. `realm='Service'` is neither: the value opens with an apostrophe, which is not a legal token character, so a conforming parser rejects the whole challenge rather than reading around it. All three parameters are quoted this way.

**And the error code is the wrong one.** RFC 6750 section 3.1 says that when a request carries no authentication information at all, the resource server SHOULD NOT include an error code. `invalid_request` is defined for a request that is malformed. This request was not malformed; it was empty.

**The description is the specification's own definition, pasted whole.** Those 240 characters are RFC 6750's text for `invalid_request`, listing five distinct ways a request can be malformed. None of the five is what happened, and the caller has to read all of them to find that out.

**`Content-Type: text/html` on zero bytes.** The header names a media type for a body that is not there.

**The retirement notice rides on the failure.** The wrong-credential body carries a `deprecation` object with a sentence, a `sunset` date of 2027-02-01, and a link to the migration guide:

```json
"deprecation":{"message":"The Buffer legacy REST API is deprecated and will be
retired on 1 February 2027. Please migrate to the GraphQL API before then.",
"sunset":"2027-02-01","link":"https:\/\/developers.buffer.com\/guides\/rest-migration.html"}
```

So the announcement that this whole API is going away is delivered **only** to callers who got their credential wrong. A working integration never sees it. RFC 8594 defines a `Sunset` *header* for exactly this purpose and it is not sent — the date is in a body that a successful request never receives.

(The forward slashes are escaped, which is PHP's `json_encode` default and a reliable fingerprint of what serialised the response.)

**The documentation is gone.** `developers.buffer.com` now describes the GraphQL API and nothing else. The only page still describing this surface is the migration guide the failure above links to, and it describes it by contrast: base URL, HTTP methods, `page`/`count` offset paging, and "Errors: HTTP status codes (401, 404, etc.)" — which is not what the 401 above does.

## Modelling limits

- **One route.** Profiles. Updates, schedules, links and the analytics endpoints each want their own evidence, and none of them has a reference left to check against.
- **The record shape is three fields, deliberately.** The legacy reference no longer exists; `id`, `service` and `avatar` are what the migration guide's own field mapping names for the equivalent object, and inventing the rest of a record from memory is exactly the guess this collection forbids.
- **Nothing is mapped in detection.** The packages that call this API are wrappers of an interface being retired in February 2027.
- **No `spec:`.** And no reference either.
