# upcloud

Emulates the UpCloud account API for local development and tests.

**6 conformance cases, 3 checked against the live API on 2026-09-13.**

Read from UpCloud's own API documentation, and struck live against `api.upcloud.com` on 2026-09-13 with no credential, with a deliberately invalid one, and on a path that does not exist.

## What this Recipe found

**Two media types and two formatters for one status.**

```
(no header)     401 application/problem+json
                {   "error" : {      "error_code" : "AUTHENTICATION_REQUIRED",
                     "error_message" : "The use of this API requires authentication.
                     Create your account at https://www.upcloud.com/."   }}

Basic (wrong)   401 application/json
                {"error":{"error_code":"AUTHENTICATION_FAILED","error_message":
                 "authentication failed using the given username and password"}}
```

The same status, from the same endpoint, in two different media types — and the first is pretty-printed with spaces on both sides of its colons while the second is compact. Two serialisers, and the label changes with them.

**And neither is a problem document.** RFC 9457 defines `application/problem+json` as an object with `type`, `title`, `status` and `detail` at the top level. This one has none of those: it nests `error_code` and `error_message` under `error`.

A client that sees the media type and reaches for `body.title` finds nothing. The media type is the only part of the response that follows the standard it names.

**The unauthenticated message is marketing.** "The use of this API requires authentication. Create your account at https://www.upcloud.com/." — a sign-up link, in an error body, to a caller who may well have an account and have forgotten a header.

**A path that does not exist answers the credential failure.** `/1.3/cauldron-nope` with no credential answers `AUTHENTICATION_REQUIRED`, in problem+json, exactly as a real path does. The credential is checked before the route is resolved, so a mistyped path is indistinguishable from a missing header.

**Money is a float with four decimal places and no currency.** `"credits": 9972.2324` — and nothing in the response says what a credit is worth or in what currency. The documentation's own explanation is that credits "are used to pay for cloud resources" and "can be purchased from the UpCloud website".

**Fifteen limits, and one of them names its unit.** `resource_limits` holds counts and sizes side by side — `cores: 200`, `memory: 1048576`, `storage_hdd: 10240`, `public_ipv4: 100` — with no unit on any of them except `ntp_excess_gib`, which is about network time traffic and is the one place a unit appears in a key. A client reading `memory` has to know from somewhere else whether it counts megabytes or gigabytes.

## Sources

- Live: `api.upcloud.com`, struck 2026-09-13.
- [Accounts](https://developers.upcloud.com/1.3/3-accounts/) — the record shape and the explanation of credits.

## Modelling limits

- **One route.** The account. Servers, storages, networks, floating IPs, load balancers, managed databases, Kubernetes and object storage each want their own evidence.
- **The whitespace difference is recorded, not served.** Cauldron serialises every body the same way, so the two 401s here differ in their media type and not in their formatting. Live they differ in both.
- **The credential is a control-panel password.** UpCloud's Basic auth takes the account's own username and password rather than a token minted for the API, which is why the username is also a field on the record this endpoint returns.
- **No `spec:`.** UpCloud documents this API as prose with worked request and response transcripts. There is no OpenAPI document at any address this Recipe could find.
- **Nothing is mapped in detection.** UpCloud is driven by its own CLI and by a Terraform provider; `upcloud-go-api` is a Go module a project holds only when it is building tooling, so a dependency list mostly shows neither. Checked 2026-09-13.
