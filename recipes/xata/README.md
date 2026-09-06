# Xata

Emulates the Xata API (1.0), for local development and tests.

**5 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

Checked live: `api.xata.io` and every documented Xata host return NXDOMAIN or a bare Cloudflare "origin DNS error" -- the REST database API this Recipe was meant to describe has been retired, and the company's own GitHub organisation confirms the pivot to a different product. What survives is the generated TypeScript client, still on GitHub, and its own generated types are more precise than the vanished docs would have been: creating a record is typed as a three-way union, and two of the three shapes carry none of the data the caller actually submitted back -- not by omission, but as one of three documented possibilities from the vendor's own code generator.

Every record's metadata also comes in two competing naming conventions on the same object -- dotted (`xata.version`) and flattened with underscores (`xata_version`) -- read as evidence of the same schema migration that killed the product this Recipe was asked to model.

## Sources

- Documentation: https://xata.io/docs/api-reference
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve xata     # run it
cauldron verify xata -v # check every claim
```
