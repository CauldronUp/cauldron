# PubNub

Emulates the PubNub API (v2), for local development and tests.

**14 conformance cases, all of them checked against the live API on 2026-08-31.**

## What this Recipe found

**The cursor the whole API runs on is a bare number
too large for JavaScript.** `/time/0` answers `[17882284253083560]`: seventeen
digits, unquoted, where `Number.MAX_SAFE_INTEGER` is sixteen. It arrives already
rounded, the corruption happens at `JSON.parse` before any code runs, and it is
the value every history call and every replay is addressed by. Steam, written
the same day, quotes its 64-bit ids for exactly this reason. PubNub also
disagrees with itself about what an error is: across four responses from one
host, `error` is a string, a boolean, and a number.

## Sources

- Documentation: https://www.pubnub.com/docs/sdks/rest-api
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve pubnub     # run it
cauldron verify pubnub -v # check every claim
```
