# Mailchimp

Emulates the Mailchimp API (3.0), for local development and tests.

**15 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Mailchimp has no sandbox. Every account is a real account, every audience a real audience, and a test that adds a subscriber can trigger a real welcome automation to a real address; the free tier also counts test contacts against the plan limit, so a CI suite that adds a member per run eventually starts failing for billing reasons rather than code reasons.

Failures follow RFC 7807 (`type`, `title`, `status`, `detail`) with no `code` field anywhere, so code written for a provider-specific error shape has nothing to read, and the human-readable text is under `detail`, not `message`. A member is addressed by the MD5 hash of their lower-cased email rather than an opaque id, so an integration that stores the id Mailchimp returned works fine until the person changes their email address.

## Sources

- Documentation: https://mailchimp.com/developer/marketing/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mailchimp     # run it
cauldron verify mailchimp -v # check every claim
```
