# ElevenLabs

Emulates the ElevenLabs API (v1), for local development and tests.

**16 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The API key is documented as optional on every endpoint that needs one. ElevenLabs' own OpenAPI document declares `xi-api-key` as an ordinary header parameter, `required: false`, on 383 separate operations, with the real rule buried only in prose -- "required by most endpoints." A client generated straight from the schema has no type-level signal that the key is needed anywhere, and 385 of 387 operations document only a 422 validation failure, so there's barely a documented branch for a bad key either. This Recipe enforces the key and answers 401 regardless of what the schema promises.

Two more things are easy to get wrong. Text to Speech and over a dozen other operations answer with raw audio or a zip file, not JSON, so calling `.json()` on the endpoint most integrations actually want throws. And a history item's cost is a subtraction rather than a field: `character_count_change_from` and `character_count_change_to` are two running-total readings, and the price of one generation is the difference between them.

## Sources

- Documentation: https://elevenlabs.io/docs/api-reference/introduction
- Machine-readable description: https://api.elevenlabs.io/openapi.json, last checked 2026-08-31
  `cauldron drift elevenlabs` compares it against what this Recipe claims.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve elevenlabs     # run it
cauldron verify elevenlabs -v # check every claim
```
