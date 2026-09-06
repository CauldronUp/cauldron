# aiven

Emulates the Aiven API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Written against Aiven's own OpenAPI document — `api.aiven.io/doc/openapi.json`, 321 paths — fetched on 2026-09-06. The document is Aiven's; the responses are this emulator's, and no case here claims otherwise.

## What this Recipe found

**Aiven's own spec contradicts itself about the credential.** The `securitySchemes` entry declares `"scheme": "bearer"`, which is what every OpenAPI generator reads and every HTTP helper implements. Its `description`, in the same object, says the header must be `authorization: aivenv1 <TOKEN>`. So a generated client sends `Authorization: Bearer …`, is refused, and the thing that told it to is the same file that says not to. A person reading the prose gets it right and a tool reading the schema gets it wrong.

**Every successful response is entitled to carry an `errors` array.** Listing services answers with `services`, `errors` and `message`, and only `services` is required; fetching one answers with `service`, `errors` and `message`, and only `service` is required. The reflex guard — `if (body.errors) throw new Error(body.message)` — is therefore testing for a key the *success* envelope may legitimately include, and whether it fires depends on whether Aiven sent an empty array or omitted it. Both are within the spec, and neither is the failure case.

**A service has four states and three of them are "up".** `state` is `POWEROFF`, `REBALANCING`, `REBUILDING` or `RUNNING`. `REBUILDING` is the expensive one: the service is serving, the URI resolves, and the plan is being moved underneath it. Code that treats anything other than `RUNNING` as down reports an outage that is not happening; code that treats anything other than `POWEROFF` as ready connects to a node that is going away. A powered-off service also keeps listing, with the plan and the disk it is still billed for.

**`termination_protection` is a field, not a permission.** A delete that comes back refused is this boolean, and a client that reads the refusal as an authorisation problem goes looking in the wrong place entirely.

## Modelling limits

- **One service type.** A real project holds PostgreSQL, Kafka and OpenSearch, each with its own `metadata`, `connection_info` and `components` shapes. Modelling one honestly is a Recipe of its own; modelling all of them from a schema would be inventing.
- **States, not transitions.** Creating a service here answers `RUNNING` rather than `REBUILDING` settling over minutes. The interesting thing is that a client must branch on four values, not how long each lasts.
- **Stated and not served: the top-level `message` on a failure.** Aiven's error envelope is a `message` *and* an `errors` array. The format here expresses one envelope, so the array is served — it is where Aiven's own schema puts `status` and `more_info`. A client reading `body.message` off a failure gets `undefined` here and a sentence from the real API, which is the wrong direction for a difference to run in, and is why it is written down rather than left to be found.
- **Plans, clouds and service types are not enumerated.** Those lists are an account's entitlements and vary per customer.
