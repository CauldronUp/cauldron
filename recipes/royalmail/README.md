# Royal Mail

Emulates the Royal Mail API (v1), for local development and tests.

**9 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://api.parcel.royalmail.com/
- Machine-readable description: https://api.parcel.royalmail.com/swagger/v1/swagger.json, last checked 2026-08-31
  `cauldron drift royalmail` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve royalmail     # run it
cauldron verify royalmail -v # check every claim
```
