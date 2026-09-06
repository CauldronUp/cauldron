# cilium

Emulates the Cilium agent API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Cilium's own Swagger 2.0 document, published in its own repository — `cilium/cilium`, `api/v1/openapi.yaml`, 33 paths, version v1beta1 — read on 2026-09-06. The agent API listens on a unix socket on each node and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**One path segment, four identifier namespaces, and colons inside it.** `/endpoint/{id}` takes, in the document's own words, "a string describing an endpoint with the format `[prefix:]id`":

```
cilium-local:3389595
cilium-global:cluster1:nodeX:452343
cni-attachment-id:22222:eth0
cep-name:default:foobar-net1
```

Three colons in the second one. A colon is a legal path character so nothing needs escaping — but a router that splits on `:`, a proxy that treats a colon as a port separator, or a client that builds the URL by joining parts will all disagree about where the identifier ends.

**And if you send no prefix, one is assumed.** "If no prefix is specified, a prefix of `cilium-local:` is assumed." So `/endpoint/3389595` silently means one namespace of the four, and a caller holding a CNI attachment ID that happens to be all digits addresses a different endpoint than they meant.

**The document's own example contradicts the prefix it illustrates:**

> - cep-name: cep name for this container if K8s is enabled, e.g. `pod-name:default:foobar-net1`

The bullet names the prefix `cep-name` and its example uses `pod-name`. One of the two is wrong and the document does not say which.

**Not every endpoint answers to every prefix.** "Not all endpoints will be addressable by all endpoint ID prefixes with the exception of the local Cilium UUID." So a 404 means either "no such endpoint" or "that endpoint exists and is not reachable by the namespace you used", and nothing distinguishes them.

**The entire error schema is a bare JSON string:**

```yaml
Error:
  type: string
```

A 400 answers `"Invalid endpoint ID format for specified type"` — with the quotes, as a JSON string literal, not an object. `body.message` is undefined; `body` **is** the message. A client calling `.json()` succeeds and gets a string, so nothing throws and code reading `body.error` reports the reason as "undefined". And the 404 and 429 declare no schema at all, so the one failure whose shape is described is described as a string.

**There is no authentication, and here that is correct.** The document declares no `securityDefinitions` and no `security` — exactly as [argocd](../argocd)'s does. But this API listens on a root-owned unix socket, so the absence is the design rather than an omission. Two documents identical on this point, meaning opposite things: which is why a scan for "APIs with no declared auth" needs a human at the end of it.

**A DELETE on the collection.** `DELETE /endpoint` takes a batch request and removes many endpoints — a destructive verb at a path with no identifier in it, the same shape [portainer](../portainer) has and reached independently.

**Nothing in a response spells an endpoint the way its own URL does.** `Endpoint.id` is a bare integer; the prefixed form is not a field on the record at all. A client that stores what it read cannot reconstruct the URL it read it from without knowing which namespace it used.

**An endpoint always has an identity, and sometimes it is the one meaning none.** Policy is written against the security identity, not the endpoint — so the number that decides what an endpoint may do is not the number that names it — and identity 5 is the reserved "not resolved yet", a real number rather than a null.

**Health is four axes and one of them is a boolean.** `overallHealth`, `bpf` and `policy` are strings; `connected` is a boolean, in the same object. "Is this endpoint healthy" has four answers of two types, and the boolean is the one a client reaches for.

## Modelling limits

- **Nothing here is verified against a live API.** The agent API is a unix socket on each node.
- **Endpoints, listed and fetched.** 33 paths is a datapath agent: identities, policy, services, IPAM, prefilter, BGP, the map surface and the debug surface each want their own evidence.
- **One prefix is served.** Routing all four namespaces to one store would mean inventing which endpoints are addressable by which prefix — the thing the document explicitly says varies. A caller using another prefix gets this emulator's unrouted-path 404 rather than Cilium's "Endpoint not found", which is the one place a wrong prefix is *more* legible here than upstream.
- **`Endpoint.status` is a deep tree** — eleven keys, several of them their own objects. Four are modelled.
