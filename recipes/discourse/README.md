# Discourse

Emulates the Discourse API (v3), for local development and tests.

**8 conformance cases, 7 checked against the live API on 2026-08-23.**

## What this Recipe found

It is the first provider added to this list since the
list was written. Any Discourse forum answers `/latest.json` and `/t/{id}.json`
without a key, so meta.discourse.org settles all five: the topics are two
levels down and there is no `topics` or `per_page` at the top level at all;
`topic_list.per_page` is thirty, the number this Recipe declares; a topic
carries `last_poster_username` as a string and neither a `last_poster` nor a
`user` object, while the names and avatars sit in a `users` array beside the
topics; `title` and `fancy_title` are the same words with the apostrophe as
`&rsquo;`, on a topic nobody wrote for the purpose; and `/t/{id}.json` answers
sixty-four keys of which `bumped_at` and `last_poster_username` are not two.

That last one is the case a description could not have settled. A description
lists what a response may carry. It does not say what a response leaves out,
and the whole point of that case is a field the listing has and the topic does
not.

## Sources

- Documentation: https://docs.discourse.org/openapi.json
- Machine-readable description: https://docs.discourse.org/openapi.json, last checked 2026-08-31
  `cauldron drift discourse` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve discourse     # run it
cauldron verify discourse -v # check every claim
```
