# Jira Service Management

Emulates the Jira Service Management API (servicedeskapi), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://developer.atlassian.com/cloud/jira/service-desk/rest/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve jiraservicemanagement     # run it
cauldron verify jiraservicemanagement -v # check every claim
```
