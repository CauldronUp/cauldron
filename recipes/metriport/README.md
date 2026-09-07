# metriport

Emulates the Metriport medical API for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-07.**

Written against Metriport's reference at `docs.metriport.com` and struck live against `api.metriport.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**Two words, three situations, one status.**

```
GET /medical/v1/patient          (no credential)
GET /medical/v1/patient          x-api-key: notreal
GET /medical/v1/cauldron-nope    x-api-key: notreal

403 {"message":"Forbidden"}
```

Byte-identical all three times. A missing credential, a wrong credential and an unknown path are one response, and the response is a single word in a single field.

That sentence is AWS API Gateway's default for a request its usage-plan check rejected — so **the gateway is answering and Metriport's own application is never reached**. The API's whole surface is opaque from outside: a caller cannot discover a route, cannot tell a typo from a permissions problem, and cannot tell whether their key was rejected or never read.

It sits at one end of a range this collection now has three points on:

| | Answer to all three mistakes |
|---|---|
| [statsig](../statsig) | eleven bytes, **no content type** |
| **metriport** | two words, with a content type |
| [healthgorilla](../healthgorilla) | a full `OperationOutcome` — and 42 resource types enumerated |

The last two are both clinical APIs, at opposite ends of the same decision.

**403, not 401**, for a request carrying nothing at all — so the status that tells a client to authenticate never arrives, and the one that does means "you are known and not allowed", which nobody was.

**Identifiers are UUIDv7.** The version nibble is 7, so the first 48 bits are a millisecond timestamp: these identifiers sort chronologically and **leak when a patient record was created**, which is useful and is also a fact about a person in a clinical system.

**A patient can be held by several facilities**, so "which facility is this patient at" is a list and a client showing one has picked an index.

**Dates of birth are whole dates** — unlike [healthgorilla](../healthgorilla) next door, where FHIR permits a year alone. A record crossing between the two has to invent or discard a day.

## Modelling limits

- **One route.** Patients. Documents, consolidated data, facilities, organisations, FHIR conversion and the whole webhook surface each want their own evidence.
- **`fromItem` is served exclusive and is inclusive on the real API.** This format has no inclusive-cursor option; the difference is one record, and it is recorded rather than approximated — the same call the [telegram](../telegram) Recipe makes.
- **Nothing is mapped in detection.** Metriport's SDKs are generated per language and named for the company rather than the surface.
- **No `spec:`.** Metriport publishes a rendered reference site.
