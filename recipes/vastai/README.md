# vastai

Emulates the Vast.ai instance listing for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-13.**

Read from Vast.ai's own published CLI ([`vast-python`](https://github.com/vast-ai/vast-python)'s `vast.py`), and struck live against `cloud.vast.ai` on 2026-09-13 with no credential, with a wrong bearer, with the key in the query string, on a path that does not exist, and with a method the path does not take.

## What this Recipe found

**A wrong key is a 404.**

```
(no header)      403  {"success":false,"error":"auth_error","msg":"This action requires login."}
Bearer notreal   404  {"success":false,"error":"auth_error","msg":"Invalid user key"}
```

The body says `auth_error` and the status line says Not Found. A client branching on the status retries the path; one branching on the body fixes the key; the same response tells them different things. And the missing credential is Forbidden, which is the status for a caller who *has* been identified.

**The key is accepted in the query string.** `?api_key=notarealkey` is read and judged, so the credential can travel in the URL — into access logs, browser history and `Referer` headers — and nothing about the API stops it.

**The discriminator disappears exactly when the path is wrong.** An unrouted path answers `{"success":false,"msg":"Not found"}` — two keys where every other failure has three, with `error` simply absent.

**And the address without the trailing slash is a permanent redirect to another hostname.** `GET https://cloud.vast.ai/api/v0/instances` answers a 301 to `https://console.vast.ai/api/v0/instances/`, as an HTML page from the proxy. The host changes, the path gains a slash, and the body is not JSON.

**`duration` means two different things on two records.** In the vendor's own CLI it is divided by 86400 and labelled `Max_Days` on an offer, and divided by 3600 and labelled `age(hours)` on an instance. One field name, one unit on the wire, two meanings and two tables. And the instance carries a second time field that does say its unit — `uptime_mins`, beside the one without a name for them.

**`reliability2` is a fraction the same tool renders two ways.** The instance and offer tables multiply it by 100 and call it `R`; the volume table prints it raw to four decimals and calls it `reliab`. A field with a version number in its name, and no unit anywhere in the response.

**Three status fields.** `actual_status`, `intended_status` and `status_msg` — so a record says what it is, what it was meant to be, and a sentence about the difference, and the field a client reads first is the one called "actual".

**`cpu_ram` is megabytes**, which the CLI divides by 1000 rather than 1024 to show gigabytes; `dph_total` is dollars per hour abbreviated to eight characters, with no currency field anywhere; and `success: false` sits in the body of every failure beside a status line that already said so.

## Sources

- [`vast-ai/vast-python`](https://github.com/vast-ai/vast-python) — `vast.py`'s `instance_fields`, `displayable_fields` and `volume_fields` tables, with their own display transforms.
- Live: `cloud.vast.ai`, struck 2026-09-13 with no credential, a wrong bearer, a key in the query string, an unrouted path, and a wrong method.

## Modelling limits

- **One route.** The instance listing. Offers, machines, volumes, clusters, overlays, network disks, teams and the billing surface are the rest.
- **The 301 is not served.** The slashless address redirects to a different hostname with an HTML body from a proxy; this Recipe records that rather than emulating another host.
- **No `spec:`.** Vast.ai publishes a Python CLI rather than a machine-readable description, and no OpenAPI document is served at any address this Recipe could find — so there is nothing for `cauldron drift` to record.
- **The success fixture is CLI-derived.** Listing instances needs a real key; the records here are `instance_fields`' own names with values of the shapes the CLI's format strings imply, and every case reading them is marked documentation-only.
- **Nothing is mapped in detection.** Vast.ai is reached through `vastai` on PyPI or a plain HTTP call carrying a bearer token, and neither resolves to this host through a dependency file. Checked 2026-09-13.
