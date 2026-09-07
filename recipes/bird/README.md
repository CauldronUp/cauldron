# bird

Emulates the Bird API for local development and tests.

**11 conformance cases, 5 checked against the live API on 2026-09-07.**

Written against Bird's reference at `docs.bird.com` and struck live against `api.bird.com` on 2026-09-07 with no credential and then with a deliberately invalid one.

## What this Recipe found

**The code is the status phrase and the message is the status phrase in a sentence.**

```
401 {"code":"Unauthorized","message":"The request is not authorized"}
```

`code` is the reason phrase, capitalised, in the field a client branches on. `message` is the same fact written out. Two fields, one of them machine-readable in shape, and neither carrying anything the status line did not.

It is the passive voice again, as in [revai](../revai): "The request is not authorized" names no actor and no reason. **Nothing says whether a credential was sent, whether it was readable, or which header it belonged in.**

**One response to three mistakes** — a missing credential, a wrong credential and an unknown path, byte-identical. So route existence is not discoverable and the two credential failures cannot be told apart.

**The scheme is `AccessKey`.** `Authorization: AccessKey <key>` — which Bird inherited from MessageBird before the rename. So the scheme name is a product that no longer exists, no HTTP library has a helper for it, and the refusal does not mention it.

**The company was renamed and the wire format kept the old name.** MessageBird became Bird; the host is `api.bird.com`; the credential scheme still says MessageBird's word. The third provider in this collection whose name and wire format disagree, after [maxio](../maxio) at `chargify.com` and [terraformcloud](../terraformcloud) at `app.terraform.io` — and here it is the **credential** rather than the hostname carrying the old name, which is the harder one to notice.

**The request and the response spell the cursor differently** — `pageToken` going out and `nextPageToken` coming back, so a client cannot pass the field straight through.

**The region is per workspace rather than per organisation**, so where messages are processed varies within one account and a client cannot assume one jurisdiction.

## Modelling limits

- **One route.** Workspaces. Channels, contacts, conversations, messages, templates and the whole Flows surface each want their own evidence.
- **Nothing is mapped in detection.** The published clients are MessageBird's, for the old API, and do not call `api.bird.com`.
- **No `spec:`.** Bird publishes a rendered reference site.
