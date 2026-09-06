# Lemon Squeezy

Emulates the Lemon Squeezy API (v1), for local development and tests.

**25 conformance cases, 1 checked against the live API.**

Struck live 2026-09-05 against api.lemonsqueezy.com, no account and no key. `/v1/orders` specifically sits behind a Cloudflare bot challenge, reproduced three times; every other modelled path answered a real 401 on the first try. Checked against `/v1/subscriptions` instead: the failure's `detail` matched exactly, and `title` did not -- the real word is "Unauthorized", not a second spelling of "Unauthenticated". A `{"jsonapi":{"version":"1.0"}}` wrapper on every error was also unmodelled. Both fixed.

## What this Recipe found

Most of what makes Lemon Squeezy awkward is shared with FastSpring elsewhere in this collection -- merchant-of-record tax, every amount doubled as an integer and a preformatted string, a cancelled subscription that keeps running until its period ends, test and live records sharing one API distinguished only by a boolean -- and this Recipe asserts all four rather than skip them just because another Recipe already does.

What's specific to Lemon Squeezy: a license key is the one credential in this API that somebody outside your company actually holds, and it carries an activation limit with a running count of machines -- the count is what runs out, not the subscription, which stays fully active while every seat is taken and nobody new can install anything. A partial refund is also not "refunded": `status` reads `partial_refund` while the `refunded` boolean stays `false`, both correctly, so a report keyed on the boolean counts the order at its full original total after a fifth of it has already gone back. And totals arrive in two different currencies on the same record, not just two formats -- a euro sale carries both `total` (euro cents) and `total_usd` (dollar cents), with nothing in either field name saying which one your books actually use.

## Sources

- Documentation: https://docs.lemonsqueezy.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve lemonsqueezy     # run it
cauldron verify lemonsqueezy -v # check every claim
```
