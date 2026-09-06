# Slack

Emulates the Slack API (2026-05), for local development and tests.

**21 conformance cases, 2 checked against the live API on 2026-09-05.**

Most of this Recipe cites documentation rather than a workspace, because a
real one needs an app and a bot token. The credential shape does not, and
checking it live found this Recipe's own case wrong.

## What writing this Recipe changed

Its failures arrive with HTTP 200 and its identifier lives in a query
parameter rather than a path. Writing this Recipe is what taught the router RPC
shapes, which several later Recipes depend on.

## What checking it live found

A case named "a bad token is also an HTTP 200" sent no `Authorization` header
at all and asserted `invalid_auth` -- the code a present-but-wrong token gets.
Struck live against `slack.com`, sending nothing gets `not_authed` instead: a
different code, still on a 200. Split into two cases now, with
`auth.absent_error` carrying the distinction.

Two more were observed and left as prose rather than cases, because modelling
them would reorder routing against auth for every route in the file: an
unrecognised method name (`/api/not.a.real.method`) answers its own
`unknown_method` code with no `Authorization` header sent at all, so an
unknown method is resolved before Slack asks who is calling -- the opposite
order from every declared method. And the HTTP verb is not checked: `DELETE`
against `/api/conversations.list`, unauthenticated, gets the same
`not_authed` a `GET` does rather than a 405.

## Sources

- Documentation: https://api.slack.com/web
- Machine-readable description: https://raw.githubusercontent.com/slackapi/slack-api-specs/master/web-api/slack_web_openapi_v2.json, last checked 2026-09-05
  `cauldron drift slack` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve slack     # run it
cauldron verify slack -v # check every claim
```
