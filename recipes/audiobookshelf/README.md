# audiobookshelf

Emulates the Audiobookshelf API, for local development and tests.

**10 conformance cases, none checked against a live API.**

Written against Audiobookshelf's own OpenAPI 3.0 document, published in its own repository — `advplyr/audiobookshelf`, `docs/openapi.json`, 31 paths — read on 2026-09-06. Audiobookshelf is self-hosted and there is no public instance, so every claim here is the provider's own document.

## What this Recipe found

**A listing returns the server's filesystem layout.** `libraryItemBase` — the shape a library listing answers with — declares:

| field | what it is |
|---|---|
| `path` | the absolute path on the server |
| `relPath` | the path relative to the library folder |
| `ino` | the inode number, as a string |
| `mtimeMs` | modified time |
| `ctimeMs` | inode-change time |
| `birthtimeMs` | creation time |
| `isFile` | — |

That is `stat(2)`, on every record, in an API about audiobooks. A client rendering a library shows absolute server paths; a logged response body contains the machine's directory structure; and `ino` is the filesystem inode, exposed as a **string** because an inode number does not survive JavaScript's 53-bit floats.

Three filesystem timestamps sit beside two application ones (`addedAt`, `updatedAt`), so a record carries five different notions of *when* and only two are about the library. `birthtimeMs: 0` is a real POSIX instant, not an absence — not every filesystem records a birth time.

**Paging is zero-indexed, and inert until you ask for a limit.** The document's own words:

> `page` — "The page number (zero indexed) to return. **If no limit is specified, then page will have no effect.**"
> `limit` — "The number of items to return." `default: 0`

So `page=1` is the **second** page, and a client that starts counting at one — which is what almost every other provider in this collection wants — silently skips the first page and never sees an error.

And `limit` defaults to 0, meaning unlimited. A client that sends only `page` gets the whole library every time, and the parameter it did send does nothing. Both failures are quiet: the response is a valid page of valid records.

**Booleans are integers.** `desc` is `type: integer, default: 0`; `minified` is `type: integer, minimum: 0`. So `desc=true` is not what the document asks for, and a client sending a JSON boolean through a query-string serialiser sends the string `"true"` to a parameter declared as a number.

**And one of those integers changes the response schema.** `minified` switches the answer from `libraryItemBase` to `libraryItemMinified` — one endpoint, two record shapes, selected by a parameter that reads like a display preference.

**Every record carries the identifier it used to have.** `id` is a uuid, "the ID of library items after 2.3.0". Beside it, `oldLibraryItemId` is nullable and holds "the ID of library items on server version 2.2.23 and before", shaped `li_` plus eighteen characters. A record from before the upgrade carries two identifiers; one created after carries one and a null — and a client cannot tell "no old id" from "this field is not populated".

**The constraints are in `format`, where nothing enforces them:**

```
oldLibraryItemId   format: "li_[a-z0-9]{18}"
ino                format: "[0-9]*"
```

Those are regular expressions in OpenAPI's `format` field. `format` is an open string and tooling ignores values it does not recognise, so neither constraint reaches a validator — where `pattern`, which exists for exactly this, would have been checked.

**Missing and invalid are separate states.** Present on disk and unreadable is a third condition beside *there* and *gone*, and it takes two booleans to express — so a client checking only `isMissing` shows a broken item as fine.

**Adjacent measurements, different number types.** `durationSec` is a float (33854.905); `size` is an integer (268824228).

## Modelling limits

- **Nothing here is verified against a live API.** Audiobookshelf is self-hosted and there is no public instance.
- **Library items, listed.** 31 paths is a media server: authors, series, collections, playlists, sessions, podcasts, the email surface and the whole player API each want their own evidence.
- **`minified` is declared and not served.** Serving two schemas from one route would mean deciding what the minified shape drops, and the document's `libraryItemMinified` declares no properties of its own.
- **The filesystem fields are served, because they are the finding.** The paths in the fixture are plainly synthetic.
