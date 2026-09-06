# Intercom

Emulates the Intercom API (2.11), for local development and tests.

**19 conformance cases, 2 checked against the live API.**

Struck live against api.intercom.io on 2026-09-05, no account and no key. The auth failure this file declared was wrong: it named a bad token's error as `token_unauthorized` / "Not authorized to access resource", and the case meant to prove it sent no Authorization header at all, so nothing here had ever actually asked Intercom what either failure looks like. Both are fixed now -- a missing credential answers `missing_authorization` / "No authorization was provided", a present-and-wrong one answers `unauthorized` / "Access Token Invalid" -- and split under `auth.absent_error` / `auth.rejected_error` because they are genuinely two different sentences.

Also found and deliberately not modelled: a path with no route, and `/contacts` addressed with `PATCH`, both answer a 404 without ever checking the credential -- but so do `PUT /contacts` and `DELETE /contacts`, which is the opposite of what a routing-before-auth rule would predict, checked the same way. Intercom's real router disagrees with itself by verb on paths this Recipe does not model any of; the header says so rather than picking a side to force into a case.

## What writing this Recipe changed

Its paging state is nested, and modelling that exposed a real bug in the
runtime: declared constants were silently overwriting computed values.

## Sources

- Documentation: https://developers.intercom.com/docs/references/rest-api
- Machine-readable description: https://raw.githubusercontent.com/intercom/Intercom-OpenAPI/main/descriptions/2.11/api.intercom.io.yaml, last checked 2026-08-30
  `cauldron drift intercom` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve intercom     # run it
cauldron verify intercom -v # check every claim
```
