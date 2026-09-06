# SonarCloud

Emulates the SonarCloud API (v1), for local development and tests.

**12 conformance cases, 2 checked against the live API on 2026-09-05.**

Most of this Recipe was read out of SonarCloud's own self-description at `api/webservices/list`. The credential shape needed no organisation to check, and checking it live found this Recipe's own auth model wrong.

## What this Recipe found

The quality gate has three outcomes and one of them means nobody set a gate. `project_status` can answer OK, WARN, ERROR, or NONE, and NONE means no gate is attached to the project at all -- so `if (status !== 'OK') fail()` fails a build that was never gated, and `if (status === 'ERROR') fail()` passes it. Every number inside the gate result is also a string: threshold, actual value, and comparator are all text, so evaluating a condition means parsing two numbers out of strings and reading an operator out of a third field.

There's also no verb but GET and POST across the API's 156 actions -- deleting a comment is `POST api/issues/delete_comment` -- so the path is the verb, and the HTTP method only says whether anything changed.

## What checking it live found

`project_status` and `ce/task` answer a request with no credential at all -- anonymous access to public data, which SonarCloud genuinely allows -- and only refuse one that presents something wrong, empty-bodied, with a 401. This Recipe had one `authentication_error` for every failure; the wrong-credential case is now its own error, the two GET routes are marked `public: when-absent`, and the two POST mutations (which really do refuse an absent credential) are unaffected. The not-found sentence also quotes the key that was actually asked for -- `"Component key 'cauldron-fixture' not found"` -- where this Recipe's copy had dropped it for a fixed sentence.

## Sources

- Documentation: https://sonarcloud.io/web_api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve sonarcloud     # run it
cauldron verify sonarcloud -v # check every claim
```
