# civo

Emulates the Civo instances API for local development and tests.

**10 conformance cases, 3 checked against the live API on 2026-09-13.**

Struck live against `api.civo.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist. Record shapes come from [civogo](https://github.com/civo/civogo), Civo's own Go client, which is open.

## What this Recipe found

**Neither credential failure is a 401.**

```
(no header)      403 {"code":"authentication_access_denied","reason":"Access denied, you don't have permission on this resource"}
Bearer notreal   404 {"code":"database_account_not_found","reason":"Failed to find the account within the internal database"}
/cauldron-nope   404 404 page not found
```

A caller who has said nothing is told they lack permission. A caller whose token was read and rejected is told the **account** was not found — with a 404, a code naming a database, and a sentence naming it again: "the internal database".

So the rejection of a credential is reported as a missing row in storage the caller has never heard of, and 404 is the status for it. A client cannot tell any of this apart from a genuinely missing resource: the wrong-credential case and an unrouted path both answer 404, in two different shapes, one of them not JSON at all.

**The field is `reason`, not `message`.** Two keys, `code` and `reason`, and neither is the name almost every other API in this collection uses.

**An instance carries two passwords.** `initial_password` and `rescue_password` are fields on the record a listing returns, beside `initial_user` and an `ssh_key`.

**And a second credential.** `civostatsd_token` sits on the same record — the metrics agent's key, on the object that describes the machine.

**The cloud underneath shows through.** `openstack_server_id` and `openstack_image_id` are the identifiers of the same machine in the system Civo is built on, beside Civo's own `id` and `template_id`. Four identifiers for two things.

**Five fields are the address.** `public_ip`, `private_ip`, `ipv6`, `pseudo_ip` and `reserved_ip`, with `reserved_ip_id` and `reserved_ip_name` beside them, and nothing on the record saying which one a caller should connect to.

**The envelope counts pages and never records.** `{page, per_page, pages, items}` — how many pages there are at the size you asked for, and no total anywhere. A client that wants to know how many instances exist has to walk to the end.

**And the client omits every zero value it sends.** Every field on civogo's `Instance` is tagged `omitempty`, including `volume_backed bool` and the four integer counts. Go's encoder drops those when they are false or zero, so a request meaning `volume_backed: false` or `gpu_count: 0` leaves the field out entirely and the server is told nothing rather than told the value. The boolean can be sent as true and cannot be sent as false.

## Sources

- Live: `api.civo.com`, struck 2026-09-13.
- [`instance.go`](https://github.com/civo/civogo/blob/master/instance.go) — `Instance` and `PaginatedInstanceList`, in Civo's own Go client.

## Modelling limits

- **One route.** Listing instances. Kubernetes clusters, volumes, networks, firewalls, load balancers, DNS, object stores and the rest each want their own evidence.
- **Both passwords and the metrics token are fixture values.** They are named and shaped like the real fields and are not credentials for anything; the finding is that the record has the fields at all.
- **The `omitempty` behaviour is the client's, not the server's.** Go's `omitempty` affects encoding, so it governs what civogo *sends*; nothing here claims the API drops zero values on the way out.
- **No `spec:`.** Civo documents this API as prose. There is no OpenAPI document at any address this Recipe could find; the record shape comes from the Go client, which is the artefact a description would have produced.
- **Nothing is mapped in detection.** Civo is driven by its own CLI and by a Terraform provider, and `civogo` is a Go module a project would hold only if it were building tooling on top — a dependency list mostly shows the CLI, which is not a dependency. Checked 2026-09-13.
