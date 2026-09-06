# Google Drive

Emulates the Google Drive API (v3), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A file listing can be quietly incomplete, and the only sign is a boolean nobody reads. `incompleteSearch` is documented as true when some search results might be missing -- typically when searching across multiple shared drives -- and everything else in the response looks completely normal: the array is internally consistent, paging works, and a loop that runs until `nextPageToken` disappears terminates cleanly having seen an unknown fraction of what it asked for. It is not an error, not a retry signal, not an unusual status code, just a flag next to a complete-looking answer.

Drive is explicitly not a filesystem: a file's `name` isn't unique within a folder, so two files can share a name in one location, and "find the file called X" is a query that returns a list rather than a single result. `parents` is an array that can only ever hold one entry -- multiple parents are documented as unsupported even though the field is still shaped like a list. And shared-drive items are invisible by default: `includeItemsFromAllDrives` and `supportsAllDrives` both default to false, so an integration built and tested against a personal account works perfectly and then silently returns a fraction of the results against an account that uses shared drives, with no error at all.

## Sources

- Documentation: https://developers.google.com/workspace/drive/api/reference/rest/v3
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve googledrive     # run it
cauldron verify googledrive -v # check every claim
```
