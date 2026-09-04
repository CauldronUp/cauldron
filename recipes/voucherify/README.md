# Voucherify

Emulates the Voucherify API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A failed redemption is a stored object, not an error. A `Redemption` with `result: FAILURE` carries its own `failure_code` and `failure_message`, and it's listed, counted, and returned by the ordinary redemptions endpoint like any successful one -- the default view of "how many times was this code redeemed" includes the times it wasn't. The redeem call that produces one also answers 200: codes that were refused come back inside two separate arrays in the body, `inapplicable_redeemables` and `skipped_redeemables`, so checking `response.ok` or only one of the two arrays reports a discount the customer never got.

A `Redemption` also carries two verdicts that aren't synonyms -- `result` (SUCCESS/FAILURE) and `status` (SUCCEEDED/FAILED/ROLLED_BACK) -- so a rolled-back redemption still reads `result: SUCCESS`, and `result === 'SUCCEEDED'` is false forever.

## Sources

- Documentation: https://docs.voucherify.io/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve voucherify     # run it
cauldron verify voucherify -v # check every claim
```
