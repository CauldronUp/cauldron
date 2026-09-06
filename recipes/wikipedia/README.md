# Wikipedia

Emulates the Wikipedia API (rest_v1), for local development and tests.

**10 conformance cases, 8 checked against the live API on 2026-08-28.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

A page that does not exist is reported as an
internal error. `GET /api/rest_v1/page/summary/Zzzznotarealpagexyz` answers
`{"status": 404, "type": "Internal error"}` and that is the whole body: nothing
names the page, and the `type` says the same thing for a title nobody created as
for a genuine failure inside the service, so a client branching on it cannot tell
a 404 it should expect from a 500 it should alert on. **And the title arrives
five times in four fields, two of them HTML** -- `title`, `displaytitle`, and
`titles` holding `canonical`, `normalized` and `display`, where `displaytitle`
and `titles.display` are the same markup twice and nothing says which a heading
should use.

A disambiguation page is a 200 whose `extract` is a flattened list:
`"Mercury most commonly refers to:Mercury (planet), the closest planet to the
Sun\nMercury (element), a chemical element"` -- no space after the colon,
newlines between the entries, neither a sentence nor a list. The spelling you
asked for is nowhere in the answer, so asking for `toronto` gives `Toronto` in
every field that could have recorded the correction. `originalimage.source` ends
`3840px-...jpg` while `originalimage.width` says 6632, so laying out from the
declared size and loading the URL gets a picture forty per cent smaller -- and
those image URLs carry `utm_source`, `utm_campaign` and `utm_content`
parameters. `content_urls` has `desktop` and `mobile` with the same four keys,
three of them byte-identical. `revision` is the string `"1370998502"` where
`pageid` is the number 64646. And the main namespace's `text` is the empty
string.

## Sources

- Documentation: https://en.wikipedia.org/api/rest_v1/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wikipedia     # run it
cauldron verify wikipedia -v # check every claim
```
