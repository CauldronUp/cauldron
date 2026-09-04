# Column

Emulates the Column API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A Column notification of change is not a failure -- a C01 means the money arrived and the account number on file is wrong, and treating it as a return reverses a payment that actually worked, losing a customer's money for a week while the corrected data (which exists only in the notification) goes unread. Return codes themselves are three groups sharing one two-character shape: R01 may be retried, R02 may not because the account is closed, and R07 must never be retried because the customer revoked authorisation and retrying becomes a regulatory problem rather than a technical one.

The deadline for a return also isn't fixed -- an unauthorised return can arrive up to sixty days later, an administrative one within two business days, so how long a payment must stay reconcilable depends on the code. Amounts are integer minor units with direction in a separate field rather than the sign, so summing them without reading direction gives total volume rather than net.

No money actually moves here -- there's no transfer-creation endpoint at all -- and return-window deadlines are described in this header rather than enforced, since Cauldron has no notion of a banking day.

## Sources

- Documentation: https://column.com/docs/api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve column     # run it
cauldron verify column -v # check every claim
```
