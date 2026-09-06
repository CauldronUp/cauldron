# Xero

Emulates the Xero API (2.0), for local development and tests.

**22 conformance cases, 3 checked against the live API on 2026-09-05.**

The invoice and contact cases still cite documentation, since a real organisation needs a real connection. The credential and routing shapes needed no organisation at all, and checking them live found this Recipe's own error shape and sentence wrong.

## What writing this Recipe changed

It answers a request for one invoice with a list of one, so code that expects a
single object finds an array where it did not look for one.

## What checking it live found

No credential at all, a present garbage bearer token, and a path nothing declares all answer byte-identical: RFC 7807 problem details -- `{Type, Title, Status, Detail, Instance, Extensions}`, PascalCase -- from the gateway in front of the Accounting API, not this Recipe's own `{Type, ErrorNumber, Message}` shape. The sentence lives in `Detail` and reads `"AuthenticationUnsuccessful"`, not `"Token is invalid or has expired"`. All three are checked before routing: an unrouted path gets the identical gateway answer rather than a 404.

## Sources

- Documentation: https://developer.xero.com/documentation/api/accounting/overview
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xero     # run it
cauldron verify xero -v # check every claim
```
