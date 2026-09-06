# ConfigCat

Emulates the ConfigCat API (v1), for local development and tests.

**14 conformance cases, 1 checked against the live API.**

Struck live against api.configcat.com on 2026-09-05, on the exact route this Recipe models: an unauthorized request carries no body at all, byte for byte the same whether the header is absent or a made-up Basic credential. This file had claimed a JSON body, {"message":"Unauthorized"}; ConfigCat sends none of it.

## What this Recipe found

ConfigCat serves the same flag's value at /v1/settings/{key}/value and /v2/settings/{key}/value, and the two paths answer with different documents: v1 calls the field value, v2 calls it defaultValue, and the targeting rules are two arrays in one version and one array in the other. Nothing in either path names its own vocabulary except the version digit, so a client that upgrades the URL without upgrading the reader gets undefined back from a 200.

The same path segment also addresses two different identifier spaces at once -- settingKeyOrId accepts either a human-chosen string key or an int32 settingId, so a flag whose key happens to be "42" is addressed by the same URL as setting number 42. And "can I save this" has four separate answers in v2 (readOnly, approveRequired, canBypassApproval, reasonRequired) but only one in v1, so the same flag looks freely editable on the old path and gated on the new one.

Only six of the API's seventy-six operations take the X-CONFIGCAT-SDKKEY header at all -- exactly the value-reading endpoints -- so the call that actually reads a flag needs a credential the rest of the API has no use for.

## Sources

- Documentation: https://api.configcat.com/docs/
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve configcat     # run it
cauldron verify configcat -v # check every claim
```
