# ionos

Emulates the IONOS Cloud datacentre listing for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from IONOS's own generated Go SDK ([`ionos-cloud/sdk-go`](https://github.com/ionos-cloud/sdk-go)), and struck live against `api.ionos.com` on 2026-09-13 with no header, with a wrong bearer, with a Basic credential, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**One number is quoted and the other is not.**

```json
{"httpStatus":401,"messages":[{"errorCode":"315","message":"Unauthorized"}]}
```

`httpStatus` is a JSON number and `errorCode` beside it is a JSON string, and both are integers. `messages` is a plural array carrying one object.

**That is the answer to every request.** A missing credential, a wrong bearer, a Basic header, an unrouted path and a wrong method all produce it, because the credential is checked before the path is. Nothing distinguishes a typo from an anonymous request.

**`state` is documented as a prose list, typed `string`, and one of its members is listed twice.** The description runs through sixteen values, and `FAILED_HIBERNATING` appears at both ends of it with the same explanatory clause each time. There is no enum — the list of legal values is a sentence.

**And twelve of those states are about Kubernetes.** On a *datacentre's* metadata type, most of the documented states are annotated "relevant for Kubernetes cluster/nodepool", so the state field of every resource is the union of every resource's states and almost none of them can apply here.

**The type enum names twenty-six values and accepts thirty-five.** The SDK declares constants for `datacenter`, `server`, `collection` and the rest, and its own validator then checks against a longer list containing `firewall-rule`, `flow-log`, `request-status`, `s3key`, `k8s`, `forwarding-rule`, `natgateway-rule`, `target-group` and `security-group`. Every one of the nine extra values has a hyphen in it, and a Go constant cannot — so the naming convention quietly dropped a quarter of the enum.

**A creation timestamp is documented as "The last time the resource was created."** On `createdDate`, in the metadata of every object.

**Four fields for two facts.** `createdBy` and `createdByUserId`, `lastModifiedBy` and `lastModifiedByUserId` — the name and the identifier of each, side by side.

**The offset and the limit are floating point.** `Offset` and `Limit` are `float32` in the listing envelope, so the position in a collection is a real number.

**And the collection is itself a resource.** The envelope carries `id`, `type` and `href` of its own, beside `items`, so a page of datacentres has an identity and an address the same way a datacentre does.

**`id` on a record is `[optional]`**, so the identifier may be absent; the data lives under `properties` while the identity sits beside it and the child collections sit under `entities`, all in one object; `etag` is a record field described by a link to RFC 2616 — obsoleted by RFC 7232 in 2014 — over plain HTTP; and the `Type` model's own documentation page has an empty properties table.

## Sources

- [`ionos-cloud/sdk-go`](https://github.com/ionos-cloud/sdk-go) — `docs/models/Datacenter.md`, `DatacenterProperties.md`, `DatacenterElementMetadata.md`, `Datacenters.md`, `PaginationLinks.md`, and `model_type.go`.
- Live: `api.ionos.com`, struck 2026-09-13 with no header, a wrong bearer, a Basic credential, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The datacentre listing. Servers, volumes, NICs, LANs, load balancers, NAT gateways, Kubernetes, backup units, contracts, users and groups are the rest.
- **`entities` carries one child collection.** Live a datacentre's `entities` holds servers, volumes, LANs and more; the fixture carries the servers collection to show the shape.
- **No `spec:`.** `api.ionos.com/cloudapi/v6/swagger.json` sits behind the credential check, and no machine-readable description is served at any address this Recipe could reach anonymously — so there is nothing for `cauldron drift` to record.
- **The success fixture is SDK-derived.** Listing datacentres needs a real token; the records here are the SDK's own model fields with values of the declared shapes, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** IONOS Cloud is reached through `ionos-cloud/sdk-go`, `ionoscloud` on PyPI, or a plain HTTP call carrying a bearer token, and none of those resolve to this host through a dependency file. Checked 2026-09-13.
