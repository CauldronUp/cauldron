# razorpay

Emulates the razorpay API (v1), for local development and tests.

**13 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://razorpay.com/docs/api/
- Machine-readable description: https://razorpay.com/openapi.json, last checked 2026-08-31
  `cauldron drift razorpay` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve razorpay     # run it
cauldron verify razorpay -v # check every claim
```
