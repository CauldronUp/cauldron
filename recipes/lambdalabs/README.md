# lambdalabs

Emulates the Lambda Cloud instances API for local development and tests.

**9 conformance cases, 2 checked against the live API on 2026-09-13.**

Read from the OpenAPI document Lambda serves without a credential at [`cloud.lambda.ai/api/v1/openapi.json`](https://cloud.lambda.ai/api/v1/openapi.json), and struck live on 2026-09-13 with no credential and with a deliberately invalid one.

## What this Recipe found

**The code cannot tell the two credential failures apart and the message can.**

```
(no header)      401 {"error": {"code":"global/invalid-api-key",
                                "message":"No API key was provided.", …}}
Bearer notreal   401 {"error": {"code":"global/invalid-api-key",
                                "message":"API key was invalid, expired, or deleted.", …}}
```

One code, two situations. A client branching on `error.code` — the field put there to be branched on — cannot distinguish a header it forgot to send from a key that has been revoked, and has to match on English to do it.

**The code is namespaced with a slash.** `global/invalid-api-key`: a scope and a name, joined by a character that is also a path separator, in a string a client will put in a switch.

**And the message names three causes at once.** "API key was invalid, expired, **or** deleted" — three states with three different remedies, in one sentence, under one code.

**A secret is a field on the listing.** `jupyter_token` is on the instance record, described as "the secret token used to log into the JupyterLab server hosted on the instance", with `jupyter_url` beside it carrying the same token in a query string. Listing your instances hands you a working login for every one of them.

**The same resource is at two paths.** The document declares `/api/v1/file-systems` and `/api/v1/filesystems` — one hyphenated and one not — while the fields on an instance spell it `file_system_names` and `file_system_mounts`. Three spellings of one word across one API.

**And there are two production servers, with two brands.** `cloud.lambda.ai` and `cloud.lambdalabs.com`, both declared and both answering — a rename that reached the new domain without retiring the old one.

**Twenty-one error types, one of which is about a support ticket.** The document defines a separate schema per failure — `ApiErrorQuotaExceeded`, `ApiErrorInsufficientCapacity`, `ApiErrorFileSystemInWrongRegion`, `ApiErrorInvalidBillingAddress` — and among them `ApiErrorTicketCreationFailed` and `ApiErrorTicketNotFound`. Raising a support ticket is part of the API surface and can fail like anything else.

**`name` may be the empty string.** `minLength: 0, maxLength: 64` on the field described as "the user-provided name of the instance", so a named thing may be named nothing.

## Sources

- Live: `cloud.lambda.ai`, struck 2026-09-13.
- [`openapi.json`](https://cloud.lambda.ai/api/v1/openapi.json) — served without a credential, 24 paths.

## Modelling limits

- **One route.** Listing instances. Launching, restarting and terminating them, filesystems, firewall rulesets, images, SSH keys, audit events and the ticket surface each want their own evidence.
- **The tokens are fixture values.** They are shaped like Lambda's and are credentials for nothing; the finding is that the field is on a listing at all.
- **One of the two hosts.** `cloud.lambdalabs.com` is declared beside `cloud.lambda.ai` and differs only in name.
- **Nothing is mapped in detection.** Lambda Cloud is reached through its own CLI or a generic HTTP client holding a key; the published SDKs are generated and a project holding one is pointed at whichever account the key belongs to. Checked 2026-09-13.
