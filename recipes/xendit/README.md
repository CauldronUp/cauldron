# xendit

Emulates the xendit API (v2), for local development and tests.

**20 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

A virtual account answers 201 with a working-looking bank account number, and for a few minutes after creation that number doesn't actually work -- the bank hasn't registered it yet, and a customer who pays immediately is told by their own bank that the account doesn't exist, someplace no integration's logs will ever see. Nothing in the response says this; the `status` field reads `PENDING`, and most integrations treat a 201 as done.

There's also no polling endpoint worth having: Xendit's model is that a webhook announces a payment as a separate object, and fetching the virtual account back tells you nothing about whether money has arrived.

## Sources

- Documentation: https://developers.xendit.co/api-reference/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xendit     # run it
cauldron verify xendit -v # check every claim
```
