# kratos

Emulates the kratos API (v1), for local development and tests.

**25 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Kratos inverts the usual login API: instead of posting credentials to an endpoint you hardcoded, you fetch a flow, and the flow hands back the URL to submit to, the method, and the full list of fields to render -- their names, types, values, and which are required. A client that posts `{email, password}` to a path it wrote down works against one deployment and silently breaks on the next, because the next one added a field or turned on a second factor, and none of that was a code change on the server.

An expired flow answers 410, and Kratos's own docs are explicit that this isn't a retry situation -- a new flow has to be started, and whatever the person had typed is gone. The CSRF token also travels as an ordinary field in that same list of nodes, indistinguishable from the ones a person actually fills in, and omitting it produces a 403 that says a security violation was detected without saying which one. And `identity.state` does nothing at all -- the specification says so in as many words -- yet it still has an `active`/`inactive` enum, so code that "deactivates" a user by setting it succeeds, reads back the new value, and has changed nothing about whether they can log in.

## Sources

- Documentation: https://www.ory.sh/docs/kratos/reference/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve kratos     # run it
cauldron verify kratos -v # check every claim
```
