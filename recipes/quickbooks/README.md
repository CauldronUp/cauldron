# QuickBooks

Emulates the QuickBooks API (v3), for local development and tests.

**16 conformance cases, 2 checked against the live API.**

Everything past the credential check still cites documentation rather than an observation, because reaching it needs a real company. The credential check itself was verified directly against quickbooks.api.intuit.com on 2026-09-05.

## What this Recipe found

With no Accept header at all, QuickBooks answers XML, not JSON, for every failure -- checked with no credential, a fictitious bearer, an unrouted path, and a wrong method, all four landing on the identical authentication failure. Adding `Accept: application/json` changes more than the wire format: the envelope this file had modelled as `Fault.Error`, read straight from the XML element names without being tested against the JSON body, is actually bottom-cased on the wire -- `fault.error`, `message`, `code`, all lowercase, with `queryResponse` (also lowercase) sitting right there as a sibling field in the same body. Whether that same rule reaches deeper field names like `Customer.DisplayName` is now an open question this file flags rather than guesses at, since confirming it needs a real company. The failure's own `code` is also a string of digits ("3200"), not the number this file had been sending.

Failures arrive under `fault.error` as an array, and the HTTP status can still be 200 for a batch request, so code reading a top-level error message finds nothing there either. `SyncToken`, the optimistic-concurrency check on every record, is a string of digits -- sending a stale one is refused, and sending none at all silently overwrites whatever changed since the record was last read.

## Sources

- Documentation: https://developer.intuit.com/app/developer/qbo/docs/api/accounting/most-commonly-used/account
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve quickbooks     # run it
cauldron verify quickbooks -v # check every claim
```
