# percy

Emulates the Percy build listing for local development and tests.

**8 conformance cases, 2 checked against the live API on 2026-09-13.**

Read from Percy's own published CLI client ([`percy/cli`](https://github.com/percy/cli)'s `packages/client/src/client.js`), and struck live against `percy.io` on 2026-09-13 with no credential and with a wrong token.

## What this Recipe found

**A missing credential is Forbidden and a wrong one is Unauthorized.**

```
(no header)            403  {"errors":[{"status":"forbidden"}]}
Token notarealtoken    401  {"errors":[{"status":"unauthorized","detail":"Invalid API token."}]}
```

The two statuses are the wrong way round — 401 is for a caller who has not authenticated — and the prose field appears on only one of them, so the failure with nothing to say is the one you get for sending nothing.

**And `status` is a word.** JSON:API defines the member as "the HTTP status code applicable to this problem, expressed as a string value". Percy sends `"forbidden"` and `"unauthorized"`, which are reason phrases, lowercased. There is no `title` and no `code` on either.

**Every multi-word attribute is kebab-case.** `target-branch`, `commit-sha`, `commit-author-email`, `pull-request-number`, `parallel-nonce`, `cli-start-time`, `skip-base-build` — so in the JavaScript client that reads them, every access is `attributes['commit-sha']` and never a property.

**A build carries two people's email addresses.** `commit-author-email` and `commit-committer-email` are sent on every build and stored on the record, so a visual-testing run holds the git identities of whoever wrote and whoever applied the commit.

**An absent state and a pending one are the same thing.** The vendor's own poller reads:

```js
let pending = !state || state === 'pending' || state === 'processing';
```

A build whose `state` is missing entirely is treated as still running, so the field's absence is a value.

**Another product's identifiers are on the record.** `testhub-build-uuid` and `testhub-build-run-id` travel with every build, and one of the `source` values a build can carry is `bstack_sdk_created` — the parent company's initials, in an enum, on a Percy build.

**And the client ships a comment saying one field does nothing.** From `client.js`:

> PER-9724: internal-only priority request… percy-api only honors this for eligible internal build types, so it is a no-op for customer builds.

A ticket number, an internal product name, and an admission that the attribute is inert for everyone reading it, in published code.

## Sources

- [`percy/cli`](https://github.com/percy/cli) — `packages/client/src/client.js`: the build attributes it sends, and the state values it polls on.
- Live: `percy.io`, struck 2026-09-13 with no credential and with a wrong token.

## Modelling limits

- **One route.** The build listing. Snapshots, resources, comparisons, projects, finalize, and the wait-for-build polling loop are the rest.
- **The record is the write-side attribute set.** Percy publishes a client rather than a description of its responses, so the fixture carries the attributes that client sends on a build plus `state`, the one attribute it reads back. Fields the API may add on read are not invented here.
- **No `spec:`.** No OpenAPI document is served at any address this Recipe could find, so there is nothing for `cauldron drift` to record.
- **Nothing is mapped in detection.** Percy is reached through `@percy/cli` on npm or a plain HTTP call carrying `Authorization: Token …`, and neither resolves to this host through a dependency file. Checked 2026-09-13.
