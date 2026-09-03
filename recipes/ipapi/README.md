# ip-api

Emulates the ip-api API (ipapi), for local development and tests.

**8 conformance cases, all of them checked against the live API on 2026-08-30.**

## What this Recipe found

**TLS is a paid feature and the message saying so is
delivered over TLS.** `https://ip-api.com/json/1.1.1.1` answers 403 with
`{"status": "fail", "message": "SSL unavailable for this endpoint, order a key
at https://members.ip-api.com/"}`. The connection is accepted, the certificate
is valid, the handshake completes, and the body says SSL is unavailable --
everything TLS is for worked, and what is missing is the entitlement.

**And every other failure is a 200.** A malformed address and a string that is
not an address at all both answer HTTP 200 with `{"status": "fail", "message":
"invalid query"}`, so the status line is not a signal and a client checking
`response.ok` gets true for both. A format the API does not recognise is the
marketing website: `/nosuchformat/` answers 200 with eleven kilobytes of HTML,
meta keywords and all. The rate limit is `X-Rl` and `X-Ttl`, two abbreviations
and neither standard. And the autonomous system is a number and a name in one
string, `"AS13335 Cloudflare, Inc."`, comma included, beside an `isp` with no
full stop and an `org` that is a different organisation again.

## Sources

- Documentation: https://ip-api.com/docs
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve ipapi     # run it
cauldron verify ipapi -v # check every claim
```
