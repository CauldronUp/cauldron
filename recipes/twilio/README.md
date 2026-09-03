# Twilio

Emulates the Twilio API (2010-04-01), for local development and tests.

**6 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

Its identifier is `sid`, not `id`, and its account identifier travels as the
Basic username -- which is why the runtime compared only the username for years,
and why Mailgun later found that this was wrong for every provider that puts its
key in the password.

An early round also found that single-object routes matched nothing at all here,
so every fetch of one record answered 404.

## Sources

- Documentation: https://www.twilio.com/docs/usage/api
- Machine-readable description: https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_api_v2010.json, last checked 2026-08-30
  `cauldron drift twilio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve twilio     # run it
cauldron verify twilio -v # check every claim
```
