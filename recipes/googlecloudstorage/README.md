# Google Cloud Storage

Emulates the Google Cloud Storage API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Folders come back in a different array from the files. There are no real directories in a bucket, only object names with slashes in them, and a listing with a `delimiter` set puts folder-like names in `prefixes` and only the top-level objects in `items`. A bucket with ten thousand files in folders and two at the top level answers with two items and a handful of prefixes -- and nothing in the response is wrong, it's just not the answer most code assumes it's getting when it reads `items` and concludes the bucket is nearly empty.

Paging counts both arrays together, so a short page is not necessarily the last page: `maxResults` caps items plus prefixes combined, and duplicate prefixes get dropped, so a page can come back short of what was asked for while more results remain. The common `while (items.length === pageSize)` loop stops early on exactly that page.

Numbers also arrive as strings -- `size`, `generation`, `metageneration` -- because a bucket can hold objects larger than a double can represent safely, so a client that adds up sizes with `+` concatenates them instead of summing them. And `id` is not the object's stable identity; it includes the generation number, so overwriting a file changes its id. `name` is the stable field, which is the reverse of what every other provider in this collection calls its identifier.

## Sources

- Documentation: https://cloud.google.com/storage/docs/json_api/v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googlecloudstorage     # run it
cauldron verify googlecloudstorage -v # check every claim
```
