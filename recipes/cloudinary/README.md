# Cloudinary

Emulates the Cloudinary API (v1_1), for local development and tests.

**10 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Cloudinary hands back two URLs for the same asset, url (http) and secure_url (https), always both present -- using the first on an https page gets blocked as mixed content by every modern browser, and it's the one listed first alphabetically and in most of Cloudinary's own examples, which makes it the most common Cloudinary bug. Deleting something that doesn't exist also answers HTTP 200 with {"result": "not found"} in the body rather than a 404, so code that checks only the status code treats a no-op as a successful deletion.

public_id carries no file extension -- the delivery URL needs a format appended separately, so a URL built from public_id alone 404s unless the transformation uses f_auto. And an asset pending moderation is fully uploaded, stored and billed, but not deliverable: it has a public_id and a secure_url that returns nothing, so an upload that reported success produces a broken image.

Cauldron doesn't transform images or verify upload signatures (the credential is checked instead of the signature algorithm), and which of the two "deleted"/"not found" outcomes a destroy call gets is something the emulator is told to produce rather than something it works out.

## Sources

- Documentation: https://cloudinary.com/documentation/image_upload_api_reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve cloudinary     # run it
cauldron verify cloudinary -v # check every claim
```
