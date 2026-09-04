# SonarCloud

Emulates the SonarCloud API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The quality gate has three outcomes and one of them means nobody set a gate. `project_status` can answer OK, WARN, ERROR, or NONE, and NONE means no gate is attached to the project at all -- so `if (status !== 'OK') fail()` fails a build that was never gated, and `if (status === 'ERROR') fail()` passes it. Every number inside the gate result is also a string: threshold, actual value, and comparator are all text, so evaluating a condition means parsing two numbers out of strings and reading an operator out of a third field.

There's also no verb but GET and POST across the API's 156 actions -- deleting a comment is `POST api/issues/delete_comment` -- so the path is the verb, and the HTTP method only says whether anything changed.

## Sources

- Documentation: https://sonarcloud.io/web_api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sonarcloud     # run it
cauldron verify sonarcloud -v # check every claim
```
