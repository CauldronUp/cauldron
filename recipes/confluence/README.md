# Confluence

Emulates the Confluence API (v2), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Asking Confluence for a page does not get you the page's text -- the body-format query parameter has no default, so a request that doesn't name one gets a full 200 with title, id, status, version, links and everything else a page has except what it actually says. And the format you do name becomes the field's key: ask for storage and the text sits at body.storage.value, ask for atlas_doc_format and it sits at body.atlas_doc_format.value with body.storage simply gone -- code written against either format breaks on the other with a property access on undefined, not an API error.

There's also no partial update: PUT /pages/{id} requires id, status, title, body and version all five, so changing one word means resending the whole page, and a client that read the page without specifying body-format has nothing to resend. Two ways to lose a draft entirely are documented in the field descriptions themselves rather than surfaced as errors -- a significantly diverged update can silently override what was in a draft, and changing status from draft to current deletes an existing draft in favor of it. Confluence's default listing also includes archived pages, the inverse of the shape most content APIs have.

The version-conflict write protocol here is the identical shape commercetools has -- a contract test suites structurally can't exercise, because the wrong or missing version passes every local check and only overwrites silently in production.

## Sources

- Documentation: https://developer.atlassian.com/cloud/confluence/rest/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve confluence     # run it
cauldron verify confluence -v # check every claim
```
