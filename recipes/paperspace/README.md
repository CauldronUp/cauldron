# paperspace

Emulates the Paperspace machine listing for local development and tests.

**13 conformance cases, 4 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Paperspace serves without a credential at [`api.paperspace.com/v1/openapi.json`](https://api.paperspace.com/v1/openapi.json), and struck live on 2026-09-13 with no header, with a wrong bearer, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**An API key is told it is not logged in.**

```json
{"message":"You must be logged in to access this method.","code":"UNAUTHORIZED"}
```

"Logged in", to a caller holding a bearer token and no session; and "this method", for an HTTP endpoint. The same sentence answers a missing credential and a wrong one, so nothing tells them apart.

**A wrong method is the same answer as a wrong path.** `PUT /v1/machines` and `GET /v1/cauldron-nope` both answer `{"message":"Not found","code":"NOT_FOUND"}` — the route that exists reporting that it does not.

**Thirty-one required fields, and twelve of them are nullable.** A machine `required`s `dtDeleted`, `reservation`, `restorePointSnapshotId`, `autoShutdownTimeout`, `autoShutdownForce`, `autoSnapshotFrequency`, `autoSnapshotSaveCount`, `privateIp`, `networkId`, `publicIp`, `accelerators` and `region` — every one declared nullable. The key is guaranteed and the value is not, and a generated client's non-optional field is null.

**`dtDeleted` is one of them.** Every machine record carries the time it was deleted, required, on a listing of machines that exist.

**A frequency whose only value is an event.** `restorePointFrequency` is `enum: ["shutdown"]`, nullable — two states, `null` and the word "shutdown", and neither of them is a frequency. `autoSnapshotFrequency` beside it is `["hourly", "daily", "weekly", "monthly"]`, which is what the other field's name promised.

**"No public address" is expressible twice.** `publicIp` is `string, nullable` and `publicIpType` is `["static", "dynamic", "none"]`, so a machine with no public address has a null in one field and the string `"none"` in the other, and nothing says what it means for them to disagree.

**Two readiness states, spelled as one word each.** `state` is `["off", "starting", "stopping", "restarting", "serviceready", "ready", "upgrading", "provisioning"]`. Six are ordinary participles; `serviceready` is two words run together, and it is not `ready`.

**Three timestamps carry a type prefix and twenty-eight fields do not.** `dtCreated`, `dtModified`, `dtDeleted` — Hungarian notation on three fields of a record whose other names are plain.

**The bill is two bare numbers.** `usageRate` and `storageRate` are `type: number` with no currency field and no unit field anywhere on the record or in the envelope. And `cpus` is `type: number` while `ram` is `type: integer`, so one of the two is divisible.

**No status but 200 is named anywhere.** The listing declares a 200 and a `default`, and `default` is a shared response called "error" carrying `message`, `code` and an untyped `details`. The 401 every anonymous caller meets and the 404 every typo meets are both the catch-all — which is why `cauldron drift` reports all three of this Recipe's failures as unbacked.

**And the field that says where to go next is the optional one.** The envelope is `{hasMore, nextPage, items}` with `additionalProperties: false` and `required: [hasMore, items]`, so `nextPage` comes and goes while the boolean that duplicates its presence is guaranteed.

## Sources

- [`api.paperspace.com/v1/openapi.json`](https://api.paperspace.com/v1/openapi.json) — served without a credential; recorded by `cauldron drift`, `spec_seen: 2026-09-13`.
- Live: `api.paperspace.com`, struck 2026-09-13 with no header, a wrong bearer, an unrouted path, and a wrong method.

## Modelling limits

- **One route of sixty-one.** The machine listing. Deployments, datasets, notebooks, models, container registries, custom templates, storage providers, secrets and the project surface are the rest.
- **The success fixture is document-derived.** Listing machines needs a real key; the records here are the operation's own `items` schema with values of the declared shapes, and every case reading them is marked documentation-only.
- **`region` is served as a string.** The schema declares it `anyOf`, and every branch is a region name.
- **Nothing is mapped in detection.** Paperspace is reached through the `gradient` CLI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
