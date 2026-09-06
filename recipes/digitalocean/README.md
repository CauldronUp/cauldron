# DigitalOcean

Emulates the DigitalOcean API (v2), for local development and tests.

**25 conformance cases, 4 checked against the live API on 2026-09-05.**

The resource shapes cite documentation rather than an account nobody has; the refusal cases were struck live, unauthenticated, against api.digitalocean.com.

## What this Recipe found

Creating a DigitalOcean droplet is asynchronous and the response says so plainly if you read it: 202, status "new", empty networks -- the IP address doesn't exist yet, so an integration that creates a droplet and immediately reads its address gets nothing, and only discovers the timing mismatch in production. The work itself is tracked by a separate action resource with three outcomes (in-progress, completed, errored), and polling until the status merely stops being in-progress -- without checking which of the other two it landed on -- treats a failed provision as a working machine.

The error object's code field is confusingly named id, so a client reading response.id expecting a resource identifier instead gets an error slug like "not_found", a string that looks entirely reasonable until it's used as a lookup key.

The live probe found this file wrong about its own authentication error: it claimed a lowercase `unauthorized` code with a trailing period on the message, and the real API sends `Unauthorized`, capitalised, with no period -- the same 401 for a missing token, an invented one, and one with no Bearer prefix. It also found that an unmatched path answers 404 before DigitalOcean's application ever looks at authentication, a different sentence from the 401 every real route gives an unauthenticated caller, and that a wrong method on a real path answers a genuine 405 with zero bytes and no Allow header at all.

## Sources

- Documentation: https://docs.digitalocean.com/reference/api/api-reference/
- Machine-readable description: https://raw.githubusercontent.com/digitalocean/openapi/main/specification/DigitalOcean-public.v2.yaml, last checked 2026-08-30
  `cauldron drift digitalocean` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve digitalocean     # run it
cauldron verify digitalocean -v # check every claim
```
