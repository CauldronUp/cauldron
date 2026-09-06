# coder

Emulates the Coder development-platform API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Written against Coder's own generated Swagger 2.0 document, published in its own repository — `coder/coder`, `coderd/apidoc/swagger.json`, 344 paths, "Coder API 2.0" — read on 2026-09-06. Coder is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**The search parameter is a language, and one of its keys is not a filter.** `GET /api/v2/workspaces` takes `q`, "Search query in the format `key:value`". Twelve keys are listed:

```
owner  template  name  status  has-agent  dormant  healthy  has-ai-task
last_used_after  last_used_before  has_external_agent  include_agent_metadata
```

The last one is documented, in the parameter's own description, as "expands each agent with the named metadata keys **rather than filtering**; repeat the key for multiple items".

So eleven of the twelve narrow the result and the twelfth changes its *shape*, sharing one parameter, one syntax and one description. A caller assembling `q` from user input can widen the response by accident, and a client typing the response has to know that one search key adds fields to it.

**And the keys do not agree about separators.** `has-agent` and `has-ai-task` use hyphens; `has_external_agent`, `last_used_after` and `include_agent_metadata` use underscores. Two conventions in one vocabulary, in the same string.

**A workspace's state is whether a timestamp is null.** From the document's own descriptions:

> `dormant_at` — "DormantAt being non-nil indicates a workspace that is dormant. A dormant workspace is no longer accessible must be activated. It is subject to deletion…"
> `deleting_at` — "DeletingAt indicates the time at which the workspace will be permanently deleted."

There is no `status` field carrying "dormant". The state *is* the nullness of a date — so `if (workspace.dormant_at)` is the check, and a client that renders the value as a date shows the user when they stopped being able to reach it.

`deleting_at` is sharper: a timestamp in the **future**, naming when the workspace will be destroyed. A listing of workspaces includes ones scheduled for deletion, and nothing but a non-null date in a field says so.

**`/api/experimental/` and `/api/v2/` are in the same document.** Six experimental paths — AI model prices, chat streaming, MCP server OAuth callbacks, user skills, a watch stream — sit beside 300-odd stable ones, with no marking in the operations themselves. A generated client gets methods for all of them and no signal about which are stable beyond a path segment a reader has to notice.

**Three credentials, and one names a product feature:**

| scheme | header |
|---|---|
| `Authorization` | `Authorization` |
| `CoderSessionToken` | `Coder-Session-Token` |
| `AIGatewayKey` | `X-AI-Governance-Gateway-Key` |

Three header credentials in one API, so which one authenticates depends on what the request is *for* rather than on who is calling.

**A record carries the owner twice and neither is stable.** `owner_id` and `owner_name`, and the same for the organisation — a rename changes one and not the other, so a client displaying the name and keying on the id shows a stale label until it refetches.

**A missing workspace and a forbidden one say the same thing:** "Resource not found or you do not have access to this resource". A deliberate choice not to leak existence, and it means a caller cannot tell a typo from a permission problem.

**One duration and four instants.** `ttl_ms` is milliseconds; `last_used_at`, `created_at`, `dormant_at` and `deleting_at` are RFC 3339 strings.

## Modelling limits

- **Nothing here is verified against a live API.** Coder is self-hosted and there is no public instance.
- **Workspaces, listed and fetched.** 344 paths is a development platform: templates, builds, agents, organisations, groups, audit, provisioners, the OAuth2 surface and the whole experimental block each want their own evidence.
- **`q` is served as a filter on one key.** Reproducing the mini-language would mean implementing a parser; what the finding needs is that one of its keys does not filter, which is recorded rather than served.
- **Ten of the 33 workspace fields are modelled.**
