# exoscale

Emulates the Exoscale cloud API (v2), for local development and tests.

**8 conformance cases, none checked against a live API.**

Written against Exoscale's own OpenAPI 3.0 document, served by Exoscale at `openapi-v2.exoscale.com/source.json` — 261 paths, version 2.0.0 — read on 2026-09-06. Every endpoint needs a signed request against a real account, so every claim here is the provider's own document.

## What this Recipe found

**The word "security" appears zero times in the whole document.** No `securitySchemes`, no top-level `security`, no per-operation `security` — not one occurrence of the string in a megabyte of JSON describing an API that provisions machines, networks and firewall rules.

Exoscale does authenticate, with an `EXO2-HMAC-SHA256` signature over the request, and a request signature is genuinely awkward to express in OpenAPI — which has no vocabulary for "sign these headers with this key". So this is a *third* reason for the same absence, and the three are worth holding together:

| provider | declares | needs | verdict |
|---|---|---|---|
| [argocd](../argocd) | nothing | a bearer token | an omission |
| [cilium](../cilium) | nothing | nothing — a root-owned unix socket | correct |
| **exoscale** | nothing | a signature OpenAPI cannot describe | unrepresentable |

Three documents identical on this point, meaning three different things. A scan for "APIs with no declared authentication" finds all three and can rank none of them.

**286 of 303 schemas have a hyphen in their name, and so do their fields.** `security-group`, `security-group-rule`, `external-sources`, `start-port`, `end-port`, `flow-direction`, and the listing key `security-groups`.

`body.security-groups` is a *subtraction* in JavaScript, not a property access. It has to be `body["security-groups"]` everywhere, and in every language whose struct fields cannot contain a hyphen a generated client renames all of them — differently per generator.

**`writeOnly` appears exactly once, and it is on the firewall rules.** One occurrence in a megabyte, against 169 uses of `readOnly`:

```yaml
rules:
  type: array
  writeOnly: true
  description: Security Group rules
```

`writeOnly` means "may be sent, is never returned". So a listing of security groups tells you the groups exist and **nothing about what they permit**. Auditing what a firewall actually allows takes a second request per group, and the field sits in the schema to be read as though it were not.

This is the exact opposite of [coolify](../coolify) in this collection, whose document uses `writeOnly` nowhere and declares a password on a response. Both are one keyword away from the other's problem.

**A GET returns a machine's password.**

```
GET /instance/{id}:password
"Reveal the password used during instance creation or the latest password reset."
```

The password is not in the URL, so the URL itself is safe. The response body is not: a GET is the method every cache, proxy and client library treats as storable by default, and this one's body is a credential. Whether Exoscale sends `Cache-Control: no-store` is not in the document — the document describes the body and says nothing about caching it.

**Thirty-four paths put the verb after a colon** — `/instance-pool/{id}:scale`, `/instance/{id}:create-snapshot`, `/block-storage/{id}:attach`. Google's AIP convention. The colon separates the identifier from the action *inside one path segment*, so a client that URL-encodes the whole segment encodes the colon too and asks for a resource whose id ends in `%3Ascale`. [cilium](../cilium) has colons in the same position meaning the opposite thing — there they separate a namespace from an identifier.

**There is no single base URL.** The server is `https://api-{zone}.exoscale.com/v2` with eight zones in the enum and no default host, so the zone is in the *hostname*: a client must know it before it can make any request, and resources in different zones are on different origins with no cross-zone listing.

## Modelling limits

- **Nothing here is verified against a live API.** Every endpoint needs a signed request against a real account.
- **Security groups, listed and fetched.** 261 paths is a cloud: instances, pools, load balancers, block storage, private networks, DNS, SKS clusters and the whole DBaaS surface each want their own evidence.
- **The credential is checked as a bearer token, which Exoscale does not use.** Reproducing EXO2-HMAC-SHA256 would mean implementing a signature scheme to demonstrate that the document does not describe it; what the Recipe needs is that a credential is required at all.
- **`rules` is served nowhere**, which is the finding — `writeOnly` means never returned, so this Recipe has no route that shows a rule and the rule shape is recorded in the Recipe header instead.
