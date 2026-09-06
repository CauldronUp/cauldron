# Mailjet

Emulates the Mailjet API (v3), for local development and tests.

**11 conformance cases, 3 checked against the live API on 2026-09-01.**

## What this Recipe found

Mailjet authenticates with HTTP Basic carrying two real secrets, an API Key as the username and a Secret Key as the password, and probing it live with no header, a malformed header, and a well-formed header carrying junk credentials returned the exact same 401 with an empty body all three times -- missing, malformed and wrong are indistinguishable here because Mailjet itself does not distinguish them. That response is `text/html`, not JSON, so calling `.json()` on it throws.

The more interesting behaviour cannot be reproduced as ordinary validation: Mailjet's Send API packages multiple messages in one request and can answer 200 while some of them failed and some succeeded, mixing `"Status": "success"` and `"Status": "error"` entries in a single response. Compare SendGrid or Postmark, which always report one outcome per call -- Mailjet is the only email provider in this collection where one 200 can mean two different things for two different recipients. Cauldron has no way to loop over a caller's array and decide per-entry success, so the mixed-outcome case here is a fixed, documented response rather than something derived from what was actually sent, and it is marked as documented, not verified against a live account.

## Sources

- Documentation: https://dev.mailjet.com/email/reference/overview/introduction/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mailjet     # run it
cauldron verify mailjet -v # check every claim
```
