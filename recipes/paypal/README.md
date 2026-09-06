# PayPal

Emulates the PayPal API (v2), for local development and tests.

**19 conformance cases, 2 checked against the live API.**

Everything that needs a sandbox account still cites documentation rather than an observation -- the Recipe's own header says why, and that reason is the finding as often as not. The two auth failures were checked directly against api-m.paypal.com, unauthenticated, on 2026-09-05.

## What this Recipe found

A missing Authorization header and a fictitious bearer token are not the same failure, checked live: the first answers `{"name":"AUTHENTICATION_FAILURE",...,"links":[...]}` and the second answers a bare OAuth2 shape, `{"error":"invalid_token","error_description":"..."}`, with no field in common between them. Code that learns PayPal's error shape from the easier case to write -- a missing header -- cannot read the one it will actually hit in production: a token that expired.

An order that comes back `APPROVED` has not been paid -- the buyer returned from PayPal, the order says APPROVED, and no money has moved until a separate capture call happens. An order left approved and never captured is a sale that never happened while looking, in every log, exactly like one that did; this is the single most expensive misunderstanding in the API. The `links` array is the actual flow control: a CREATED order carries an approve link and an APPROVED one does not, because approving is already done, so code that indexes `links[1]` instead of matching on `rel` reads a different link at each step and eventually sends a buyer somewhere meaningless.

Even a successful capture can still not be settled -- `status: PENDING` with a reason attached is a real HTTP 201 where the money is under review, and code that branches on the call having succeeded ships the goods anyway. Money is always a string paired with a currency code, and the number of decimal places depends on the currency (a yen amount has none), so parsing it as a float or multiplying by a hundred to get minor units is wrong for JPY by two orders of magnitude. What arrives in a completed capture also is not what was charged: gross, fee, and net are all reported separately, and reconciling against gross alone leaves a shortfall every time.

Capturing, authorizing, voiding and refunding are deliberately not modelled -- the same position Mercury, Bill.com and Deel take elsewhere in this collection, on the reasoning that an emulator making an irreversible operation easy to exercise invites the mistake it exists to warn against.

## Sources

- Documentation: https://developer.paypal.com/docs/api/orders/v2/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve paypal     # run it
cauldron verify paypal -v # check every claim
```
