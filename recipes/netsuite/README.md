# NetSuite

Emulates the NetSuite API (v1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

NetSuite's REST API does not have a shared host to test against at all -- it is addressed at `{accountId}.suitetalk.api.netsuite.com`, a hostname that only exists once an account provisions it, and this was confirmed rather than assumed: three fabricated account ids all came back NXDOMAIN, and the bare wildcard host has no DNS record either. Every route, error body and field in this Recipe is modelled from Oracle's published documentation rather than observed, with no verified date on any case, because there was no live host reachable to check one against.

The documented distinction between deleting and voiding a transaction is sharper here than the other accounting providers in this group: deleting really does remove the record entirely (NetSuite keeps only an audit-trail note that it happened), while voiding zeroes out the total and every line item but leaves the record in place, explicitly the preferred move because it preserves the audit trail. Deletion is refused only when something else references the record, like a paid invoice with an associated payment, not because the transaction itself has posted -- an open, unpaid invoice with nothing pointing at it is not documented as blocked at all.

Only get and delete are routed; no worked example of a NetSuite list response, create, or update body could be found in Oracle's own documentation, so those are not modelled rather than guessed at. The error envelope is RFC 7807-shaped, which none of the other accounting Recipes in this collection share -- Xero is flat PascalCase and QuickBooks nests under `Fault.Error`.

## Sources

- Documentation: https://docs.oracle.com/en/cloud/saas/netsuite/ns-online-help/section_156570709583.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve netsuite     # run it
cauldron verify netsuite -v # check every claim
```
