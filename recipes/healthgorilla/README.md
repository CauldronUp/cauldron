# healthgorilla

Emulates the Health Gorilla FHIR R4 API for local development and tests.

**12 conformance cases, 6 checked against the live API on 2026-09-07.**

Written against Health Gorilla's reference at `developer.healthgorilla.com` and struck live against `api.healthgorilla.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**An unauthenticated 404 lists the entire API surface.**

```
GET /fhir/R4/CauldronNope
404 {"resourceType":"OperationOutcome",
     "issue":[{"severity":"error","code":"processing",
               "diagnostics":"HAPI-0302: Unknown resource type 'CauldronNope' - Server knows how
                              to handle: [AllergyIntolerance, Bundle, CarePlan, CodeSystem, …
                              StructureDefinition, Subscription, User, ValueSet]"}]}
```

**42 resource types**, in a 656-character string, to a request carrying no credential at all.

That is HAPI FHIR's default handler being helpful: it names everything the server can serve so a caller can see what they mistyped. It is also a complete map of a clinical data API, free, before authenticating — and it arrives from a framework default rather than from a decision anybody made.

The mitigation elsewhere is real: reading a `Patient` needs a credential, so this discloses the shape rather than the data. It is still the most thorough enumeration in this collection.

**The statuses are inverted.**

| Sent | Status | Body |
|---|---|---|
| nothing | **403** | a full `OperationOutcome` |
| `Bearer notreal` | **401** | *empty* |

403 for sending nothing and 401 for sending something wrong — the opposite of the convention, and the opposite way round from every other provider here that distinguishes them at all.

And the richer response is on the anonymous request. **The failure with something to say is the one nobody had to authenticate for, and the failure a real integration will hit has no body at all.**

**The route is resolved before the credential is judged**, which is what makes the enumeration reachable anonymously and is the opposite of most providers here.

**The failure is a FHIR resource.** `OperationOutcome` has a `resourceType` discriminator and an `issue` array, so a client parses failures with the same code that parses records. That is the best-designed error model in this collection — and it is not Health Gorilla's, it is the FHIR specification, and every FHIR server has it.

**`severity` is `fatal` on the credential failure and `error` on the routing one.** FHIR defines fatal, error, warning and information; the ordering is defensible and a client switching on severity behaves differently for each.

**The media type is `application/fhir+json`** — the third distinct `+json` suffix in this collection after `vnd.api+json` and `problem+json`.

**A birth date can be a year and a month.** FHIR dates are deliberately partial — `YYYY`, `YYYY-MM` or `YYYY-MM-DD` — so `birthDate` is not a date, and every parser that treats it as one either throws or invents a day.

## Modelling limits

- **One resource type of 42.** Patients. Conditions, observations, documents, coverage, care plans and the rest each want their own evidence.
- **The resource is named `resource` in this Recipe, deliberately.** A FHIR Bundle wraps every entry under the literal key `resource` whatever type it holds, because one Bundle can hold several — where Chargebee and Maxio wrap each item under its own type name. This format takes the wrapper key from the resource's name, so the name is chosen to put the right word on the wire.
- **The 42-type list is truncated in the served fixture.** The full string is 656 characters; the case asserts its shape and the count is recorded here.
- **Nothing is mapped in detection.** FHIR clients are FHIR clients: a project talking to Health Gorilla depends on a generic FHIR library and a base URL, which no dependency name distinguishes.
- **No `spec:`.** The FHIR R4 specification describes the resources, and a `CapabilityStatement` would describe this server's own subset — but that needs a credential.
