# dialpad

Emulates the Dialpad API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Dialpad's reference at `developers.dialpad.com` and struck live against `dialpad.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The error envelope is Google's, on an API that is not Google's.**

```json
{"error": {"code": 401,
           "errors": [{"domain": "global",
                       "message": "A valid API key must be provided.",
                       "reason": "required"}],
           "message": "A valid API key must be provided."}}
```

That is the Google JSON-C error format, field for field: a nested `error` object, a numeric `code`, an `errors` array whose entries carry `domain`, `message` and `reason`, and the message repeated at both levels.

Google deprecated it years ago in favour of `google.rpc.Status`. [firebaseauth](../firebaseauth) in this collection serves the newer shape, from the same company. So a client written against Firebase's errors does not parse Dialpad's, and a client written against Dialpad's has learned a decade-old convention it will not meet anywhere else.

**`domain` is `"global"` and `reason` is `"required"`.** Those are Google's own enumerations, not Dialpad's, and neither says anything about Dialpad. The message is the only Dialpad-specific string in the document, and it appears twice.

**`reason: "required"` is the answer to supplying a credential.** A missing key and a wrong key produce the same body, so the closest thing here to a machine-readable cause reports the wrong one.

**An unknown path is `text/plain`, and it is App Engine's.**

```
GET /api/v2/cauldron-nope
404  Content-Type: text/plain
Requested URL /api/v2/cauldron-nope not found
```

Between that and the JSON-C envelope, **both halves of this API's failure behaviour are Google infrastructure showing through** rather than anything Dialpad wrote.

**Identifiers are nineteen-digit Datastore ids sent as strings**, because they do not survive a double — the same trade [etcd](../etcd) and [waypoint](../waypoint) make, here without a word about it anywhere.

**Contact details are arrays even when there is one.** One email, one number, both in arrays — so there is no primary, and a client showing "the" address is choosing an index. A user with no number sends `[]` rather than omitting the key.

## Modelling limits

- **One route.** Users. Calls, contacts, rooms, transcripts, SMS, webhooks and the whole call-routing surface each want their own evidence.
- **Nothing is mapped in detection.** No client for this API turned up on npm, Packagist or the Go module proxy under an obvious name on 2026-09-07.
- **No `spec:`.** Dialpad publishes a rendered reference site.
