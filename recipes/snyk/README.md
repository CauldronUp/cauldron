# Snyk

Emulates the Snyk API (v1), for local development and tests.

**11 conformance cases, 4 checked against the live API on 2026-09-05.**

Snyk has no sandbox, so the vulnerability and issue cases still cite documentation. What an unauthenticated request gets back needs no account, and checking it live found this Recipe's own credential case wrong.

## What this Recipe found

A vulnerability with no fix is the interesting one: severity says how bad it is, but `isUpgradable` and `isPatchable` say whether anything can be done about it today, and a gate that fails on severity alone blocks every build until an upstream maintainer ships a release. An ignored issue is also still returned -- Snyk records the ignore as a decision rather than removing the finding, so a naive count re-reports something a human already assessed.

## What checking it live found

Every unauthenticated v1 request -- no token, a garbage one, or a Bearer-prefixed one instead of Snyk's own `token` scheme, all byte-identical -- answers JSON:API shaped: `{"jsonapi":{"version":"1.0"},"errors":[{"status":"401","details":"Unauthorized"}]}`, `application/vnd.api+json`. That is not the flat `{code, message, userMessage}` shape the rest of this file uses for errors that need an account to reach, and a case here had asserted the flat shape for the Bearer-prefixed request without ever having run it. It was wrong; fixed now, with the `userMessage` field moved off the shared envelope and onto the five errors that actually carry it, since the one failure this Recipe could check does not.

A path nothing declares answers a third way: 404, zero bytes, no `Content-Type`, with no credential sent -- so routing is resolved before the credential is judged at all, which this Recipe did not model either.

## Sources

- Documentation: https://docs.snyk.io/snyk-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve snyk     # run it
cauldron verify snyk -v # check every claim
```
