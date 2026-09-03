# JFrog Artifactory

Emulates the JFrog Artifactory API (7), for local development and tests.

**7 conformance cases, all of them checked against the live API on 2026-09-01.**

## What this Recipe found

It **names its own internal token type in a refusal**
-- "Props Authentication Token not found" -- and sends a file's `size` as a JSON
string.

## Sources

- Documentation: https://jfrog.com/help/r/jfrog-rest-apis/artifactory-rest-apis
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve artifactory     # run it
cauldron verify artifactory -v # check every claim
```
