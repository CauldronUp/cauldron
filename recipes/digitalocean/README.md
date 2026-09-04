# DigitalOcean

Emulates the DigitalOcean API (v2), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Creating a DigitalOcean droplet is asynchronous and the response says so plainly if you read it: 202, status "new", empty networks -- the IP address doesn't exist yet, so an integration that creates a droplet and immediately reads its address gets nothing, and only discovers the timing mismatch in production. The work itself is tracked by a separate action resource with three outcomes (in-progress, completed, errored), and polling until the status merely stops being in-progress -- without checking which of the other two it landed on -- treats a failed provision as a working machine.

The error object's code field is confusingly named id, so a client reading response.id expecting a resource identifier instead gets an error slug like "not_found", a string that looks entirely reasonable until it's used as a lookup key.

## Sources

- Documentation: https://docs.digitalocean.com/reference/api/api-reference/
- Machine-readable description: https://raw.githubusercontent.com/digitalocean/openapi/main/specification/DigitalOcean-public.v2.yaml, last checked 2026-08-30
  `cauldron drift digitalocean` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve digitalocean     # run it
cauldron verify digitalocean -v # check every claim
```
