# Supabase

Emulates the Supabase API (v1), for local development and tests.

**8 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://api.supabase.com/api/v1-json
- Machine-readable description: https://supabase.com/openapi.json, last checked 2026-08-31
  `cauldron drift supabase` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve supabase     # run it
cauldron verify supabase -v # check every claim
```
