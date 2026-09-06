# crowdsec

Emulates the CrowdSec Local API, for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against CrowdSec's own Swagger 2.0 document, published in its own repository — `crowdsecurity/crowdsec`, `pkg/models/localapi_swagger.yaml`, 13 paths, version 1.0.0 — read on 2026-09-06. The Local API runs on the operator's own machine and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**The block list contains entries marked "do not block", and the marker is optional.** A `Decision` carries `simulated`, described as "true if the decision result from a scenario in simulation mode". Simulation is how an operator tries a scenario without acting on it — so the decisions list, which is the list a bouncer reads to decide who to turn away, contains entries deliberately not meant to be enforced.

`simulated` is `readOnly` and not in the schema's `required` list, so it may be absent, and absent is falsy. That fails in the safe direction for truthiness — an absent flag reads as "not simulated", which enforces — and in the wrong direction for anyone *counting* simulated decisions. Either way, a bouncer that never reads the field enforces every simulation the operator was testing.

**The credential's header name is a header name and a value prefix.**

```yaml
JWTAuthorizer:
  type: apiKey
  in: header
  name: "Authorization: Bearer"
```

Under `in: header`, `name` is the header's name. `Authorization: Bearer` is not one — it contains a colon and a space, neither legal in a header name. A generator that does what the document says produces a header called `Authorization: Bearer` with the token as its value, which most HTTP libraries reject outright and the rest send as a malformed line. The other scheme, `APIKeyAuthorizer` (`X-Api-Key`), is correct — so one of the two credentials in this document cannot be used by anything generated from it.

**A decision's lifetime is a string, and there are two of them.** `duration` is `type: string` — `"4h"`, `"168h"`, Go's own duration syntax, with no day unit — and `until` is a separate string holding an absolute date. Both describe when a decision stops applying, in different units, and nothing says which wins if they disagree. Arithmetic on either needs a parser.

**Two fields whose meaning depends on which way the call is going:**

| field | the document's words |
|---|---|
| `id` | "(only relevant for GET ops) the unique id" |
| `uuid` | "only relevant for LAPI->CAPI, ignored for cscli->LAPI and crowdsec->LAPI" |

So a decision's identity depends on the operation and on which of four directions the traffic is flowing, said in prose rather than in the schema. There is no shape describing a decision as it exists — only a shape with two identifiers and a note about when each is real.

**The type is an open enum with no members.** `type` is "the type of decision, might be 'ban', 'captcha' or something custom" — three named possibilities, none declared as an enum, and an explicit invitation to invent more. No generated client can have a type for it and no switch over it can be exhaustive.

**`scope` says what `value` is measured in.** An address and a CIDR block sit in the same field, so code parsing `value` as an IP throws on every `Range` decision.

**A plural field holding one string.** `ErrorResponse` has `message` (required) and `errors`, which is `type: string` and described as "more detail on individual errors". Individual, plural, and a single string — so a client iterating it iterates characters.

**Nine filters and no pagination.** `GET /decisions` takes `scope`, `value`, `type`, `ip`, `range`, `contains`, `origins`, `scenarios_containing` and `scenarios_not_containing` — and no limit, no page, no cursor. On a busy installation this is the whole blocklist. The streaming endpoint beside it exists for exactly that reason and is a different shape.

**No credential is 403, not 401**, so a client branching on 401 to fetch a credential never fetches one.

**HEAD is declared beside GET on every read** — unusual in an OpenAPI document, and worth knowing that a client may probe with it.

## Modelling limits

- **Nothing here is verified against a live API.** The Local API runs on the operator's own machine.
- **Decisions, listed.** 13 paths is a security agent: alerts, watchers, the decisions stream, allowlists and usage metrics each want their own evidence — and the stream in particular is a different shape from the listing.
- **The credential checked here is `X-Api-Key`,** the scheme this document gets right. The JWT scheme is recorded and not served: serving a header named `Authorization: Bearer` would mean emitting a malformed header line to demonstrate a defect in a document.
- **`duration` is served as the Go duration string the provider sends.** Nothing here converts it, and nothing reconciles it with `until`.
