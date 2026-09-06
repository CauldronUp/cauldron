# boundary

Emulates the HashiCorp Boundary controller API, for local development and tests.

**12 conformance cases, none checked against a live API.**

Written against Boundary's own generated Swagger 2.0 document, published in its own repository — `hashicorp/boundary`, `internal/gen/controller.swagger.json`, 95 paths, version 0.21.0 — read on 2026-09-06. Boundary is self-hosted and there is no public instance, so every claim here is the provider's own document.

## This one is mostly a counterexample, and that is why it is here

Most Recipes in this collection record a provider doing something that will surprise a caller. Boundary's document is the other case: it states its own rules, including the uncomfortable ones, in enough detail that a client can be written correctly from it. Recording that is worth as much as recording a trap — it shows the traps are choices rather than inevitabilities.

## What this Recipe found

**It documents an information leak and argues for it.** From `info.description`:

> "Boundary returns `404` if a resource cannot be found. Note that this happens *prior* to authentication/authorization checking in nearly all cases… As a result, an action against a resource that does not exist returns a `404` instead of a `401` or `403`. **While this could be considered an information leak**, since IDs are randomly generated and this only discloses whether an ID is valid, **it's tolerable** as it allows for far simpler and more robust client implementation."

So a caller *can* distinguish "no such resource" from "not yours", and the document says why that was chosen and what it costs.

[coder](../coder) makes the opposite choice — "Resource not found or you do not have access to this resource", one sentence for both — and does not explain it. Two access-control products, opposite decisions, and only one tells you which it made.

**A missing credential can succeed:**

> "A token that is invalid or missing, but where the anonymous user (`u_anon`) is able to successfully perform the action, will not return a `401` but instead will return the result of the action."

So 401 is not a property of the request. It depends on what `u_anon` has been granted, which is deployment configuration — and a client that treats "no token" as "will fail" is wrong on any instance where anonymous access is permitted, in the safe-looking direction.

**Two statuses for being rate limited, meaning different things:**

| status | cause |
|---|---|
| 429 | "any of the API rate limit quotas have been exhausted for the resource and action" |
| 503 | "unable to store a quota due to the API rate limit being exceeded" |

The first is the caller's quota; the second is the server's own accounting failing under load. Both carry `Retry-After`, and a client that retries on 429 and gives up on 503 gives up on the one most likely to clear.

**A listing tells you what to evict from your cache.**

> `removed_ids` — "A list of item IDs that have been removed since they were returned as part of a pagination. They should be dropped from any client cache."

That is rare and it is right: a paginated listing that only ever adds cannot tell a client about deletions, so every cache built on one drifts. Boundary sends the tombstones — the ordinary outcome everywhere else, here made avoidable.

**The same endpoint is a listing or a change feed.** `response_type` is `"delta"` or `"complete"` — delta means this is part of a pagination *or an update to one already finished*; complete means it is the last page. So the meaning of the response is a string, and the two meanings are "here is some data" and "here is what changed".

**The count is honestly named.** `est_item_count` — "an estimate at the total items available. This may change during pagination." Most providers here call that `total` and let a reader assume it is exact.

**Permissions travel on every record.** `authorized_actions` is "the available actions on this resource for this user", inline on each item — so a listing answers what you may do without a second call, which [vikunja](../vikunja) needs a header for and only on single items.

**Mutations carry a version.** `version` is "used in mutation requests, after the initial creation, to ensure this resource has not changed" — optimistic concurrency declared on the resource rather than left to an ETag, so it survives a client that strips headers.

**Identifiers carry their type as a prefix** — `acctpw_` for a password account, `ampw_` for its auth method — so a client that mixes two up is refused by shape rather than by lookup.

**Custom actions hang off a colon** — `/v1/accounts/{id}:set-password`. The same convention [exoscale](../exoscale) uses, and the opposite of [cilium](../cilium)'s, which puts a namespace *before* the colon rather than a verb after it.

## Modelling limits

- **Nothing here is verified against a live API.** Boundary is self-hosted and there is no public instance.
- **Accounts, listed and fetched.** 95 paths is an access proxy: targets, hosts, host sets, credential stores, sessions, scopes, roles, users, groups and the whole authentication surface each want their own evidence.
- **`removed_ids` is served as an empty array**, because a fixture has no deletions to report. What is reproduced is that the key is there to be read.
- **The `u_anon` behaviour is recorded and not served.** Reproducing it would mean modelling a grant system, and the finding is that 401 depends on one.
