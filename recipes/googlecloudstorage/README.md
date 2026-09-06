# Google Cloud Storage

Emulates the Google Cloud Storage API (v1), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real project; the refusal cases were struck live, unauthenticated, against storage.googleapis.com.

## What this Recipe found

Folders come back in a different array from the files. There are no real directories in a bucket, only object names with slashes in them, and a listing with a `delimiter` set puts folder-like names in `prefixes` and only the top-level objects in `items`. A bucket with ten thousand files in folders and two at the top level answers with two items and a handful of prefixes -- and nothing in the response is wrong, it's just not the answer most code assumes it's getting when it reads `items` and concludes the bucket is nearly empty.

Paging counts both arrays together, so a short page is not necessarily the last page: `maxResults` caps items plus prefixes combined, and duplicate prefixes get dropped, so a page can come back short of what was asked for while more results remain. The common `while (items.length === pageSize)` loop stops early on exactly that page.

Numbers also arrive as strings -- `size`, `generation`, `metageneration` -- because a bucket can hold objects larger than a double can represent safely, so a client that adds up sizes with `+` concatenates them instead of summing them. And `id` is not the object's stable identity; it includes the generation number, so overwriting a file changes its id. `name` is the stable field, which is the reverse of what every other provider in this collection calls its identifier.

The live probe found this file's error envelope claiming a "status" string field that neither a bad-credential nor a bucket-existence failure actually carries -- the older v1 JSON API this Recipe emulates uses a plainer `{code, message, errors[]}` shape than the newer convention some other Google APIs in this collection use. It also found that an anonymous request against a bucket name that does not exist gets told so directly (404), while the identical request with an invented token fails on the credential first (401) -- not modelled as a general rule, since this Recipe's own fixture bucket has real seeded objects and importing "check existence first" here would misrepresent it as empty.

## Sources

- Documentation: https://cloud.google.com/storage/docs/json_api/v1
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googlecloudstorage     # run it
cauldron verify googlecloudstorage -v # check every claim
```
