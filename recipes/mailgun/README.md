# Mailgun

Emulates the Mailgun API (v3), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

This Recipe found the worst bug in the runtime's history. Basic authentication
had only ever compared the username -- correct for Twilio, where the account
identifier is the username, and wrong here, where the username is the constant
`api` and the key is the password.

So a request carrying the right username and a completely wrong key returned
200. The failure path a test most wants to exercise could not be reached at
all.

## Sources

- Documentation: https://documentation.mailgun.com/docs/mailgun/api-reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve mailgun     # run it
cauldron verify mailgun -v # check every claim
```
