# Gumroad

Emulates the Gumroad API (v2), for local development and tests.

**13 conformance cases, 3 checked against the live API on 2026-09-05.**

The resource cases cite Gumroad's own open-source controllers rather than an observation on a real seller; the refusal cases were struck live, unauthenticated, against api.gumroad.com.

## What this Recipe found

Verifying a license spends one. `POST /v2/licenses/verify` reads like a question and is actually an update: Gumroad's own controller defaults `increment_uses_count` to true unless the caller explicitly sends `false`, then increments the use count on every call. An app that checks its license on every launch burns a use on every launch, and a five-seat key runs out on the fifth boot of the same computer -- and the strongest evidence this actually bites people is that Gumroad ships a separate endpoint, `decrement_uses_count`, whose only job is undoing it.

Failures come in three incompatible shapes from the same API: license endpoints send 404/500 with `{success: false, message}`, a bad page key sends 400 with a completely different `{status, error}` shape, and everything else -- Gumroad's own source admits this is a known bug it means to fix -- answers 200 with `{success: false, message}`. A client can't unwrap all three the same way, and a missing `product_id` answers 500 on purpose, so a client that retries on 5xx retries forever.

Verifying also doesn't tell you which license it verified -- the response merges in only the use count, no key, no id -- and the paginated `next_page_url` Gumroad hands back has the access token stripped out of it, so following the link Gumroad gives you is itself unauthenticated.

The live probe found that "unauthenticated" itself answers nothing like this file assumed: a missing or invalid `access_token` gets zero bytes, not the JSON this file declared, with the real refusal carried in a `WWW-Authenticate` header instead -- Gumroad's API auth is Doorkeeper (OAuth2), answering the RFC 6750 way. An unrouted path and a wrong method on a real path both answer Gumroad's own branded marketing-site 404 page before authentication is ever consulted.

## Sources

- Documentation: https://app.gumroad.com/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve gumroad     # run it
cauldron verify gumroad -v # check every claim
```
