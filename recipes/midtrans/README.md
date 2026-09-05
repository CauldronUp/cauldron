# midtrans

Emulates the midtrans API (v2), for local development and tests.

**14 conformance cases, 2 checked against the live API.**

Struck live 2026-09-05 against api.sandbox.midtrans.com (api.midtrans.com itself timed out from this network; the sandbox host answered normally). This file declared one authentication_error for every failure, "Unauthorized transaction, please check client or server key" -- the real API sends two different sentences, "Operation is not allowed due to unauthorized payload." for a missing credential and "Unknown Merchant server_key/id" for one nobody issued, and this file's own existing case for the second one had invented text that passed only because the emulator agreed with itself. Fixed below.

## What this Recipe found

A Midtrans payment is only safe when two separate fields agree: `transaction_status` says what happened to the money and `fraud_status` says what Midtrans thinks of it, and a transaction can be `capture` and `challenge` at the same time -- charged, but held pending a human fraud review. Code that reads only `transaction_status` and sees a success word ships the goods before the review resolves, and this never shows up in a happy-path test because a test card is never challenged.

There are also two status codes that are not the same thing: the HTTP status and a `status_code` string inside the body carrying Midtrans's own meaning, so a 200 with `status_code: "404"` is a real, valid response that a client branching on HTTP status alone reads as a successful empty lookup. `gross_amount` is a decimal string in a currency that has no cents (`"10000.00"` is ten thousand rupiah), so a client that divides by a hundred out of habit reports every amount as a hundredth of itself. `settlement` does not follow the same path every time either -- a card payment goes capture then settlement, but a bank transfer goes straight from pending to settlement, skipping capture entirely, so code that waits for capture waits forever on half its payments.

The notification signature, the entire security boundary for webhooks, is not computed or verified here on purpose: faking verification that always passes would teach the opposite lesson to the one worth learning.

## Sources

- Documentation: https://docs.midtrans.com/reference/backend-integration
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve midtrans     # run it
cauldron verify midtrans -v # check every claim
```
