# Voucherify

Emulates the Voucherify API (v1), for local development and tests.

**16 conformance cases, 4 checked against the live API on 2026-09-05.**

The redemption and voucher cases still cite documentation, since a real redemption needs a real project. The credential shape needed no project at all, and checking it live found this Recipe's own error wrong.

## What this Recipe found

A failed redemption is a stored object, not an error. A `Redemption` with `result: FAILURE` carries its own `failure_code` and `failure_message`, and it's listed, counted, and returned by the ordinary redemptions endpoint like any successful one -- the default view of "how many times was this code redeemed" includes the times it wasn't. The redeem call that produces one also answers 200: codes that were refused come back inside two separate arrays in the body, `inapplicable_redeemables` and `skipped_redeemables`, so checking `response.ok` or only one of the two arrays reports a discount the customer never got.

A `Redemption` also carries two verdicts that aren't synonyms -- `result` (SUCCESS/FAILURE) and `status` (SUCCEEDED/FAILED/ROLLED_BACK) -- so a rolled-back redemption still reads `result: SUCCESS`, and `result === 'SUCCEEDED'` is false forever.

## What checking it live found

No credential, half a credential (an existing case already sent exactly this), and a present, wrong `X-App-Id`/`X-App-Token` pair are all byte-identical: key `"unauthorized"`, message `"Unauthorized"`, no `details` field. Not the `"authentication_error"`/`"Invalid credentials."` this Recipe had invented. A path nothing declares answers a third way, with no `key` field at all, resolved before any credential is judged.

## Sources

- Documentation: https://docs.voucherify.io/reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve voucherify     # run it
cauldron verify voucherify -v # check every claim
```
