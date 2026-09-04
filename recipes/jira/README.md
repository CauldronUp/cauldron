# Jira

Emulates the Jira API (v3), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

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
