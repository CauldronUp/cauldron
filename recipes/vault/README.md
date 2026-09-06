# HashiCorp Vault

Emulates the HashiCorp Vault API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

The same read gives you the secret directly, or a box containing the secret, depending on a decision made when the Vault mount was created. On KV v1, `response.data` is the secret. On KV v2, `response.data.data` is the secret, and `response.data` is a wrapper holding it beside a metadata block -- same client, same call, one extra layer, and nothing in the URL says which version you're talking to. `secret = response.data` is correct against one mount and returns `{data, metadata}` against the next, without throwing.

A write also never echoes back what was written -- the v2 response is pure metadata (created_time, version, and so on), and the v1 equivalent answers 204 with no body at all. And "not deleted" is two separate fields: `deletion_time` is an empty string, not null, and `destroyed` is `false` -- a soft-deleted version has a deletion_time with `destroyed` still false, and code checking `!= null` gets it wrong.

## Sources

- Documentation: https://developer.hashicorp.com/vault/api-docs/secret/kv/kv-v2
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vault     # run it
cauldron verify vault -v # check every claim
```
