# RubyGems

Emulates the RubyGems API (v1), for local development and tests.

**11 conformance cases, 9 checked against the live API on 2026-08-25.** The unchecked one is the paging case: it sends the parameter names this Recipe declares, read from the provider's own description rather than struck against it.

## What this Recipe found

The same reason again: no OpenAPI, no credential for
reads, every claim checked live. The flag that decides whether publishing needs
two-factor authentication is the string `"true"`, because a gemspec's metadata
is a map of strings to strings -- so `if (metadata.rubygems_mfa_required)` is
true for `"true"` and also for `"false"`. A gem nobody has published answers a
404 whose entire body is `This rubygem could not be found.`, in plain text, from
a path ending in `.json`. `authors` is prose and `licenses` is a list, so
nokogiri's four authors arrive as one comma-joined string. A dependency
requirement is a sentence -- `{"name": "actioncable", "requirements": "=
8.1.3.1"}`. And the gemspec's URIs are promoted to the top level and left in
place, agreeing three times out of four: `homepage_uri` is at the top and absent
from `metadata`, so neither copy is authoritative.

## Sources

- Documentation: https://guides.rubygems.org/rubygems-org-api/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve rubygems     # run it
cauldron verify rubygems -v # check every claim
```
