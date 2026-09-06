# Twilio

Emulates the Twilio API (2010-04-01), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

Twilio's sandbox sends real messages to real phones, so the message-sending cases still cite documentation. The credential and routing shapes needed no account at all, and checking them live found this Recipe's own message incomplete.

## What writing this Recipe changed

Its identifier is `sid`, not `id`, and its account identifier travels as the
Basic username -- which is why the runtime compared only the username for years,
and why Mailgun later found that this was wrong for every provider that puts its
key in the password.

An early round also found that single-object routes matched nothing at all here,
so every fetch of one record answered 404.

## What checking it live found

No credential at all and a present, wrong Basic one share code `20003` and disagree on the sentence -- `"No credentials provided"` against `"invalid username"` -- and this Recipe had only ever claimed the second for both. A path nothing declares answers XML, not the JSON every declared route sends: Twilio negotiates format by the `.json` suffix on a real route, and there is no suffix to read on a path that was never a route at all -- checked before the credential is judged, not after.

## Sources

- Documentation: https://www.twilio.com/docs/usage/api
- Machine-readable description: https://raw.githubusercontent.com/twilio/twilio-oai/main/spec/json/twilio_api_v2010.json, last checked 2026-08-30
  `cauldron drift twilio` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve twilio     # run it
cauldron verify twilio -v # check every claim
```
