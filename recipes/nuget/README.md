# NuGet

Emulates the NuGet API (v3), for local development and tests.

**14 conformance cases, 12 checked against the live API on 2026-08-26.** The 2 unchecked ones are the paging cases: they send the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

It is the one where the same package answers two
different ways at once. The service index advertises the registration resource
under four versions of its type pointing at three base URLs, and a client picks
by asking for a version of the *type*: resolve the unversioned name and you get
`registration5-semver1`, which for `Microsoft.Extensions.Logging` is 98
versions. The `semver2` base URL has 175. The seventy-seven missing are the
SemVer 2.0.0 ones -- every preview of the current major and the one before it --
because SemVer 1 cannot express a dotted prerelease identifier, and nothing in
the response says anything was left out. Whether the index arrives whole or in
pieces is decided by size rather than by endpoint: 84 versions come back
inlined with fragment links, 600 come back as pages with real URLs, and
`Microsoft.Extensions.Logging` is on both sides of that line at once. And asking
for a package that does not exist reports that a *blob* does not exist, in
Azure Storage's XML, from a JSON API.

## Sources

- Documentation: https://learn.microsoft.com/en-us/nuget/api/registration-base-url-resource
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve nuget     # run it
cauldron verify nuget -v # check every claim
```
