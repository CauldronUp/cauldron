# Jira

Emulates the Jira API (v3), for local development and tests.

**20 conformance cases, 7 checked against the live API.**

Struck live on 2026-09-05 against `https://your-domain.atlassian.net` -- Atlassian's own documentation placeholder, and a real, provisioned Jira Cloud site rather than a bare wildcard. Two things it answered contradicted this file outright: the retired `/rest/api/3/search` endpoint's message was invented text, and the real one is longer and cites a changelog number; and the field list (`/rest/api/3/field`) had inherited the same `issues`-wrapped envelope the search endpoint uses and was never checked -- it is a bare array on the real API, confirmed on twenty-eight built-in fields. A third finding was an omission rather than an error: `isLast`, the field marking a search's last page, was declared absent, and a real empty search answered `{"issues":[],"isLast":true}`. All three are fixed. Two other claims matched byte for byte: the 401 body for a route that always requires identity, and the 404 body for an issue that does not exist.

## What this Recipe found

Jira's API is hard not because it's large but because almost every value in it is configurable by an administrator, and the parts that look stable often aren't. A custom field is `customfield_10023` on the wire, with the number assigned per instance in creation order -- the same logical field has a different id in staging than in production, so every integration ends up building its own field-name mapping at runtime, and the day it caches the wrong one is the day a story point estimate gets written into a field called Severity.

A status name is likewise whatever an administrator typed -- `"In Review"` isn't a value the API defines, and code branching on it breaks the moment someone renames a column. The stable field is `statusCategory.key`, which has exactly three values. Creating an issue also returns only an id, a key, and a URL -- none of them the issue -- so a test that asserts against the create response is asserting against almost nothing, and has to fetch the issue separately to see any of its fields.

The old search endpoint doesn't return 404 when you hit it, it returns 410 Gone -- `/rest/api/3/search` was removed in favor of `/rest/api/3/search/jql`, and the replacement pages by opaque token with no `startAt` and no `total` at all. A pagination UI built around a result count can't be rebuilt against the new endpoint, because that number simply isn't available anymore, at any price.

## Sources

- Documentation: https://developer.atlassian.com/cloud/jira/platform/rest/v3/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve jira     # run it
cauldron verify jira -v # check every claim
```
