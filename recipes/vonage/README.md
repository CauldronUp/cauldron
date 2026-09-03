# Vonage

Emulates the Vonage API (1), for local development and tests.

**7 conformance cases, none checked against a live API.**

Every case here cites documentation rather than an observation. The Recipe's own header says why, and that reason is the finding as often as not.

## What writing this Recipe changed

It reports a successful send with the string `"0"`, which is what made the
conformance checker compare a scalar's kind as well as its value. Before that,
`"0"` and `0` were indistinguishable and this Recipe's case passed whichever the
emulator sent.

## Sources

- Documentation: https://developer.vonage.com/en/api/sms
- No machine-readable description is recorded. The Recipe's header says whether one exists and could not be read, or does not exist.

Every case cites where it came from. The Recipe itself, [`recipe.yaml`](recipe.yaml), carries the full notes: what was probed, what was deliberately not modelled, and why.

```bash
cauldron serve vonage     # run it
cauldron verify vonage -v # check every claim
```
