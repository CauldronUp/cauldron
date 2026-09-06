# EasyPost

Emulates the EasyPost API (v2), for local development and tests.

**19 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An EasyPost rate isn't a price you can look up later -- it's a value handed to you from a specific shipment's response, addressable nowhere on its own (there's no GET for a rate id), so code that stores a rate id and buys with it an hour later is buying with something that no longer exists. Every object also carries a mode field saying test or production, and a test label is a picture of a label that scans as nothing -- easy to read past because both modes answer 200 with what looks like a real postage_label and tracking code.

An address can also be created successfully and still be undeliverable: verification is a nested object with its own success flag, so a 201 only means the record was stored, not that anyone actually lives there. And a failure's error.message is sometimes a plain string and sometimes an array of field errors, so code that renders it straight into a sentence works fine until the first validation failure, which is the one time anyone actually reads it.

Rates don't expire here and buying with a rate from the wrong shipment isn't refused -- both are real EasyPost failures this format can't check because it addresses everything by id alone.

## Sources

- Documentation: https://docs.easypost.com/docs/shipments
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve easypost     # run it
cauldron verify easypost -v # check every claim
```
