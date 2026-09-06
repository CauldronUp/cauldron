# Spaceflight News

Emulates the Spaceflight News API (spaceflightnews), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-08-30.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

News's, where **one record links the same site twice, once
securely and once not.** `url` is `https://spaceflightnow.com/...` and
`image_url` is `http://spaceflightnow.com/...`, on the same host, in the same
object -- so a page served over TLS that renders the picture gets a
mixed-content block and shows nothing, while the link beside it works.

**And the API's own 404 page carries a third-party analytics beacon**, a script
tag pointing at `static.cloudflareinsights.com/beacon.min.js` inside an error
response. Two timestamps in one record carry two precisions: `published_at` to
the second and `updated_at` to the microsecond. A `limit` that is not a number
is silently replaced with the default, so a client with a broken parameter never
learns. Absence is spelled two ways four fields apart, `socials: null` and
`events: []`. And the two failures do not agree on a format at all: an article
that does not exist answers JSON with the model name capitalised mid-sentence,
and a path that does not exist answers HTML.

## Sources

- Documentation: https://api.spaceflightnewsapi.net/v4/documentation
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve spaceflightnews     # run it
cauldron verify spaceflightnews -v # check every claim
```
