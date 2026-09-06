# WordPress

Emulates the WordPress API (wp/v2), for local development and tests.

**16 conformance cases, 6 checked against the live API on 2026-08-23.**

## What this Recipe found

It serves any public site's posts to anybody:
`title` is an object whose only key is `rendered` and `id` beside it is a JSON
number; `categories` and `tags` are arrays of term ids with no names anywhere
in the post; a missing post answers exactly three keys -- `code`, `message`,
`data` -- with the status nested at `data.status` and nothing called `status`,
`error` or `errors` at the top; and asking for `page=0` is refused with "page
must be greater than or equal to 1", which settles the counting-from-one half
outright rather than by inference.

Two of those six were left without a date until a site could answer them, and
what the notes asked for turned out to exist. `date` and `date_gmt` are both
there with no zone marker on either, which is the shape -- and on
wordpress.org/news they are identical, because that site runs UTC, so nothing
of the trap was visible. wptavern.com runs on another timezone: its post
185079 carries `2025-01-08T22:48:02` and `2025-01-09T03:48:02`, five hours
apart and on different days, so a client parsing the first as UTC files that
post under the wrong date. Likewise, wordpress.org/news publishes two pages
and neither has a parent, which cannot show a child naming one; ma.tt's page
2545 carries `parent: 2536`, and 2536 carries `parent: 0`, so the id is an
integer pointing at another page and zero is what a top-level page says
rather than the field being absent.

One is still without a date. Eight of twenty posts carry `featured_media: 0`,
so the zero is real and common, but the case claims a sticky post as well --
and wordpress.org/news, wptavern.com, ma.tt and techcrunch.com were each
asked for `?sticky=true` and every one answered with nothing. Half a case
watched is not a case watched.

## Sources

- Documentation: https://developer.wordpress.org/rest-api/reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve wordpress     # run it
cauldron verify wordpress -v # check every claim
```
