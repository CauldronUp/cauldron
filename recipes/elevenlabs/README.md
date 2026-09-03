# ElevenLabs

Emulates the ElevenLabs API (v1), for local development and tests.

**11 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## Sources

- Documentation: https://elevenlabs.io/docs/api-reference/introduction
- Machine-readable description: https://api.elevenlabs.io/openapi.json, last checked 2026-08-31
  `cauldron drift elevenlabs` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve elevenlabs     # run it
cauldron verify elevenlabs -v # check every claim
```
