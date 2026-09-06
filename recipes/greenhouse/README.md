# Greenhouse

Emulates the Greenhouse API (v1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource cases cite documentation rather than an observation on a real organisation, because a test candidate would be a fictional person in a real recruiter's pipeline. The refusal cases were struck live, unauthenticated, against harvest.greenhouse.io.

## What this Recipe found

A candidate is not an application. `GET /candidates` returns people, but the hiring decision -- status, stage, outcome -- lives on the application underneath, and one person can hold several applications at once, to different jobs, at different stages. "The candidate's status" isn't a real thing, and code that reads one field gets whichever application happened to come back first.

Rejection is not deletion: a rejected application still shows up in a plain list and still sits in the candidate's applications array, so a pipeline count that doesn't filter on status includes everyone ever turned down. A prospect -- someone sourced but never actually applied -- is a candidate with `is_prospect: true` and an application carrying `prospect: true`, no job attached until someone converts it, so an unfiltered funnel counts sourcing leads as applicants.

Writes also need an `On-Behalf-Of` header naming a Greenhouse user id, while reads don't need it at all -- an integration whose tests only read meets this requirement for the first time in production, and the refusal names the header rather than the missing permission.

The live probe found this file's declared authentication error was already exactly right: Greenhouse checks the credential before it looks at the path or the method at all, so a missing key, a wrong key, an unrouted path, and a wrong method on a real path all answer the identical 401.

## Sources

- Documentation: https://developers.greenhouse.io/harvest.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve greenhouse     # run it
cauldron verify greenhouse -v # check every claim
```
