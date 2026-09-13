# runpod

Emulates the Runpod pod listing for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Runpod serves without a credential at [`rest.runpod.io/v1/openapi.json`](https://rest.runpod.io/v1/openapi.json), and struck live on 2026-09-13 with no credential, with a wrong one, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A refusal typed as JSON with nothing in it.**

```
HTTP/1.1 401 Unauthorized
Content-Type: application/json
Content-Length: 0
```

Zero bytes, announced as JSON. `response.json()` throws on it, `response.text()` returns the empty string, and there is no code, no message and no identifier anywhere in the body — because there is no body. The only handle on the failure is a header.

**Two of those headers, for one request.** `x-request-id: req_1b210a59-…` and `x-trace-id: trace_e6aa3fb2-…`, both fresh uuids, differing only in their prefix. A support thread has to be told which one is wanted.

**A path that does not exist is a 400, and the lecture is the body.**

```
GET /v1/cauldron-nope   400 text/plain

[{"error":"At #/paths/get for GET https://rest.runpod.io/v1/cauldron-nope, The GET
request contains a path of '/v1/cauldron-nope' however that path, or the GET method
for that path does not exist in the specification. Suggestion: Check the path is
correct, and check that the correct HTTP method has been used (e.g. GET, POST, PUT,
DELETE)"}]
```

Bad Request, for a request that was not bad — it was addressed to nothing. The message is a spec-validating middleware talking to its own author: it names an internal JSON pointer, quotes the path back, and then lists the HTTP verbs for the reader. `#/paths/get` is not a pointer into any OpenAPI document; the wrong-method failure beside it says `#/paths/pods/put`, which is not one either, and the two are not even wrong the same way.

**And a JSON array labelled `text/plain`, next to an empty body labelled `application/json`.** One API, two adjacent requests, both content types wrong, in opposite directions.

**One record prints the same instant twice, once in JavaScript.**

```
lastStartedAt:     "2024-07-12T19:14:40.144Z"
lastStatusChange:  "Rented by User: Fri Jul 12 2024 15:14:40 GMT-0400 (Eastern Daylight Time)"
```

The second is `Date.prototype.toString()` with a sentence glued to the front, in a field named like a timestamp — carrying the *server's* timezone, spelled out in English, on a record delivered to callers anywhere.

**Two prices, one exampled as a number and one as a string.** `adjustedCostPerHr` is `type: number, example: 0.69`; `costPerHr` beside it is `type: number, example: "0.74"`. The declarations agree and the examples do not, and a generated client's fixtures come from the examples.

**The ports are on the record twice, in two encodings.** `ports` is `["8888/http", "22/tcp"]` and `portMappings` is `{"22": 10341}` — an object whose keys are port numbers written as strings and whose values are port numbers written as integers.

**Three string fields are exampled as null and never declared nullable.** `aiApiId`, `endpointId` and `templateId` are each `type: string` with `example: null`, while `publicIp` and `portMappings` on the same schema do say `nullable: true`. Nullability is declared twice and demonstrated three more times.

**A pod says what it is supposed to be.** `desiredStatus` is the only status field on the record; there is nothing saying what it currently is.

**And `ListPods` declares no 401, and two failures that are about one pod.** Its responses are 200, `400 "Invalid ID supplied."` and `404 "Pod not found."` — on an operation that takes no identifier, from a host that answers 401 to every anonymous request. `cauldron drift` reports the 401 as unbacked, and it is right: the only failure every caller meets is the one no operation declares. The listing also takes sixteen query parameters, six of them `include*` flags that change what a record contains.

## Sources

- [`rest.runpod.io/v1/openapi.json`](https://rest.runpod.io/v1/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `rest.runpod.io`, struck 2026-09-13 with no credential, a wrong credential, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The pod listing. Endpoints, templates, network volumes, container registry auth, billing and the serverless surface are the rest.
- **The `include*` flags are not modelled.** Six query parameters decide whether a record carries its machine, template, network volume, workers or savings plans. This Recipe serves the record without them, which is what an unadorned `GET /v1/pods` returns.
- **`x-request-id` and `x-trace-id` are constants.** Live they are fresh per request; the live cases assert their shape with a regex rather than their value.
- **The success fixture is document-derived.** Listing pods needs a real key. The records here are `Pod`'s own properties with the document's example values, including both timestamp formats verbatim, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Runpod is reached through `runpod` on PyPI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
