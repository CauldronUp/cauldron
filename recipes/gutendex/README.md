# Gutendex

Emulates the Gutendex API (v1), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**A filter it cannot parse is dropped, and you get
the whole catalogue.** `?ids=84` answers one book; `?ids=abc` answers 79,296 --
every book Project Gutenberg has -- at 200, with nothing anywhere saying the
filter was ignored. **And the invalid parameter is carried faithfully into the
next link**, `?ids=abc&page=2`, so a client that pages by following the link it
was given walks the entire catalogue still asking for "abc" every time. The
parameter is preserved exactly where it does nothing and dropped exactly where
it would matter.

`bookshelves` holds four entries prefixed `Category: ` and then four that are
not, in one array, with nothing marking where one convention ends. The keys of
`formats` are MIME types, so they carry slashes, plus signs and -- in
`"text/plain; charset=utf-8"` -- a semicolon and a space; not one can be written
after a dot. And the same title appears twice in one record capitalised two
ways: `title` is `"Frankenstein; or, the modern prometheus"` while the summary
beside it opens by quoting `"Frankenstein; Or, The Modern Prometheus"`, so the
field a client renders disagrees with the prose it renders underneath.

## Sources

- Documentation: https://gutendex.com
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gutendex     # run it
cauldron verify gutendex -v # check every claim
```
