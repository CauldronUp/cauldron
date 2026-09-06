# Akamai

Emulates the Akamai API (v1), for local development and tests.

**6 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Akamai has no shared, generic API hostname at all -- unlike every other CDN in this collection (Bunny, Cloudflare, Fastly, KeyCDN), its management API lives at a customer-specific host provisioned only after an account exists, so nothing in this Recipe could be probed live and every case is sourced from Akamai's own published documentation.

The interesting shape is that "saved" and "serving" are visibly different questions, in Akamai's own worked examples: a property carries latestVersion (newest draft), stagingVersion and productionVersion as three separate numbers, and a version's own record carries independent productionStatus and stagingStatus enum fields -- so the exact same version can be live on staging and not production at once. The activation record's own status enum goes further, with ZONE_1/ZONE_2/ZONE_3 as a real documented in-between "progressive activation" state that neither Fastly's flat boolean pair nor KeyCDN's missing field has any equivalent for.

EdgeGrid, the auth scheme, packs five parameters (client token, access token, timestamp, nonce, HMAC signature) into one Authorization header -- and Cauldron does not verify the signature itself, so a request with a correctly shaped but wrongly signed header is accepted. What each malformation of that header should return could not be established from Akamai's own documentation at all, which is recorded as an honest unknown rather than guessed at.

## Sources

- Documentation: https://techdocs.akamai.com/property-mgr/reference/property-manager-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve akamai     # run it
cauldron verify akamai -v # check every claim
```
