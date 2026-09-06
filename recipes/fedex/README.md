# FedEx

Emulates the FedEx API (v1), for local development and tests.

**12 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What this Recipe found

An unknown tracking number comes back as a 200. FedEx lets you track up to thirty parcels in one call, so a failure can't be the status of the whole response -- it's an error object buried two arrays deep, at `output.completeTrackResults[0].trackResults[0]`, indistinguishable in shape from a real result unless something actually looks inside it. Almost every client hardcodes that `[0][0]` path and stays right only until a duplicated tracking number returns more than one match.

A handful of fields hide a similar trap. A delivery date is found by scanning a `dateAndTimes` array for a `type` of `ACTUAL_DELIVERY`, and a parcel that hasn't arrived yet has an `ESTIMATED_DELIVERY` entry sitting in the identical shape in the same array. Status comes three separate ways -- a machine code, a finer derived code, and a locale-dependent prose string that reads the most like an answer of the three. Weight arrives as a string with its unit attached, so adding two weights concatenates them.

The access token expires in an hour with no refresh modelled here; FedEx really does expire it, and this Recipe accepts its fixture token forever, so an integration that fetches a token once at boot works for an hour after every deploy and nothing here can catch that gap.

## Sources

- Documentation: https://developer.fedex.com/api/en-us/catalog/track/v1/docs.html
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve fedex     # run it
cauldron verify fedex -v # check every claim
```
