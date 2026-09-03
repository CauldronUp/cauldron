# Wikimedia

Emulates the Wikimedia API (v1), for local development and tests.

**6 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The error message is inside a map keyed
by language.** Every failure carries `messageTranslations`, and none carries
`message` -- so `err.message` is undefined on all of them and the sentence a
client wanted is one level down under a language code it has to choose.
Choosing wrongly is silent: ask for `fr` on a response carrying only `en` and
you get undefined again, from a field that exists.

**And three failures share no key set at all.** A missing page sends four keys;
a missing path sends two, with no message in any language -- not an empty map,
no map; and an empty search sends nine, two of which name the same failure and
disagree, `error` saying `parameter-validation-failed` while `errorKey` says
`missingparam`. The record itself links to a different host running the same
API: served from `api.wikimedia.org/core/v1`, its `html_url` points at
`en.wikipedia.org/w/rest.php/v1`. Its two identifiers are four orders of
magnitude apart, `id` 9228 for the page and `latest.id` 1372042770 for the
revision, with nothing in either name saying which is which. And the licence URL
ends in a language, `.../by-sa/4.0/deed.en`.

## Sources

- Documentation: https://api.wikimedia.org/wiki/Core_REST_API
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wikimedia     # run it
cauldron verify wikimedia -v # check every claim
```
