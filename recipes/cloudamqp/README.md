# cloudamqp

Emulates the CloudAMQP instance API for local development and tests.

**10 conformance cases, 4 checked against the live API on 2026-09-09.**

Written against CloudAMQP's published OpenAPI document at `docs.cloudamqp.com` and struck live against `api.cloudamqp.com` on 2026-09-09 with no credential and with a deliberately invalid one.

## What this Recipe found

**The challenge header is sent only to callers who already know the scheme.**

```
GET /api/nodes    (no credential)
401  Content-Type: application/json
(no WWW-Authenticate)
<empty body>

GET /api/nodes    Authorization: Basic bm90cmVhbDpub3RyZWFs
401  Content-Type: application/json
WWW-Authenticate: Basic realm="CloudAMQP API"
{"error":"Not authorized"}
```

That is backwards.

RFC 9110 §15.5.2 says a 401 **MUST** carry a `WWW-Authenticate` header, and the entire reason the field exists is to tell a caller that has *not* authenticated which scheme to use. CloudAMQP withholds it from exactly that caller, and sends it to the one who has already demonstrated they know — a request carrying a Basic header gets told to use Basic.

The body follows the same pattern: nothing at all for the caller who sent nothing, and a sentence for the caller who sent something wrong. And `Content-Type: application/json` is named over the empty bytes.

**CloudAMQP's other API does it correctly.** The customer API at `customer.cloudamqp.com`, probed the same day, answers `WWW-Authenticate: Basic realm="Customer API"` to a missing credential and a wrong one alike, with `{"error":"Not authorized"}` in both.

One vendor, two APIs, two behaviours on one header — and the half that gets it wrong is wrong in the least useful direction available.

Both are well formed, which is worth saying. [buffer](../buffer) single-quotes its auth-params, which RFC 9110 does not permit for a quoted-string, so a conforming parser rejects the whole challenge. CloudAMQP's is correct wherever it appears; the fault here is *when* it appears rather than how.

**The credential is Basic with an empty username.** The API key is the **password** and the username is blank — the mirror image of [close](../close) and [insightly](../insightly), which both put the key in the username and leave the password empty.

Three providers, one scheme, two opposite conventions, and nothing on the wire says which way round a given API wants it. Get it backwards and the key is simply never checked.

OpenAPI has no way to express "leave the username blank", so the document invents two vendor extensions to say it — `x-username-optional: true` and `x-default-username: ""` — beside a description that spells out the curl invocation by hand. A generator reading only the standard keywords produces a client that asks for a username nobody has.

**The document declares one response and it is the success.** `GET /nodes` lists only `200`, on an API that answers 401 to every request without a credential. That is the same silence [middesk](../middesk) publishes, now found on a second provider two days later, and it means a generated client has no failure type for the response every new integration meets first.

Nothing on a node is guaranteed either: the item schema has no `required` array at all.

**Two of everything on a node.**

- `disk_size` and `additional_disk_size` are both integers in gigabytes, and the total is neither of them — a client wanting the capacity has to add them.
- `hostname` and `hostname_internal` are both addresses for the same machine, distinguished only by the suffix.
- `running` and `configured` are both booleans that read as "is it up", and a node can be configured and not running — which is the distinction one boolean could not carry.

**There is no paging anywhere.** The operation declares no parameters at all, and the response is a bare array.

## Modelling limits

- **One route.** Listing nodes. Plugins, firewall rules, custom domains and certificates, the VPC surface and the whole account section each want their own evidence.
- **This is the instance API, not the customer API.** They are different hosts with different credentials: `api.cloudamqp.com` acts on one instance, `customer.cloudamqp.com` manages the instances themselves. The customer API's behaviour is recorded above as the contrast and is not served here.
- **Nothing is mapped in detection.** A project talking to CloudAMQP holds an AMQP client — `amqplib`, `php-amqplib`, `streadway/amqp` — pointed at a broker URL, not a client of this management API, so a dependency name says nothing about whether this surface is used.
