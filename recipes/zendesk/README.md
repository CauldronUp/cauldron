# Zendesk

Emulates the Zendesk API (v2), for local development and tests.

**20 conformance cases, 3 checked against the live API on 2026-09-05.**

Most helpdesks need an account, but Zendesk runs its own support desk on Zendesk, and support.zendesk.com needed no account of its own to check an unauthenticated request against. Doing so found this Recipe's own credential model incomplete.

## What this Recipe found

The error envelope reuses `error` as a bare string code rather than an object -- a failure is `{"error": "RecordNotFound", "description": "..."}`, so code that reads `error.message` finds a string where it expected an object and typically throws on the property access instead of reporting the actual failure. A collection listing also reports its true total separately from the page length, so a pagination UI can't be built from the page alone.

## What checking it live found

No credential at all and a present, wrong bearer token are different failures, not one: absence gets the short `{"error":"Couldn't authenticate you"}` this Recipe already had -- an existing case had sent exactly this while calling it "a bad token" -- but a present, wrong token gets `{"error":"invalid_token","error_description":"..."}` instead, a different code and a field this Recipe's `message_field` never reached. A path nothing declares is a third shape again, `{"error":"InvalidEndpoint","description":"Not found"}`, resolved before either credential question.

## Sources

- Documentation: https://developer.zendesk.com/api-reference/ticketing/introduction
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve zendesk     # run it
cauldron verify zendesk -v # check every claim
```
