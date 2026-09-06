# imgix

Emulates the imgix Management API for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-09-06.**

Written against what `api.imgix.com` actually answers, struck live on 2026-09-06 with no credential and then with a deliberately invalid one. imgix's Management API reference is a JavaScript-rendered page with nothing in its served HTML, which is why the attributes of a source are recorded as a gap rather than guessed.

## What this Recipe found

**The 401 tells you the server's build, its health, and whether it is in read-only or tombstone state.**

Struck live, with no credential at all:

```json
{"errors": [{"detail": "A valid API key is required and could not be determined.",
             "id": "87e18b65a7a846c1aa113cca1a7e986b",
             "meta": {}, "status": "401", "title": "Unauthorized"}],
 "jsonapi": {"version": "1.0"},
 "meta": {"authentication": {"authorized": false, "clientId": null,
                             "mode": "PUBLIC_APIKEY",
                             "modeTitle": "Public API Key",
                             "tag": null, "user": null},
          "server": {"commit": "production",
                     "status": {"healthy": true, "read_only": false, "tombstone": false},
                     "version": "3.213.15"}}}
```

Two things at once, pulling in opposite directions.

**The disclosure**: an anonymous request learns the running version — `3.213.15`, to the patch — and the deployment's own health flags. Anyone can ask, with no account.

**The usefulness**, which is rarer and worth naming as clearly: `read_only` and `tombstone` ride *every* response, so a client can tell a maintenance window from a bug **without a status page, on the same request that failed**. Almost nothing else in this collection does that. [boundary](../boundary) is the nearest and it tells you about its pagination rather than about itself.

And `meta.authentication` explains the refusal in structured fields rather than in the sentence: `authorized: false`, `mode: "PUBLIC_APIKEY"`, `clientId: null`. So the failure says **which scheme it tried** — the thing a caller holding two kinds of key most needs and almost never gets.

**A missing credential and a wrong one are the same response.** Struck live both ways: identical status, identical `detail`, identical `title`. Only the per-request id differs. So a client cannot tell "I never configured this" from "the key I configured is wrong" — the two failures that need entirely different fixes, and the pair every other Recipe in this collection has found a provider distinguishing.

**`status` is a string inside the error and a number outside it.** The entry carries `"status": "401"`, quoted, while the HTTP status beside it is `401`. That is JSON:API's specification rather than imgix's choice — and the consequence is still that one fact appears twice in one exchange with two types.

**The error carries its own correlation id, and so does a header, and they disagree.** `errors[0].id` is 32 hex characters, fresh on every call. The same response carries `x-cloud-trace-context` — also 32 hex characters, a different value. So a failure hands back **two unrelated correlation ids** and nothing says which one support wants.

**The media type is `application/vnd.api+json`.** Not `application/json`. A client checking the header before parsing does not parse, and a framework that decodes only `application/json` hands the handler an empty body and reports nothing — the same trap [unit](../unit)'s JSON:API responses set, and the reason this format has a `content_type` key at all.

## Detection

**Nothing is mapped, and that is the finding.**

`@imgix/js-core` is imgix's flagship package, it is official and current, and it names `imgix.net` **75 times** and `api.imgix.com` **zero**. It is a URL builder for the image-rendering service — a different host and a different product from the Management API this Recipe serves. `imgix/imgix-php` and `imgix-core-js` are the same tool in other languages.

So a project can depend on imgix, use imgix every day, and touch nothing this Recipe emulates. Mapping the URL builder here would report the wrong API for the overwhelming majority of imgix users. It is the [streamchat](../streamchat) shape again — the package named for the company belongs to a sibling product — and here it goes further: **the sibling product is the one nearly everybody means.**

## Modelling limits

- **A source's attributes are not modelled.** JSON:API requires `type` and `id` on every resource object and this Recipe serves both; the attribute names live in a reference page rendered by JavaScript, and inventing them would teach a client field names imgix may not use. The gap is smaller than the guess.
- **Paging is not served.** JSON:API pages with `page[number]` and `page[size]`, and this format has no way to send bracketed parameter names. Recorded rather than approximated.
- **The per-request id is a constant here.** A fake minting fresh hex on every call would be asserting a format it cannot check. What the cases pin is that the field exists, is 32 hex characters, and belongs to the entry rather than to the envelope.
- **One route.** Sources, and the failure envelope everything shares.
