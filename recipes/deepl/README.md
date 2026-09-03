# DeepL

Emulates the DeepL API (v2), for local development and tests.

**13 conformance cases, 9 checked against the live API on 2026-09-02.**

## What this Recipe found

**One wrong key gets two different sentences** -- the
usage endpoint names the key as invalid and the translate endpoint answers the
single word Forbidden, to the identical credential.

## Sources

- Documentation: https://developers.deepl.com/docs/getting-started/auth
- Machine-readable description: https://raw.githubusercontent.com/DeepLcom/openapi/main/openapi.yaml, last checked 2026-09-02
  `cauldron drift deepl` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve deepl     # run it
cauldron verify deepl -v # check every claim
```
